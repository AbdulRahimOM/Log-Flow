package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

type SupabaseStorage struct {
	SupabaseURL string
	SupabaseKey string
	BucketName  string
}

func NewSupabaseStorage(url, key, bucketName string) Storage {
	return &SupabaseStorage{
		SupabaseURL: url,
		SupabaseKey: key,
		BucketName:  bucketName,
	}
}

func (sb *SupabaseStorage) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	client := resty.New()

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Upload file to Supabase Storage
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+sb.SupabaseKey).
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(file).
		Put(fmt.Sprintf("%s/storage/v1/object/%s/%s", sb.SupabaseURL, sb.BucketName, fileHeader.Filename))

	if err != nil {
		return "", fmt.Errorf("failed to upload: %v", err)
	}

	// Check for errors
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("failed to upload. Status not 200. resp: %s", resp.String())
	}

	// Return the file URL
	// fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", sb.SupabaseURL, sb.BucketName, fileHeader.Filename)
	fileURL, err := sb.generateSignedURL(fileHeader.Filename, 3600)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	return fileURL, nil
}

func (sb *SupabaseStorage) generateSignedURL(filePath string, expirySeconds int) (string, error) {
	client := resty.New()

	// Make API request
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+sb.SupabaseKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{"expiresIn": expirySeconds}).
		Post(fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", sb.SupabaseURL, sb.BucketName, filePath))

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed: %s", resp.String())
	}

	var signedURLResponse struct {
		SignedURL string `json:"signedURL"`
	}

	err = json.Unmarshal(resp.Body(), &signedURLResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	signedURL := fmt.Sprintf("%s/storage/v1%s", sb.SupabaseURL, signedURLResponse.SignedURL)

	return signedURL, nil
}

// DownloadFile fetches a file from Supabase Storage and saves it locally
func (sb *SupabaseStorage) DownloadFile(fileName string, savePath string) error {
	// Generate a signed URL (valid for 1 hour)
	fileURL, err := sb.generateSignedURL(fileName, 3600)
	if err != nil {
		return fmt.Errorf("failed to generate signed URL: %v", err)
	}

	resp, err := resty.New().R().Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to download file. Status code: %d", resp.StatusCode())
	}

	outFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.RawBody())
	if err != nil {
		return fmt.Errorf("failed to save file locally: %v", err)
	}

	fmt.Printf("✅ File downloaded successfully: %s\n", savePath)
	return nil
}

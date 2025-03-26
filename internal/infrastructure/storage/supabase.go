package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	supabaseUrlExpiryInSeconds = 360000
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

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("failed to upload. Status not 200. resp: %s", resp.String())
	}

	// fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", sb.SupabaseURL, sb.BucketName, fileHeader.Filename)	//Permanent URL
	fileURL, err := sb.generateSignedURL(fileHeader.Filename, supabaseUrlExpiryInSeconds)
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

func (sb *SupabaseStorage) StreamLogs(fileURL string) (io.ReadCloser, error) {
	resp, err := resty.New().R().SetDoNotParseResponse(true).Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file stream: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file. Status: %d", resp.StatusCode())
	}

	return resp.RawBody(), nil
}

func (sb *SupabaseStorage) GetFileSize(fileURL string) (int64, error) {
	resp, err := resty.New().R().Head(fileURL)
	if err != nil || resp.RawResponse == nil {
		log.Errorf("Warning: Unable to determine total file size. %v", err)
		return 0, err
	}
	return resp.RawResponse.ContentLength, nil
}

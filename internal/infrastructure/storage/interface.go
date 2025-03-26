package storage

import (
	"io"
	"mime/multipart"
)

type Storage interface {
	UploadFile(fileHeader *multipart.FileHeader) (url string, err error)
	StreamLogs(fileURL string) (io.ReadCloser, error)
	GetFileSize(fileURL string) (int64, error)
}

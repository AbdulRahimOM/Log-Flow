package storage

import "mime/multipart"

type Storage interface {
	UploadFile(fileHeader *multipart.FileHeader) (url string, err error)
}

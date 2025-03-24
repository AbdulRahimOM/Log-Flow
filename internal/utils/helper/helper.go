package helper

import (
	"os"
	"strings"
)

func IsValidLogFile(filename string) bool {
	return strings.HasSuffix(filename, ".log")
}

func EnsureUploadsDir(uploadPath string) error {
	// uploadPath := "./uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		// Create the directory with proper permissions
		err := os.Mkdir(uploadPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

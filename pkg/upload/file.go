package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var headers = map[string]string{
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type FileUpload struct {
	req          *http.Request
	field        string
	maxSize      int64
	allowedTypes []string
}

func NewUpload(r *http.Request, field string) *FileUpload {
	return &FileUpload{req: r, field: field}
}

func (fu *FileUpload) SetMaxSize(max int64) {
	fu.maxSize = max
}

func (fu *FileUpload) SetAllowedTypes(types ...string) {
	fu.allowedTypes = types
}

func imageType(value string, options []string) bool {
	for _, option := range options {
		h, ok := headers[option]

		if !ok {
			return false
		}

		if h == value {
			return true
		}
	}

	return false
}

func (fu *FileUpload) validateFile(h *multipart.FileHeader) error {
	if len(fu.allowedTypes) > 0 {
		if !imageType(h.Header.Get("Content-Type"), fu.allowedTypes) {
			return fmt.Errorf("File is not in: %v", fu.allowedTypes)
		}
	}

	if fu.maxSize > 0 {
		if h.Size > fu.maxSize {
			return fmt.Errorf("File is too big: %d", fu.maxSize)
		}
	}

	return nil
}

func (fu *FileUpload) ValidateAndUpload(dir string) (*FileInfo, error) {
	file, handler, err := fu.req.FormFile(fu.field)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err := fu.validateFile(handler); err != nil {
		return nil, err
	}

	fileExt := strings.Split(handler.Filename, ".")[1]

	// Create a temporary file within our temp-images directory that follows a particular naming pattern
	tempFile, err := os.CreateTemp(dir, fmt.Sprintf("%s-*.%s", "image", fileExt))
	if err != nil {
		return nil, err
	}

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	if fileBytes, err := io.ReadAll(file); err != nil {
		return nil, err
	} else {
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		return &FileInfo{Name: tempFile.Name(), Size: handler.Size}, nil
	}
}

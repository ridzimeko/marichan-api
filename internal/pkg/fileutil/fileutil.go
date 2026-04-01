package fileutil

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func SaveUploadedFile(file *multipart.FileHeader, dir string) (string, string, error) {
	if err := EnsureDir(dir); err != nil {
		return "", "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	name := uuid.NewString() + ext
	dst := filepath.Join(dir, name)

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	if err != nil {
		return "", "", err
	}

	return dst, name, nil
}

func BuildFileUrl(baseUrl, folder, filename string) string {
	folder = strings.Trim(folder, "/")
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(baseUrl, "/"), folder, filename)
}

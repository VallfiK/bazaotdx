// internal/utils/file_utils.go
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateGuestFolder(rootPath, fullName string) (string, error) {
	// Нормализуем имя для использования в пути
	normalizedName := strings.ReplaceAll(strings.TrimSpace(fullName), " ", "_")
	normalizedName = strings.ReplaceAll(normalizedName, "/", "-")

	// Формат: ГГГГ-ММ-ДД_ЧЧ-ММ-СС_ФИО
	folderName := fmt.Sprintf("%s_%s",
		time.Now().Format("2006-01-02_15-04-05"),
		normalizedName,
	)

	path := filepath.Join(rootPath, folderName)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("failed to create folder: %w", err)
	}
	return path, nil
}

func CopyDocument(sourcePath, destDir string) (string, error) {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	destPath := filepath.Join(destDir, filepath.Base(sourcePath))
	dstFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return "", err
	}
	return destPath, nil
}

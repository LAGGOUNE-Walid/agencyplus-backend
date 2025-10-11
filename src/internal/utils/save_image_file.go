package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrNotAnImage   = errors.New("uploaded file is not a valid image")
	ErrFileTooLarge = errors.New("file size exceeds limit")
)

func SaveFile(file multipart.File, header *multipart.FileHeader, baseFolder string, maxSize int64) (string, error) {
	defer file.Close()

	if header.Size > maxSize {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))

	// Create subfolder by year-month
	subfolder := time.Now().Format("2006-01") // e.g. "2025-05"
	destFolder := filepath.Join(baseFolder, subfolder)
	if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create folder: %w", err)
	}

	// Generate random 8-character prefix
	prefix, err := generateRandomString(8)
	if err != nil {
		return "", fmt.Errorf("failed to generate random prefix: %w", err)
	}

	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	filePath := filepath.Join(destFolder, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path: "2025-05/filename.jpg"
	relativePath := filepath.Join(subfolder, filename)
	return relativePath, nil
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func GeneratePDFThumbnail(pdfPath string, outputImagePath string) error {
	cmd := exec.Command(
		"pdftoppm",
		"-f", "1",
		"-l", "1",
		"-png",
		"-singlefile",
		pdfPath,
		outputImagePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}

	return nil
}

package service

import (
	"backend/internal/storage"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// StorageService defines the interface for file storage operations
type StorageService interface {
	// UploadFile uploads a file to MinIO and returns the URL
	UploadFile(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, error)

	// UploadImage uploads an image file with validation
	UploadImage(ctx context.Context, file io.Reader, filename string, size int64) (string, error)

	// DeleteFile deletes a file from storage
	DeleteFile(ctx context.Context, objectName string) error

	// GetFileURL generates a presigned URL for file access
	GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)

	// ListFiles lists files with a given prefix
	ListFiles(ctx context.Context, prefix string) ([]minio.ObjectInfo, error)
}

// storageService implements StorageService
type storageService struct {
	client  *storage.MinIOClient
	baseURL string
}

// NewStorageService creates a new storage service
func NewStorageService(client *storage.MinIOClient, baseURL string) StorageService {
	return &storageService{
		client:  client,
		baseURL: baseURL,
	}
}

// UploadFile uploads a file to MinIO storage
func (s *storageService) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, error) {
	// Generate unique object name
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	// Upload to MinIO
	_, err := s.client.GetClient().PutObject(ctx, s.client.GetBucket(), objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return the public URL
	return fmt.Sprintf("%s/%s/%s", s.baseURL, s.client.GetBucket(), objectName), nil
}

// UploadImage uploads an image with validation
func (s *storageService) UploadImage(ctx context.Context, file io.Reader, filename string, size int64) (string, error) {
	// Validate file size (max 10MB)
	const maxSize = 10 * 1024 * 1024
	if size > maxSize {
		return "", fmt.Errorf("file size exceeds 10MB limit")
	}

	// Validate file extension
	ext := filepath.Ext(filename)
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !validExts[ext] {
		return "", fmt.Errorf("invalid image format. Allowed: jpg, jpeg, png, gif, webp")
	}

	return s.UploadFile(ctx, file, filename, "image/"+ext[1:], size)
}

// DeleteFile deletes a file from storage
func (s *storageService) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.GetClient().RemoveObject(ctx, s.client.GetBucket(), objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFileURL generates a presigned URL for temporary access
func (s *storageService) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := s.client.GetClient().PresignedGetObject(ctx, s.client.GetBucket(), objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// ListFiles lists files with a given prefix
func (s *storageService) ListFiles(ctx context.Context, prefix string) ([]minio.ObjectInfo, error) {
	var files []minio.ObjectInfo

	for obj := range s.client.GetClient().ListObjects(ctx, s.client.GetBucket(), minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		files = append(files, obj)
	}

	return files, nil
}

// UploadResult represents the result of a file upload
type UploadResult struct {
	URL        string `json:"url"`
	ObjectName string `json:"object_name"`
	Size       int64  `json:"size"`
}

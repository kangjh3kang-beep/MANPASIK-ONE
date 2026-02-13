// Package storage는 S3 호환 오브젝트 스토리지 클라이언트를 제공합니다.
package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Client wraps MinIO client for S3-compatible storage.
type S3Client struct {
	client *minio.Client
	bucket string
}

// NewS3Client creates a new S3 client and ensures the bucket exists.
func NewS3Client(endpoint, accessKey, secretKey, bucket, region string, useSSL bool) (*S3Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("s3 client creation failed: %w", err)
	}

	s3 := &S3Client{client: client, bucket: bucket}

	// Ensure bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s3.ensureBucket(ctx); err != nil {
		return nil, err
	}

	return s3, nil
}

func (s *S3Client) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
	}
	return nil
}

// Upload stores a file in S3.
func (s *S3Client) Upload(ctx context.Context, path string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, path, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Download retrieves a file from S3.
func (s *S3Client) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Delete removes a file from S3.
func (s *S3Client) Delete(ctx context.Context, path string) error {
	return s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
}

// GetPresignedURL generates a pre-signed URL for temporary access.
func (s *S3Client) GetPresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.bucket, path, expiry, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// Exists checks if a file exists in S3.
func (s *S3Client) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, path, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Health checks if S3 is reachable.
func (s *S3Client) Health(ctx context.Context) error {
	_, err := s.client.BucketExists(ctx, s.bucket)
	return err
}

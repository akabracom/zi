// internal/storage/minio/minio.go
package minio

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioClient() *MinioClient {
	endpoint := "minio:9000"
	accessKey := "minioadmin"
	secretKey := "minioadmin123"
	bucket := "deposits"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("MinIO client: %v", err)
	}

	m := &MinioClient{Client: client, BucketName: bucket}
	m.EnsureBucket()
	return m
}

func (m *MinioClient) EnsureBucket() {
	ctx := context.Background()
	exists, err := m.Client.BucketExists(ctx, m.BucketName)
	if err != nil {
		log.Printf("Bucket check failed: %v", err)
		return
	}
	if !exists {
		if err := m.Client.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{}); err != nil {
			log.Printf("Bucket creation failed: %v", err)
		} else {
			log.Printf("Bucket '%s' created", m.BucketName)
		}
	}
}

func (m *MinioClient) UploadImage(objectName, filePath string) error {
	ctx := context.Background()
	_, err := m.Client.FPutObject(ctx, m.BucketName, objectName, filePath, minio.PutObjectOptions{})
	return err
}

func (m *MinioClient) RemoveImage(objectName string) error {
	ctx := context.Background()
	return m.Client.RemoveObject(ctx, m.BucketName, objectName, minio.RemoveObjectOptions{})
}
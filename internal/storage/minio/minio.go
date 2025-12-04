// internal/storage/minio/minio.go
package minio

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioClient() *MinioClient {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin123"
	}
	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "deposits"
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("MinIO client: %v", err)
	}

	return &MinioClient{
		Client:     client,
		BucketName: bucket,
	}
}

func (m *MinioClient) EnsureBucket() {
	ctx := context.Background()
	exists, err := m.Client.BucketExists(ctx, m.BucketName)
	if err != nil {
		log.Fatalf("Check bucket: %v", err)
	}
	if !exists {
		if err := m.Client.MakeBucket(ctx, m.BucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("Create bucket: %v", err)
		}
		log.Printf("Bucket '%s' created", m.BucketName)
	}
}

func (m *MinioClient) UploadLocalImages(dir string) error {
	ctx := context.Background()
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(dir, file.Name())
		objectName := "images/" + file.Name()

		_, err := m.Client.FPutObject(ctx, m.BucketName, objectName, path, minio.PutObjectOptions{})
		if err != nil {
			return err
		}
		log.Printf("Uploaded %s to MinIO", file.Name())
	}
	return nil
}
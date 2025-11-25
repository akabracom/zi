package minio

import (
    "context"
    "io"
    "log"
    "os"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
    Client     *minio.Client
    BucketName string
}

func NewMinioClient() (*MinioClient, error) {
    endpoint := os.Getenv("MINIO_ENDPOINT")
    if endpoint == "" {
        endpoint = "minio:9000"
    }
    accessKey := os.Getenv("MINIO_ACCESS_KEY")
    if accessKey == "" {
        accessKey = "minioadmin"
    }
    secretKey := os.Getenv("MINIO_SECRET_KEY")
    if secretKey == "" {
        secretKey = "minioadmin123"
    }
    bucketName := os.Getenv("MINIO_BUCKET")
    if bucketName == "" {
        bucketName = "deposits"
    }

    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: false,
    })
    if err != nil {
        return nil, err
    }

    ctx := context.Background()
    exists, err := client.BucketExists(ctx, bucketName)
    if err != nil {
        return nil, err
    }

    if !exists {
        err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
        if err != nil {
            return nil, err
        }
        log.Printf("Bucket '%s' created", bucketName)
    }

    log.Printf("MinIO connected: endpoint=%s, bucket=%s", endpoint, bucketName)
    return &MinioClient{
        Client:     client,
        BucketName: bucketName,
    }, nil
}

func (mc *MinioClient) UploadImageDirect(objectName string, reader io.Reader, size int64) error {
    _, err := mc.Client.PutObject(context.Background(), mc.BucketName, objectName, reader, size, minio.PutObjectOptions{
        ContentType: "image/png",
    })
    return err
}

func (mc *MinioClient) RemoveImage(objectName string) error {
    return mc.Client.RemoveObject(context.Background(), mc.BucketName, objectName, minio.RemoveObjectOptions{})
}

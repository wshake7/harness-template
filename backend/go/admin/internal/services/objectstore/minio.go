package objectstore

import (
	"admin/internal/config"
	"context"
	"errors"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioClient interface {
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
	PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	PresignedPutObject(ctx context.Context, bucketName string, objectName string, expires time.Duration) (*url.URL, error)
	PresignedGetObject(ctx context.Context, bucketName string, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error)
	StatObject(ctx context.Context, bucketName string, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error
}

type MinIOEngine struct {
	conf   config.MinIOStorageConfig
	client minioClient
}

func NewMinIOEngine(conf config.MinIOStorageConfig) (*MinIOEngine, error) {
	client, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Secure: conf.UseSSL,
		Region: conf.Region,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOEngine{
		conf:   conf,
		client: client,
	}, nil
}

func (m *MinIOEngine) Name() string {
	return "minio"
}

func (m *MinIOEngine) PutObject(ctx context.Context, input PutInput) (PutResult, error) {
	info, err := m.client.PutObject(ctx, input.Bucket, input.ObjectKey, input.Reader, input.Size, minio.PutObjectOptions{
		ContentType: input.ContentType,
	})
	if err != nil {
		return PutResult{}, err
	}

	return PutResult{
		Bucket:    info.Bucket,
		ObjectKey: info.Key,
		ETag:      info.ETag,
		Size:      info.Size,
	}, nil
}

func (m *MinIOEngine) PresignedGetObject(ctx context.Context, input PresignInput) (string, time.Time, error) {
	params := make(url.Values)
	if input.Disposition != "" {
		params.Set("response-content-disposition", input.Disposition)
	}

	u, err := m.client.PresignedGetObject(ctx, input.Bucket, input.ObjectKey, input.Expiry, params)
	if err != nil {
		return "", time.Time{}, err
	}

	return u.String(), time.Now().Add(input.Expiry), nil
}

func (m *MinIOEngine) PresignedPutObject(ctx context.Context, input PresignPutInput) (string, time.Time, error) {
	u, err := m.client.PresignedPutObject(ctx, input.Bucket, input.ObjectKey, input.Expiry)
	if err != nil {
		return "", time.Time{}, err
	}

	return u.String(), time.Now().Add(input.Expiry), nil
}

func (m *MinIOEngine) StatObject(ctx context.Context, bucket string, objectKey string) (ObjectInfo, error) {
	info, err := m.client.StatObject(ctx, bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return ObjectInfo{}, err
	}
	return ObjectInfo{
		Bucket:    bucket,
		ObjectKey: info.Key,
		Size:      info.Size,
	}, nil
}

func (m *MinIOEngine) GetObject(ctx context.Context, bucket string, objectKey string) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
}

func (m *MinIOEngine) RemoveObject(ctx context.Context, bucket string, objectKey string) error {
	err := m.client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{})
	if err == nil {
		return nil
	}

	var resp minio.ErrorResponse
	if errors.As(err, &resp) && (resp.Code == "NoSuchKey" || resp.Code == "NoSuchObject" || resp.Code == "NoSuchVersion") {
		return nil
	}
	return err
}

func (m *MinIOEngine) Health(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.conf.Bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if !m.conf.AutoCreateBucket {
		return errors.New("minio bucket does not exist: " + m.conf.Bucket)
	}
	return m.client.MakeBucket(ctx, m.conf.Bucket, minio.MakeBucketOptions{Region: m.conf.Region})
}

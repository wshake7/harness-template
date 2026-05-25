package objectstore

import (
	"admin/internal/config"
	"bytes"
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-common/viperc"

	"github.com/minio/minio-go/v7"
)

// loadTestMinIOConfig tries to load MinIO config from etc/config.local.yaml first,
// then falls back to etc/config.yaml. Returns the config and true if Storage is enabled.
func loadTestMinIOConfig(t *testing.T) (config.MinIOStorageConfig, bool) {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Skipf("cannot get working directory: %v", err)
	}

	// Walk up to find the admin project root (where etc/ lives).
	root := cwd
	for {
		if _, err := os.Stat(filepath.Join(root, "etc")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Skip("cannot find admin project root (etc/ not found)")
		}
		root = parent
	}

	// Try config.local.yaml first, then config.yaml.
	var cfg config.Config
	for _, name := range []string{"config.local.yaml", "config.yaml"} {
		path := filepath.Join(root, "etc", name)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		_, err := viperc.ParseFile(path, &cfg)
		if err != nil {
			t.Fatalf("parse config %s: %v", path, err)
		}
		break
	}

	if !cfg.Storage.Enabled {
		return config.MinIOStorageConfig{}, false
	}
	return cfg.Storage.MinIO, true
}

func TestMinIOEnginePutObjectPassesParameters(t *testing.T) {
	fake := &fakeMinIOClient{
		putObjectResult: minio.UploadInfo{
			Bucket: "admin-files",
			Key:    "uploads/demo.txt",
			ETag:   "etag",
			Size:   5,
		},
	}
	engine := &MinIOEngine{
		conf:   config.MinIOStorageConfig{Bucket: "admin-files"},
		client: fake,
	}

	result, err := engine.PutObject(context.Background(), PutInput{
		Bucket:      "admin-files",
		ObjectKey:   "uploads/demo.txt",
		Reader:      bytes.NewBufferString("hello"),
		Size:        5,
		ContentType: "text/plain",
	})
	if err != nil {
		t.Fatalf("put object failed: %v", err)
	}
	if fake.putBucket != "admin-files" || fake.putObjectKey != "uploads/demo.txt" {
		t.Fatalf("unexpected put target: bucket=%q key=%q", fake.putBucket, fake.putObjectKey)
	}
	if fake.putSize != 5 {
		t.Fatalf("expected size 5, got %d", fake.putSize)
	}
	if fake.putContentType != "text/plain" {
		t.Fatalf("expected content type text/plain, got %q", fake.putContentType)
	}
	if result.ETag != "etag" {
		t.Fatalf("expected etag, got %q", result.ETag)
	}
}

func TestMinIOEnginePresignedGetObjectSetsDisposition(t *testing.T) {
	fake := &fakeMinIOClient{
		presignedURL: mustParseURL(t, "https://minio.local/admin-files/uploads/demo.txt"),
	}
	engine := &MinIOEngine{
		conf:   config.MinIOStorageConfig{Bucket: "admin-files"},
		client: fake,
	}

	_, _, err := engine.PresignedGetObject(context.Background(), PresignInput{
		Bucket:      "admin-files",
		ObjectKey:   "uploads/demo.txt",
		Expiry:      time.Hour,
		Disposition: "attachment",
	})
	if err != nil {
		t.Fatalf("presigned get object failed: %v", err)
	}
	if got := fake.presignedParams.Get("response-content-disposition"); got != "attachment" {
		t.Fatalf("expected attachment disposition, got %q", got)
	}
}

func TestMinIOEngineRemoveObjectIgnoresMissingObject(t *testing.T) {
	engine := &MinIOEngine{
		conf: config.MinIOStorageConfig{Bucket: "admin-files"},
		client: &fakeMinIOClient{
			removeObjectErr: minio.ErrorResponse{Code: "NoSuchKey"},
		},
	}

	if err := engine.RemoveObject(context.Background(), "admin-files", "uploads/demo.txt"); err != nil {
		t.Fatalf("expected missing object to be ignored, got %v", err)
	}
}

func TestMinIOEngineHealthReturnsErrorWhenBucketMissingAndAutoCreateDisabled(t *testing.T) {
	engine := &MinIOEngine{
		conf: config.MinIOStorageConfig{
			Bucket:           "admin-files",
			AutoCreateBucket: false,
		},
		client: &fakeMinIOClient{
			bucketExists: false,
		},
	}

	err := engine.Health(context.Background())
	if err == nil {
		t.Fatal("expected bucket missing error")
	}
	if !strings.Contains(err.Error(), "admin-files") {
		t.Fatalf("expected bucket name in error, got %v", err)
	}
}

func TestMinIOEngineHealthAutoCreatesBucket(t *testing.T) {
	fake := &fakeMinIOClient{bucketExists: false}
	engine := &MinIOEngine{
		conf: config.MinIOStorageConfig{
			Bucket:           "admin-files",
			Region:           "us-east-1",
			AutoCreateBucket: true,
		},
		client: fake,
	}

	if err := engine.Health(context.Background()); err != nil {
		t.Fatalf("expected bucket to be auto created, got %v", err)
	}
	if fake.makeBucketName != "admin-files" {
		t.Fatalf("expected make bucket for admin-files, got %q", fake.makeBucketName)
	}
	if fake.makeBucketRegion != "us-east-1" {
		t.Fatalf("expected region us-east-1, got %q", fake.makeBucketRegion)
	}
}

type fakeMinIOClient struct {
	bucketExists    bool
	bucketExistsErr error

	makeBucketName   string
	makeBucketRegion string
	makeBucketErr    error

	putBucket       string
	putObjectKey    string
	putBody         []byte
	putSize         int64
	putContentType  string
	putObjectErr    error
	putObjectResult minio.UploadInfo

	presignedURL    *url.URL
	presignedErr    error
	presignedParams url.Values

	removeObjectErr error
}

func (f *fakeMinIOClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return f.bucketExists, f.bucketExistsErr
}

func (f *fakeMinIOClient) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	f.makeBucketName = bucketName
	f.makeBucketRegion = opts.Region
	return f.makeBucketErr
}

func (f *fakeMinIOClient) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	f.putBucket = bucketName
	f.putObjectKey = objectName
	f.putSize = objectSize
	f.putContentType = opts.ContentType
	if reader != nil {
		body, err := io.ReadAll(reader)
		if err != nil {
			return minio.UploadInfo{}, err
		}
		f.putBody = body
	}
	return f.putObjectResult, f.putObjectErr
}

func (f *fakeMinIOClient) PresignedGetObject(ctx context.Context, bucketName string, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	f.presignedParams = reqParams
	if f.presignedURL == nil {
		return nil, f.presignedErr
	}
	return f.presignedURL, f.presignedErr
}

func (f *fakeMinIOClient) RemoveObject(ctx context.Context, bucketName string, objectName string, opts minio.RemoveObjectOptions) error {
	return f.removeObjectErr
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()

	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse url failed: %v", err)
	}
	return parsed
}

// TestMinIOIntegration runs against a real MinIO instance using config from
// etc/config.local.yaml (or etc/config.yaml). It is skipped when Storage is
// not enabled or the config file is missing.
func TestMinIOIntegration(t *testing.T) {
	minioCfg, ok := loadTestMinIOConfig(t)
	if !ok {
		t.Skip("Storage not enabled in config, skipping integration test")
	}

	engine, err := NewMinIOEngine(minioCfg)
	if err != nil {
		t.Fatalf("new minio engine: %v", err)
	}

	ctx := context.Background()

	if err := engine.Health(ctx); err != nil {
		t.Fatalf("health check: %v", err)
	}

	objKey := "test/integration-test.txt"
	content := []byte("minio integration test")

	_, err = engine.PutObject(ctx, PutInput{
		Bucket:      minioCfg.Bucket,
		ObjectKey:   objKey,
		Reader:      bytes.NewReader(content),
		Size:        int64(len(content)),
		ContentType: "text/plain",
	})
	if err != nil {
		t.Fatalf("put object: %v", err)
	}
	t.Cleanup(func() {
		if err := engine.RemoveObject(context.Background(), minioCfg.Bucket, objKey); err != nil {
			t.Logf("cleanup remove object: %v", err)
		}
	})

	url, expiresAt, err := engine.PresignedGetObject(ctx, PresignInput{
		Bucket:      minioCfg.Bucket,
		ObjectKey:   objKey,
		Expiry:      5 * time.Minute,
		Disposition: "attachment",
	})
	if err != nil {
		t.Fatalf("presigned get object: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty presigned URL")
	}
	if expiresAt.Before(time.Now()) {
		t.Fatal("expiresAt should be in the future")
	}

	if err := engine.RemoveObject(ctx, minioCfg.Bucket, objKey); err != nil {
		t.Fatalf("remove object: %v", err)
	}
}

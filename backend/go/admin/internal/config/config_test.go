package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorageConfigDefaults(t *testing.T) {
	configPath := writeConfigFile(t, "Host: 127.0.0.1\n")

	cfg := Init(configPath)

	if cfg.Storage.Enabled {
		t.Fatal("expected storage disabled by default")
	}
	if cfg.Storage.Engine != "minio" {
		t.Fatalf("expected default engine minio, got %q", cfg.Storage.Engine)
	}
	if cfg.Storage.MaxUploadBytes != 10*1024*1024 {
		t.Fatalf("expected default max upload bytes 10485760, got %d", cfg.Storage.MaxUploadBytes)
	}
	if cfg.Storage.PresignedExpiresSeconds != 3600 {
		t.Fatalf("expected default presigned expiry 3600, got %d", cfg.Storage.PresignedExpiresSeconds)
	}
	if cfg.Storage.MinIO.Bucket != "admin-files" {
		t.Fatalf("expected default bucket admin-files, got %q", cfg.Storage.MinIO.Bucket)
	}
}

func TestStorageConfigParseMinIO(t *testing.T) {
	configPath := writeConfigFile(t, `
Storage:
  Enabled: true
  Engine: minio
  MinIO:
    Endpoint: minio.local:9000
    AccessKeyID: minioadmin
    SecretAccessKey: hidden-secret
    Bucket: uploads
    Region: cn-shanghai-1
    UseSSL: true
    AutoCreateBucket: false
`)

	cfg := Init(configPath)

	if !cfg.Storage.Enabled {
		t.Fatal("expected storage enabled")
	}
	if cfg.Storage.MinIO.Endpoint != "minio.local:9000" {
		t.Fatalf("expected endpoint minio.local:9000, got %q", cfg.Storage.MinIO.Endpoint)
	}
	if cfg.Storage.MinIO.AccessKeyID != "minioadmin" {
		t.Fatalf("expected access key minioadmin, got %q", cfg.Storage.MinIO.AccessKeyID)
	}
	if cfg.Storage.MinIO.SecretAccessKey != "hidden-secret" {
		t.Fatal("expected secret access key to be parsed")
	}
	if cfg.Storage.MinIO.Bucket != "uploads" {
		t.Fatalf("expected bucket uploads, got %q", cfg.Storage.MinIO.Bucket)
	}
	if cfg.Storage.MinIO.Region != "cn-shanghai-1" {
		t.Fatalf("expected region cn-shanghai-1, got %q", cfg.Storage.MinIO.Region)
	}
	if !cfg.Storage.MinIO.UseSSL {
		t.Fatal("expected use ssl true")
	}
	if cfg.Storage.MinIO.AutoCreateBucket {
		t.Fatal("expected auto create bucket false")
	}
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config file failed: %v", err)
	}
	return path
}

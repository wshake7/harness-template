package objectstore

import (
	"strings"
	"testing"
	"time"

	"admin/internal/config"
)

func TestNewEngineDisabled(t *testing.T) {
	engine, err := NewEngine(config.StorageConfig{})
	if err != nil {
		t.Fatalf("expected no error for disabled storage, got %v", err)
	}
	if engine == nil {
		t.Fatal("expected disabled engine")
	}
	if engine.Name() != "disabled" {
		t.Fatalf("expected disabled engine name, got %q", engine.Name())
	}
}

func TestClampPresignedExpiry(t *testing.T) {
	got := ClampPresignedExpiry(-1, 3600)
	if got != time.Hour {
		t.Fatalf("expected default expiry 1h, got %s", got)
	}

	got = ClampPresignedExpiry(999999, 3600)
	if got != 7*24*time.Hour {
		t.Fatalf("expected max expiry 7d, got %s", got)
	}
}

func TestBuildObjectKey(t *testing.T) {
	now := time.Date(2026, time.May, 24, 10, 0, 0, 0, time.UTC)
	key := BuildObjectKey("uploads", "../demo.PNG", now, "018f6e7e-4a77-7a2e-a91c-6cf114e2f52a")

	if !strings.HasPrefix(key, "uploads/2026/05/24/") {
		t.Fatalf("expected dated prefix, got %q", key)
	}
	if !strings.HasSuffix(key, ".png") {
		t.Fatalf("expected normalized extension .png, got %q", key)
	}
	if strings.Contains(key, "..") {
		t.Fatalf("expected sanitized key, got %q", key)
	}
}

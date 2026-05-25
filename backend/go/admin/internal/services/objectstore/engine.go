package objectstore

import (
	"admin/internal/config"
	"context"
	"errors"
	"io"
	"path"
	"strings"
	"time"
)

const maxPresignedExpiry = 7 * 24 * time.Hour

var ErrStorageDisabled = errors.New("object storage disabled")
var ErrObjectNotFound = errors.New("object storage object not found")

type Engine interface {
	Name() string
	PutObject(ctx context.Context, input PutInput) (PutResult, error)
	PresignedGetObject(ctx context.Context, input PresignInput) (string, time.Time, error)
	RemoveObject(ctx context.Context, bucket string, objectKey string) error
	Health(ctx context.Context) error
}

type PutInput struct {
	Bucket      string
	ObjectKey   string
	Reader      io.Reader
	Size        int64
	ContentType string
}

type PutResult struct {
	Bucket    string
	ObjectKey string
	ETag      string
	Size      int64
}

type PresignInput struct {
	Bucket      string
	ObjectKey   string
	Expiry      time.Duration
	Disposition string
}

type ObjectInfo struct {
	Bucket    string
	ObjectKey string
	Size      int64
}

type disabledEngine struct{}

func NewEngine(conf config.StorageConfig) (Engine, error) {
	if !conf.Enabled {
		return disabledEngine{}, nil
	}

	switch strings.ToLower(strings.TrimSpace(conf.Engine)) {
	case "", "minio":
		return NewMinIOEngine(conf.MinIO)
	default:
		return nil, errors.New("unsupported object storage engine: " + conf.Engine)
	}
}

func ClampPresignedExpiry(requestedSeconds int, defaultSeconds int) time.Duration {
	effective := requestedSeconds
	if effective <= 0 {
		effective = defaultSeconds
	}
	if effective <= 0 {
		effective = int(time.Hour / time.Second)
	}

	expiry := time.Duration(effective) * time.Second
	if expiry > maxPresignedExpiry {
		return maxPresignedExpiry
	}
	return expiry
}

func BuildObjectKey(prefix string, originalName string, now time.Time, objectID string) string {
	cleanPrefix := strings.Trim(strings.ReplaceAll(prefix, "\\", "/"), "/")
	fileName := path.Base(strings.ReplaceAll(originalName, "\\", "/"))
	ext := strings.ToLower(path.Ext(fileName))
	if ext == "." {
		ext = ""
	}

	segments := []string{}
	if cleanPrefix != "" {
		segments = append(segments, cleanPrefix)
	}
	segments = append(segments, now.Format("2006"), now.Format("01"), now.Format("02"))
	return path.Join(append(segments, objectID+ext)...)
}

func (disabledEngine) Name() string {
	return "disabled"
}

func (disabledEngine) PutObject(ctx context.Context, input PutInput) (PutResult, error) {
	return PutResult{}, ErrStorageDisabled
}

func (disabledEngine) PresignedGetObject(ctx context.Context, input PresignInput) (string, time.Time, error) {
	return "", time.Time{}, ErrStorageDisabled
}

func (disabledEngine) RemoveObject(ctx context.Context, bucket string, objectKey string) error {
	return ErrStorageDisabled
}

func (disabledEngine) Health(ctx context.Context) error {
	return nil
}

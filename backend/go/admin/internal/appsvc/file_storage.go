package appsvc

import (
	"admin/internal/services/orm/models"
	"context"
	"io"
	"time"
)

//go:generate mockgen -source=file_storage.go -destination=../mock/mock_file_storage.go -package=mock -typed

type FileStorage interface {
	Upload(ctx context.Context, input UploadInput) (*models.FileAsset, error)
	Detail(ctx context.Context, id uint64) (*models.FileAsset, error)
	PresignedURL(ctx context.Context, input PresignedURLInput) (*PresignedURLResult, error)
	Delete(ctx context.Context, id uint64, operatorID uint64) error
}

type UploadInput struct {
	OperatorID   uint64
	Reader       io.Reader
	Size         int64
	OriginalName string
	ContentType  string
	BizType      string
	BizID        string
	Metadata     string
	Remark       string
}

type PresignedURLInput struct {
	ID             uint64
	ExpiresSeconds int
	Disposition    string
}

type PresignedURLResult struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expiresAt"`
}

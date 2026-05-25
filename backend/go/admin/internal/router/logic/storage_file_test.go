package logic

import (
	"admin/internal/appsvc"
	"admin/internal/auth"
	"admin/internal/config"
	"admin/internal/fiberc/handler"
	"admin/internal/services/orm/models"
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func TestStorageFileHandlerUploadSuccess(t *testing.T) {
	storage := &stubFileStorage{
		uploadResult: &models.FileAsset{
			Engine:       "minio",
			Bucket:       "admin-files",
			ObjectKey:    "uploads/2026/05/24/demo.txt",
			OriginalName: "demo.txt",
			ContentType:  "text/plain",
			Size:         5,
			SHA256:       "sha",
		},
	}
	app := newStorageTestApp(t, storage, 1024)

	req := newMultipartRequest(t, map[string]string{
		"bizType":  "attachment",
		"metadata": `{"source":"test"}`,
	}, "file", "demo.txt", []byte("hello"))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("upload request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if storage.uploadInput == nil {
		t.Fatal("expected upload service to be called")
	}
	if storage.uploadInput.BizType != "attachment" {
		t.Fatalf("expected bizType attachment, got %q", storage.uploadInput.BizType)
	}
}

func TestStorageFileHandlerUploadRejectsEmptyFile(t *testing.T) {
	storage := &stubFileStorage{}
	app := newStorageTestApp(t, storage, 1024)

	req := newMultipartRequest(t, nil, "file", "empty.txt", []byte{})
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("upload request failed: %v", err)
	}
	defer resp.Body.Close()

	if storage.uploadInput != nil {
		t.Fatal("expected empty file not to call upload service")
	}
}

func TestStorageFileHandlerUploadRejectsTooLargeFile(t *testing.T) {
	storage := &stubFileStorage{}
	app := newStorageTestApp(t, storage, 4)

	req := newMultipartRequest(t, nil, "file", "large.txt", []byte("hello"))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("upload request failed: %v", err)
	}
	defer resp.Body.Close()

	if storage.uploadInput != nil {
		t.Fatal("expected oversized file not to call upload service")
	}
}

func TestStorageFileHandlerPresignedRejectsInvalidDisposition(t *testing.T) {
	h := NewStorageFileHandler(&stubFileStorage{}, 1024)
	_, err := h.Presigned(newTestCtx(t), &ReqStorageFilePresigned{
		ID:          1,
		Disposition: "bad",
	})
	if err == nil {
		t.Fatal("expected invalid disposition error")
	}
}

func newStorageTestApp(t *testing.T, storage appsvc.FileStorage, maxUploadBytes int64) *fiber.App {
	t.Helper()

	app := fiber.NewWithCustomCtx(func(app *fiber.App) fiber.CustomCtx {
		ctx := &handler.Ctx{
			DefaultCtx: *fiber.NewDefaultCtx(app),
			Config: &config.Config{
				Storage: config.StorageConfig{MaxUploadBytes: maxUploadBytes},
			},
			SessionInfo: &auth.SessionInfo{Id: 7},
		}
		ctx.SetLogger(zap.NewNop())
		return ctx
	})
	h := NewStorageFileHandler(storage, maxUploadBytes)
	app.Post("/upload", handler.CtxFunc(h.Upload))
	return app
}

func newMultipartRequest(t *testing.T, fields map[string]string, fileField string, fileName string, body []byte) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write field failed: %v", err)
		}
	}
	part, err := writer.CreateFormFile(fileField, fileName)
	if err != nil {
		t.Fatalf("create form file failed: %v", err)
	}
	if _, err = part.Write(body); err != nil {
		t.Fatalf("write form file failed: %v", err)
	}
	if err = writer.Close(); err != nil {
		t.Fatalf("close multipart writer failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

type stubFileStorage struct {
	uploadInput     *appsvc.UploadInput
	uploadResult    *models.FileAsset
	uploadErr       error
	presignedResult *appsvc.PresignedURLResult
	presignedErr    error
}

func (s *stubFileStorage) Upload(ctx context.Context, input appsvc.UploadInput) (*models.FileAsset, error) {
	copied := input
	if input.Reader != nil {
		data, err := io.ReadAll(input.Reader)
		if err != nil {
			return nil, err
		}
		copied.Reader = bytes.NewReader(data)
	}
	s.uploadInput = &copied
	if s.uploadResult != nil {
		return s.uploadResult, s.uploadErr
	}
	return &models.FileAsset{Engine: "minio"}, s.uploadErr
}

func (s *stubFileStorage) Detail(ctx context.Context, id uint64) (*models.FileAsset, error) {
	return &models.FileAsset{}, nil
}

func (s *stubFileStorage) PresignedURL(ctx context.Context, input appsvc.PresignedURLInput) (*appsvc.PresignedURLResult, error) {
	if s.presignedResult != nil || s.presignedErr != nil {
		return s.presignedResult, s.presignedErr
	}
	return &appsvc.PresignedURLResult{
		URL:       "https://minio.local/demo",
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (s *stubFileStorage) Delete(ctx context.Context, id uint64, operatorID uint64) error {
	return nil
}

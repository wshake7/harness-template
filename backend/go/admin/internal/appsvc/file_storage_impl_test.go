package appsvc

import (
	"admin/internal/services/objectstore"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFileStorageUploadSuccess(t *testing.T) {
	q := newFileStorageTestQuery(t)
	engine := &fakeStorageEngine{}
	svc := newTestFileStorage(t, q, engine)

	result, err := svc.Upload(context.Background(), UploadInput{
		OperatorID:   7,
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Size:         5,
		Reader:       bytes.NewBufferString("hello"),
		BizType:      "attachment",
		Metadata:     `{"source":"test"}`,
	})
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if engine.putInput == nil {
		t.Fatal("expected engine put object to be called")
	}
	if result.SHA256 != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Fatalf("unexpected sha256: %q", result.SHA256)
	}

	item, err := q.FileAsset.Where(q.FileAsset.ID.Eq(result.ID)).First()
	if err != nil {
		t.Fatalf("expected file asset stored, got %v", err)
	}
	if item.BizType != "attachment" {
		t.Fatalf("expected biz type attachment, got %q", item.BizType)
	}
	if item.Metadata["source"] != "test" {
		t.Fatalf("expected metadata source=test, got %#v", item.Metadata)
	}
}

func TestFileStorageUploadRejectsInvalidMetadata(t *testing.T) {
	q := newFileStorageTestQuery(t)
	engine := &fakeStorageEngine{}
	svc := newTestFileStorage(t, q, engine)

	_, err := svc.Upload(context.Background(), UploadInput{
		OperatorID:   7,
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Size:         5,
		Reader:       bytes.NewBufferString("hello"),
		Metadata:     "{bad json",
	})
	if err == nil {
		t.Fatal("expected invalid metadata error")
	}
	if engine.putInput != nil {
		t.Fatal("expected upload to fail before put object")
	}
}

func TestFileStorageUploadCompensatesWhenCreateFails(t *testing.T) {
	q := newFileStorageTestQuery(t)
	engine := &fakeStorageEngine{}
	svc := newTestFileStorage(t, q, engine)
	objectKey := objectstore.BuildObjectKey("uploads", "demo.txt", svc.(*fileStorageImpl).now(), "fixed-object-id")

	if err := q.FileAsset.Create(&models.FileAsset{
		Engine:       "minio",
		Bucket:       "admin-files",
		ObjectKey:    objectKey,
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Extension:    ".txt",
		Size:         5,
		SHA256:       "existing",
		Status:       models.FileAssetStatusActive,
		Metadata:     datatypes.JSONMap{},
	}); err != nil {
		t.Fatalf("seed file asset failed: %v", err)
	}

	_, err := svc.Upload(context.Background(), UploadInput{
		OperatorID:   7,
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Size:         5,
		Reader:       bytes.NewBufferString("hello"),
	})
	if err == nil {
		t.Fatal("expected duplicate create to fail")
	}
	if engine.removeCalls != 1 {
		t.Fatalf("expected compensation remove once, got %d", engine.removeCalls)
	}
}

func TestFileStoragePresignedURLCapsExpiry(t *testing.T) {
	q := newFileStorageTestQuery(t)
	engine := &fakeStorageEngine{}
	svc := newTestFileStorage(t, q, engine)

	if err := q.FileAsset.Create(newStoredFileAsset()); err != nil {
		t.Fatalf("seed file asset failed: %v", err)
	}

	result, err := svc.PresignedURL(context.Background(), PresignedURLInput{
		ID:             1,
		ExpiresSeconds: 999999,
		Disposition:    "attachment",
	})
	if err != nil {
		t.Fatalf("presigned url failed: %v", err)
	}
	if engine.presignInput == nil {
		t.Fatal("expected engine presign to be called")
	}
	if engine.presignInput.Expiry != 7*24*time.Hour {
		t.Fatalf("expected capped expiry 7d, got %s", engine.presignInput.Expiry)
	}
	if result.URL == "" {
		t.Fatal("expected non-empty presigned url")
	}
}

func TestFileStoragePrepareDirectUploadCreatesPendingAsset(t *testing.T) {
	q := newFileStorageTestQuery(t)
	engine := &fakeStorageEngine{}
	svc := newTestFileStorage(t, q, engine)

	result, err := svc.PrepareDirectUpload(context.Background(), PrepareDirectUploadInput{
		OperatorID:   7,
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Size:         5,
		BizType:      "knowledge-document",
		Metadata:     `{"source":"test"}`,
	})
	if err != nil {
		t.Fatalf("prepare direct upload failed: %v", err)
	}
	if engine.presignPutInput == nil {
		t.Fatal("expected engine presigned put to be called")
	}
	if result.Asset.Status != models.FileAssetStatusPendingUpload {
		t.Fatalf("expected pending upload status, got %q", result.Asset.Status)
	}
}

func TestFileStorageCompleteDirectUploadActivatesAfterStat(t *testing.T) {
	q := newFileStorageTestQuery(t)
	engine := &fakeStorageEngine{
		statInfo: objectstore.ObjectInfo{
			Bucket:    "admin-files",
			ObjectKey: "uploads/2026/05/24/fixed-object-id.txt",
			Size:      5,
		},
	}
	svc := newTestFileStorage(t, q, engine)

	if err := q.FileAsset.Create(&models.FileAsset{
		Engine:       "minio",
		Bucket:       "admin-files",
		ObjectKey:    "uploads/2026/05/24/fixed-object-id.txt",
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Extension:    ".txt",
		Size:         5,
		SHA256:       "",
		Status:       models.FileAssetStatusPendingUpload,
		Metadata:     datatypes.JSONMap{},
	}); err != nil {
		t.Fatalf("seed pending file asset failed: %v", err)
	}

	result, err := svc.CompleteDirectUpload(context.Background(), CompleteDirectUploadInput{
		ID:         1,
		OperatorID: 9,
	})
	if err != nil {
		t.Fatalf("complete direct upload failed: %v", err)
	}
	if result.Status != models.FileAssetStatusActive {
		t.Fatalf("expected active status, got %q", result.Status)
	}
}

func TestFileStorageDeleteIsIdempotentForMissingObject(t *testing.T) {
	q := newFileStorageTestQuery(t)
	engine := &fakeStorageEngine{removeErr: objectstore.ErrObjectNotFound}
	svc := newTestFileStorage(t, q, engine)

	item := newStoredFileAsset()
	if err := q.FileAsset.Create(item); err != nil {
		t.Fatalf("seed file asset failed: %v", err)
	}

	if err := svc.Delete(context.Background(), item.ID, 9); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	count, err := q.FileAsset.Where(q.FileAsset.ID.Eq(item.ID)).Count()
	if err != nil {
		t.Fatalf("count active file asset failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected active file asset removed, got %d", count)
	}
}

func newFileStorageTestQuery(t *testing.T) *query.Query {
	t.Helper()

	dbName := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "_")
	db, err := gorm.Open(sqlite.Open("file:"+dbName+"?mode=memory&cache=shared"), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.AutoMigrate(&models.FileAsset{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	return query.Use(db)
}

func newTestFileStorage(t *testing.T, q *query.Query, engine objectstore.Engine) FileStorage {
	t.Helper()

	return &fileStorageImpl{
		q:      q,
		engine: engine,
		conf: storageRuntimeConfig{
			engineName:              "minio",
			bucket:                  "admin-files",
			objectKeyPrefix:         "uploads",
			maxUploadBytes:          10 * 1024 * 1024,
			presignedExpiresSeconds: 3600,
		},
		now: func() time.Time {
			return time.Date(2026, time.May, 24, 12, 0, 0, 0, time.UTC)
		},
		newObjectID: func() string {
			return "fixed-object-id"
		},
	}
}

func newStoredFileAsset() *models.FileAsset {
	return &models.FileAsset{
		Engine:       "minio",
		Bucket:       "admin-files",
		ObjectKey:    "uploads/2026/05/24/fixed-object-id.txt",
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Extension:    ".txt",
		Size:         5,
		SHA256:       "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		Status:       models.FileAssetStatusActive,
		Metadata:     datatypes.JSONMap{},
	}
}

type fakeStorageEngine struct {
	putInput        *objectstore.PutInput
	presignInput    *objectstore.PresignInput
	presignPutInput *objectstore.PresignPutInput
	statInfo        objectstore.ObjectInfo
	removeErr       error
	removeCalls     int
}

func (f *fakeStorageEngine) Name() string {
	return "minio"
}

func (f *fakeStorageEngine) PutObject(ctx context.Context, input objectstore.PutInput) (objectstore.PutResult, error) {
	copied := input
	if input.Reader != nil {
		body, err := io.ReadAll(input.Reader)
		if err != nil {
			return objectstore.PutResult{}, err
		}
		copied.Reader = bytes.NewReader(body)
	}
	f.putInput = &copied
	return objectstore.PutResult{
		Bucket:    input.Bucket,
		ObjectKey: input.ObjectKey,
		Size:      input.Size,
	}, nil
}

func (f *fakeStorageEngine) PresignedGetObject(ctx context.Context, input objectstore.PresignInput) (string, time.Time, error) {
	copied := input
	f.presignInput = &copied
	return "https://minio.local/presigned", time.Now().Add(input.Expiry), nil
}

func (f *fakeStorageEngine) PresignedPutObject(ctx context.Context, input objectstore.PresignPutInput) (string, time.Time, error) {
	copied := input
	f.presignPutInput = &copied
	return "https://minio.local/presigned-put", time.Now().Add(input.Expiry), nil
}

func (f *fakeStorageEngine) StatObject(ctx context.Context, bucket string, objectKey string) (objectstore.ObjectInfo, error) {
	if f.statInfo.Bucket == "" {
		return objectstore.ObjectInfo{Bucket: bucket, ObjectKey: objectKey}, nil
	}
	return f.statInfo, nil
}

func (f *fakeStorageEngine) GetObject(ctx context.Context, bucket string, objectKey string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}

func (f *fakeStorageEngine) RemoveObject(ctx context.Context, bucket string, objectKey string) error {
	f.removeCalls++
	return f.removeErr
}

func (f *fakeStorageEngine) Health(ctx context.Context) error {
	return nil
}

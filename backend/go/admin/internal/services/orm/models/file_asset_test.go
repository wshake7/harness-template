package models

import (
	"testing"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFileAssetAutoMigrate(t *testing.T) {
	db := openFileAssetTestDB(t, "file_asset_auto_migrate")
	if err := db.AutoMigrate(&FileAsset{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	asset := &FileAsset{
		Engine:       "minio",
		Bucket:       "admin-files",
		ObjectKey:    "uploads/2026/05/24/demo.txt",
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Extension:    ".txt",
		Size:         5,
		SHA256:       "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		BizType:      "attachment",
		BizID:        "biz-1",
		Status:       "active",
		Metadata:     datatypes.JSONMap{"source": "test"},
	}
	if err := db.Create(asset).Error; err != nil {
		t.Fatalf("create file asset failed: %v", err)
	}

	var got FileAsset
	if err := db.First(&got, asset.ID).Error; err != nil {
		t.Fatalf("query file asset failed: %v", err)
	}
	if got.BizType != "attachment" || got.BizID != "biz-1" {
		t.Fatalf("unexpected biz binding: %#v", got)
	}
	if got.Metadata["source"] != "test" {
		t.Fatalf("expected metadata source=test, got %#v", got.Metadata)
	}
	if got.DeletedAt != 0 {
		t.Fatalf("expected active record not deleted, got %d", got.DeletedAt)
	}
}

func TestFileAssetObjectKeyUniqueWhenActive(t *testing.T) {
	db := openFileAssetTestDB(t, "file_asset_unique")
	if err := db.AutoMigrate(&FileAsset{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	first := newTestFileAsset("uploads/2026/05/24/demo.txt")
	if err := db.Create(first).Error; err != nil {
		t.Fatalf("create first file asset failed: %v", err)
	}

	duplicate := newTestFileAsset("uploads/2026/05/24/demo.txt")
	if err := db.Create(duplicate).Error; err == nil {
		t.Fatal("expected duplicate active object key to fail")
	}

	if err := db.Delete(first).Error; err != nil {
		t.Fatalf("soft delete first file asset failed: %v", err)
	}

	recreated := newTestFileAsset("uploads/2026/05/24/demo.txt")
	if err := db.Create(recreated).Error; err != nil {
		t.Fatalf("expected recreate after soft delete to succeed, got %v", err)
	}
}

func openFileAssetTestDB(t *testing.T, name string) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	return db
}

func newTestFileAsset(objectKey string) *FileAsset {
	return &FileAsset{
		Engine:       "minio",
		Bucket:       "admin-files",
		ObjectKey:    objectKey,
		OriginalName: "demo.txt",
		ContentType:  "text/plain",
		Extension:    ".txt",
		Size:         5,
		SHA256:       "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		Status:       "active",
		Metadata:     datatypes.JSONMap{},
	}
}

package appsvc

import (
	"admin/internal/config"
	"admin/internal/services/objectstore"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestKnowledgeDocumentIndexerIndexesManualDocument(t *testing.T) {
	q := newKnowledgeIndexerTestQuery(t)
	mustCreateKnowledgeCollection(t, q, 1)
	mustCreateKnowledgeDocument(t, q, &models.KnowledgeDocument{
		CollectionID: 1,
		DocumentID:   "manual-doc",
		Title:        "Manual Doc",
		Content:      "# Title\nmanual content",
		ContentType:  "markdown",
		Source:       "manual://manual-doc",
		VectorStatus: "pending",
		Metadata:     datatypes.JSONMap{},
	})

	runner := &fakeKnowledgeIndexRunner{ids: []string{"vec-1", "vec-2"}}
	cleaner := &fakeVectorStoreCleaner{}
	indexer := newKnowledgeDocumentIndexerForTest(t, q, &fakeKnowledgeStorageEngine{}, runner, cleaner)

	err := indexer.IndexDocument(context.Background(), 1)
	if err != nil {
		t.Fatalf("index manual document failed: %v", err)
	}
	if runner.lastCollection == nil || runner.lastCollection.CollectionName != "biz-1" {
		t.Fatalf("expected collection biz-1, got %#v", runner.lastCollection)
	}
	if runner.lastMetadata["document_db_id"] != "1" {
		t.Fatalf("expected document_db_id metadata, got %#v", runner.lastMetadata)
	}
	item, _ := q.KnowledgeDocument.Where(q.KnowledgeDocument.ID.Eq(1)).First()
	if item.VectorStatus != "indexed" || item.VectorID != "vec-1,vec-2" {
		t.Fatalf("expected indexed vector status, got status=%q vectorID=%q", item.VectorStatus, item.VectorID)
	}
}

func TestKnowledgeDocumentIndexerImportsActiveFileAsset(t *testing.T) {
	q := newKnowledgeIndexerTestQuery(t)
	mustCreateKnowledgeCollection(t, q, 1)
	if err := q.FileAsset.Create(&models.FileAsset{
		Engine:       "minio",
		Bucket:       "admin-files",
		ObjectKey:    "uploads/test/runbook.md",
		OriginalName: "runbook.md",
		ContentType:  "text/markdown",
		Extension:    ".md",
		Size:         12,
		Status:       models.FileAssetStatusActive,
		Metadata:     datatypes.JSONMap{},
	}); err != nil {
		t.Fatalf("seed file asset failed: %v", err)
	}

	runner := &fakeKnowledgeIndexRunner{ids: []string{"vec-file-1"}}
	indexer := newKnowledgeDocumentIndexerForTest(t, q, &fakeKnowledgeStorageEngine{
		getBody: "# Runbook\nhello world",
	}, runner, &fakeVectorStoreCleaner{})

	doc, err := indexer.ImportFile(context.Background(), ImportKnowledgeFileInput{
		CollectionID: 1,
		FileAssetID:  1,
		Title:        "Runbook",
		ContentType:  "markdown",
		Remark:       "imported",
		OperatorID:   7,
	})
	if err != nil {
		t.Fatalf("import file failed: %v", err)
	}
	if doc.DocumentID != "file-1" {
		t.Fatalf("expected generated document id file-1, got %q", doc.DocumentID)
	}
	if !strings.Contains(doc.Content, "hello world") {
		t.Fatalf("expected imported content, got %q", doc.Content)
	}
	if doc.VectorStatus != "pending" {
		t.Fatalf("expected pending vector status after import, got %q", doc.VectorStatus)
	}
	if runner.lastMetadata != nil {
		t.Fatalf("expected import to defer indexing, got metadata %#v", runner.lastMetadata)
	}
}

func TestKnowledgeDocumentIndexerRejectsInactiveFileAsset(t *testing.T) {
	q := newKnowledgeIndexerTestQuery(t)
	mustCreateKnowledgeCollection(t, q, 1)
	if err := q.FileAsset.Create(&models.FileAsset{
		Engine:       "minio",
		Bucket:       "admin-files",
		ObjectKey:    "uploads/test/runbook.md",
		OriginalName: "runbook.md",
		ContentType:  "text/markdown",
		Extension:    ".md",
		Size:         12,
		Status:       models.FileAssetStatusPendingUpload,
		Metadata:     datatypes.JSONMap{},
	}); err != nil {
		t.Fatalf("seed file asset failed: %v", err)
	}

	indexer := newKnowledgeDocumentIndexerForTest(t, q, &fakeKnowledgeStorageEngine{}, &fakeKnowledgeIndexRunner{}, &fakeVectorStoreCleaner{})
	_, err := indexer.ImportFile(context.Background(), ImportKnowledgeFileInput{
		CollectionID: 1,
		FileAssetID:  1,
		Title:        "Runbook",
		OperatorID:   7,
	})
	if err == nil {
		t.Fatal("expected inactive file asset to be rejected")
	}
}

func TestKnowledgeDocumentIndexerDeletesExistingVectorsByMetadata(t *testing.T) {
	q := newKnowledgeIndexerTestQuery(t)
	mustCreateKnowledgeCollection(t, q, 1)
	mustCreateKnowledgeDocument(t, q, &models.KnowledgeDocument{
		CollectionID: 1,
		DocumentID:   "manual-doc",
		Title:        "Manual Doc",
		Content:      "hello",
		ContentType:  "text",
		Source:       "manual://manual-doc",
		VectorStatus: "indexed",
		Metadata:     datatypes.JSONMap{},
	})

	cleaner := &fakeVectorStoreCleaner{}
	indexer := newKnowledgeDocumentIndexerForTest(t, q, &fakeKnowledgeStorageEngine{}, &fakeKnowledgeIndexRunner{ids: []string{"vec-1"}}, cleaner)

	err := indexer.DeleteDocumentVectors(context.Background(), 1)
	if err != nil {
		t.Fatalf("delete document vectors failed: %v", err)
	}
	if cleaner.lastDocumentID != 1 || cleaner.lastCollectionName != "biz-1" {
		t.Fatalf("unexpected cleaner call: %#v", cleaner)
	}
}

func TestKnowledgeDocumentIndexerMarksFailedWhenIndexingFails(t *testing.T) {
	q := newKnowledgeIndexerTestQuery(t)
	mustCreateKnowledgeCollection(t, q, 1)
	mustCreateKnowledgeDocument(t, q, &models.KnowledgeDocument{
		CollectionID: 1,
		DocumentID:   "manual-doc",
		Title:        "Manual Doc",
		Content:      "hello",
		ContentType:  "text",
		Source:       "manual://manual-doc",
		VectorStatus: "pending",
		Metadata:     datatypes.JSONMap{},
	})

	indexer := newKnowledgeDocumentIndexerForTest(t, q, &fakeKnowledgeStorageEngine{}, &fakeKnowledgeIndexRunner{
		err: errors.New("runner failed"),
	}, &fakeVectorStoreCleaner{})

	err := indexer.IndexDocument(context.Background(), 1)
	if err == nil {
		t.Fatal("expected indexing failure")
	}
	item, _ := q.KnowledgeDocument.Where(q.KnowledgeDocument.ID.Eq(1)).First()
	if item.VectorStatus != "failed" {
		t.Fatalf("expected failed vector status, got %q", item.VectorStatus)
	}
	if item.IndexingError == "" {
		t.Fatal("expected indexing error to be recorded")
	}
}

func newKnowledgeIndexerTestQuery(t *testing.T) *query.Query {
	t.Helper()

	dbName := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "_")
	db, err := gorm.Open(sqlite.Open("file:"+dbName+"?mode=memory&cache=shared"), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.AutoMigrate(&models.KnowledgeCollection{}, &models.KnowledgeDocument{}, &models.FileAsset{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	return query.Use(db)
}

func newKnowledgeDocumentIndexerForTest(t *testing.T, q *query.Query, storage objectstore.Engine, runner KnowledgeIndexRunner, cleaner VectorStoreCleaner) *knowledgeDocumentIndexer {
	t.Helper()

	return &knowledgeDocumentIndexer{
		q:       q,
		conf:    &config.Config{AI: config.AIConfig{Embedding: config.AIEmbeddingConfig{Dimensions: 2048, Model: "test-embedding"}}},
		storage: storage,
		runner:  runner,
		cleaner: cleaner,
		now: func() time.Time {
			return time.Date(2026, time.May, 25, 12, 0, 0, 0, time.UTC)
		},
	}
}

func mustCreateKnowledgeCollection(t *testing.T, q *query.Query, id uint64) {
	t.Helper()

	err := q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		CollectionName:   "biz-1",
		DisplayName:      "Biz 1",
		MetricType:       "COSINE",
		IndexType:        "auto",
	})
	if err != nil {
		t.Fatalf("seed knowledge collection failed: %v", err)
	}
	item, err := q.KnowledgeCollection.Where(q.KnowledgeCollection.CollectionName.Eq("biz-1")).First()
	if err != nil {
		t.Fatalf("query seeded knowledge collection failed: %v", err)
	}
	if item.ID != id {
		t.Fatalf("expected collection id %d, got %d", id, item.ID)
	}
}

func mustCreateKnowledgeDocument(t *testing.T, q *query.Query, doc *models.KnowledgeDocument) {
	t.Helper()
	if doc.Metadata == nil {
		doc.Metadata = datatypes.JSONMap{}
	}
	if err := q.KnowledgeDocument.Create(doc); err != nil {
		t.Fatalf("seed knowledge document failed: %v", err)
	}
}

type fakeKnowledgeStorageEngine struct {
	getBody string
}

func (f *fakeKnowledgeStorageEngine) Name() string { return "minio" }
func (f *fakeKnowledgeStorageEngine) PutObject(ctx context.Context, input objectstore.PutInput) (objectstore.PutResult, error) {
	return objectstore.PutResult{}, nil
}
func (f *fakeKnowledgeStorageEngine) PresignedPutObject(ctx context.Context, input objectstore.PresignPutInput) (string, time.Time, error) {
	return "", time.Time{}, nil
}
func (f *fakeKnowledgeStorageEngine) PresignedGetObject(ctx context.Context, input objectstore.PresignInput) (string, time.Time, error) {
	return "", time.Time{}, nil
}
func (f *fakeKnowledgeStorageEngine) StatObject(ctx context.Context, bucket string, objectKey string) (objectstore.ObjectInfo, error) {
	return objectstore.ObjectInfo{Bucket: bucket, ObjectKey: objectKey, Size: int64(len(f.getBody))}, nil
}
func (f *fakeKnowledgeStorageEngine) GetObject(ctx context.Context, bucket string, objectKey string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewBufferString(f.getBody)), nil
}
func (f *fakeKnowledgeStorageEngine) RemoveObject(ctx context.Context, bucket string, objectKey string) error {
	return nil
}
func (f *fakeKnowledgeStorageEngine) Health(ctx context.Context) error { return nil }

type fakeKnowledgeIndexRunner struct {
	ids            []string
	err            error
	lastCollection *models.KnowledgeCollection
	lastMetadata   map[string]any
}

func (f *fakeKnowledgeIndexRunner) Index(ctx context.Context, conf *config.Config, collection *models.KnowledgeCollection, sourcePath string, metadata map[string]any) ([]string, error) {
	f.lastCollection = collection
	f.lastMetadata = metadata
	if f.err != nil {
		return nil, f.err
	}
	return f.ids, nil
}

type fakeVectorStoreCleaner struct {
	lastCollectionName string
	lastDocumentID     uint64
	lastSource         string
}

func (f *fakeVectorStoreCleaner) DeleteByDocument(ctx context.Context, collectionName string, documentDBID uint64, source string) error {
	f.lastCollectionName = collectionName
	f.lastDocumentID = documentDBID
	f.lastSource = source
	return nil
}

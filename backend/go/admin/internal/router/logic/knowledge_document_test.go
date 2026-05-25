package logic

import (
	"context"
	"fmt"
	"testing"

	"admin/internal/appsvc"
	"admin/internal/mock"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"admin/internal/services/temporaljob"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/client"
	"go.uber.org/mock/gomock"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func setupKnowledgeDocumentHandler(t *testing.T) *KnowledgeDocumentHandler {
	t.Helper()
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeDocument.UnderlyingDB())
	return NewKnowledgeDocumentHandler(q, nil, nil, "")
}

func TestKnowledgeDocumentHandler_List_Empty(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
	assert.Len(t, result.Items, 0)
}

func TestKnowledgeDocumentHandler_List_WithData(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)
	q := h.Q

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "doc-001",
		Title:           "Test Document",
		Content:         "This is test content",
		ContentType:     "text",
		VectorStatus:    "pending",
	})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "doc-001", result.Items[0].DocumentID)
	assert.Equal(t, "Test Document", result.Items[0].Title)
}

func TestKnowledgeDocumentHandler_ListByCollection_Empty(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	result, err := h.ListByCollection(newTestCtx(t), &ReqKnowledgeDocumentListByCollection{
		CollectionID: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestKnowledgeDocumentHandler_ListByCollection_WithData(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)
	q := h.Q

	// Create documents for collection 1
	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "doc-001",
		Title:           "Doc 1",
		Content:         "Content 1",
		ContentType:     "text",
		VectorStatus:    "pending",
	})
	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 2},
		CollectionID:    1,
		DocumentID:      "doc-002",
		Title:           "Doc 2",
		Content:         "Content 2",
		ContentType:     "text",
		VectorStatus:    "indexed",
	})
	// Document for collection 2
	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 3},
		CollectionID:    2,
		DocumentID:      "doc-003",
		Title:           "Doc 3",
		Content:         "Content 3",
		ContentType:     "text",
		VectorStatus:    "pending",
	})

	result, err := h.ListByCollection(newTestCtx(t), &ReqKnowledgeDocumentListByCollection{
		CollectionID: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), result.Total)
	assert.Len(t, result.Items, 2)
}

func TestKnowledgeDocumentHandler_Detail_Success(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)
	q := h.Q

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "detail-doc",
		Title:           "Detail Document",
		Content:         "Detail content",
		ContentType:     "markdown",
		VectorStatus:    "indexed",
	})

	result, err := h.Detail(newTestCtx(t), &ReqKnowledgeDocumentID{ID: 1})
	assert.NoError(t, err)
	assert.Equal(t, "detail-doc", result.DocumentID)
	assert.Equal(t, "Detail Document", result.Title)
}

func TestKnowledgeDocumentHandler_Detail_NotFound(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	_, err := h.Detail(newTestCtx(t), &ReqKnowledgeDocumentID{ID: 999})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "文档不存在")
}

func TestKnowledgeDocumentHandler_Create_Success(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	err := h.Create(newTestCtx(t), &ReqKnowledgeDocumentCreate{
		CollectionID: 1,
		DocumentID:   "new-doc",
		Title:        "New Document",
		Content:      "New content",
		IsEnabled:    true,
	})
	assert.NoError(t, err)

	item, err := h.Q.KnowledgeDocument.Where(h.Q.KnowledgeDocument.DocumentID.Eq("new-doc")).First()
	assert.NoError(t, err)
	assert.Equal(t, "New Document", item.Title)
	assert.Equal(t, "text", item.ContentType)
	assert.Equal(t, "pending", item.VectorStatus)
}

func TestKnowledgeDocumentHandler_Create_WithMetadata(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	err := h.Create(newTestCtx(t), &ReqKnowledgeDocumentCreate{
		CollectionID: 1,
		DocumentID:   "meta-doc",
		Title:        "Meta Document",
		Content:      "Content",
		Metadata:     `{"key":"value","num":123}`,
		IsEnabled:    true,
	})
	assert.NoError(t, err)

	item, _ := h.Q.KnowledgeDocument.Where(h.Q.KnowledgeDocument.DocumentID.Eq("meta-doc")).First()
	assert.Equal(t, "value", item.Metadata["key"])
	// SQLite stores JSON numbers as json.Number
	numStr := fmt.Sprintf("%v", item.Metadata["num"])
	assert.Equal(t, "123", numStr)
}

func TestKnowledgeDocumentHandler_Create_InvalidMetadata(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	err := h.Create(newTestCtx(t), &ReqKnowledgeDocumentCreate{
		CollectionID: 1,
		DocumentID:   "bad-meta",
		Title:        "Bad Meta",
		Content:      "Content",
		Metadata:     `not-json`,
		IsEnabled:    true,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "元数据格式错误")
}

func TestKnowledgeDocumentHandler_Update_Success(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)
	q := h.Q

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "update-doc",
		Title:           "Original Title",
		Content:         "Original content",
		ContentType:     "text",
		VectorStatus:    "pending",
	})

	newTitle := "Updated Title"
	err := h.Update(newTestCtx(t), &ReqKnowledgeDocumentUpdate{
		ID:    1,
		Title: &newTitle,
	})
	assert.NoError(t, err)

	item, _ := q.KnowledgeDocument.Where(q.KnowledgeDocument.ID.Eq(1)).First()
	assert.Equal(t, "Updated Title", item.Title)
}

func TestKnowledgeDocumentHandler_Update_NotFound(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	newTitle := "Updated"
	err := h.Update(newTestCtx(t), &ReqKnowledgeDocumentUpdate{
		ID:    999,
		Title: &newTitle,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "文档不存在")
}

func TestKnowledgeDocumentHandler_Update_WithMetadata(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)
	q := h.Q

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "update-meta",
		Title:           "Title",
		Content:         "Content",
		VectorStatus:    "pending",
	})

	meta := `{"updated":true}`
	err := h.Update(newTestCtx(t), &ReqKnowledgeDocumentUpdate{
		ID:       1,
		Metadata: &meta,
	})
	assert.NoError(t, err)

	item, _ := q.KnowledgeDocument.Where(q.KnowledgeDocument.ID.Eq(1)).First()
	assert.Equal(t, true, item.Metadata["updated"])
}

func TestKnowledgeDocumentHandler_Del_Success(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)
	q := h.Q

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "delete-doc",
		Title:           "Delete Me",
		Content:         "Content",
		VectorStatus:    "pending",
	})

	err := h.Del(newTestCtx(t), &ReqKnowledgeDocumentID{ID: 1})
	assert.NoError(t, err)

	_, err = q.KnowledgeDocument.Where(q.KnowledgeDocument.ID.Eq(1)).First()
	assert.Error(t, err)
}

func TestKnowledgeDocumentHandler_Del_NotFound(t *testing.T) {
	h := setupKnowledgeDocumentHandler(t)

	err := h.Del(newTestCtx(t), &ReqKnowledgeDocumentID{ID: 999})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "文档不存在")
}

func TestKnowledgeDocumentHandler_ImportFileDelegatesToIndexer(t *testing.T) {
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeDocument.UnderlyingDB())
	indexer := &stubKnowledgeDocumentIndexer{
		importResult: &models.KnowledgeDocument{
			AutoIncrementID: mixin.AutoIncrementID{ID: 1},
			CollectionID:    1,
			DocumentID:      "file-1",
			Title:           "Imported",
			VectorStatus:    "indexed",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	temporalSvc := mock.NewMockTemporalService(ctrl)
	temporalSvc.EXPECT().
		IsConnected().
		Return(true)
	temporalSvc.EXPECT().
		ExecuteWorkflow(gomock.Any(), gomock.Any(), temporaljob.DispatchWorkflowName, gomock.Any()).
		Return((client.WorkflowRun)(nil), nil)
	h := NewKnowledgeDocumentHandler(q, indexer, temporalSvc, "admin")

	result, err := h.ImportFile(newTestCtx(t), &ReqKnowledgeDocumentImportFile{
		CollectionID: 1,
		FileAssetID:  1,
		Title:        "Imported",
		ContentType:  "markdown",
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, uint64(1), indexer.lastImport.CollectionID)
	assert.Equal(t, uint64(1), indexer.lastImport.FileAssetID)
}

func TestKnowledgeDocumentHandler_CreateIndexesDocument(t *testing.T) {
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeDocument.UnderlyingDB())
	indexer := &stubKnowledgeDocumentIndexer{}
	h := NewKnowledgeDocumentHandler(q, indexer, nil, "")

	err := h.Create(newTestCtx(t), &ReqKnowledgeDocumentCreate{
		CollectionID: 1,
		DocumentID:   "new-doc",
		Title:        "New Document",
		Content:      "New content",
		IsEnabled:    true,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), indexer.lastIndexedID)
}

func TestKnowledgeDocumentHandler_UpdateReindexesWhenContentChanges(t *testing.T) {
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeDocument.UnderlyingDB())
	indexer := &stubKnowledgeDocumentIndexer{}
	h := NewKnowledgeDocumentHandler(q, indexer, nil, "")

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "update-doc",
		Title:           "Original",
		Content:         "Original content",
		ContentType:     "text",
		VectorStatus:    "pending",
	})

	nextContent := "Updated content"
	err := h.Update(newTestCtx(t), &ReqKnowledgeDocumentUpdate{
		ID:      1,
		Content: &nextContent,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), indexer.lastIndexedID)
}

func TestKnowledgeDocumentHandler_UpdateDoesNotReindexWhenOnlyRemarkChanges(t *testing.T) {
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeDocument.UnderlyingDB())
	indexer := &stubKnowledgeDocumentIndexer{}
	h := NewKnowledgeDocumentHandler(q, indexer, nil, "")

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "update-doc",
		Title:           "Original",
		Content:         "Original content",
		ContentType:     "text",
		VectorStatus:    "pending",
	})

	remark := "Remark only"
	err := h.Update(newTestCtx(t), &ReqKnowledgeDocumentUpdate{
		ID:     1,
		Remark: &remark,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), indexer.lastIndexedID)
}

func TestKnowledgeDocumentHandler_DelDeletesVectorsBeforeDBRecord(t *testing.T) {
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeDocument.UnderlyingDB())
	indexer := &stubKnowledgeDocumentIndexer{}
	h := NewKnowledgeDocumentHandler(q, indexer, nil, "")

	q.KnowledgeDocument.Create(&models.KnowledgeDocument{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionID:    1,
		DocumentID:      "delete-doc",
		Title:           "Delete Me",
		Content:         "Content",
		VectorStatus:    "indexed",
	})

	err := h.Del(newTestCtx(t), &ReqKnowledgeDocumentID{ID: 1})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), indexer.lastDeletedID)
}

type stubKnowledgeDocumentIndexer struct {
	importResult  *models.KnowledgeDocument
	lastImport    appsvc.ImportKnowledgeFileInput
	lastIndexedID uint64
	lastDeletedID uint64
	indexErr      error
	deleteErr     error
	importErr     error
}

func (s *stubKnowledgeDocumentIndexer) ImportFile(ctx context.Context, input appsvc.ImportKnowledgeFileInput) (*models.KnowledgeDocument, error) {
	s.lastImport = input
	if s.importResult != nil || s.importErr != nil {
		return s.importResult, s.importErr
	}
	return &models.KnowledgeDocument{}, nil
}

func (s *stubKnowledgeDocumentIndexer) IndexDocument(ctx context.Context, documentID uint64) error {
	s.lastIndexedID = documentID
	return s.indexErr
}

func (s *stubKnowledgeDocumentIndexer) DeleteDocumentVectors(ctx context.Context, documentID uint64) error {
	s.lastDeletedID = documentID
	return s.deleteErr
}

func (s *stubKnowledgeDocumentIndexer) EnsureCollection(ctx context.Context, collection *models.KnowledgeCollection) error {
	return nil
}

func (s *stubKnowledgeDocumentIndexer) DropCollection(ctx context.Context, collectionName string) error {
	return nil
}

package logic

import (
	"context"
	"testing"

	"admin/internal/appsvc"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

type mockCollectionIndexer struct {
	ensureCalled []*models.KnowledgeCollection
	dropCalled   []string
	ensureErr    error
	dropErr      error
}

func (m *mockCollectionIndexer) ImportFile(ctx context.Context, input appsvc.ImportKnowledgeFileInput) (*models.KnowledgeDocument, error) {
	return nil, nil
}
func (m *mockCollectionIndexer) IndexDocument(ctx context.Context, documentID uint64) error { return nil }
func (m *mockCollectionIndexer) DeleteDocumentVectors(ctx context.Context, documentID uint64) error {
	return nil
}
func (m *mockCollectionIndexer) EnsureCollection(ctx context.Context, collection *models.KnowledgeCollection) error {
	m.ensureCalled = append(m.ensureCalled, collection)
	return m.ensureErr
}
func (m *mockCollectionIndexer) DropCollection(ctx context.Context, collectionName string) error {
	m.dropCalled = append(m.dropCalled, collectionName)
	return m.dropErr
}

func setupKnowledgeCollectionHandler(t *testing.T) *KnowledgeCollectionHandler {
	t.Helper()
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeCollection.UnderlyingDB())
	return NewKnowledgeCollectionHandler(q, nil)
}

func setupKnowledgeCollectionHandlerWithMock(t *testing.T) (*KnowledgeCollectionHandler, *mockCollectionIndexer) {
	t.Helper()
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeCollection.UnderlyingDB())
	mock := &mockCollectionIndexer{}
	return NewKnowledgeCollectionHandler(q, mock), mock
}

func TestKnowledgeCollectionHandler_List_Empty(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
	assert.Len(t, result.Items, 0)
}

func TestKnowledgeCollectionHandler_List_WithData(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)
	q := h.Q

	q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		AutoIncrementID:  mixin.AutoIncrementID{ID: 1},
		CollectionName:   "test-collection",
		DisplayName:      "Test Collection",
		MetricType:       "COSINE",
		IndexType:        "auto",
		IsEnabled:        mixin.IsEnabled{IsEnabled: true},
	})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "test-collection", result.Items[0].CollectionName)
	assert.True(t, result.Items[0].CanWrite)
	assert.True(t, result.Items[0].CanDelete)
}

func TestKnowledgeCollectionHandler_Detail_Success(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)
	q := h.Q

	q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionName:  "detail-test",
		DisplayName:     "Detail Test",
	})

	result, err := h.Detail(newTestCtx(t), &ReqKnowledgeCollectionID{ID: 1})
	assert.NoError(t, err)
	assert.Equal(t, "detail-test", result.CollectionName)
	assert.Equal(t, "Detail Test", result.DisplayName)
}

func TestKnowledgeCollectionHandler_Detail_NotFound(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)

	_, err := h.Detail(newTestCtx(t), &ReqKnowledgeCollectionID{ID: 999})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "集合不存在")
}

func TestKnowledgeCollectionHandler_Create_Success(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)

	err := h.Create(newTestCtx(t), &ReqKnowledgeCollectionCreate{
		CollectionName:  "new-collection",
		DisplayName:     "New Collection",
		IsEnabled:       true,
	})
	assert.NoError(t, err)

	// Verify created
	item, err := h.Q.KnowledgeCollection.Where(h.Q.KnowledgeCollection.CollectionName.Eq("new-collection")).First()
	assert.NoError(t, err)
	assert.Equal(t, "New Collection", item.DisplayName)
	assert.Equal(t, "COSINE", item.MetricType)
	assert.Equal(t, "auto", item.IndexType)
}

func TestKnowledgeCollectionHandler_Create_DefaultValues(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)

	err := h.Create(newTestCtx(t), &ReqKnowledgeCollectionCreate{
		CollectionName:  "defaults-test",
		DisplayName:     "Defaults Test",
		IsEnabled:       false,
	})
	assert.NoError(t, err)

	item, _ := h.Q.KnowledgeCollection.Where(h.Q.KnowledgeCollection.CollectionName.Eq("defaults-test")).First()
	assert.Equal(t, "COSINE", item.MetricType)
	assert.Equal(t, "auto", item.IndexType)
}

func TestKnowledgeCollectionHandler_Update_Success(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)
	q := h.Q

	q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionName:  "update-test",
		DisplayName:     "Original",
	})

	newName := "Updated Name"
	err := h.Update(newTestCtx(t), &ReqKnowledgeCollectionUpdate{
		ID:        1,
		DisplayName: &newName,
	})
	assert.NoError(t, err)

	item, _ := q.KnowledgeCollection.Where(q.KnowledgeCollection.ID.Eq(1)).First()
	assert.Equal(t, "Updated Name", item.DisplayName)
}

func TestKnowledgeCollectionHandler_Update_NotFound(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)

	newName := "Updated"
	err := h.Update(newTestCtx(t), &ReqKnowledgeCollectionUpdate{
		ID:          999,
		DisplayName: &newName,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "集合不存在")
}

func TestKnowledgeCollectionHandler_Del_Success(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)
	q := h.Q

	q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionName:  "delete-test",
		DisplayName:     "Delete Test",
	})

	err := h.Del(newTestCtx(t), &ReqKnowledgeCollectionID{ID: 1})
	assert.NoError(t, err)

	_, err = q.KnowledgeCollection.Where(q.KnowledgeCollection.ID.Eq(1)).First()
	assert.Error(t, err) // Should be not found
}

func TestKnowledgeCollectionHandler_Del_NotFound(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)

	err := h.Del(newTestCtx(t), &ReqKnowledgeCollectionID{ID: 999})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "集合不存在")
}

func TestKnowledgeCollectionHandler_Update_IsEnabled(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)
	q := h.Q

	q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionName:  "status-test",
		DisplayName:     "Status Test",
		IsEnabled:       mixin.IsEnabled{IsEnabled: true},
	})

	nextEnabled := false
	err := h.Update(newTestCtx(t), &ReqKnowledgeCollectionUpdate{ID: 1, IsEnabled: &nextEnabled})
	assert.NoError(t, err)

	item, _ := q.KnowledgeCollection.Where(q.KnowledgeCollection.ID.Eq(1)).First()
	assert.False(t, item.IsEnabled.IsEnabled)
}

func TestKnowledgeCollectionHandler_Create_EnsuresMilvusCollection(t *testing.T) {
	h, mock := setupKnowledgeCollectionHandlerWithMock(t)

	err := h.Create(newTestCtx(t), &ReqKnowledgeCollectionCreate{
		CollectionName:  "milvus-test",
		DisplayName:     "Milvus Test",
		MetricType:      "IP",
		IndexType:       "hnsw",
		IsEnabled:       true,
	})
	require.NoError(t, err)

	require.Len(t, mock.ensureCalled, 1)
	c := mock.ensureCalled[0]
	assert.Equal(t, "milvus-test", c.CollectionName)
	assert.Equal(t, "IP", c.MetricType)
	assert.Equal(t, "hnsw", c.IndexType)
}

func TestKnowledgeCollectionHandler_Del_DropsMilvusCollection(t *testing.T) {
	h, mock := setupKnowledgeCollectionHandlerWithMock(t)
	q := h.Q

	q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionName:  "drop-test",
		DisplayName:     "Drop Test",
	})

	err := h.Del(newTestCtx(t), &ReqKnowledgeCollectionID{ID: 1})
	require.NoError(t, err)

	require.Len(t, mock.dropCalled, 1)
	assert.Equal(t, "drop-test", mock.dropCalled[0])
}

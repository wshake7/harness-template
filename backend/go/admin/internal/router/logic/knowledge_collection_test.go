package logic

import (
	"testing"

	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"

	"github.com/stretchr/testify/assert"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func setupKnowledgeCollectionHandler(t *testing.T) *KnowledgeCollectionHandler {
	t.Helper()
	q := mustMigrateKnowledge(t)
	query.SetDefault(q.KnowledgeCollection.UnderlyingDB())
	return NewKnowledgeCollectionHandler(q)
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
		EmbeddingModel:   "text-embedding-3-small",
		VectorDimension:  1536,
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
		EmbeddingModel:  "test-model",
		VectorDimension: 768,
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
		EmbeddingModel:  "text-embedding-ada-002",
		VectorDimension: 1536,
		IsEnabled:       true,
	})
	assert.NoError(t, err)

	// Verify created
	item, err := h.Q.KnowledgeCollection.Where(h.Q.KnowledgeCollection.CollectionName.Eq("new-collection")).First()
	assert.NoError(t, err)
	assert.Equal(t, "New Collection", item.DisplayName)
	assert.Equal(t, "COSINE", item.MetricType)
	assert.Equal(t, "auto", item.IndexType)
	assert.Equal(t, 255, item.IDMaxLength)
	assert.Equal(t, 8192, item.ContentMaxLength)
}

func TestKnowledgeCollectionHandler_Create_DefaultValues(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)

	err := h.Create(newTestCtx(t), &ReqKnowledgeCollectionCreate{
		CollectionName:  "defaults-test",
		DisplayName:     "Defaults Test",
		EmbeddingModel:  "model",
		VectorDimension: 512,
		IsEnabled:       false,
	})
	assert.NoError(t, err)

	item, _ := h.Q.KnowledgeCollection.Where(h.Q.KnowledgeCollection.CollectionName.Eq("defaults-test")).First()
	assert.Equal(t, "COSINE", item.MetricType)
	assert.Equal(t, "auto", item.IndexType)
	assert.Equal(t, 255, item.IDMaxLength)
	assert.Equal(t, 8192, item.ContentMaxLength)
}

func TestKnowledgeCollectionHandler_Update_Success(t *testing.T) {
	h := setupKnowledgeCollectionHandler(t)
	q := h.Q

	q.KnowledgeCollection.Create(&models.KnowledgeCollection{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		CollectionName:  "update-test",
		DisplayName:     "Original",
		EmbeddingModel:  "model",
		VectorDimension: 768,
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
		EmbeddingModel:  "model",
		VectorDimension: 768,
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
		EmbeddingModel:  "model",
		VectorDimension: 768,
		IsEnabled:       mixin.IsEnabled{IsEnabled: true},
	})

	nextEnabled := false
	err := h.Update(newTestCtx(t), &ReqKnowledgeCollectionUpdate{ID: 1, IsEnabled: &nextEnabled})
	assert.NoError(t, err)

	item, _ := q.KnowledgeCollection.Where(q.KnowledgeCollection.ID.Eq(1)).First()
	assert.False(t, item.IsEnabled.IsEnabled)
}

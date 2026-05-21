package logic

import (
	"admin/internal/mock"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func mustMigrateResourceAPI(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysResourceApi{},
	)
}

func TestSysResourceApiHandler_List_Empty(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysResourceApiHandler(q, mock.NewMockCasbinService(ctrl))

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysResourceApiHandler_List_WithData(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysResourceApiHandler(q, mock.NewMockCasbinService(ctrl))

	q.SysResourceApi.Create(&models.SysResourceApi{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Path: "/api/test", Method: "GET"})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
}

func TestSysResourceApiHandler_Create_Success(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysResourceApiHandler(q, mock.NewMockCasbinService(ctrl))

	err := h.Create(newTestCtx(t), &ReqResourceApiCreate{Path: "/api/test", Method: "GET"})
	assert.NoError(t, err)
}

func TestSysResourceApiHandler_Create_Duplicate(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysResourceApiHandler(q, mock.NewMockCasbinService(ctrl))

	assert.NoError(t, h.Create(newTestCtx(t), &ReqResourceApiCreate{Path: "/api/test", Method: "GET"}))
	err := h.Create(newTestCtx(t), &ReqResourceApiCreate{Path: "/api/test", Method: "GET"})
	assert.Error(t, err)
}

func TestSysResourceApiHandler_Update_NotFound(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysResourceApiHandler(q, mock.NewMockCasbinService(ctrl))

	err := h.Update(newTestCtx(t), &ReqResourceApiUpdate{ID: 99})
	assert.Error(t, err)
}

func TestSysResourceApiHandler_Update_Success_WithCasbinSync(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockCasbin := mock.NewMockCasbinService(ctrl)
	h := NewSysResourceApiHandler(q, mockCasbin)

	q.SysResourceApi.Create(&models.SysResourceApi{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Path: "/api/test", Method: "GET", IsEnabled: mixin.IsEnabled{IsEnabled: true}})
	newPath := "/api/changed"
	mockCasbin.EXPECT().SyncAPIResourcePolicies(gomock.Any(), gomock.Any()).Return(nil)

	err := h.Update(newTestCtx(t), &ReqResourceApiUpdate{ID: 1, Path: &newPath})
	assert.NoError(t, err)
}

func TestSysResourceApiHandler_Del_Success(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockCasbin := mock.NewMockCasbinService(ctrl)
	h := NewSysResourceApiHandler(q, mockCasbin)

	q.SysResourceApi.Create(&models.SysResourceApi{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Path: "/api/test", Method: "GET"})
	mockCasbin.EXPECT().RemoveAPIResourcePolicies(gomock.Any()).Return(nil)

	err := h.Del(newTestCtx(t), &ReqResourceApiBatchDelete{IDs: []uint64{1}})
	assert.NoError(t, err)
}

func TestSysResourceApiHandler_Del_NotFound(t *testing.T) {
	q := mustMigrateResourceAPI(t)
	query.SetDefault(q.SysResourceApi.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysResourceApiHandler(q, mock.NewMockCasbinService(ctrl))

	err := h.Del(newTestCtx(t), &ReqResourceApiBatchDelete{IDs: []uint64{99}})
	assert.Error(t, err)
}

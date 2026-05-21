package logic

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func mustMigrateApiLog(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysApiLog{},
		&models.SysUser{},
	)
}

func TestSysApiLogHandler_List_Empty(t *testing.T) {
	q := mustMigrateApiLog(t)
	query.SetDefault(q.SysApiLog.UnderlyingDB())
	h := NewSysApiLogHandler(q)

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysApiLogHandler_List_WithUserReload(t *testing.T) {
	q := mustMigrateApiLog(t)
	query.SetDefault(q.SysApiLog.UnderlyingDB())
	h := NewSysApiLogHandler(q)

	userID := uint64(1)
	q.SysUser.Create(&models.SysUser{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin"})
	q.SysApiLog.Create(&models.SysApiLog{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, RequestID: "r1", Path: "/api/test", Method: "GET", SysUserID: &userID})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
	assert.NotNil(t, result.Items[0].SysUser)
	assert.Equal(t, "admin", result.Items[0].SysUser.Username)
}

func TestSysApiLogHandler_Detail_Success(t *testing.T) {
	q := mustMigrateApiLog(t)
	query.SetDefault(q.SysApiLog.UnderlyingDB())
	h := NewSysApiLogHandler(q)

	userID := uint64(1)
	q.SysUser.Create(&models.SysUser{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin"})
	q.SysApiLog.Create(&models.SysApiLog{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, RequestID: "r1", Path: "/api/test", Method: "GET", SysUserID: &userID})

	result, err := h.Detail(newTestCtx(t), &ReqLogDetail{ID: 1})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.ID)
	assert.NotNil(t, result.SysUser)
}

func TestSysApiLogHandler_Detail_NotFound(t *testing.T) {
	q := mustMigrateApiLog(t)
	query.SetDefault(q.SysApiLog.UnderlyingDB())
	h := NewSysApiLogHandler(q)

	_, err := h.Detail(newTestCtx(t), &ReqLogDetail{ID: 99})
	assert.Error(t, err)
}

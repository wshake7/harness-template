package logic

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func mustMigrateLoginLog(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysLoginLog{},
		&models.SysUser{},
	)
}

func TestSysLoginLogHandler_List_Empty(t *testing.T) {
	q := mustMigrateLoginLog(t)
	query.SetDefault(q.SysLoginLog.UnderlyingDB())
	h := NewSysLoginLogHandler(q)

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysLoginLogHandler_List_WithUserReload(t *testing.T) {
	q := mustMigrateLoginLog(t)
	query.SetDefault(q.SysLoginLog.UnderlyingDB())
	h := NewSysLoginLogHandler(q)

	userID := uint64(1)
	q.SysUser.Create(&models.SysUser{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin", Nickname: "Admin"})
	q.SysLoginLog.Create(&models.SysLoginLog{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin", SysUserID: &userID})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
	assert.NotNil(t, result.Items[0].SysUser)
	assert.Equal(t, "Admin", result.Items[0].SysUser.Nickname)
}

func TestSysLoginLogHandler_Detail_Success(t *testing.T) {
	q := mustMigrateLoginLog(t)
	query.SetDefault(q.SysLoginLog.UnderlyingDB())
	h := NewSysLoginLogHandler(q)

	userID := uint64(1)
	q.SysUser.Create(&models.SysUser{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin", Nickname: "Admin"})
	q.SysLoginLog.Create(&models.SysLoginLog{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin", SysUserID: &userID})

	result, err := h.Detail(newTestCtx(t), &ReqLoginLogDetail{ID: 1})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.ID)
	assert.NotNil(t, result.SysUser)
}

func TestSysLoginLogHandler_Detail_NotFound(t *testing.T) {
	q := mustMigrateLoginLog(t)
	query.SetDefault(q.SysLoginLog.UnderlyingDB())
	h := NewSysLoginLogHandler(q)

	_, err := h.Detail(newTestCtx(t), &ReqLoginLogDetail{ID: 99})
	assert.Error(t, err)
}

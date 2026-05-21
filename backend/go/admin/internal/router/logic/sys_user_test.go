package logic

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func mustMigrateUser(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysUser{},
	)
}

func TestSysUserHandler_List_Empty(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysUserHandler_List_WithData(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	q.SysUser.Create(&models.SysUser{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin"})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
}

func TestSysUserHandler_Create_Success(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	err := h.Create(newTestCtx(t), &ReqSysUserCreate{Username: "admin", Password: "123456"})
	assert.NoError(t, err)
}

func TestSysUserHandler_Create_Duplicate(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	assert.NoError(t, h.Create(newTestCtx(t), &ReqSysUserCreate{Username: "admin", Password: "123456"}))
	err := h.Create(newTestCtx(t), &ReqSysUserCreate{Username: "admin", Password: "123456"})
	assert.Error(t, err)
}

func TestSysUserHandler_Update_NotFound(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	nickname := "NewName"
	err := h.Update(newTestCtx(t), &ReqSysUserUpdate{ID: 99, Nickname: &nickname})
	assert.Error(t, err)
}

func TestSysUserHandler_Update_Success(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	q.SysUser.Create(&models.SysUser{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin"})
	nickname := "NewName"
	err := h.Update(newTestCtx(t), &ReqSysUserUpdate{ID: 1, Nickname: &nickname})
	assert.NoError(t, err)
}

func TestSysUserHandler_Del_Success(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	q.SysUser.Create(&models.SysUser{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Username: "admin"})
	err := h.Del(newTestCtx(t), &ReqSysUserBatchDelete{IDs: []uint64{1}})
	assert.NoError(t, err)
}

func TestSysUserHandler_Del_NotFound(t *testing.T) {
	q := mustMigrateUser(t)
	query.SetDefault(q.SysUser.UnderlyingDB())
	h := NewSysUserHandler(q)

	err := h.Del(newTestCtx(t), &ReqSysUserBatchDelete{IDs: []uint64{99}})
	assert.Error(t, err)
}

func TestReqSysUserCreate_normalize(t *testing.T) {
	req := &ReqSysUserCreate{Username: " admin ", Password: " 123456 ", Nickname: " Admin ", LanguageCode: " zh ", Remark: " test "}
	req.normalize()
	assert.Equal(t, "admin", req.Username)
	assert.Equal(t, "123456", req.Password)
	assert.Equal(t, "Admin", req.Nickname)
	assert.Equal(t, "zh", req.LanguageCode)
	assert.Equal(t, "test", req.Remark)
}

func TestReqSysUserUpdate_normalize(t *testing.T) {
	username := " admin "
	nickname := " User "
	lang := " zh "
	remark := " test "
	req := &ReqSysUserUpdate{Username: &username, Nickname: &nickname, LanguageCode: &lang, Remark: &remark}
	req.normalize()
	assert.Equal(t, "admin", *req.Username)
	assert.Equal(t, "User", *req.Nickname)
	assert.Equal(t, "zh", *req.LanguageCode)
	assert.Equal(t, "test", *req.Remark)
}

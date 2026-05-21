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

func TestSysRoleHandler_List_Empty(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysRoleHandler_List_WithData(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	q.SysRole.Create(&models.SysRole{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Admin",
		Code:            "admin",
	})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
}

func TestSysRoleHandler_Tree_Empty(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	result, err := h.Tree(newTestCtx(t))
	assert.NoError(t, err)
	assert.Len(t, *result, 0)
}

func TestSysRoleHandler_Create_Success(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	err := h.Create(newTestCtx(t), &ReqSysRoleCreate{
		Name: "Admin", Code: "admin",
	})
	assert.NoError(t, err)
}

func TestSysRoleHandler_Permissions_NotFound(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	_, err := h.Permissions(newTestCtx(t), &ReqSysRolePermissionQuery{ID: 99})
	assert.Error(t, err)
}

func TestSysRoleHandler_Permissions_Empty(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	q.SysRole.Create(&models.SysRole{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Admin",
		Code:            "admin",
	})

	result, err := h.Permissions(newTestCtx(t), &ReqSysRolePermissionQuery{ID: 1})
	assert.NoError(t, err)
	assert.Len(t, result.MenuIDs, 0)
	assert.Len(t, result.ApiIDs, 0)
}

func TestSysRoleHandler_Del_HasChildren(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	q.SysRole.Create(&models.SysRole{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Parent",
		Code:            "parent",
	})
	parentID := uint64(1)
	q.SysRole.Create(&models.SysRole{
		AutoIncrementID: mixin.AutoIncrementID{ID: 2},
		Name:            "Child",
		Code:            "child",
		ParentID:        &parentID,
	})

	err := h.Del(newTestCtx(t), &ReqSysRoleBatchDelete{IDs: []uint64{1}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "子角色")
}

// --- Pure function tests ---

func TestReqSysRoleCreate_normalize(t *testing.T) {
	req := &ReqSysRoleCreate{
		Name: " Admin ", Code: " admin ", Remark: " test ",
	}
	req.normalize()
	assert.Equal(t, "Admin", req.Name)
	assert.Equal(t, "admin", req.Code)
	assert.Equal(t, "test", req.Remark)
}

func TestReqSysRoleUpdate_normalize(t *testing.T) {
	name := " Admin "
	code := " admin "
	remark := " test "
	req := &ReqSysRoleUpdate{
		Name: &name, Code: &code, Remark: &remark,
	}
	req.normalize()
	assert.Equal(t, "Admin", *req.Name)
	assert.Equal(t, "admin", *req.Code)
	assert.Equal(t, "test", *req.Remark)
}

func TestValidateSysRoleValues(t *testing.T) {
	assert.Error(t, validateSysRoleValues("", "code"))
	assert.Error(t, validateSysRoleValues("name", ""))
	assert.NoError(t, validateSysRoleValues("name", "code"))
}

// --- ensureSysRoleCanDelete ---

func TestEnsureSysRoleCanDelete_HasUsers(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Test", Code: "test"})
	q.SysUserRole.Create(&models.SysUserRole{RoleID: 1, UserID: 1})

	err := ensureSysRoleCanDelete(q, []uint64{1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "角色已绑定用户")
}

func TestEnsureSysRoleCanDelete_HasMenus(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Test", Code: "test"})
	q.SysRoleMenu.Create(&models.SysRoleMenu{RoleID: 1, MenuID: 1})

	err := ensureSysRoleCanDelete(q, []uint64{1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "角色已绑定菜单权限")
}

func TestEnsureSysRoleCanDelete_HasApis(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Test", Code: "test"})
	q.SysRoleApi.Create(&models.SysRoleApi{RoleID: 1, ApiID: 1})

	err := ensureSysRoleCanDelete(q, []uint64{1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "角色已绑定API权限")
}

// --- SavePermissions ---

func TestSysRoleHandler_SavePermissions_Success(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockCasbin := mock.NewMockCasbinService(ctrl)
	h := NewSysRoleHandler(q, mockCasbin)

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Admin", Code: "admin", IsEnabled: mixin.IsEnabled{IsEnabled: true}})

	mockCasbin.EXPECT().SyncRoleAPIPermissions("admin", gomock.Any(), gomock.Any()).Return(nil)

	err := h.SavePermissions(newTestCtx(t), &ReqSysRolePermissionSave{
		ID: 1, MenuIDs: []uint64{}, ApiIDs: []uint64{},
	})
	assert.NoError(t, err)
}

func TestSysRoleHandler_SavePermissions_NotFound(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	err := h.SavePermissions(newTestCtx(t), &ReqSysRolePermissionSave{
		ID: 99,
	})
	assert.Error(t, err)
}

// --- Update ---

func TestSysRoleHandler_Update_Success(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Admin", Code: "admin", IsEnabled: mixin.IsEnabled{IsEnabled: true}})

	name := "Admin Updated"
	err := h.Update(newTestCtx(t), &ReqSysRoleUpdate{ID: 1, Name: &name})
	assert.NoError(t, err)
}

func TestSysRoleHandler_Update_NotFound(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	name := "X"
	err := h.Update(newTestCtx(t), &ReqSysRoleUpdate{ID: 99, Name: &name})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "角色不存在")
}

func TestSysRoleHandler_Update_TriggersCasbinSync(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockCasbin := mock.NewMockCasbinService(ctrl)
	h := NewSysRoleHandler(q, mockCasbin)

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Admin", Code: "admin", IsEnabled: mixin.IsEnabled{IsEnabled: true}})
	code := "admin2"
	mockCasbin.EXPECT().SyncRoleState(uint64(1), "admin", true, "admin2", true).Return(nil)

	err := h.Update(newTestCtx(t), &ReqSysRoleUpdate{ID: 1, Code: &code})
	assert.NoError(t, err)
}

func TestEnsureSysRoleParentNotDescendant_SelfOrChild(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	parentID := uint64(1)
	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Parent", Code: "parent"})
	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 2}, Name: "Child", Code: "child", ParentID: &parentID})

	err := ensureSysRoleParentNotDescendant(q, 1, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "父级角色不能选择自身或子角色")
}

func TestEnsureSysRoleParentNotDescendant_MissingParent(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	err := ensureSysRoleParentNotDescendant(q, 1, 99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "父级角色不存在")
}

func TestEnsureSysRoleMenuIDsExist(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	q.SysResourceMenu.Create(&models.SysResourceMenu{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Menu", MenuType: MenuTypeCatalog, Path: "/menu"})
	assert.NoError(t, ensureSysRoleMenuIDsExist(q, []uint64{1}))
	assert.Error(t, ensureSysRoleMenuIDsExist(q, []uint64{1, 99}))
}

func TestEnsureSysRoleAPIIDsExist(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	q.SysResourceApi.Create(&models.SysResourceApi{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Path: "/api/test", Method: "GET"})
	assert.NoError(t, ensureSysRoleAPIIDsExist(q, []uint64{1}))
	assert.Error(t, ensureSysRoleAPIIDsExist(q, []uint64{1, 99}))
}

func TestNormalizeSysRoleParentID_Nil(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	parentID, err := normalizeSysRoleParentID(q, 0, nil)
	assert.NoError(t, err)
	assert.Nil(t, parentID)
}

func TestNormalizeSysRoleParentID_Zero(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	zero := uint64(0)
	parentID, err := normalizeSysRoleParentID(q, 0, &zero)
	assert.NoError(t, err)
	assert.Nil(t, parentID)
}

func TestNormalizeSysRoleParentID_Self(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	self := uint64(1)
	_, err := normalizeSysRoleParentID(q, 1, &self)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "父级角色不能选择自身")
}

func TestNormalizeSysRoleParentID_Missing(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	missing := uint64(99)
	_, err := normalizeSysRoleParentID(q, 0, &missing)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "父级角色不存在")
}

func TestNormalizeSysRoleParentID_Success(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 2}, Name: "Parent", Code: "parent"})
	parent := uint64(2)
	got, err := normalizeSysRoleParentID(q, 0, &parent)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, uint64(2), *got)
}

func TestSysRoleHandler_Del_NotFoundRowsAffected(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "Only", Code: "only"})
	err := h.Del(newTestCtx(t), &ReqSysRoleBatchDelete{IDs: []uint64{1, 99}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "角色不存在")
}

func TestSysRoleHandler_Del_Success(t *testing.T) {
	q := mustMigrateRole(t)
	query.SetDefault(q.SysRole.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewSysRoleHandler(q, mock.NewMockCasbinService(ctrl))

	q.SysRole.Create(&models.SysRole{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Name: "A", Code: "a"})
	err := h.Del(newTestCtx(t), &ReqSysRoleBatchDelete{IDs: []uint64{1}})
	assert.NoError(t, err)
}

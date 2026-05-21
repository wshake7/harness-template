package logic

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func TestSysResourceMenuHandler_List_Empty(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	ctx := newTestCtx(t)
	result, err := h.List(ctx, &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysResourceMenuHandler_List_WithData(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	q.SysResourceMenu.Create(&models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Dashboard",
		MenuType:        MenuTypeCatalog,
		Path:            "/dashboard",
	})

	ctx := newTestCtx(t)
	result, err := h.List(ctx, &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
}

func TestSysResourceMenuHandler_Tree_NoRoles(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	ctx := newTestCtx(t)
	ctx.SessionInfo.RoleIDs = nil
	result, err := h.Tree(ctx)
	assert.NoError(t, err)
	assert.Len(t, *result, 0)
}

func TestSysResourceMenuHandler_Tree_WithData(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	// Create a catalog menu
	q.SysResourceMenu.Create(&models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Dashboard",
		MenuType:        MenuTypeCatalog,
		Path:            "/dashboard",
		IsEnabled:       mixin.IsEnabled{IsEnabled: true},
		SortOrder:       mixin.SortOrder{SortOrder: 1},
	})
	// Give role 1 access
	q.SysRoleMenu.Create(&models.SysRoleMenu{
		RoleID: 1,
		MenuID: 1,
	})

	ctx := newTestCtx(t)
	ctx.SessionInfo.RoleIDs = []uint64{1}
	result, err := h.Tree(ctx)
	assert.NoError(t, err)
	assert.Len(t, *result, 1)
	assert.Equal(t, "Dashboard", (*result)[0].Name)
}

func TestSysResourceMenuHandler_Del_Success(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	q.SysResourceMenu.Create(&models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Test",
		MenuType:        MenuTypeButton,
		Path:            "/test",
		ParentID:        nil,
	})

	err := h.Del(newTestCtx(t), &ReqResourceMenuBatchDelete{IDs: []uint64{1}})
	assert.NoError(t, err)
}

func TestSysResourceMenuHandler_Del_HasChildren(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	q.SysResourceMenu.Create(&models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Parent",
		MenuType:        MenuTypeCatalog,
		Path:            "/parent",
	})
	parentID := uint64(1)
	q.SysResourceMenu.Create(&models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: 2},
		Name:            "Child",
		MenuType:        MenuTypeMenu,
		Path:            "/child",
		ParentID:        &parentID,
	})

	err := h.Del(newTestCtx(t), &ReqResourceMenuBatchDelete{IDs: []uint64{1}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "子节点")
}

func TestSysResourceMenuHandler_Create_Catalog(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	err := h.Create(newTestCtx(t), &ReqResourceMenuCreate{
		Name:     "Dashboard",
		MenuType: MenuTypeCatalog,
	})
	assert.NoError(t, err)
}

func TestSysResourceMenuHandler_Create_InvalidType(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	err := h.Create(newTestCtx(t), &ReqResourceMenuCreate{
		Name: "Test", MenuType: "INVALID",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "菜单类型无效")
}

func TestSysResourceMenuHandler_Create_NoParentForMenu(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	err := h.Create(newTestCtx(t), &ReqResourceMenuCreate{
		Name: "Users", MenuType: MenuTypeMenu, Path: "/users", Component: "Users",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "上级菜单不能为空")
}

func TestSysResourceMenuHandler_Create_WithParent(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	// Create parent catalog first
	q.SysResourceMenu.Create(&models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "System",
		MenuType:        MenuTypeCatalog,
		Path:            "/system",
	})

	parentID := uint64(1)
	err := h.Create(newTestCtx(t), &ReqResourceMenuCreate{
		Name: "Users", MenuType: MenuTypeMenu, Path: "/users", Component: "Users",
		ParentID: &parentID,
	})
	assert.NoError(t, err)
}

func TestSysResourceMenuHandler_Update_NotFound(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	err := h.Update(newTestCtx(t), &ReqResourceMenuUpdate{ID: 99})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "菜单资源不存在")
}

func TestSysResourceMenuHandler_Update_Success(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	q.SysResourceMenu.Create(&models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		Name:            "Dashboard",
		MenuType:        MenuTypeCatalog,
		Path:            "/dashboard",
	})

	name := "Dashboard v2"
	err := h.Update(newTestCtx(t), &ReqResourceMenuUpdate{ID: 1, Name: &name})
	assert.NoError(t, err)
}

func TestSysResourceMenuHandler_Del_NotFound(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	h := NewSysResourceMenuHandler(q)

	err := h.Del(newTestCtx(t), &ReqResourceMenuBatchDelete{IDs: []uint64{99}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "菜单资源不存在")
}

// --- Pure function tests ---

func TestValidateResourceMenuValues(t *testing.T) {
	assert.NoError(t, validateResourceMenuValues(MenuTypeCatalog, "Test", "", ""))
	assert.NoError(t, validateResourceMenuValues(MenuTypeButton, "Btn", "", ""))
	assert.Error(t, validateResourceMenuValues("INVALID", "N", "", ""))
	assert.Error(t, validateResourceMenuValues(MenuTypeMenu, "", "", ""))
	assert.Error(t, validateResourceMenuValues(MenuTypeMenu, "Test", "", ""))
	assert.Error(t, validateResourceMenuValues(MenuTypeMenu, "Test", "/test", ""))
}

func TestCanAssociateResourceMenuAPIs(t *testing.T) {
	assert.True(t, canAssociateResourceMenuAPIs(MenuTypeMenu))
	assert.True(t, canAssociateResourceMenuAPIs(MenuTypeButton))
	assert.False(t, canAssociateResourceMenuAPIs(MenuTypeCatalog))
	assert.False(t, canAssociateResourceMenuAPIs(MenuTypeEmbedded))
}

func TestCompactStrings(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, compactStrings([]string{" a ", " b ", "", " a "}))
	assert.Equal(t, []string{}, compactStrings(nil))
}

func TestBoolFromAny(t *testing.T) {
	assert.True(t, boolFromAny(true))
	assert.True(t, boolFromAny("true"))
	assert.True(t, boolFromAny("1"))
	assert.False(t, boolFromAny(false))
	assert.False(t, boolFromAny("false"))
	assert.False(t, boolFromAny(nil))
}

func TestStringFromAny(t *testing.T) {
	assert.Equal(t, "hello", stringFromAny("hello"))
	assert.Equal(t, "", stringFromAny(123))
}

func TestSameOptionalUint64(t *testing.T) {
	a := uint64(1)
	b := uint64(1)
	c := uint64(2)
	assert.True(t, sameOptionalUint64(nil, nil))
	assert.True(t, sameOptionalUint64(&a, &b))
	assert.False(t, sameOptionalUint64(&a, nil))
	assert.False(t, sameOptionalUint64(&a, &c))
}

// --- Additional pure function tests ---

func TestInt32FromAny(t *testing.T) {
	assert.Equal(t, int32(42), int32FromAny(42))
	assert.Equal(t, int32(42), int32FromAny(int32(42)))
	assert.Equal(t, int32(42), int32FromAny(int64(42)))
	assert.Equal(t, int32(42), int32FromAny(float64(42.9)))
	assert.Equal(t, int32(42), int32FromAny("42"))
	assert.Equal(t, int32(0), int32FromAny("abc"))
	assert.Equal(t, int32(0), int32FromAny(nil))
	assert.Equal(t, int32(0), int32FromAny(true))
}

func TestStringsFromAny_More(t *testing.T) {
	assert.Equal(t, []string{}, stringsFromAny(nil))
	assert.Equal(t, []string{"hello"}, stringsFromAny("hello"))
	assert.Equal(t, []string{"a", "b"}, stringsFromAny([]string{"a", "b"}))
	assert.Equal(t, []string{"x"}, stringsFromAny([]any{"x"}))
	assert.Equal(t, []string{}, stringsFromAny(""))
}

func TestBoolFromAny_More(t *testing.T) {
	assert.True(t, boolFromAny("true"))
	assert.True(t, boolFromAny("1"))
	assert.True(t, boolFromAny(true))
	assert.True(t, boolFromAny(float64(3.14)))
	assert.True(t, boolFromAny(int(1)))
	assert.False(t, boolFromAny("false"))
	assert.False(t, boolFromAny("0"))
	assert.False(t, boolFromAny(false))
	assert.False(t, boolFromAny(int32(0)))
	assert.False(t, boolFromAny(nil))
}

func TestNormalizeResourceMenuMetadata(t *testing.T) {
	m := datatypes.JSONMap{"hidden": true, "authorities": []string{"admin"}}
	normalizeResourceMenuMetadata(m)
	assert.Equal(t, true, boolFromAny(m["hidden"]))
	assert.NotNil(t, m["authorities"])
}

func TestNormalizeResourceMenuMetadata_Nil(t *testing.T) {
	normalizeResourceMenuMetadata(nil) // no panic
}

func TestEnsureResourceMenuParentNotDescendant_NilParent(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	assert.NoError(t, ensureResourceMenuParentNotDescendant(q, 1, nil))
}

func TestEnsureResourceMenuParentNotDescendant_MissingParent(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	parentID := uint64(99)
	err := ensureResourceMenuParentNotDescendant(q, 1, &parentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "父级菜单不存在")
}

func TestEnsureResourceMenuParentNotDescendant_Descendant(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	parentID := uint64(3)
	treePath := "/1/2/3/"
	q.SysResourceMenu.Create(&models.SysResourceMenu{AutoIncrementID: mixin.AutoIncrementID{ID: 3}, Name: "Child", MenuType: MenuTypeMenu, Path: "/child", TreePath: &treePath})
	err := ensureResourceMenuParentNotDescendant(q, 2, &parentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "父级不能选择自身或子节点")
}

func TestEnsureResourceMenuAPIIDsExist(t *testing.T) {
	q := mustMigrateResourceMenu(t)
	query.SetDefault(q.SysResourceMenu.UnderlyingDB())
	q.SysResourceApi.Create(&models.SysResourceApi{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, Path: "/api/test", Method: "GET"})
	assert.NoError(t, ensureResourceMenuAPIIDsExist(q, nil))
	assert.NoError(t, ensureResourceMenuAPIIDsExist(q, []uint64{1}))
	assert.Error(t, ensureResourceMenuAPIIDsExist(q, []uint64{1, 99}))
}

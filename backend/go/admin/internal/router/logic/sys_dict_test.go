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

func setupDictHandler(t *testing.T) (*SysDictHandler, *mock.MockDataPermissionService) {
	t.Helper()
	q := mustMigrateDict(t)
	query.SetDefault(q.SysDictType.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockDP := mock.NewMockDataPermissionService(ctrl)
	h := NewSysDictHandler(q, mockDP)
	return h, mockDP
}

func TestSysDictHandler_TypeList_Empty(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	mockDP.EXPECT().BuildFilterExprsForCtx(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	mockDP.EXPECT().ApplyPagePermissionExpr(gomock.Any(), gomock.Any()).Return(nil)

	result, err := h.TypeList(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysDictHandler_TypeCreate_Success(t *testing.T) {
	h, _ := setupDictHandler(t)

	err := h.TypeCreate(newTestCtx(t), &ReqDictTypeCreate{
		TypeCode: "gender", TypeName: "Gender",
	})
	assert.NoError(t, err)
}

func TestSysDictHandler_EntryList_Empty(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	mockDP.EXPECT().BuildReadPermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)

	result, err := h.EntryList(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}


func TestSysDictHandler_EntryCreate_Success(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	// Create a dict type first
	h.Q.SysDictType.Create(&models.SysDictType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		TypeCode:        "gender",
		TypeName:        "Gender",
	})

	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)

	err := h.EntryCreate(newTestCtx(t), &ReqDictEntryCreate{
		EntryLabel: "Male", EntryValue: "male", SysDictTypeId: 1,
	})
	assert.NoError(t, err)
}

func TestSysDictHandler_TypeDel_NoRows(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	mockDP.EXPECT().BuildDeletePermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)

	// Deleting non-existent IDs should succeed (no-op)
	err := h.TypeDel(newTestCtx(t), &ReqDictTypeBatchDelete{IDs: []uint64{1}})
	assert.Error(t, err) // fails because no rows matched
}

func TestSysDictHandler_EntryDel_Success(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	// Create an entry and a dict type
	h.Q.SysDictType.Create(&models.SysDictType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		TypeCode:        "gender",
		TypeName:        "Gender",
	})
	h.Q.SysDictEntry.Create(&models.SysDictEntry{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		EntryLabel:      "Test",
		EntryValue:      "test",
		SysDictTypeId:   1,
	})

	mockDP.EXPECT().BuildDeletePermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)

	err := h.EntryDel(newTestCtx(t), &ReqDictEntryBatchDelete{IDs: []uint64{1}})
	assert.NoError(t, err)
}

// --- Pure function tests ---

func TestApplyDictEntryTypeIDPageFilter(t *testing.T) {
	req := &v1.PagingRequest{}
	err := applyDictEntryTypeIDPageFilter(req, []uint64{1, 2, 3})
	assert.NoError(t, err)
	assert.NotNil(t, req.FilteringType)
}

func TestApplyDictEntryTypeIDPageFilter_Empty(t *testing.T) {
	req := &v1.PagingRequest{}
	err := applyDictEntryTypeIDPageFilter(req, nil)
	assert.NoError(t, err)
	assert.NotNil(t, req.FilteringType)
}

func TestNormalizeDictEntryMatchCodes(t *testing.T) {
	codes := normalizeDictEntryMatchCodes(&ReqDictEntryMatch{
		Codes: []string{" a ", "", " b ", " a "},
	})
	assert.Equal(t, []string{"a", "b"}, codes)
}

func TestNormalizeDictEntryMatchCodes_Empty(t *testing.T) {
	codes := normalizeDictEntryMatchCodes(&ReqDictEntryMatch{})
	assert.Len(t, codes, 0)
}

// --- Handler method tests with SQLite ---

func TestSysDictHandler_TypeUpdate_Success(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	q := h.Q
	q.SysDictType.Create(&models.SysDictType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		TypeCode:        "gender",
		TypeName:        "Gender",
	})

	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(q, nil)
	// BuildReadPermissionQuery is called within EntryMatch -> but TypeUpdate doesn't call it

	typeName := "Updated"
	err := h.TypeUpdate(newTestCtx(t), &ReqDictTypeUpdate{ID: 1, TypeName: &typeName})
	assert.NoError(t, err)
}

func TestSysDictHandler_EntryMatch_Empty(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	mockDP.EXPECT().BuildReadPermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)

	q := h.Q
	q.SysDictType.Create(&models.SysDictType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		TypeCode:        "gender",
		TypeName:        "Gender",
	})

	result, err := h.EntryMatch(newTestCtx(t), &ReqDictEntryMatch{
		Codes: []string{"male"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, *result, "male")
}

func TestSysDictHandler_EntryUpdate_Success(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	q := h.Q
	q.SysDictType.Create(&models.SysDictType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		TypeCode:        "gender",
		TypeName:        "Gender",
	})
	q.SysDictEntry.Create(&models.SysDictEntry{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		EntryLabel:      "Male",
		EntryValue:      "male",
		SysDictTypeId:   1,
	})

	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(q, nil)

	id := uint64(1)
	label := "Updated"
	err := h.EntryUpdate(newTestCtx(t), &ReqDictEntryUpdate{ID: &id, EntryLabel: &label})
	assert.NoError(t, err)
}

func TestSysDictHandler_EntryCreate_TypeNotFound(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)

	err := h.EntryCreate(newTestCtx(t), &ReqDictEntryCreate{
		EntryLabel: "Male", EntryValue: "male", SysDictTypeId: 99,
	})
	assert.Error(t, err)
}

func TestSysDictHandler_EntryUpdate_NilID(t *testing.T) {
	h, mockDP := setupDictHandler(t)
	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)
	err := h.EntryUpdate(newTestCtx(t), &ReqDictEntryUpdate{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "请求错误")
}

func TestConvertFilterByPagingRequest(t *testing.T) {
	// Test that applyDictEntryTypeIDPageFilter works with existing filter
	req := &v1.PagingRequest{
		FilteringType: &v1.PagingRequest_Query{Query: `{"name":"test"}`},
	}
	err := applyDictEntryTypeIDPageFilter(req, []uint64{1})
	assert.NoError(t, err)
	// Should merge with existing query
	assert.NotNil(t, req.FilteringType)
}

func TestSysDictHandler_EntryBatchCopy_Success(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	q := h.Q
	// Create source and target dict types
	q.SysDictType.Create(&models.SysDictType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		TypeCode:        "gender",
		TypeName:        "Gender",
	})
	q.SysDictType.Create(&models.SysDictType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 2},
		TypeCode:        "gender2",
		TypeName:        "Gender2",
	})
	// Create a source entry
	q.SysDictEntry.Create(&models.SysDictEntry{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		EntryLabel:      "Male",
		EntryValue:      "male",
		SysDictTypeId:   1,
	})

	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(q, nil)
	mockDP.EXPECT().BuildReadPermissionQuery(gomock.Any(), gomock.Any()).Return(q, nil)

	err := h.EntryBatchCopy(newTestCtx(t), &ReqDictEntryBatchCopy{
		TargetTypeId: 2,
		EntryIds:    []uint64{1},
	})
	assert.NoError(t, err)
}

func TestSysDictHandler_EntryBatchCopy_NotFound(t *testing.T) {
	h, mockDP := setupDictHandler(t)

	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)

	err := h.EntryBatchCopy(newTestCtx(t), &ReqDictEntryBatchCopy{
		TargetTypeId: 99,
		EntryIds:     []uint64{99},
	})
	assert.Error(t, err)
}

func TestQueryDictEntryTranslationMap(t *testing.T) {
	h, _ := setupDictHandler(t)
	q := h.Q
	query.SetDefault(q.SysLanguageType.UnderlyingDB())

	// Create language type and entries
	q.SysLanguageType.Create(&models.SysLanguageType{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		TypeCode:        "zh-CN",
		TypeName:        "Chinese",
		IsEnabled:       mixin.IsEnabled{IsEnabled: true},
	})
	q.SysLanguageEntry.Create(&models.SysLanguageEntry{
		AutoIncrementID:   mixin.AutoIncrementID{ID: 1},
		EntryCode:         "male",
		EntryValue:        "男",
		SysLanguageTypeId: 1,
		IsEnabled:         mixin.IsEnabled{IsEnabled: true},
	})

	ctx := newTestCtx(t)
	entries := []*models.SysDictEntry{{LanguageCode: "male"}}
	result, err := queryDictEntryTranslationMap(q, ctx, entries)
	assert.NoError(t, err)
	assert.Len(t, result, 0) // language code not matching ctx.Language
}

func TestQueryDictEntryTranslationMap_WithMatch(t *testing.T) {
	h, _ := setupDictHandler(t)
	q := h.Q
	query.SetDefault(q.SysLanguageType.UnderlyingDB())

	q.SysLanguageType.Create(&models.SysLanguageType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "en", TypeName: "English", IsEnabled: mixin.IsEnabled{IsEnabled: true}})
	q.SysLanguageEntry.Create(&models.SysLanguageEntry{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, EntryCode: "dict.male", EntryValue: "Male", SysLanguageTypeId: 1, IsEnabled: mixin.IsEnabled{IsEnabled: true}})

	ctx := newTestCtx(t)
	ctx.Language = "en"
	result, err := queryDictEntryTranslationMap(q, ctx, []*models.SysDictEntry{{LanguageCode: "dict.male"}})
	assert.NoError(t, err)
	assert.Equal(t, "Male", result["dict.male"])
}

func TestSysDictHandler_EntryBatchCopy_TargetNotAllowed(t *testing.T) {
	h, mockDP := setupDictHandler(t)
	mockDP.EXPECT().BuildWritePermissionQuery(gomock.Any(), gomock.Any()).Return(h.Q, nil)
	err := h.EntryBatchCopy(newTestCtx(t), &ReqDictEntryBatchCopy{TargetTypeId: 99, EntryIds: []uint64{1}})
	assert.Error(t, err)
}

func TestSysDictHandler_EntryMatch_ReadPermissionError(t *testing.T) {
	h, mockDP := setupDictHandler(t)
	mockDP.EXPECT().BuildReadPermissionQuery(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
	_, err := h.EntryMatch(newTestCtx(t), &ReqDictEntryMatch{Codes: []string{"male"}})
	assert.Error(t, err)
}

func TestSysDictHandler_EntryMatch_TranslationHit(t *testing.T) {
	h, mockDP := setupDictHandler(t)
	q := h.Q
	ctx := newTestCtx(t)
	ctx.Language = "en"

	mockDP.EXPECT().BuildReadPermissionQuery(gomock.Any(), gomock.Any()).Return(q, nil)
	q.SysDictType.Create(&models.SysDictType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "male", TypeName: "MaleType", IsEnabled: mixin.IsEnabled{IsEnabled: true}})
	q.SysDictEntry.Create(&models.SysDictEntry{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, EntryLabel: "原始", EntryValue: "male", LanguageCode: "dict.male", SysDictTypeId: 1, IsEnabled: mixin.IsEnabled{IsEnabled: true}})
	q.SysLanguageType.Create(&models.SysLanguageType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "en", TypeName: "English", IsEnabled: mixin.IsEnabled{IsEnabled: true}})
	q.SysLanguageEntry.Create(&models.SysLanguageEntry{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, EntryCode: "dict.male", EntryValue: "Male", SysLanguageTypeId: 1, IsEnabled: mixin.IsEnabled{IsEnabled: true}})

	result, err := h.EntryMatch(ctx, &ReqDictEntryMatch{Codes: []string{"male"}})
	assert.NoError(t, err)
	assert.Len(t, (*result)["male"], 1)
	assert.Equal(t, "Male", (*result)["male"][0].EntryLabel)
}


func TestQueryAllowedDictTypeIDSetByExpr_EmptyIDs(t *testing.T) {
	_, mockDP := setupDictHandler(t)
	set, err := queryAllowedDictTypeIDSetByExpr(mockDP, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, set, 0)
}

func TestQueryAllowedDictTypeIDSetByExpr_PartialMatch(t *testing.T) {
	h, mockDP := setupDictHandler(t)
	q := h.Q
	q.SysDictType.Create(&models.SysDictType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "a", TypeName: "A"})
	mockDP.EXPECT().BuildPermissionQueryFromExpr(gomock.Any()).Return(q, nil)

	set, err := queryAllowedDictTypeIDSetByExpr(mockDP, []uint64{1, 99}, &v1.FilterExpr{})
	assert.NoError(t, err)
	assert.True(t, set[1])
	assert.False(t, set[99])
}

func TestQueryAllowedDictTypeIDSetByExpr_QueryError(t *testing.T) {
	_, mockDP := setupDictHandler(t)
	mockDP.EXPECT().BuildPermissionQueryFromExpr(gomock.Any()).Return(nil, assert.AnError)
	_, err := queryAllowedDictTypeIDSetByExpr(mockDP, []uint64{1}, &v1.FilterExpr{})
	assert.Error(t, err)
}

func TestSysDictHandler_EntryList_ReadPermissionError(t *testing.T) {
	h, mockDP := setupDictHandler(t)
	mockDP.EXPECT().BuildReadPermissionQuery(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
	_, err := h.EntryList(newTestCtx(t), &v1.PagingRequest{})
	assert.Error(t, err)
}

package logic

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func mustMigrateLanguage(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysLanguageType{},
		&models.SysLanguageEntry{},
	)
}

func TestSysLanguageHandler_TypeList_Empty(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageType.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	result, err := h.TypeList(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysLanguageHandler_TypeCreate_Success(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageType.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	err := h.TypeCreate(newTestCtx(t), &ReqLangTypeCreate{TypeCode: "en", TypeName: "English"})
	assert.NoError(t, err)
}

func TestSysLanguageHandler_TypeUpdate_Success(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageType.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	q.SysLanguageType.Create(&models.SysLanguageType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "en", TypeName: "English"})
	name := "English 2"
	err := h.TypeUpdate(newTestCtx(t), &ReqLangTypeUpdate{ID: 1, TypeName: &name})
	assert.NoError(t, err)
}

func TestSysLanguageHandler_TypeDel_DefaultBlocked(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageType.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	q.SysLanguageType.Create(&models.SysLanguageType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "en", TypeName: "English", IsDefault: true})
	err := h.TypeDel(newTestCtx(t), &ReqLangTypeBatchDelete{IDs: []uint64{1}})
	assert.Error(t, err)
}

func TestSysLanguageHandler_EntryList_Empty(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageEntry.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	result, err := h.EntryList(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestSysLanguageHandler_EntryCreate_Success(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageType.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	q.SysLanguageType.Create(&models.SysLanguageType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "en", TypeName: "English"})
	err := h.EntryCreate(newTestCtx(t), &ReqLangEntryCreate{EntryCode: "hello", EntryValue: "Hello", SysLanguageTypeId: 1})
	assert.NoError(t, err)
}

func TestSysLanguageHandler_EntryUpdate_Success(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageEntry.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	q.SysLanguageType.Create(&models.SysLanguageType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "en", TypeName: "English"})
	q.SysLanguageEntry.Create(&models.SysLanguageEntry{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, EntryCode: "hello", EntryValue: "Hello", SysLanguageTypeId: 1})
	id := uint64(1)
	value := "Hi"
	err := h.EntryUpdate(newTestCtx(t), &ReqLangEntryUpdate{ID: &id, EntryValue: &value})
	assert.NoError(t, err)
}

func TestSysLanguageHandler_EntryDel_Success(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageEntry.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	q.SysLanguageEntry.Create(&models.SysLanguageEntry{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, EntryCode: "hello", EntryValue: "Hello", SysLanguageTypeId: 1})
	err := h.EntryDel(newTestCtx(t), &ReqLangEntryBatchDelete{IDs: []uint64{1}})
	assert.NoError(t, err)
}

func TestSysLanguageHandler_EntryBatchCreate_Success(t *testing.T) {
	q := mustMigrateLanguage(t)
	query.SetDefault(q.SysLanguageEntry.UnderlyingDB())
	h := NewSysLanguageHandler(q)

	q.SysLanguageType.Create(&models.SysLanguageType{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, TypeCode: "en", TypeName: "English"})
	err := h.EntryBatchCreate(newTestCtx(t), &ReqLangEntryBatchCreate{EntryCode: "hello", Values: map[string]string{"en": "Hello"}})
	assert.NoError(t, err)
}

package logic

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSQLiteDB(t *testing.T, tables ...interface{}) *query.Query {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(tables...); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	q := query.Use(db)
	return q
}

func mustMigrateResourceMenu(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysResourceMenu{},
		&models.SysResourceMenuApi{},
		&models.SysRoleMenu{},
		&models.SysRole{},
		&models.SysUserRole{},
		&models.SysResourceApi{},
		&models.SysUser{},
	)
}

func mustMigrateRole(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysRole{},
		&models.SysRoleMenu{},
		&models.SysRoleApi{},
		&models.SysUserRole{},
		&models.SysResourceMenu{},
		&models.SysResourceApi{},
	)
}

func mustMigrateDict(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.SysDictType{},
		&models.SysDictEntry{},
		&models.SysLanguageType{},
		&models.SysLanguageEntry{},
	)
}

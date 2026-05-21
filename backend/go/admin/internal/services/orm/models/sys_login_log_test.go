package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSysLoginLogAllowsNilUserID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:sys_login_log_nil_user?mode=memory&cache=shared"), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.AutoMigrate(&SysUser{}, &SysLoginLog{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	loginLog := &SysLoginLog{
		Username:   "missing-user",
		LoginIP:    "127.0.0.1",
		StatusCode: 100,
		Success:    false,
		Reason:     "invalid username or password",
	}
	if err := db.Create(loginLog).Error; err != nil {
		t.Fatalf("create login log with nil user id failed: %v", err)
	}
	if loginLog.SysUserID != nil {
		t.Fatalf("expected nil sys_user_id, got %d", *loginLog.SysUserID)
	}
	if loginLog.LoginTime == nil {
		t.Fatal("expected login time to be filled")
	}
}

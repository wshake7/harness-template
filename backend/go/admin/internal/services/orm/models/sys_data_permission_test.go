package models

import (
	"testing"

	"gorm.io/datatypes"
)

func TestSysDataPermissionNormalizeAction(t *testing.T) {
	permission := &SysDataPermission{
		SubjectType:   DataPermissionSubjectRole,
		SubjectID:     1,
		ResourceTable: "sys_dict_type",
		Action:        datatypes.JSON([]byte(`["delete","read","write"]`)),
		ScopeType:     DataPermissionScopeAll,
	}

	if err := permission.NormalizeAndValidate(); err != nil {
		t.Fatalf("NormalizeAndValidate returned error: %v", err)
	}
	if string(permission.Action) != `["read","write","delete"]` {
		t.Fatalf("expected canonical action JSON, got %s", permission.Action)
	}
	if permission.ActionKey != "read,write,delete" {
		t.Fatalf("expected canonical action key, got %q", permission.ActionKey)
	}
}

func TestSysDataPermissionNormalizeActionAllWins(t *testing.T) {
	permission := &SysDataPermission{
		SubjectType:   DataPermissionSubjectAnyRole,
		SubjectID:     0,
		ResourceTable: "sys_dict_type",
		Action:        datatypes.JSON([]byte(`["write","all"]`)),
		ScopeType:     DataPermissionScopeAll,
	}

	if err := permission.NormalizeAndValidate(); err != nil {
		t.Fatalf("NormalizeAndValidate returned error: %v", err)
	}
	if string(permission.Action) != `["all"]` {
		t.Fatalf("expected all action JSON, got %s", permission.Action)
	}
	if permission.ActionKey != "all" {
		t.Fatalf("expected all action key, got %q", permission.ActionKey)
	}
}

func TestSysDataPermissionRejectsInvalidScopeValues(t *testing.T) {
	permission := &SysDataPermission{
		SubjectType:   DataPermissionSubjectRole,
		SubjectID:     1,
		ResourceTable: "sys_dict_type",
		Action:        datatypes.JSON([]byte(`["read"]`)),
		ScopeType:     DataPermissionScopeInclude,
		ScopeValues:   datatypes.JSON([]byte(`[]`)),
	}

	if err := permission.NormalizeAndValidate(); err == nil {
		t.Fatal("expected include scope with empty values to be rejected")
	}
}

func TestSysDataPermissionRejectsInvalidAction(t *testing.T) {
	permission := &SysDataPermission{
		SubjectType:   DataPermissionSubjectRole,
		SubjectID:     1,
		ResourceTable: "sys_dict_type",
		Action:        datatypes.JSON([]byte(`["read","export"]`)),
		ScopeType:     DataPermissionScopeAll,
	}

	if err := permission.NormalizeAndValidate(); err == nil {
		t.Fatal("expected invalid action to be rejected")
	}
}

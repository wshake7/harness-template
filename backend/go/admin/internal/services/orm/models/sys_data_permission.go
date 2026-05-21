package models

import (
	"errors"
	"slices"
	"strings"

	"orm-crud/gormc/mixin"

	"github.com/bytedance/sonic"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const (
	DataPermissionSubjectUser    = "USER"
	DataPermissionSubjectRole    = "ROLE"
	DataPermissionSubjectAnyUser = "ANY_USER"
	DataPermissionSubjectAnyRole = "ANY_ROLE"

	DataPermissionScopeAll     = "all"
	DataPermissionScopeNone    = "none"
	DataPermissionScopeInclude = "include"
	DataPermissionScopeExclude = "exclude"
	DataPermissionScopeCustom  = "custom"

	DataPermissionActionAll    = "all"
	DataPermissionActionRead   = "read"
	DataPermissionActionWrite  = "write"
	DataPermissionActionDelete = "delete"
)

var (
	validDataPermissionSubjects = map[string]struct{}{
		DataPermissionSubjectUser:    {},
		DataPermissionSubjectRole:    {},
		DataPermissionSubjectAnyUser: {},
		DataPermissionSubjectAnyRole: {},
	}
	validDataPermissionScopes = map[string]struct{}{
		DataPermissionScopeAll:     {},
		DataPermissionScopeNone:    {},
		DataPermissionScopeInclude: {},
		DataPermissionScopeExclude: {},
		DataPermissionScopeCustom:  {},
	}
	dataPermissionActionOrder = []string{
		DataPermissionActionAll,
		DataPermissionActionRead,
		DataPermissionActionWrite,
		DataPermissionActionDelete,
	}
)

func init() {
	Models = append(Models, &SysDataPermission{})
}

type SysDataPermission struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.Remark
	mixin.OperatorID
	mixin.IsEnabled
	DeletedAt     soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_data_permission_deleted_at" json:"deletedAt"`
	SubjectType   string                `gorm:"column:subject_type;type:varchar(16);not null;index:idx_sys_data_permission_subject;uniqueIndex:idx_sys_data_permission_subject_resource_action_active,priority:1,where:deleted_at = 0;comment:主体类型(USER/ROLE/ANY_USER/ANY_ROLE)" json:"subjectType"`
	SubjectID     uint64                `gorm:"column:subject_id;type:bigint;not null;index:idx_sys_data_permission_subject;uniqueIndex:idx_sys_data_permission_subject_resource_action_active,priority:2,where:deleted_at = 0;comment:主体ID，ANY_*时为0" json:"subjectID"`
	ResourceTable string                `gorm:"column:resource_table;type:varchar(32);not null;uniqueIndex:idx_sys_data_permission_subject_resource_action_active,priority:3,where:deleted_at = 0;comment:资源表名" json:"resourceTable"`
	Action        datatypes.JSON        `gorm:"column:action;not null;default:'[\"read\"]';comment:操作列表(all/read/write/delete)" json:"action"`
	ActionKey     string                `gorm:"column:action_key;type:varchar(64);not null;default:'read';uniqueIndex:idx_sys_data_permission_subject_resource_action_active,priority:4,where:deleted_at = 0;comment:规范化操作列表" json:"actionKey"`
	ScopeType     string                `gorm:"column:scope_type;type:varchar(32);not null;default:'none';comment:作用域类型(all/none/include/exclude/custom)" json:"scopeType"`
	ScopeField    string                `gorm:"column:scope_field;type:varchar(64);not null;default:'id';comment:用于匹配scope_values的字段" json:"scopeField"`
	ScopeValues   datatypes.JSON        `gorm:"column:scope_values;not null;default:'[]';comment:作用域值列表" json:"scopeValues"`
	Conditions    datatypes.JSONMap     `gorm:"column:conditions;not null;default:'{}';comment:行过滤条件" json:"conditions"`
	Priority      int                   `gorm:"column:priority;type:integer;not null;default:0;comment:多角色冲突时的优先级" json:"priority"`
}

func (*SysDataPermission) TableName() string {
	return "sys_data_permission"
}

func (m *SysDataPermission) BeforeSave(tx *gorm.DB) error {
	return m.NormalizeAndValidate()
}

func (m *SysDataPermission) NormalizeAndValidate() error {
	if _, ok := validDataPermissionSubjects[m.SubjectType]; !ok {
		return errors.New("invalid data permission subject_type")
	}
	if strings.HasPrefix(m.SubjectType, "ANY_") && m.SubjectID != 0 {
		return errors.New("data permission ANY_* subject_id must be 0")
	}
	if !strings.HasPrefix(m.SubjectType, "ANY_") && m.SubjectID == 0 {
		return errors.New("data permission USER/ROLE subject_id must be greater than 0")
	}

	if m.ScopeType == "" {
		m.ScopeType = DataPermissionScopeNone
	}
	if _, ok := validDataPermissionScopes[m.ScopeType]; !ok {
		return errors.New("invalid data permission scope_type")
	}
	if m.ScopeField == "" {
		m.ScopeField = "id"
	}
	if m.Conditions == nil {
		m.Conditions = datatypes.JSONMap{}
	}
	if len(m.ScopeValues) == 0 {
		m.ScopeValues = []byte("[]")
	}
	if err := m.normalizeAction(); err != nil {
		return err
	}
	if m.ScopeType == DataPermissionScopeInclude || m.ScopeType == DataPermissionScopeExclude {
		var values []any
		if err := sonic.Unmarshal(m.ScopeValues, &values); err != nil {
			return errors.New("invalid data permission scope_values")
		}
		if len(values) == 0 {
			return errors.New("data permission include/exclude scope_values cannot be empty")
		}
	}
	return nil
}

func (m *SysDataPermission) normalizeAction() error {
	if len(m.Action) == 0 {
		m.Action = []byte(`["read"]`)
	}
	var rawActions []string
	if err := sonic.Unmarshal(m.Action, &rawActions); err != nil {
		return errors.New("invalid data permission action")
	}
	if len(rawActions) == 0 {
		return errors.New("data permission action cannot be empty")
	}

	seen := make(map[string]struct{}, len(rawActions))
	for _, action := range rawActions {
		action = strings.TrimSpace(action)
		if action == "" {
			return errors.New("data permission action cannot contain empty value")
		}
		seen[action] = struct{}{}
	}
	for action := range seen {
		if !isValidDataPermissionAction(action) {
			return errors.New("invalid data permission action")
		}
	}

	canonical := make([]string, 0, len(seen))
	if _, ok := seen[DataPermissionActionAll]; ok {
		canonical = append(canonical, DataPermissionActionAll)
	} else {
		for _, action := range dataPermissionActionOrder[1:] {
			if _, ok := seen[action]; ok {
				canonical = append(canonical, action)
			}
		}
	}

	actionBytes, err := sonic.Marshal(canonical)
	if err != nil {
		return err
	}
	m.Action = actionBytes
	m.ActionKey = strings.Join(canonical, ",")
	return nil
}

func isValidDataPermissionAction(action string) bool {
	return slices.Contains(dataPermissionActionOrder, action)
}

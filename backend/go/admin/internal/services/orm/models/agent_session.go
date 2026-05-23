package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &AgentSession{})
}

// AgentSession agent 会话，对应用户的一次对话
type AgentSession struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.Remark
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_agent_session_deleted_at" json:"deletedAt"`
	SessionID string                `gorm:"column:session_id;type:varchar(64);not null;uniqueIndex:idx_agent_session_session_id_active,where:deleted_at = 0;comment:会话唯一标识" json:"sessionID"`
	Title     string                `gorm:"column:title;type:varchar(255);default:'';comment:会话标题" json:"title"`
	Status    string                `gorm:"column:status;type:varchar(32);not null;default:'active';comment:会话状态 active/archived" json:"status"`
	Metadata  datatypes.JSONMap     `gorm:"column:metadata;default:'{}';comment:扩展元数据" json:"metadata"`
}

func (*AgentSession) TableName() string {
	return "agent_session"
}

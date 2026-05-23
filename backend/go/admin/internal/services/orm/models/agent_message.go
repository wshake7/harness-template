package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &AgentMessage{})
}

// AgentMessage agent 会话中的单条消息
type AgentMessage struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	DeletedAt  soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_agent_message_deleted_at" json:"deletedAt"`
	SessionID  uint64                `gorm:"column:session_id;type:bigint;not null;index:idx_agent_message_session_id;comment:所属会话ID" json:"sessionID"`
	Role       string                `gorm:"column:role;type:varchar(32);not null;comment:角色 user/assistant/system/tool" json:"role"`
	Content    string                `gorm:"column:content;type:text;not null;comment:消息内容" json:"content"`
	TokenCount int                   `gorm:"column:token_count;type:int;not null;default:0;comment:Token 数量" json:"tokenCount"`
	Metadata   datatypes.JSONMap     `gorm:"column:metadata;default:'{}';comment:扩展元数据" json:"metadata"`

	Session *AgentSession `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"session"`
}

func (*AgentMessage) TableName() string {
	return "agent_message"
}

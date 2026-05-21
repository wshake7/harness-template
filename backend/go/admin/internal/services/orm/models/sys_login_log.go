package models

import (
	"time"

	"orm-crud/gormc/mixin"

	"gorm.io/gorm"
)

func init() {
	Models = append(Models, &SysLoginLog{})
}

type SysLoginLog struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	Username       string     `gorm:"column:username;type:varchar(64);not null;default:'';index:idx_sys_login_log_username;comment:登录账号" json:"username"`
	LoginIP        string     `gorm:"column:login_ip;type:varchar(64);not null;default:'';index:idx_sys_login_log_login_ip;comment:登录IP地址" json:"loginIP"`
	LoginMAC       string     `gorm:"column:login_mac;type:varchar(128);not null;default:'';comment:登录MAC地址" json:"loginMAC"`
	LoginTime      *time.Time `gorm:"column:login_time;index:idx_sys_login_log_login_time;index:idx_sys_login_log_success_login_time,priority:2;index:idx_sys_login_log_status_login_time,priority:2;comment:登录时间" json:"loginTime"`
	UserAgent      string     `gorm:"column:user_agent;type:text;default:'';comment:浏览器的用户代理信息" json:"userAgent"`
	BrowserName    string     `gorm:"column:browser_name;type:varchar(128);default:'';comment:浏览器名称" json:"browserName"`
	BrowserVersion string     `gorm:"column:browser_version;type:varchar(128);default:'';comment:浏览器版本" json:"browserVersion"`
	ClientID       string     `gorm:"column:client_id;type:varchar(128);default:'';comment:客户端ID" json:"clientID"`
	ClientName     string     `gorm:"column:client_name;type:varchar(128);default:'';comment:客户端名称" json:"clientName"`
	OSName         string     `gorm:"column:os_name;type:varchar(128);default:'';comment:操作系统名称" json:"osName"`
	OSVersion      string     `gorm:"column:os_version;type:varchar(128);default:'';comment:操作系统版本" json:"osVersion"`
	SysUserID      *uint64    `gorm:"column:sys_user_id;type:bigint;default:null;index:idx_sys_login_log_sys_user_id;index:idx_sys_login_log_user_login_time,priority:1;comment:登录用户ID" json:"sysUserID"`
	StatusCode     int32      `gorm:"column:status_code;type:int;default:0;index:idx_sys_login_log_status_code;index:idx_sys_login_log_status_login_time,priority:1;comment:状态码" json:"statusCode"`
	Success        bool       `gorm:"column:success;type:boolean;default:false;index:idx_sys_login_log_success_login_time,priority:1;comment:登录成功" json:"success"`
	Reason         string     `gorm:"column:reason;type:varchar(255);default:'';comment:登录失败原因" json:"reason"`
	Location       string     `gorm:"column:location;type:varchar(255);default:'';comment:登录地理位置" json:"location"`

	SysUser *SysUser `gorm:"foreignKey:SysUserID;references:ID" json:"sysUser"`
}

func (m *SysLoginLog) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt.CreatedAt == nil {
		now := time.Now()
		m.CreatedAt.CreatedAt = &now
	}
	if m.LoginTime == nil {
		now := time.Now()
		m.LoginTime = &now
	}
	return nil
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}

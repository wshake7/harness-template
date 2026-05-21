package models

import (
	"time"

	"orm-crud/gormc/mixin"

	"gorm.io/gorm"
)

func init() {
	Models = append(Models, &SysApiLog{})
}

type SysApiLog struct {
	mixin.AutoIncrementID
	CreatedAt      *time.Time `gorm:"column:created_at;index:idx_sys_api_log_created_at;index:idx_sys_api_log_user_created_at,priority:2;index:idx_sys_api_log_success_created_at,priority:2;index:idx_sys_api_log_status_created_at,priority:2" json:"createdAt"`
	RequestID      string     `gorm:"column:request_id;type:varchar(128);not null;uniqueIndex:idx_sys_api_log_request_id;comment:请求ID" json:"requestID"`
	Method         string     `gorm:"column:method;type:varchar(16);not null;index:idx_sys_api_log_method_path;comment:请求方法" json:"method"`
	Module         string     `gorm:"column:module;type:varchar(255);default:'';comment:模块" json:"module"`
	Path           string     `gorm:"column:path;type:varchar(255);not null;index:idx_sys_api_log_method_path;comment:请求路径" json:"path"`
	Referer        string     `gorm:"column:referer;type:text;default:'';comment:请求源" json:"referer"`
	BeforeChange   string     `gorm:"column:before_change;type:text;default:'';comment:变更前内容" json:"beforeChange"`
	AfterChange    string     `gorm:"column:after_change;type:text;default:'';comment:变更后内容" json:"afterChange"`
	FormatChange   string     `gorm:"column:format_change;type:text;default:'';comment:格式化变化内容" json:"formatChange"`
	RequestURI     string     `gorm:"column:request_uri;type:text;default:'';comment:请求URI" json:"requestURI"`
	RequestBody    string     `gorm:"column:request_body;type:text;default:'';comment:请求体" json:"requestBody"`
	RequestHeader  string     `gorm:"column:request_header;type:text;default:'';comment:请求头" json:"requestHeader"`
	Response       string     `gorm:"column:response;type:text;default:'';comment:响应信息" json:"response"`
	CostTime       int64      `gorm:"column:cost_time;type:bigint;default:0;comment:操作耗时" json:"costTime"`
	SysUserID      *uint64    `gorm:"column:sys_user_id;type:bigint;default:null;index:idx_sys_api_log_sys_user_id;index:idx_sys_api_log_user_created_at,priority:1;comment:操作者用户ID" json:"sysUserID"`
	ClientIP       string     `gorm:"column:client_ip;type:varchar(64);default:'';comment:操作者IP" json:"clientIP"`
	StatusCode     int        `gorm:"column:status_code;type:int;default:0;index:idx_sys_api_log_status_code;index:idx_sys_api_log_status_created_at,priority:1;comment:状态码" json:"statusCode"`
	Reason         string     `gorm:"column:reason;type:varchar(255);default:'';comment:操作失败原因" json:"reason"`
	Success        bool       `gorm:"column:success;type:boolean;default:false;index:idx_sys_api_log_success_created_at,priority:1;comment:操作成功" json:"success"`
	Location       string     `gorm:"column:location;type:varchar(255);default:'';comment:操作地理位置" json:"location"`
	UserAgent      string     `gorm:"column:user_agent;type:text;default:'';comment:浏览器的用户代理信息" json:"userAgent"`
	BrowserName    string     `gorm:"column:browser_name;type:varchar(128);default:'';comment:浏览器名称" json:"browserName"`
	BrowserVersion string     `gorm:"column:browser_version;type:varchar(128);default:'';comment:浏览器版本" json:"browserVersion"`
	ClientID       string     `gorm:"column:client_id;type:varchar(128);default:'';comment:客户端ID" json:"clientID"`
	ClientName     string     `gorm:"column:client_name;type:varchar(128);default:'';comment:客户端名称" json:"clientName"`
	OSName         string     `gorm:"column:os_name;type:varchar(128);default:'';comment:操作系统名称" json:"oSName"`
	OSVersion      string     `gorm:"column:os_version;type:varchar(128);default:'';comment:操作系统版本" json:"oSVersion"`

	SysUser *SysUser `gorm:"foreignKey:SysUserID;references:ID" json:"sysUser"`
}

func (m *SysApiLog) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt == nil {
		now := time.Now()
		m.CreatedAt = &now
	}
	return nil
}

func (SysApiLog) TableName() string {
	return "sys_api_log"
}

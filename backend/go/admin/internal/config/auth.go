package config

import (
	"github.com/click33/sa-token-go/core/config"
	"time"
)

type AuthConfig struct {
	// Timeout Token expiration time in seconds, -1 for never expire | Token超时时间（单位：秒，-1代表永不过期）
	Timeout int64 `mapstructure:"Timeout" default:"15552000"`

	// ActiveTimeout Token minimum activity frequency in seconds. If Token is not accessed for this time, it will be frozen. -1 means no limit | Token最低活跃频率（单位：秒），如果Token超过此时间没有访问，则会被冻结。-1代表不限制，永不冻结 默认2个月
	ActiveTimeout int64 `mapstructure:"ActiveTimeout" default:"5184000"`

	// TokenName Token name (also used as Cookie name) | Token名称（同时也是Cookie名称）
	TokenName string `mapstructure:"TokenName" default:"Token"`

	// MaxRefresh Threshold for triggering async token renewal (in seconds) | Token自动续期触发阈值（单位：秒，当剩余有效期低于该值时触发异步续期 -1或0代表不限制）
	MaxRefresh int64 `mapstructure:"MaxRefresh" default:"0"`

	// RenewInterval Minimum interval between token renewals (ms) | Token最小续期间隔（单位：秒，同一个Token在此时间内只会续期一次 -1或0代表不限制）
	RenewInterval int64 `mapstructure:"RenewInterval" default:"86400"`

	// IsConcurrent Allow concurrent login for the same account (true=allow concurrent login, false=new login kicks out old login) | 是否允许同一账号并发登录（为true时允许一起登录，为false时新登录挤掉旧登录）
	IsConcurrent bool `mapstructure:"IsConcurrent" default:"false"`

	// IsShare Share the same Token for concurrent logins (true=share one Token, false=create new Token for each login) | 在多人登录同一账号时，是否共用一个Token（为true时所有登录共用一个Token，为false时每次登录新建一个Token）
	IsShare bool `mapstructure:"IsShare" default:"false"`

	// MaxLoginCount Maximum number of concurrent logins for the same account, -1 means no limit (only effective when IsConcurrent=true and IsShare=false) | 同一账号最大登录数量，-1代表不限（只有在IsConcurrent=true，IsShare=false时此配置才有效）
	MaxLoginCount int `mapstructure:"MaxLoginCount" default:"-1"`

	// IsReadBody Try to read Token from request body (default: false) | 是否尝试从请求体里读取Token（默认：false）
	IsReadBody bool `mapstructure:"IsReadBody" default:"false"`

	// IsReadHeader Try to read Token from HTTP Header (default: true, recommended) | 是否尝试从Header里读取Token（默认：true，推荐）
	IsReadHeader bool `mapstructure:"IsReadHeader" default:"true"`

	// IsReadCookie Try to read Token from Cookie (default: false) | 是否尝试从Cookie里读取Token（默认：false）
	IsReadCookie bool `mapstructure:"IsReadCookie" default:"false"`

	// TokenStyle Token generation style | Token风格
	TokenStyle config.TokenStyle `mapstructure:"TokenStyle" default:"uuid"`

	// DataRefreshPeriod Auto-refresh period in seconds, -1 means no auto-refresh | 自动续签（单位：秒），-1代表不自动续签
	DataRefreshPeriod int64 `mapstructure:"DataRefreshPeriod" default:"30"`

	// TokenSessionCheckLogin Check if Token-Session is kicked out when logging in (true=check on login, false=skip check) | Token-Session在登录时是否检查（true=登录时验证是否被踢下线，false=不作此检查）
	TokenSessionCheckLogin bool `mapstructure:"TokenSessionCheckLogin" default:"true"`

	// AutoRenew Auto-renew Token expiration time on each validation | 是否自动续期（每次验证Token时，都会延长Token的有效期）
	AutoRenew bool `mapstructure:"AutoRenew" default:"true"`

	// JwtSecretKey JWT secret key (only effective when TokenStyle=JWT) | JWT密钥（只有TokenStyle=JWT时，此配置才生效）
	JwtSecretKey string `mapstructure:"JwtSecretKey"`

	// IsLog Enable API logging | 是否输出API日志
	IsLog bool `mapstructure:"IsLog" default:"false"`

	// IsPrintBanner Print startup banner (default: true) | 是否打印启动 Banner（默认：true）
	IsPrintBanner bool `mapstructure:"IsLog" default:"false"`

	// KeyPrefix Storage key prefix for Redis isolation (default: "satoken:") | 存储键前缀，用于Redis隔离（默认："satoken:"）
	// Set to empty "" to be compatible with Java sa-token default behavior | 设置为空""以兼容Java sa-token默认行为
	KeyPrefix string `mapstructure:"KeyPrefix" default:"token:"`

	// CookieConfig Cookie configuration | Cookie配置
	CookieConfig AuthCookieConfig

	// RenewPoolConfig Configuration for renewal pool manager | 续期池配置
	RenewPoolConfig AuthRenewPoolConfig
}

type AuthCookieConfig struct {
	// Domain Cookie domain | 作用域
	Domain string `mapstructure:"Domain" default:""`

	// Path Cookie path | 路径
	Path string `mapstructure:"Path" default:"/"`

	// Secure Only effective under HTTPS | 是否只在HTTPS下生效
	Secure bool `mapstructure:"Secure" default:"false"`

	// HttpOnly Prevent JavaScript access to Cookie | 是否禁止JS操作Cookie
	HttpOnly bool `mapstructure:"HttpOnly" default:"false"`

	// SameSite SameSite attribute (Strict, Lax, None) | SameSite属性（Strict、Lax、None）
	SameSite config.SameSiteMode `mapstructure:"SameSite" default:"Lax"`

	// MaxAge Cookie expiration time in seconds | 过期时间（单位：秒）
	MaxAge int `mapstructure:"MaxAge" default:"0"`
}

type AuthRenewPoolConfig struct {
	MinSize             int           `mapstructure:"MinSize" default:"100"`           // Minimum pool size | 最小协程数
	MaxSize             int           `mapstructure:"MaxSize" default:"2000"`          // Maximum pool size | 最大协程数
	ScaleUpRate         float64       `mapstructure:"ScaleUpRate" default:"0.8"`       // Scale-up threshold | 扩容阈值
	ScaleDownRate       float64       `mapstructure:"ScaleDownRate" default:"0.3"`     // Scale-down threshold | 缩容阈值
	CheckInterval       time.Duration `mapstructure:"CheckInterval" default:"60s"`     // Auto-scale check interval | 检查间隔
	Expiry              time.Duration `mapstructure:"Expiry" default:"10s"`            // Idle worker expiry duration | 空闲协程过期时间
	PrintStatusInterval time.Duration `mapstructure:"PrintStatusInterval" default:"0"` // Interval for periodic status printing (0 = disabled) | 定时打印池状态的间隔（0表示关闭）
	PreAlloc            bool          `mapstructure:"PreAlloc" default:"false"`        // Whether to pre-allocate memory | 是否预分配内存
	NonBlocking         bool          `mapstructure:"NonBlocking" default:"true"`      // Whether to use non-blocking mode | 是否为非阻塞模式
}

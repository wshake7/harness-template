package config

import (
	"github.com/gofiber/fiber/v3"
	"time"
)

type FiberConfig struct {
	//启用“服务器：值”这一 HTTP 标头。
	ServerHeader string `mapstructure:"ServerHeader" default:""`
	//启用后，路由器会把 /foo 和 /foo/ 视为不同。否则，路由器会把 /foo 和 /foo/ 当作同一个数字处理。
	StrictRouting bool `mapstructure:"StrictRouting" default:"false"`

	//启用时，/Foo 和 /foo 是不同的路由。禁用时，/Foo 和 /foo 的处理方式相同启用时，/Foo 和 /foo 是不同的路由。禁用时，/Foo 和 /foo 的处理方式相同
	CaseSensitive bool `mapstructure:"CaseSensitive" default:"false"`

	//防止光纤自动为每个GET路由注册HEAD路由，从而提供自定义的HEAD处理器;手动HEAD路线仍然会覆盖生成的路线。
	DisableHeadAutoRegister bool `mapstructure:"DisableHeadAutoRegister" default:"false"`

	//启用后，上下文方法返回的所有值都是不可变的。默认情况下，这些信息在你从处理程序返回之前有效。
	Immutable bool `mapstructure:"Immutable" default:"false"`

	//当设置为“真”时，会将路由中的所有编码字符转换回原始形式，然后再为上下文设置路径，这样路由、从上下文 `ctx.Path()` 中返回当前 URL 以及带有解码字符的参数 `ctx.Params(%key%)` 就能够正常工作了。
	UnescapePath bool `mapstructure:"UnescapePath" default:"false"`

	//设置请求体的最大允许大小。零值或负值则回落到默认限制。如果大小超过配置限制，会发送413 - 请求实体过大响应。这个限制同样适用于通过适配器中间件从网络/http运行光纤。
	BodyLimit int `mapstructure:"BodyLimit" default:"4194304"`

	//设置从范围头解析的最大范围数。零值或负值则回落到默认限制。如果超过限制，请求将以416 - 请求范围不可满足和内容范围：字节 */<size> 拒绝
	MaxRanges int `mapstructure:"MaxRanges" default:"16"`

	//最大并发连接数。
	Concurrency int `mapstructure:"Concurrency" default:"262144"`

	//视图布局是所有模板渲染的全局布局，直到在“渲染”函数中进行重新设置。
	ViewsLayout string `mapstructure:"ViewsLayout" default:""`

	//“PassLocalsToViews”功能允许将在一个纤维（Context）中设置的局部变量传递给模板引擎。
	PassLocalsToViews bool `mapstructure:"PassLocalsToViews" default:"false"`

	//“PassLocalsToContext”用于控制“StoreInContext”是否还会将值传递到请求的“上下文”中（对于基于“Fiber”的上下文而言）。对于基于“Fiber”的上下文，“ValueFromContext”总是从“c.Locals()”中读取值。
	PassLocalsToContext bool `mapstructure:"PassLocalsToContext" default:"false"`

	//读超时 默认无限
	ReadTimeout time.Duration `mapstructure:"ReadTimeout" default:"0s"`

	//写超时 默认无限
	WriteTimeout time.Duration `mapstructure:"WriteTimeout" default:"0s"`

	//启用保持生命功能时，等待下一个请求的最大时间。如果IdleTimeout为零，则使用ReadTimeout的值。
	IdleTimeout time.Duration `mapstructure:"IdleTimeout" default:"0s"`

	//请求读取的每个连接缓冲区大小。这也限制了最大头部大小。如果你的客户端发送多KB的RequestURI和/或多KB头部（例如BIG cookies），增加缓冲区。
	ReadBufferSize int `mapstructure:"ReadBufferSize" default:"4096"`

	//响应写入时的每个连接缓冲区大小。
	WriteBufferSize int `mapstructure:"WriteBufferSize" default:"4096"`

	//ProxyHeader 会使 c.IP() 函数能够返回指定的头部键的值
	//默认情况下，c.IP() 会返回 TCP 连接中的远程 IP 地址
	//如果您的服务器位于负载均衡器之后，此属性可能会很有用：X-Forwarded-*（注意：头部很容易被伪造，检测到的 IP 地址不可靠。）
	ProxyHeader string `mapstructure:"ProxyHeader" default:""`

	//如果将该选项设置为“true”，则会拒绝所有非 GET 类型的请求。
	//此选项对于仅接受 GET 请求的服务器而言，具有防拒绝服务攻击（DoS）的作用。
	//如果将 GETOnly 设置为“true”，则请求大小将受到 ReadBufferSize 的限制。
	GETOnly bool `mapstructure:"GETOnly" default:"false"`

	//禁用保持活泼连接，服务器在第一次响应后关闭每个连接。
	DisableKeepalive bool `mapstructure:"DisableKeepalive" default:"false"`

	//如果为真，则在回复中省略日期头。
	DisableDefaultDate bool `mapstructure:"DisableDefaultDate" default:"false"`

	//当为真时，会省略响应中的默认 Content-Type 头部。
	DisableDefaultContentType bool `mapstructure:"DisableDefaultContentType" default:"false"`

	//默认情况下，所有头部名称均为规范化：conteNT-tYPE -> Content-Type
	DisableHeaderNormalizing bool `mapstructure:"DisableHeaderNormalizing" default:"false"`

	//设置日志中使用的应用名称和服务器头部
	AppName string `mapstructure:"AppName" default:""`

	//StreamRequestBody 启用请求体流，并在主体大于当前限制时更早调用处理程序。
	StreamRequestBody bool `mapstructure:"StreamRequestBody" default:"false"`

	//如果设置为true，就不会预解析多部分表单数据。该选项对于希望将多部分表单数据视为二进制数据，或选择何时解析数据的服务器非常有用。
	DisablePreParseMultipartForm bool `mapstructure:"DisablePreParseMultipartForm" default:"false"`

	//如果设置为true，可以大幅降低内存占用，但代价是CPU占用率更高。
	ReduceMemoryUsage bool `mapstructure:"ReduceMemoryUsage" default:"false"`

	//如果您发现自己处于某种代理（例如负载均衡器）之后，那么某些标头信息可能会通过特殊的 X-Forwarded-* 标头或 Forwarded 标头发送给您。
	//例如，Host HTTP 标头通常用于返回请求的主机。但在您处于代理之后时，实际的主机可能存储在 X-Forwarded-Host 标头中。如果您处于代理之后，应启用 TrustProxy 以防止标头欺骗。
	//如果您启用了 TrustProxy 但未提供 TrustProxyConfig，Fiber 将跳过所有可能被欺骗的标头。
	//如果请求 IP 在 TrustProxyConfig.Proxies 允许列表中，那么：
	//1. c.Scheme() 从 X-Forwarded-Proto、X-Forwarded-Protocol、X-Forwarded-Ssl 或 X-Url-Scheme 标头获取值
	//2. 从“ProxyHeader”标头获取 c.IP() 的值。3. c.Host() 和 c.Hostname() 从 X-Forwarded-Host 标头获取值 但如果请求 IP 不在 TrustProxyConfig.Proxies 允许列表中，则：
	//	1. c.Scheme() 不会从 X-Forwarded-Proto、X-Forwarded-Protocol、X-Forwarded-Ssl 或 X-Url-Scheme 标头获取值，当应用程序处理 TLS 连接时将返回 https，否则返回 http。
	//	2. c.IP() 不会从 ProxyHeader 标头获取值，将从 fasthttp 上下文返回 RemoteIP()3. c.Host() 和 c.Hostname() 不会从 X-Forwarded-Host 标头获取值，而是使用 fasthttp.Request.URI().Host() 来获取主机名。
	//若要自动信任所有环回、链路本地或私有 IP 地址，而无需手动将其添加到 TrustProxyConfig.Proxies 允许列表中，您可以将 TrustProxyConfig.Loopback、TrustProxyConfig.LinkLocal 或 TrustProxyConfig.Private 设置为 true。
	TrustProxy bool `mapstructure:"TrustProxy" default:"false"`

	//如果设置为“真”，那么 c.IP() 和 c.IPs() 函数在返回 IP 地址之前会对其进行验证。
	//此外，c.IP() 函数只会返回第一个有效的 IP 地址，而不会仅仅返回原始的报头信息。
	//警告：这会带来一定的性能开销。
	EnableIPValidation bool `mapstructure:"EnableIPValidation" default:"false"`

	//RequestMethods 提供了对 HTTP 方法的自定义功能。您可以根据需要添加或删除方法。可选。默认值：DefaultMethods
	RequestMethods []string `mapstructure:"RequestMethods" default:"[]"`

	//“EnableSplittingOnParsers”属性若设为“true”，则会将查询/主体/头部参数以逗号分隔开来。
	//例如，您可以使用它来解析来自查询参数的多个值，如下所示：
	///api？foo=bar,baz == foo[]=bar&foo[]=baz
	EnableSplittingOnParsers bool `mapstructure:"EnableSplittingOnParsers" default:"false"`

	Services []fiber.Service
}

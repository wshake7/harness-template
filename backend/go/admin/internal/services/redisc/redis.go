package redisc

import (
	"admin/internal/config"
	"context"
	"errors"
	"go-common/utils/types"
	"net"

	"github.com/bytedance/sonic"
	"github.com/redis/rueidis"
)

type RedisClient struct {
	rueidis.Client
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

var Client *RedisClient

func New(conf config.RedisConfig) *RedisClient {
	clientOption := rueidis.ClientOption{
		InitAddress:           conf.Addr,                  // 初始地址，支持多个地址（集群）
		Username:              conf.Username,              // 用户名
		Password:              conf.Password,              // 密码
		SelectDB:              conf.SelectDB,              // 选择的数据库
		ClientName:            conf.ClientName,            // 客户端名称
		CacheSizeEachConn:     conf.CacheSizeEachConn,     // 每个连接的缓存大小
		RingScaleEachConn:     conf.RingScaleEachConn,     // 每个连接的环形缓冲区缩放比例
		ReadBufferEachConn:    conf.ReadBufferEachConn,    // 每个连接的读缓冲区大小
		WriteBufferEachConn:   conf.WriteBufferEachConn,   // 每个连接的写缓冲区大小
		BlockingPoolSize:      conf.BlockingPoolSize,      // 阻塞池大小
		ConnWriteTimeout:      conf.ConnWriteTimeout,      // 连接写入超时
		ConnLifetime:          conf.ConnLifetime,          // 连接存活时间
		MaxFlushDelay:         conf.MaxFlushDelay,         // 最大刷新延迟
		DisableTCPNoDelay:     conf.DisableTCPNoDelay,     // 是否禁用 TCP NoDelay
		ShuffleInit:           conf.ShuffleInit,           // 是否随机打乱初始地址
		DisableRetry:          conf.DisableRetry,          // 是否禁用重试
		DisableCache:          conf.DisableCache,          // 是否禁用客户端缓存
		DisableAutoPipelining: conf.DisableAutoPipeline,   // 是否禁用自动流水线
		AlwaysPipelining:      conf.AlwaysPipelining,      // 是否始终使用流水线
		AlwaysRESP2:           conf.AlwaysRESP2,           // 是否始终使用 RESP2 协议
		TLSConfig:             nil,                        // TLS 配置
		DialFn:                nil,                        // 自定义拨号函数
		DialCtxFn:             nil,                        // 带上下文的自定义拨号函数
		NewCacheStoreFn:       nil,                        // 自定义缓存存储函数
		OnInvalidations:       nil,                        // 缓存失效时的回调
		SendToReplicas:        nil,                        // 是否发送到从节点
		AuthCredentialsFn:     nil,                        // 身份验证凭据函数
		RetryDelay:            nil,                        // 重试延迟函数
		ReplicaSelector:       nil,                        // 从节点选择器
		ReadNodeSelector:      nil,                        // 读节点选择器
		Sentinel:              rueidis.SentinelOption{},   // 哨兵配置
		Dialer:                net.Dialer{},               // 拨号器配置
		ClientSetInfo:         []string{},                 // CLIENT SETINFO 参数
		ClientTrackingOptions: []string{},                 // CLIENT TRACKING 参数
		Standalone:            rueidis.StandaloneOption{}, // 单机模式配置
		BlockingPoolCleanup:   0,                          // 阻塞池清理间隔
		BlockingPoolMinSize:   0,                          // 阻塞池最小大小
		BlockingPipeline:      0,                          // 阻塞流水线配置
		PipelineMultiplex:     0,                          // 流水线多路复用配置
		ClusterOption:         rueidis.ClusterOption{},    // 集群选项
		ClientNoTouch:         false,                      // 是否不触碰（不更新 LRU）
		ForceSingleClient:     false,                      // 是否强制单客户端
		ReplicaOnly:           false,                      // 是否仅连接从节点
		ClientNoEvict:         false,                      // 是否不驱逐
		EnableReplicaAZInfo:   false,                      // 是否启用从节点可用区信息
		AZFromInfo:            false,                      // 是否从信息中获取可用区
	}
	client, err := rueidis.NewClient(clientOption)
	if err != nil {
		panic(err)
	}
	if err = client.Do(context.Background(), client.B().Ping().Build()).Error(); err != nil {
		panic("redis连接失败: " + err.Error())
	}
	Client = &RedisClient{
		Client:    client,
		marshal:   sonic.Marshal,
		unmarshal: sonic.Unmarshal,
	}
	return Client
}

func (r *RedisClient) GetJson(ctx context.Context, key string, obj any) error {
	if !types.IsPointer(obj) {
		return errors.New("obj is not a pointer")
	}
	result, err := r.Client.Do(ctx, r.Client.B().Get().Key(key).Build()).AsBytes()
	if err != nil {
		return err
	}
	return r.unmarshal(result, obj)
}

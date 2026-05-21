package rueidis

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"net"
	"time"

	"github.com/click33/sa-token-go/core/adapter"
	"github.com/redis/rueidis"
)

// Storage Redis存储实现
type Storage struct {
	client    rueidis.Client
	ctx       context.Context
	opTimeout time.Duration
}

// Config Redis配置
type Config struct {
	Host     string
	Port     int
	Password string
	Database int
	PoolSize int
	// Optional timeouts for redis client
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	// OperationTimeout applies to each single storage operation context
	OperationTimeout time.Duration
}

// NewStorage 通过Redis URL创建存储
func NewStorage(url string) (adapter.Storage, error) {
	opts, err := rueidis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	client, err := rueidis.NewClient(opts)

	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}
	// 测试连接
	ctx := context.Background()

	if err := client.Do(context.Background(), client.B().Ping().Build()).Error(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)

	}

	return &Storage{
		client:    client,
		ctx:       ctx,
		opTimeout: 3 * time.Second,
	}, nil
}

// NewStorageFromConfig 通过配置创建存储
func NewStorageFromConfig(cfg *Config) (adapter.Storage, error) {
	config := rueidis.ClientOption{
		InitAddress:           []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, // 初始地址，支持多个地址（集群）
		Password:              cfg.Password,                                       // 密码
		SelectDB:              cfg.Database,                                       // 选择的数据库
		ClientName:            "",                                                 // 客户端名称
		CacheSizeEachConn:     0,                                                  // 每个连接的缓存大小
		RingScaleEachConn:     0,                                                  // 每个连接的环形缓冲区缩放比例
		ReadBufferEachConn:    0,                                                  // 每个连接的读缓冲区大小
		WriteBufferEachConn:   0,                                                  // 每个连接的写缓冲区大小
		BlockingPoolSize:      cfg.PoolSize,                                       // 阻塞池大小
		ConnWriteTimeout:      cfg.WriteTimeout,                                   // 连接写入超时
		ConnLifetime:          0,                                                  // 连接存活时间
		MaxFlushDelay:         0,                                                  // 最大刷新延迟
		DisableTCPNoDelay:     false,                                              // 是否禁用 TCP NoDelay
		ShuffleInit:           false,                                              // 是否随机打乱初始地址
		DisableRetry:          false,                                              // 是否禁用重试
		DisableCache:          false,                                              // 是否禁用客户端缓存
		DisableAutoPipelining: false,                                              // 是否禁用自动流水线
		AlwaysPipelining:      false,                                              // 是否始终使用流水线
		AlwaysRESP2:           false,                                              // 是否始终使用 RESP2 协议
		TLSConfig:             nil,                                                // TLS 配置
		DialFn:                nil,                                                // 自定义拨号函数
		DialCtxFn:             nil,                                                // 带上下文的自定义拨号函数
		NewCacheStoreFn:       nil,                                                // 自定义缓存存储函数
		OnInvalidations:       nil,                                                // 缓存失效时的回调
		SendToReplicas:        nil,                                                // 是否发送到从节点
		AuthCredentialsFn:     nil,                                                // 身份验证凭据函数
		RetryDelay:            nil,                                                // 重试延迟函数
		ReplicaSelector:       nil,                                                // 从节点选择器
		ReadNodeSelector:      nil,                                                // 读节点选择器
		Sentinel:              rueidis.SentinelOption{},                           // 哨兵配置
		Dialer:                net.Dialer{},                                       // 拨号器配置
		ClientSetInfo:         []string{},                                         // CLIENT SETINFO 参数
		ClientTrackingOptions: []string{},                                         // CLIENT TRACKING 参数
		Standalone:            rueidis.StandaloneOption{},                         // 单机模式配置
		BlockingPoolCleanup:   0,                                                  // 阻塞池清理间隔
		BlockingPoolMinSize:   0,                                                  // 阻塞池最小大小
		BlockingPipeline:      0,                                                  // 阻塞流水线配置
		PipelineMultiplex:     0,                                                  // 流水线多路复用配置
		ClusterOption:         rueidis.ClusterOption{},                            // 集群选项
		ClientNoTouch:         false,                                              // 是否不触碰（不更新 LRU）
		ForceSingleClient:     false,                                              // 是否强制单客户端
		ReplicaOnly:           false,                                              // 是否仅连接从节点
		ClientNoEvict:         false,                                              // 是否不驱逐
		EnableReplicaAZInfo:   false,                                              // 是否启用从节点可用区信息
		AZFromInfo:            false,                                              // 是否从信息中获取可用区
	}
	client, err := rueidis.NewClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	opTimeout := cfg.OperationTimeout
	if opTimeout <= 0 {
		opTimeout = 3 * time.Second
	}

	return &Storage{
		client:    client,
		ctx:       ctx,
		opTimeout: opTimeout,
	}, nil
}

// NewStorageFromClient 从已有的Redis客户端创建存储
func NewStorageFromClient(client rueidis.Client) adapter.Storage {
	return &Storage{
		client:    client,
		ctx:       context.Background(),
		opTimeout: 3 * time.Second,
	}
}

// getKey 获取完整的键名（Storage 层不处理前缀，前缀由 Manager 层统一管理）
func (s *Storage) getKey(key string) string {
	return key
}

// Set 设置键值对
func (s *Storage) Set(key string, value any, expiration time.Duration) (err error) {
	ctx, cancel := s.withTimeout()
	defer cancel()
	var data string
	switch value.(type) {
	case string:
		data = value.(string)
	case []byte:
		data = string(value.([]byte))
	default:
		data, err = sonic.MarshalString(value)
		if err != nil {
			return err
		}
	}

	cmd := s.client.B().Set().Key(s.getKey(key)).Value(data)
	if expiration > 0 {
		return s.client.Do(ctx, cmd.Ex(expiration).Build()).Error()
	}
	return s.client.Do(ctx, cmd.Build()).Error()
}

func (s *Storage) SetKeepTTL(key string, value any) error {
	ctx, cancel := s.withTimeout()
	defer cancel()

	count, err := s.client.Do(ctx, s.client.B().Exists().Key(s.getKey(key)).Build()).AsInt64()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("key not found: %s", key)
	}
	var data string
	switch value.(type) {
	case string:
		data = value.(string)
	case []byte:
		data = string(value.([]byte))
	default:
		data, err = sonic.MarshalString(value)
		if err != nil {
			return err
		}
	}
	cmd := s.client.B().Set().Key(s.getKey(key)).Value(data).Keepttl().Build()
	return s.client.Do(ctx, cmd).Error()
}

// Get 获取值
func (s *Storage) Get(key string) (any, error) {
	ctx, cancel := s.withTimeout()
	defer cancel()

	result, err := s.client.Do(ctx, s.client.B().Get().Key(s.getKey(key)).Build()).AsBytes()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, err
	}
	return result, nil
}

// Delete 删除键
func (s *Storage) Delete(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	ctx, cancel := s.withTimeout()
	defer cancel()

	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = s.getKey(key)
	}

	return s.client.Do(ctx, s.client.B().Del().Key(fullKeys...).Build()).Error()
}

// Exists 检查键是否存在
func (s *Storage) Exists(key string) bool {
	ctx, cancel := s.withTimeout()
	defer cancel()

	result, err := s.client.Do(ctx, s.client.B().Exists().Key(s.getKey(key)).Build()).AsInt64()
	if err != nil {
		return false
	}
	return result > 0
}

// Keys 获取匹配模式的所有键
func (s *Storage) Keys(pattern string) ([]string, error) {
	ctx, cancel := s.withTimeout()
	defer cancel()

	var (
		cursor uint64
		result []string
	)
	for {
		res, err := s.client.Do(ctx, s.client.B().Scan().Cursor(cursor).Match(pattern).Count(1000).Build()).AsScanEntry()
		if err != nil {
			return nil, err
		}
		result = append(result, res.Elements...)
		cursor = res.Cursor
		if cursor == 0 {
			break
		}
	}
	return result, nil
}

// Expire 设置键的过期时间
func (s *Storage) Expire(key string, expiration time.Duration) error {
	ctx, cancel := s.withTimeout()
	defer cancel()
	seconds := int64(expiration.Seconds())
	return s.client.Do(ctx, s.client.B().Expire().Key(s.getKey(key)).Seconds(seconds).Build()).Error()
}

// TTL 获取键的剩余生存时间
func (s *Storage) TTL(key string) (time.Duration, error) {
	ctx, cancel := s.withTimeout()
	defer cancel()

	// /*go-redis v9*/
	// return s.client.TTL(ctx, s.getKey(key)).Result()

	// rueidis
	result, err := s.client.Do(ctx, s.client.B().Ttl().Key(s.getKey(key)).Build()).AsInt64()
	if err != nil {
		return 0, err
	}
	return time.Duration(result) * time.Second, nil
}

// Clear 清空所有数据（警告：会清空整个 Redis，谨慎使用！应由 Manager 层控制）
func (s *Storage) Clear() error {
	ctx, cancel := s.withTimeout()
	defer cancel()
	var cursor uint64
	for {
		res, err := s.client.Do(ctx, s.client.B().Scan().Cursor(cursor).Match("*").Count(1000).Build()).AsScanEntry()
		if err != nil {
			return err
		}
		if len(res.Elements) > 0 {
			err = s.client.Do(ctx, s.client.B().Unlink().Key(res.Elements...).Build()).Error()
			if err != nil {
				return err
			}
		}
		cursor = res.Cursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// Ping 检查连接
func (s *Storage) Ping() error {
	ctx, cancel := s.withTimeout()
	defer cancel()

	return s.client.Do(ctx, s.client.B().Ping().Build()).Error()
}

// Close 关闭连接
func (s *Storage) Close() error {
	s.client.Close()
	return nil
}

// GetClient 获取Redis客户端（用于高级操作）
func (s *Storage) GetClient() rueidis.Client {
	return s.client
}

// withTimeout returns a context with the configured per-operation timeout.
func (s *Storage) withTimeout() (context.Context, context.CancelFunc) {
	if s.opTimeout > 0 {
		return context.WithTimeout(s.ctx, s.opTimeout)
	}
	return context.WithCancel(s.ctx)
}

// Builder Redis存储构建器
type Builder struct {
	host     string
	port     int
	password string
	database int
	poolSize int
}

// NewBuilder 创建构建器
func NewBuilder() *Builder {
	return &Builder{
		host:     "localhost",
		port:     6379,
		password: "",
		database: 0,
		poolSize: 10,
	}
}

// Host 设置主机
func (b *Builder) Host(host string) *Builder {
	b.host = host
	return b
}

// Port 设置端口
func (b *Builder) Port(port int) *Builder {
	b.port = port
	return b
}

// Password 设置密码
func (b *Builder) Password(password string) *Builder {
	b.password = password
	return b
}

// Database 设置数据库
func (b *Builder) Database(database int) *Builder {
	b.database = database
	return b
}

// PoolSize 设置连接池大小
func (b *Builder) PoolSize(poolSize int) *Builder {
	b.poolSize = poolSize
	return b
}

// Build 构建存储
func (b *Builder) Build() (adapter.Storage, error) {
	return NewStorageFromConfig(&Config{
		Host:     b.host,
		Port:     b.port,
		Password: b.password,
		Database: b.database,
		PoolSize: b.poolSize,
	})
}

package config

import (
	"time"
)

type Config struct {
	AppName         string `mapstructure:"AppName"`
	Host            string `mapstructure:"Host" default:"0.0.0.0"`
	Port            int    `mapstructure:"Port" default:"3000"`
	RestPrefix      string `mapstructure:"RestPrefix" default:"/"`
	IsSwagger       bool   `mapstructure:"IsSwagger" default:"false"`
	SwaggerPrefix   string `mapstructure:"SwaggerPrefix" default:"/swagger"`
	DefaultLanguage string `mapstructure:"DefaultLanguage" default:"cn"`
	Auth            AuthConfig
	Fiber           FiberConfig
	AI              AIConfig
	Orm             OrmConfig
	Redis           RedisConfig
	Temporal        TemporalConfig
	Milvus          MilvusConfig
}

var Conf = new(Config)

type AIConfig struct {
	Embedding AIEmbeddingConfig `mapstructure:"Embedding"`
	Knowledge AIKnowledgeConfig `mapstructure:"Knowledge"`
}

type AIEmbeddingConfig struct {
	Provider   string `mapstructure:"Provider" default:"dashscope"`      // Embedding provider. Current default is DashScope.
	APIKey     string `mapstructure:"APIKey"`                            // Third-party embedding credential.
	Model      string `mapstructure:"Model" default:"text-embedding-v3"` // Embedding model name.
	Dimensions uint   `mapstructure:"Dimensions" default:"1024"`         // Dense vector dimension expected by the embedding model.
}

type AIKnowledgeConfig struct {
	Collection          string `mapstructure:"Collection" default:"biz"`                                            // Target Milvus collection.
	CollectionDesc      string `mapstructure:"CollectionDesc" default:"Knowledge documents for admin AI workflows"` // Collection description used on first create.
	IDMaxLength         uint   `mapstructure:"IDMaxLength" default:"255"`                                           // Max length for the Milvus varchar primary key field.
	ContentMaxLength    uint   `mapstructure:"ContentMaxLength" default:"8192"`                                     // Max length for the Milvus varchar content field.
	IndexType           string `mapstructure:"IndexType" default:"auto"`                                            // Dense vector index type: auto, hnsw, ivf_flat.
	IndexMetricType     string `mapstructure:"IndexMetricType" default:"COSINE"`                                    // Dense vector similarity metric.
	LoadTimeoutSeconds  int    `mapstructure:"LoadTimeoutSeconds" default:"60"`
	FlushTimeoutSeconds int    `mapstructure:"FlushTimeoutSeconds" default:"30"`
}

type OrmConfig struct {
	DriverName      string        `mapstructure:"DriverName"`
	DataSource      string        `mapstructure:"DataSource"`
	ConnMaxIdleTime time.Duration `mapstructure:"ConnMaxIdleTime" default:"60s"`
	ConnMaxLifetime time.Duration `mapstructure:"ConnMaxLifetime" default:"120s"`
	MaxIdleConn     int           `mapstructure:"MaxIdleConn" default:"10"`
	MaxOpenConn     int           `mapstructure:"MaxOpenConn" default:"20"`
	IsGenCode       bool          `mapstructure:"IsGenCode" default:"false"`
	IsAutoMigrate   bool          `mapstructure:"IsAutoMigrate" default:"false"`
	IsLog           bool          `mapstructure:"IsLog" default:"false"`
}

type RedisConfig struct {
	Addr                []string      `mapstructure:"Addr"`                 // Redis 地址，支持集群模式
	Username            string        `mapstructure:"Username"`             // Redis 用户名
	Password            string        `mapstructure:"Password"`             // Redis 密码
	SelectDB            int           `mapstructure:"SelectDB" default:"0"` // 选择的数据库索引
	ClientName          string        `mapstructure:"ClientName"`           // 客户端名称
	CacheSizeEachConn   int           `mapstructure:"CacheSizeEachConn"`    // 每个连接的客户端缓存大小 (字节)
	RingScaleEachConn   int           `mapstructure:"RingScaleEachConn"`    // 每个连接的环形缓冲区大小
	ReadBufferEachConn  int           `mapstructure:"ReadBufferEachConn"`   // 每个连接的读缓冲区大小
	WriteBufferEachConn int           `mapstructure:"WriteBufferEachConn"`  // 每个连接的写缓冲区大小
	BlockingPoolSize    int           `mapstructure:"BlockingPoolSize"`     // 阻塞操作的连接池大小
	ConnWriteTimeout    time.Duration `mapstructure:"ConnWriteTimeout"`     // 连接写入超时时间
	ConnDialTimeout     time.Duration `mapstructure:"ConnDialTimeout"`
	ConnReadTimeout     time.Duration `mapstructure:"ConnReadTimeout"`
	ConnLifetime        time.Duration `mapstructure:"ConnLifetime"`        // 连接最大存活时间
	MaxFlushDelay       time.Duration `mapstructure:"MaxFlushDelay"`       // 最大刷新延迟
	DisableTCPNoDelay   bool          `mapstructure:"DisableTCPNoDelay"`   // 是否禁用 TCP_NODELAY
	ShuffleInit         bool          `mapstructure:"ShuffleInit"`         // 是否在初始化时打乱地址顺序
	DisableRetry        bool          `mapstructure:"DisableRetry"`        // 是否禁用重试
	DisableCache        bool          `mapstructure:"DisableCache"`        // 是否禁用客户端缓存
	DisableAutoPipeline bool          `mapstructure:"DisableAutoPipeline"` // 是否禁用自动管道
	AlwaysPipelining    bool          `mapstructure:"AlwaysPipelining"`    // 是否始终使用管道
	AlwaysRESP2         bool          `mapstructure:"AlwaysRESP2"`         // 是否始终使用 RESP2 协议
}

type TemporalConfig struct {
	Enabled       bool   `mapstructure:"Enabled" default:"false"`
	HostPort      string `mapstructure:"HostPort" default:"127.0.0.1:7233"`
	Namespace     string `mapstructure:"Namespace" default:"default"`
	Identity      string `mapstructure:"Identity"`
	TaskQueue     string `mapstructure:"TaskQueue" default:"admin"`
	WorkerEnabled bool   `mapstructure:"WorkerEnabled" default:"false"`
}

type MilvusConfig struct {
	Enabled bool   `mapstructure:"Enabled" default:"true"`
	Address string `mapstructure:"Address" default:"127.0.0.1:19530"`
	APIKey  string `mapstructure:"APIKey"`
	DBName  string `mapstructure:"DBName"`
}

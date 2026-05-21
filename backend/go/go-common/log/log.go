package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"log/slog"
	"os"
	"path/filepath"
)

// Config 日志配置（重命名为更规范的 Config）
type Config struct {
	// 基础配置
	ProjectName string
	Level       zapcore.Level
	IsJson      bool // 重命名为更规范的驼峰命名

	// 编码器配置
	EncoderConfig zapcore.EncoderConfig

	// 输出配置
	Writers      []zapcore.WriteSyncer                // 重命名为更清晰的名称
	LevelWriters map[zapcore.Level]*lumberjack.Logger // 重命名

	// 文件配置
	FilePath string
	FileName string

	// 动态字段函数
	FieldFunc func() []zapcore.Field // 重命名
}

// LoggerOption 配置选项函数
type LoggerOption func(*Config)

// WithProjectName 设置项目名称
func WithProjectName(name string) LoggerOption {
	return func(c *Config) {
		c.ProjectName = name
	}
}

// WithLevel 设置日志级别
func WithLevel(level zapcore.Level) LoggerOption {
	return func(c *Config) {
		c.Level = level
	}
}

// WithJSON 设置 JSON 编码
func WithJSON(isJSON bool) LoggerOption {
	return func(c *Config) {
		c.IsJson = isJSON
	}
}

// WithWriters 设置输出写入器
func WithWriters(writers ...zapcore.WriteSyncer) LoggerOption {
	return func(c *Config) {
		c.Writers = writers
	}
}

// WithFile 设置文件输出
func WithFile(filePath, fileName string) LoggerOption {
	return func(c *Config) {
		c.FilePath = filePath
		c.FileName = fileName
		lj := DefaultLumberjack(filePath, fileName)
		c.Writers = append(c.Writers, zapcore.AddSync(lj))
	}
}

// WithFieldFunc 设置动态字段函数
func WithFieldFunc(fn func() []zapcore.Field) LoggerOption {
	return func(c *Config) {
		c.FieldFunc = fn
	}
}

// defaultZapOptions 默认 zap 选项
var defaultZapOptions = []zap.Option{
	zap.AddCaller(),
	zap.AddStacktrace(zapcore.ErrorLevel),
}

// Init 初始化全局日志器（zap 和 slog）
func Init(configs ...*Config) {
	InitZapGlobal(configs...)
	InitSlogGlobal(configs...)
}

// InitZapGlobal 初始化全局 zap 日志器
func InitZapGlobal(configs ...*Config) {
	if len(configs) == 0 {
		configs = []*Config{DevConfig()}
	}
	logger := NewZap(configs...)
	zap.ReplaceGlobals(logger)
}

// InitSlogGlobal 初始化全局 slog 日志器
func InitSlogGlobal(configs ...*Config) {
	if len(configs) == 0 {
		configs = []*Config{DevConfig()}
	}
	logger := NewSlog(configs...)
	slog.SetDefault(logger)
}

// NewZap 创建 zap 日志器
func NewZap(configs ...*Config) *zap.Logger {
	return NewZapWithOptions(defaultZapOptions, configs...)
}

// NewZapWithOptions 使用自定义选项创建 zap 日志器
func NewZapWithOptions(options []zap.Option, configs ...*Config) *zap.Logger {
	if len(configs) == 0 {
		configs = []*Config{DevConfig()}
	}

	cores := make([]zapcore.Core, 0, len(configs)*2)

	for _, conf := range configs {
		// 验证配置
		if err := validateConfig(conf); err != nil {
			panic(fmt.Errorf("invalid log config: %w", err))
		}

		// 创建编码器
		encoder := newEncoder(conf)

		// 为每个 writer 创建 core
		for _, ws := range conf.Writers {
			core := zapcore.NewCore(encoder, ws, conf.Level)
			cores = append(cores, core)
		}
	}

	if len(cores) == 0 {
		panic("no log cores configured")
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, options...)
}

// NewSlog 创建 slog 日志器
func NewSlog(configs ...*Config) *slog.Logger {
	if len(configs) == 0 {
		configs = []*Config{DevConfig()}
	}

	cores := make([]zapcore.Core, 0, len(configs)*2)

	for _, conf := range configs {
		// 验证配置
		if err := validateConfig(conf); err != nil {
			panic(fmt.Errorf("invalid log config: %w", err))
		}

		// 创建编码器
		encoder := newEncoder(conf)

		// 为每个 writer 创建 core
		for _, ws := range conf.Writers {
			core := zapcore.NewCore(encoder, ws, conf.Level)
			cores = append(cores, core)
		}
	}

	if len(cores) == 0 {
		panic("no log cores configured")
	}

	core := zapcore.NewTee(cores...)
	return slog.New(zapslog.NewHandler(core, zapslog.WithCaller(true)))
}

// validateConfig 验证配置
func validateConfig(conf *Config) error {
	if conf == nil {
		return fmt.Errorf("config is nil")
	}
	if len(conf.Writers) == 0 && len(conf.LevelWriters) == 0 {
		return fmt.Errorf("no writers configured")
	}
	return nil
}

// DevConfig 开发环境配置
func DevConfig() *Config {
	return &Config{
		ProjectName: "",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    CapitalColorLevelEncoder, // 使用彩色输出
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		IsJson:       false,
		Writers:      []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)},
		Level:        zapcore.DebugLevel,
		FieldFunc:    nil,
		LevelWriters: nil,
	}
}

// ProdConfig 生产环境配置
func ProdConfig() *Config {
	return &Config{
		ProjectName: "",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		IsJson:       true,
		Writers:      []zapcore.WriteSyncer{zapcore.AddSync(DefaultLumberjack("", ""))},
		Level:        zapcore.InfoLevel,
		FieldFunc:    nil,
		LevelWriters: nil,
	}
}

// NewDevLogger 创建开发环境日志器（便捷方法）
func NewDevLogger(opts ...LoggerOption) *zap.Logger {
	config := DevConfig()
	for _, opt := range opts {
		opt(config)
	}
	return NewZap(config)
}

// NewProdLogger 创建生产环境日志器（便捷方法）
func NewProdLogger(opts ...LoggerOption) *zap.Logger {
	config := ProdConfig()
	for _, opt := range opts {
		opt(config)
	}
	return NewZap(config)
}

// DefaultLumberjack 创建默认的文件滚动日志
func DefaultLumberjack(filePath string, fileName string) *lumberjack.Logger {
	if fileName == "" {
		fileName = "app.log"
	}

	logDir := filePath
	if filePath == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("failed to get working directory: %w", err))
		}

		logDir = filepath.Join(wd, "logs")
	}

	// 确保目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Errorf("failed to create log directory %s: %w", logDir, err))
	}

	return &lumberjack.Logger{
		Filename:   filepath.Join(logDir, fileName),
		MaxSize:    100,  // MB，增大到 100MB
		MaxAge:     30,   // 天，减少到 30 天
		MaxBackups: 10,   // 个，减少到 10 个
		LocalTime:  true, // 使用本地时间
		Compress:   true, // 压缩旧文件
	}
}

// NewLumberjack 创建自定义配置的文件滚动日志
func NewLumberjack(filePath, fileName string, maxSize, maxAge, maxBackups int, compress bool) *lumberjack.Logger {
	if fileName == "" {
		fileName = "app.log"
	}

	logDir := filePath
	if filePath == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("failed to get working directory: %w", err))
		}
		logDir = filepath.Join(wd, "logs")
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Errorf("failed to create log directory %s: %w", logDir, err))
	}

	return &lumberjack.Logger{
		Filename:   filepath.Join(logDir, fileName),
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  true,
		Compress:   compress,
	}
}

func LowercaseLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	default:
		enc.AppendString(l.String())
	}
}

func CapitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	zapcore.CapitalColorLevelEncoder(l, enc)
}

func ParseLevel(text string) (zapcore.Level, error) {
	switch text {
	default:
		atomicLevel, err := zap.ParseAtomicLevel(text)
		if err != nil {
			return zapcore.InfoLevel, err
		}
		return atomicLevel.Level(), err
	}
}

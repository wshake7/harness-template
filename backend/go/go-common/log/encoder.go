package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// customEncoder 自定义编码器
type customEncoder struct {
	zapcore.Encoder
	projectName  string
	levelWriters map[zapcore.Level]*lumberjack.Logger
	fieldFunc    func() []zapcore.Field
}

// EncodeEntry 编码日志条目
func (e *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 添加项目名称
	if e.projectName != "" {
		fields = append(fields, zap.String("project", e.projectName))
	}

	// 添加动态字段
	if e.fieldFunc != nil {
		if dynamicFields := e.fieldFunc(); len(dynamicFields) > 0 {
			fields = append(fields, dynamicFields...)
		}
	}

	// 编码条目
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, err
	}

	// 写入级别特定的文件
	if e.levelWriters != nil {
		if writer := e.levelWriters[entry.Level]; writer != nil {
			// 忽略写入错误，避免日志记录失败
			_, _ = writer.Write(buf.Bytes())
		}
	}

	return buf, nil
}

// Clone 克隆编码器（实现 zapcore.Encoder 接口）
func (e *customEncoder) Clone() zapcore.Encoder {
	return &customEncoder{
		Encoder:      e.Encoder.Clone(),
		projectName:  e.projectName,
		levelWriters: e.levelWriters,
		fieldFunc:    e.fieldFunc,
	}
}

// newEncoder 创建编码器
func newEncoder(conf *Config) zapcore.Encoder {
	var enc zapcore.Encoder
	if conf.IsJson {
		enc = zapcore.NewJSONEncoder(conf.EncoderConfig)
	} else {
		enc = zapcore.NewConsoleEncoder(conf.EncoderConfig)
	}

	return &customEncoder{
		Encoder:      enc,
		projectName:  conf.ProjectName,
		levelWriters: conf.LevelWriters,
		fieldFunc:    conf.FieldFunc,
	}
}

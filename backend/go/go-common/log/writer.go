package log

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

// customWriter 自定义写入器（保留原有功能但改进实现）
type customWriter struct {
	file  *os.File
	mutex sync.Mutex
	meta  map[string]any
	fn    func(map[string]any, []byte) (int, error)
}

// NewCustomWriter 创建自定义写入器
func NewCustomWriter(file *os.File, meta map[string]any, fn func(map[string]any, []byte) (int, error)) *customWriter {
	if fn == nil {
		panic("write function is required")
	}
	return &customWriter{
		file: file,
		meta: meta,
		fn:   fn,
	}
}

// Write 实现 io.Writer 接口
func (w *customWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.fn(w.meta, p)
}

// Sync 实现 zapcore.WriteSyncer 接口
func (w *customWriter) Sync() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭文件
func (w *customWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// LeveledWriter 按日志级别分别写入不同文件
type LeveledWriter struct {
	writers map[zapcore.Level]zapcore.WriteSyncer
	mutex   sync.RWMutex
}

// NewLeveledWriter 创建按级别写入的写入器
func NewLeveledWriter() *LeveledWriter {
	return &LeveledWriter{
		writers: make(map[zapcore.Level]zapcore.WriteSyncer),
	}
}

// AddWriter 添加级别写入器
func (lw *LeveledWriter) AddWriter(level zapcore.Level, writer zapcore.WriteSyncer) {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()
	lw.writers[level] = writer
}

// Write 实现 io.Writer 接口（这里不实际使用，只是占位）
func (lw *LeveledWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Sync 同步所有写入器
func (lw *LeveledWriter) Sync() error {
	lw.mutex.RLock()
	defer lw.mutex.RUnlock()

	var errs []error
	for _, writer := range lw.writers {
		if err := writer.Sync(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to sync writers: %v", errs)
	}
	return nil
}

// BufferedWriter 带缓冲的写入器
type BufferedWriter struct {
	writer zapcore.WriteSyncer
	buffer *buffer.Buffer
	mutex  sync.Mutex
	size   int
}

// NewBufferedWriter 创建带缓冲的写入器
func NewBufferedWriter(writer zapcore.WriteSyncer, bufferSize int) *BufferedWriter {
	if bufferSize <= 0 {
		bufferSize = 4096
	}
	return &BufferedWriter{
		writer: writer,
		buffer: buffer.NewPool().Get(),
		size:   bufferSize,
	}
}

// Write 实现 io.Writer 接口
func (bw *BufferedWriter) Write(p []byte) (n int, err error) {
	bw.mutex.Lock()
	defer bw.mutex.Unlock()

	n, err = bw.buffer.Write(p)
	if err != nil {
		return n, err
	}

	// 如果缓冲区满了，刷新到底层写入器
	if bw.buffer.Len() >= bw.size {
		return n, bw.flush()
	}

	return n, nil
}

// Sync 同步缓冲区
func (bw *BufferedWriter) Sync() error {
	bw.mutex.Lock()
	defer bw.mutex.Unlock()
	return bw.flush()
}

// flush 刷新缓冲区（需要持有锁）
func (bw *BufferedWriter) flush() error {
	if bw.buffer.Len() == 0 {
		return nil
	}

	_, err := bw.writer.Write(bw.buffer.Bytes())
	bw.buffer.Reset()

	if err != nil {
		return err
	}

	return bw.writer.Sync()
}

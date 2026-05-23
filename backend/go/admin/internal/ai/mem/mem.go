package mem

import (
	"admin/internal/config"
	"context"

	"github.com/cloudwego/eino/schema"
)

// Memory is the top-level factory that creates sessions keyed by ID.
type Memory interface {
	GetSession(id string) Session
}

// Session holds conversation messages for a single dialogue.
type Session interface {
	SetMessages(ctx context.Context, msg *schema.Message) error
	GetMessages(ctx context.Context) ([]*schema.Message, error)
}

// InMemConfig holds parameters for NewInMemMemory.
type InMemConfig struct {
	MaxWindowSize int
}

// RedisConfig holds parameters for NewRedisMemory.
type RedisConfig struct {
	MaxWindowSize int
	TTLSeconds    int
}

// DBConfig holds parameters for NewDBMemory.
type DBConfig struct {
	MaxWindowSize int
}

// New creates a Memory from config. Use the typed constructors directly when
// you need different backends in different places without coupling to config.
func New(conf config.AIMemoryConfig) Memory {
	switch conf.Type {
	case "redis":
		return NewRedisMemory(RedisConfig{
			MaxWindowSize: conf.MaxWindowSize,
			TTLSeconds:    conf.TTLSeconds,
		})
	case "db":
		return NewDBMemory(DBConfig{MaxWindowSize: conf.MaxWindowSize})
	default:
		return NewInMemMemory(InMemConfig{MaxWindowSize: conf.MaxWindowSize})
	}
}

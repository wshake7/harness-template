package mem

import (
	"admin/internal/services/redisc"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/rueidis"
)

type redisMemory struct {
	maxWindowSize int
	ttl           time.Duration
}

func NewRedisMemory(cfg RedisConfig) Memory {
	ttl := time.Duration(cfg.TTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &redisMemory{
		maxWindowSize: cfg.MaxWindowSize,
		ttl:           ttl,
	}
}

func (m *redisMemory) GetSession(id string) Session {
	return &redisSession{
		key:           fmt.Sprintf("ai:memory:%s", id),
		maxWindowSize: m.maxWindowSize,
		ttl:           m.ttl,
	}
}

type redisSession struct {
	key           string
	maxWindowSize int
	ttl           time.Duration
}

func (s *redisSession) SetMessages(ctx context.Context, msg *schema.Message) error {
	msgs, err := s.get(ctx)
	if err != nil {
		return err
	}
	msgs = append(msgs, msg)
	if len(msgs) > s.maxWindowSize {
		excess := len(msgs) - s.maxWindowSize
		if excess%2 != 0 {
			excess++
		}
		msgs = msgs[excess:]
	}
	return s.set(ctx, msgs)
}

func (s *redisSession) GetMessages(ctx context.Context) ([]*schema.Message, error) {
	return s.get(ctx)
}

func (s *redisSession) get(ctx context.Context) ([]*schema.Message, error) {
	b, err := redisc.Client.Do(ctx, redisc.Client.B().Get().Key(s.key).Build()).AsBytes()
	if err != nil {
		if errors.Is(err, rueidis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get messages: %w", err)
	}
	var msgs []*schema.Message
	if err := sonic.Unmarshal(b, &msgs); err != nil {
		return nil, fmt.Errorf("unmarshal messages: %w", err)
	}
	return msgs, nil
}

func (s *redisSession) set(ctx context.Context, msgs []*schema.Message) error {
	b, err := sonic.Marshal(msgs)
	if err != nil {
		return fmt.Errorf("marshal messages: %w", err)
	}
	return redisc.Client.Do(ctx, redisc.Client.B().Set().Key(s.key).Value(rueidis.BinaryString(b)).ExSeconds(int64(s.ttl.Seconds())).Build()).Error()
}

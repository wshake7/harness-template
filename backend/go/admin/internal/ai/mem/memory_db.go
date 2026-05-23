package mem

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/schema"
)

type dbMemory struct {
	maxWindowSize int
}

func NewDBMemory(cfg DBConfig) Memory {
	return &dbMemory{maxWindowSize: cfg.MaxWindowSize}
}

func (m *dbMemory) GetSession(id string) Session {
	return &dbSession{sessionID: id, maxWindowSize: m.maxWindowSize}
}

type dbSession struct {
	sessionID     string
	maxWindowSize int
}

func (s *dbSession) ensureSession(ctx context.Context) (*models.AgentSession, error) {
	q := query.AgentSession.WithContext(ctx)
	sess, err := q.Where(query.AgentSession.SessionID.Eq(s.sessionID)).First()
	if err == nil {
		return sess, nil
	}
	sess = &models.AgentSession{SessionID: s.sessionID, Status: "active"}
	if err := q.Create(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *dbSession) SetMessages(ctx context.Context, msg *schema.Message) error {
	sess, err := s.ensureSession(ctx)
	if err != nil {
		return err
	}

	content, err := messageToJSON(msg)
	if err != nil {
		return err
	}

	if err := query.AgentMessage.WithContext(ctx).Create(&models.AgentMessage{
		SessionID:  sess.ID,
		Role:       string(msg.Role),
		Content:    content,
		TokenCount: 0,
	}); err != nil {
		return err
	}

	s.trim(ctx, sess.ID)
	return nil
}

func (s *dbSession) GetMessages(ctx context.Context) ([]*schema.Message, error) {
	sess, err := s.ensureSession(ctx)
	if err != nil {
		return nil, err
	}

	mq := query.AgentMessage.WithContext(ctx)
	msgs, err := mq.Where(query.AgentMessage.SessionID.Eq(sess.ID)).Order(query.AgentMessage.CreatedAt).Find()
	if err != nil {
		return nil, nil
	}

	out := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		msg, err := messageFromJSON(m.Content)
		if err != nil {
			continue
		}
		out = append(out, msg)
	}
	return out, nil
}

func (s *dbSession) trim(ctx context.Context, sessionID uint64) {
	if s.maxWindowSize <= 0 {
		return
	}
	mq := query.AgentMessage.WithContext(ctx)
	count, err := mq.Where(query.AgentMessage.SessionID.Eq(sessionID)).Count()
	if err != nil {
		return
	}
	excess := int(count) - s.maxWindowSize
	if excess <= 0 {
		return
	}
	if excess%2 != 0 {
		excess++
	}
	oldest, err := mq.Where(query.AgentMessage.SessionID.Eq(sessionID)).
		Order(query.AgentMessage.CreatedAt).Limit(excess).Find()
	if err != nil || len(oldest) == 0 {
		return
	}
	for _, m := range oldest {
		_, _ = mq.Delete(m)
	}
}

func messageToJSON(msg *schema.Message) (string, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func messageFromJSON(data string) (*schema.Message, error) {
	var msg schema.Message
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

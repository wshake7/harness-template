package auth

import (
	"admin/internal/fiberc/res"
	"errors"
	"github.com/click33/sa-token-go/core/manager"
	"github.com/click33/sa-token-go/core/session"
	"github.com/click33/sa-token-go/stputil"
	"go-common/utils/trans"
	"time"
)

const (
	KeyEncryptedKey = "encrypted:key"
)

type SessionInfo struct {
	PrivateKey string   `json:"privateKey"`
	Id         uint64   `json:"id"`
	Username   string   `json:"username"`
	RoleCodes  []string `json:"roleCodes"`
	RoleIDs    []uint64 `json:"roleIDs"`
	Language   string   `json:"language"`
}

type Session struct {
	*session.Session
}

func (s *Session) GetInfo() (SessionInfo, error) {
	obj, err := trans.Map2Obj[SessionInfo](s.Data)
	if err != nil {
		return SessionInfo{}, err
	}
	return obj, err
}

func (s *Session) SaveInfo(info *SessionInfo, ttl ...time.Duration) error {
	data, err := trans.Obj2Map[any](info)
	if err != nil {
		return err
	}
	err = s.SetMulti(data, ttl...)
	if err != nil {
		return err
	}
	return nil
}

func GetSession(loginID any) (*Session, error) {
	s, err := stputil.GetSession(loginID)
	if err != nil {
		return nil, err
	}
	return &Session{s}, nil
}

func GetSessionByToken(token string) (*Session, error) {
	s, err := stputil.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}
	return &Session{s}, nil
}

func CheckLoginErr(err error) error {
	switch {
	case errors.Is(err, manager.ErrAccountDisabled):
		return res.FailAccountDisabled
	case errors.Is(err, manager.ErrNotLogin):
		return res.FailNotLogin
	case errors.Is(err, manager.ErrTokenNotFound):
		return res.FailTokenNotFound
	case errors.Is(err, manager.ErrInvalidTokenData):
		return res.FailInvalidTokenData
	case errors.Is(err, manager.ErrLoginLimitExceeded):
		return res.FailLoginLimitExceeded
	case errors.Is(err, manager.ErrTokenKickout):
		return res.FailTokenKickout
	case errors.Is(err, manager.ErrTokenReplaced):
		return res.FailTokenReplaced
	default:
		return res.FailDefault
	}
}

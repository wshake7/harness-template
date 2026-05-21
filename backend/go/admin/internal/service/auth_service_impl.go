package service

import (
	"admin/internal/auth"

	"github.com/click33/sa-token-go/stputil"
)

type authServiceImpl struct{}

func NewAuthService() AuthService {
	return &authServiceImpl{}
}

func (s *authServiceImpl) Login(id uint64) (string, error) {
	return stputil.Login(id)
}

func (s *authServiceImpl) Logout(loginID any) error {
	return stputil.Logout(loginID)
}

func (s *authServiceImpl) GetLoginID(token string) (any, error) {
	return stputil.GetLoginID(token)
}

func (s *authServiceImpl) SaveSession(loginID uint64, info *auth.SessionInfo) error {
	session, err := auth.GetSession(loginID)
	if err != nil {
		return err
	}
	return session.SaveInfo(info)
}

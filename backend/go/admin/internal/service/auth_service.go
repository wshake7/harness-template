package service

import (
	"admin/internal/auth"
)

//go:generate mockgen -source=auth_service.go -destination=../mock/mock_auth_service.go -package=mock -typed

type AuthService interface {
	Login(id uint64) (token string, err error)
	Logout(loginID any) error
	GetLoginID(token string) (any, error)
	SaveSession(loginID uint64, info *auth.SessionInfo) error
}

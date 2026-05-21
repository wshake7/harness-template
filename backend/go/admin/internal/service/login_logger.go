package service

import (
	"admin/internal/fiberc/handler"
)

//go:generate mockgen -source=login_logger.go -destination=../mock/mock_login_logger.go -package=mock -typed

type LoginLogger interface {
	RecordPwdLogin(ctx *handler.Ctx, username string, userID *uint64, statusCode int, success bool, reason string)
}

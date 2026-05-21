package service

import (
	"admin/internal/fiberc/handler"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"time"

	"github.com/mileusna/useragent"
	"go-common/utils/coroutine"
	"go-common/utils/ip_util"
	"go.uber.org/zap"
)

type loginLoggerImpl struct{}

func NewLoginLogger() LoginLogger {
	return &loginLoggerImpl{}
}

func (l *loginLoggerImpl) RecordPwdLogin(ctx *handler.Ctx, username string, userID *uint64, statusCode int, success bool, reason string) {
	if ctx == nil {
		return
	}

	logger := ctx.L()
	ip := ctx.IP()
	userAgent := ctx.UserAgent()
	loginTime := time.Now()

	coroutine.Launch(func() {
		location := ""
		if result, err := ip_util.Client.Query(ip); err == nil {
			location = result.String()
		} else if logger != nil {
			logger.Error("query login IP location fail", zap.Error(err), zap.String("ip", ip))
		}

		ua := useragent.Parse(userAgent)
		clientName := ua.Device
		if clientName == "" && ua.Desktop {
			clientName = "PC"
		}

		loginLog := &models.SysLoginLog{
			Username:       username,
			LoginIP:        ip,
			LoginMAC:       "",
			LoginTime:      &loginTime,
			UserAgent:      userAgent,
			BrowserName:    ua.Name,
			BrowserVersion: ua.Version,
			ClientID:       "",
			ClientName:     clientName,
			OSName:         ua.OS,
			OSVersion:      ua.OSVersion,
			SysUserID:      userID,
			StatusCode:     int32(statusCode),
			Success:        success,
			Reason:         reason,
			Location:       location,
		}
		if err := query.SysLoginLog.Create(loginLog); err != nil && logger != nil {
			logger.Error("SysLoginLog.Create fail", zap.Error(err), zap.String("username", username))
		}
	})
}

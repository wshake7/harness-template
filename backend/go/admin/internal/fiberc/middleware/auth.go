package middleware

import (
	"admin/internal/auth"
	"admin/internal/config"
	"admin/internal/domains"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	redisc2 "admin/internal/services/redisc"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func AuthMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		headerToken := fiber.GetReqHeader[string](ctx, config.Conf.Auth.TokenName)
		cookieToken := ctx.Cookies(config.Conf.Auth.TokenName)
		token := headerToken
		if token == "" {
			token = cookieToken
		}
		if token == "" {
			return res.FailNotLogin
		}

		session, err := auth.GetSessionByToken(token)
		if err != nil {
			ctx.L().Warn("header token failed, fallback to cookie token", zap.Error(err))
			session, err = auth.GetSessionByToken(token)
		}
		if err != nil {
			ctx.L().Error("get session error", zap.Error(err))
			return auth.CheckLoginErr(err)
		}
		ctx.SessionRaw = session
		info, err := session.GetInfo()
		if err != nil || info.PrivateKey == "" {
			ctx.L().Error("get session info error", zap.Error(err), zap.String("key", info.PrivateKey))
			return res.FailDefault
		}
		ctx.PrivateKey = info.PrivateKey
		ctx.SessionInfo = &info
		ctx.AddResLogFields(zap.Any(domains.LogFieldSessionInfo, info))
		ctx.AddLogFields(zap.Any(domains.LogFieldSessionInfo, info))
		return ctx.Next()
	})
}

func PublicMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		var keyPair redisc2.DtoKeyPair
		err := redisc2.Client.GetJson(ctx, redisc2.KeyGlobalEncryptPublicKey, &keyPair)
		if err != nil || keyPair.PrivateKey == "" {
			ctx.L().Error("get key error", zap.Error(err), zap.String("key", keyPair.PrivateKey))
			return res.FailRequestKey
		}
		ctx.PrivateKey = keyPair.PrivateKey
		return ctx.Next()
	})
}

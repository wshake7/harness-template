package middleware

import (
	domains2 "admin/internal/domains"
	"admin/internal/fiberc/handler"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"strings"
)

func LanguageMiddleware(defaultLanguage string) fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		language := strings.TrimSpace(fiber.GetReqHeader[string](ctx, domains2.HeaderXLanguage))
		if language == "" {
			language = defaultLanguage
		}
		ctx.Language = language
		info := ctx.SessionInfo
		if info != nil && info.Language != language {
			ctx.SessionInfo.Language = language
			err := ctx.SessionRaw.SaveInfo(ctx.SessionInfo)
			if err != nil {
				ctx.L().Error("save session info error", zap.Error(err), zap.String("language", language))
			}
		}
		return ctx.Next()
	})
}

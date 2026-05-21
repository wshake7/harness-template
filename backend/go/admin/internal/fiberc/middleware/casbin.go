package middleware

import (
	"fmt"
	"strings"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/casbin"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func CasbinAPIMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		info := ctx.SessionInfo
		if info == nil {
			ctx.L().Warn("casbin auth skipped because session info is empty")
			return res.FailNotLogin
		}
		if casbin.E == nil {
			ctx.L().Error("casbin enforcer is not initialized")
			return res.FailDefault
		}

		subjects := []string{fmt.Sprintf("user:%d", info.Id)}
		for _, role := range info.RoleCodes {
			role = strings.TrimSpace(role)
			if role == "" {
				continue
			}
			subjects = append(subjects, "role:"+role)
		}

		obj := ctx.Path()
		act := strings.ToUpper(ctx.Method())

		for _, sub := range subjects {
			ok, err := casbin.E.Enforce(sub, obj, act)
			if err != nil {
				if strings.Contains(err.Error(), "please make sure rule exists in policy when using eval() in matcher") {
					ctx.L().Warn("casbin auth denied because no evaluable policy matched",
						zap.String("subject", sub),
						zap.String("object", obj),
						zap.String("action", act),
					)
					continue
				}
				ctx.L().Error("casbin enforce error",
					zap.Error(err),
					zap.String("subject", sub),
					zap.String("object", obj),
					zap.String("action", act),
				)
				return res.FailDefault
			}
			if ok {
				return ctx.Next()
			}
		}

		ctx.L().Warn("casbin auth denied",
			zap.Strings("subjects", subjects),
			zap.String("object", obj),
			zap.String("action", act),
		)
		return res.FailAuthUnauthorized
	})
}

package router

import (
	"admin/internal/appsvc"
	"admin/internal/config"
	"admin/internal/fiberc/handler"
	middleware2 "admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"
	"github.com/gofiber/fiber/v3"
)

func registerAccountRouters(router fiber.Router, conf *config.Config) {
	accountHandler := logic.NewAccountHandler(
		query.Q,
		appsvc.NewAuthService(),
		appsvc.NewLoginLogger(),
	)
	router.Post("/login/pwd", middleware2.PublicMiddleware(), middleware2.EncryptMiddleware(), handler.CtxHandlerFunc(func(ctx *handler.Ctx, req *logic.ReqAccountPwdLogin) (*logic.ResAccountPwdLogin, error) {
		res, err := accountHandler.PwdLogin(ctx, req)
		if err != nil {
			return res, err
		}
		ctx.Cookie(&fiber.Cookie{
			Name:     conf.Auth.TokenName,
			Value:    res.Token,
			SameSite: "Lax",
			Secure:   false,
			HTTPOnly: true,
		})
		return res, err
	}))
	router.Post("/changePwd", middleware2.AuthMiddleware(conf.Auth.TokenName), middleware2.EncryptMiddleware(), handler.CtxHandlerNilFunc(accountHandler.ChangePwd))
	router.Get("/logout", middleware2.AuthMiddleware(conf.Auth.TokenName), handler.CtxHandlerNilFunc(accountHandler.Logout))
}

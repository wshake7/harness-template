package auth_router

import (
	"admin/internal/fiberc/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterRouters(router fiber.Router) {
	eventGroup := router.Use(middleware.AuthMiddleware(), middleware.CasbinAPIMiddleware(), middleware.LanguageMiddleware())
	registerEventRouters(eventGroup)

	group := router.Use(middleware.AuthMiddleware(), middleware.CasbinAPIMiddleware(), middleware.EncryptMiddleware(), middleware.LanguageMiddleware())
	registerSysRoleRouters(group.Group("/sys/role"))
	registerSysUserRouters(group.Group("/sys/user"))
	registerSysDictRouters(group.Group("/sys/dict"))
	registerSysLanguageRouters(group.Group("/sys/language"))
	registerSysApiLogRouters(group.Group("/sys/api/log"))
	registerSysLoginLogRouters(group.Group("/sys/login/log"))
	registerSysResourceMenuRouters(group.Group("/sys/resource/menu"))
	registerSysResourceApiRouters(group.Group("/sys/resource/api"))
	registerJobScheduleRouters(group.Group("/sys/job/schedule"))
	registerJobExecutionRouters(group.Group("/sys/job/execution"))
}

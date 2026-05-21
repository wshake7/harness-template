package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/service"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerJobScheduleRouters(router fiber.Router) {
	jobScheduleHandler := logic.NewJobScheduleHandler(query.Q, service.NewTemporalService())
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_schedule"))
	createLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_schedule"), middleware.WithChangeQuery(jobScheduleCreateChangeQuery))
	updateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_schedule"), middleware.WithChangeQuery(jobScheduleUpdateChangeQuery))
	switchLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_schedule"), middleware.WithChangeQuery(jobScheduleSwitchChangeQuery))

	router.Post("/options", handler.CtxFunc(jobScheduleHandler.Options))
	router.Post("/list", handler.CtxHandlerFunc(jobScheduleHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(jobScheduleHandler.Detail))
	router.Post("/create", createLogMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Create))
	router.Post("/update", updateLogMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Del))
	router.Post("/switch", switchLogMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Switch))
	router.Post("/sync", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Sync))
	router.Post("/trigger", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Trigger))
}

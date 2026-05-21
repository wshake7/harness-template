package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/service"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerJobExecutionRouters(router fiber.Router) {
	jobExecutionHandler := logic.NewJobExecutionHandler(query.Q, service.NewTemporalService())
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_execution"))
	changeLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_execution"), middleware.WithChangeQuery(jobExecutionChangeQuery))

	router.Post("/list", handler.CtxHandlerFunc(jobExecutionHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(jobExecutionHandler.Detail))
	router.Post("/cancel", changeLogMiddleware, handler.CtxHandlerNilFunc(jobExecutionHandler.Cancel))
	router.Post("/retry", logMiddleware, handler.CtxHandlerNilFunc(jobExecutionHandler.Retry))
}

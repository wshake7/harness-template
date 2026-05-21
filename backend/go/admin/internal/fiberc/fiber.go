package fiberc

import (
	"admin/internal/config"
	"admin/internal/domains"
	"admin/internal/fiberc/handler"
	middleware2 "admin/internal/fiberc/middleware"
	"admin/internal/fiberc/res"
	"admin/internal/lifecycle"
	"context"
	"encoding/xml"
	"errors"
	fiberzap "github.com/gofiber/contrib/v3/zap"
	"html/template"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/monitor"
	"github.com/gofiber/fiber/v3/binder"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "admin/docs"
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
)

func NewFiber(conf *config.Config) *App {
	app := initialize(conf)
	initializeHoos(app)
	logger := zap.L().With(zap.String("module", "fiberc"))
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, err any) {
			ctx := handler.Trans(c)
			ctx.IsPanic = true
			// 用 zap 打印，包含完整堆栈
			ctx.L().Error("panic recovered",
				zap.Any("error", err),
				zap.String("stack", string(debug.Stack())),
			)
		},
	}))
	app.Use(middleware2.LogMiddleware(logger))
	app.Use(middleware2.TraceMiddleware())
	app.Use(cors.New())
	app.Use(pprof.New(pprof.Config{Prefix: "/pprof"}))
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get("/metrics", monitor.New())
	app.Get("/prometheus", adaptor.HTTPHandler(promhttp.Handler()))
	if app.conf.IsSwagger {
		swgCfg := swaggo.Config{
			Title:              "Swagger UI",
			Layout:             "StandaloneLayout",
			URL:                "doc.json",
			DeepLinking:        true,
			ShowMutatedRequest: true,
			Plugins: []template.JS{
				template.JS("SwaggerUIBundle.plugins.DownloadUrl"),
			},
			Presets: []template.JS{
				template.JS("SwaggerUIBundle.presets.apis"),
				template.JS("SwaggerUIStandalonePreset"),
			},
			SyntaxHighlight: &swaggo.SyntaxHighlightConfig{Activate: true, Theme: "agate"},
		}
		app.Get(app.conf.SwaggerPrefix+"/*", swaggo.New(swgCfg))
	}
	logCfg := fiberzap.Config{
		Logger: logger,
		//Fields: []string{"ip", "latency", "status", "method", "url", "resBody", "body", "queryParams"},
		//FieldsFunc: func(c fiber.Ctx) []zap.Field {
		//	ctx := handler.Trans(c)
		//	return ctx.LogResFields
		//},
	}
	app.Use(fiberzap.New(logCfg))
	app.done = gracefulShutdown(app)
	return app
}

func initialize(conf *config.Config) *App {
	fiberConfig := conf.Fiber
	fiberApp := fiber.NewWithCustomCtx(func(app *fiber.App) fiber.CustomCtx {
		return &handler.Ctx{
			DefaultCtx: *fiber.NewDefaultCtx(app),
			Config:     conf,
		}
	}, fiber.Config{
		ServerHeader:            fiberConfig.ServerHeader,
		StrictRouting:           fiberConfig.StrictRouting,
		CaseSensitive:           fiberConfig.CaseSensitive,
		DisableHeadAutoRegister: fiberConfig.DisableHeadAutoRegister,
		Immutable:               fiberConfig.Immutable,
		UnescapePath:            fiberConfig.UnescapePath,
		BodyLimit:               fiberConfig.BodyLimit,
		MaxRanges:               fiberConfig.MaxRanges,
		Concurrency:             fiberConfig.Concurrency,
		Views:                   nil,
		ViewsLayout:             fiberConfig.ViewsLayout,
		PassLocalsToViews:       fiberConfig.PassLocalsToViews,
		PassLocalsToContext:     fiberConfig.PassLocalsToContext,
		ReadTimeout:             fiberConfig.ReadTimeout,
		WriteTimeout:            fiberConfig.WriteTimeout,
		IdleTimeout:             fiberConfig.IdleTimeout,
		ReadBufferSize:          fiberConfig.ReadBufferSize,
		WriteBufferSize:         fiberConfig.WriteBufferSize,
		CompressedFileSuffixes:  map[string]string{"gzip": ".fiber.gz", "br": ".fiber.br", "zstd": ".fiber.zst"},
		ProxyHeader:             fiberConfig.ProxyHeader,
		GETOnly:                 fiberConfig.GETOnly,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			ctx := handler.Trans(c)
			isPanic := ctx.IsPanic
			ctx.IsPanic = false
			code := fiber.StatusOK
			var fe *fiber.Error
			if errors.As(err, &fe) {
				if fe.Code == fiber.StatusNotFound {
					code = fe.Code // ← 404、405 等 Fiber 内置错误
				}
			}
			ctx.Status(code)
			if err != nil {
				var e res.Response
				switch {
				case errors.As(err, &e):
					err = ctx.JSON(e)
					if err != nil {
						ctx.L().Error("Fiber返回信息错误", zap.Error(err))
					}
				default:
					if !isPanic {
						ctx.L().Error("Fiber处理信息错误", zap.Error(err))
					}
					err = ctx.JSON(domains.JsonErr)
					if err != nil {
						ctx.L().Error("Fiber返回信息错误", zap.Error(err))
					}
				}
			}
			return nil
		},
		DisableKeepalive:             fiberConfig.DisableKeepalive,
		DisableDefaultDate:           fiberConfig.DisableDefaultDate,
		DisableDefaultContentType:    fiberConfig.DisableDefaultContentType,
		DisableHeaderNormalizing:     fiberConfig.DisableHeaderNormalizing,
		AppName:                      fiberConfig.AppName,
		StreamRequestBody:            fiberConfig.StreamRequestBody,
		DisablePreParseMultipartForm: fiberConfig.DisablePreParseMultipartForm,
		ReduceMemoryUsage:            fiberConfig.ReduceMemoryUsage,
		JSONEncoder:                  sonic.Marshal,
		JSONDecoder:                  sonic.Unmarshal,
		MsgPackEncoder:               binder.UnimplementedMsgpackMarshal,
		MsgPackDecoder:               binder.UnimplementedMsgpackUnmarshal,
		CBOREncoder:                  binder.UnimplementedCborMarshal,
		CBORDecoder:                  binder.UnimplementedCborUnmarshal,
		XMLEncoder:                   xml.Marshal,
		XMLDecoder:                   xml.Unmarshal,
		TrustProxy:                   fiberConfig.TrustProxy,
		TrustProxyConfig:             fiber.DefaultTrustProxyConfig,
		EnableIPValidation:           fiberConfig.EnableIPValidation,
		ColorScheme:                  fiber.DefaultColors,
		StructValidator:              nil,
		RequestMethods:               fiberConfig.RequestMethods,
		EnableSplittingOnParsers:     fiberConfig.EnableSplittingOnParsers,
		Services:                     fiberConfig.Services,
		ServicesStartupContextProvider: func() context.Context {
			return context.Background()
		},
		ServicesShutdownContextProvider: func() context.Context {
			return context.Background()
		},
	})
	return &App{
		App:  fiberApp,
		conf: conf,
	}
}

func initializeHoos(app *App) {
	// 在启动消息打印前执行，允许自定义横幅和信息条目
	app.Hooks().OnPreStartupMessage(func(data *fiber.PreStartupMessageData) error {
		return nil
	})
	// 在启动消息打印后执行，支持启动后逻辑
	app.Hooks().OnPostStartupMessage(func(data *fiber.PostStartupMessageData) error {
		return nil
	})
	// 在服务器开始关闭流程前执行，用于处理需要在关闭前完成的清理任务，例如停止后台任务、保存状态或关闭数据库连接
	app.Hooks().OnPreShutdown(func() error {
		lifecycle.BeginShutdown()
		zap.L().Info("服务器关闭OnPreShutdown")
		return nil
	})

	// 关闭后：释放资源，在服务器完全关闭后执行，通常用于执行最后的收尾工作，例如记录日志、释放资源或通知其他服务
	app.Hooks().OnPostShutdown(func(err error) error {
		zap.L().Info("服务器关闭OnPostShutdown")
		if err != nil {
			zap.L().Error("服务器异常关闭", zap.Error(err))
			return nil
		}
		return nil
	})
}

func gracefulShutdown(app *App) chan struct{} {
	done := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		zap.L().Info("收到服务器关闭信号")
		err := app.ShutdownWithTimeout(5 * time.Second)
		if err != nil {
			zap.L().Error("服务器关闭异常")
		} else {
			zap.L().Info("服务器关闭成功")
		}
		_ = zap.L().Sync()
		close(done)
	}()
	return done
}

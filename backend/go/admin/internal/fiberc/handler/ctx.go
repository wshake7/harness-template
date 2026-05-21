package handler

import (
	"admin/internal/auth"
	"admin/internal/config"
	"admin/internal/domains"
	"admin/internal/fiberc/res"
	"admin/internal/validator"
	"errors"
	"go-common/utils/types"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Trans(ctx fiber.Ctx) *Ctx {
	c := ctx.(*Ctx)
	return c
}

func CtxNilMiddlewareFunc(fn func(ctx *Ctx) error) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		c := Trans(ctx)
		err := fn(c)
		if err != nil {
			c.ErrCode = &domains.StatusFail
			c.ErrMsg = err.Error()
			return err
		}
		return nil
	}
}

func CtxNilFunc(fn func(ctx *Ctx) error) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		c := Trans(ctx)
		err := fn(c)
		if err != nil {
			var e res.Response
			switch {
			case errors.As(err, &e):
				c.ErrCode = &e.Code
				c.ErrMsg = e.Msg
				return e
			default:
				c.ErrCode = &domains.StatusFail
				c.ErrMsg = e.Error()
				return c.FailErr(err)
			}
		}
		return c.Ok0()
	}
}

func CtxFunc[R any](fn func(ctx *Ctx) (*R, error)) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		c := Trans(ctx)
		r, err := fn(c)
		if err != nil {
			var e res.Response
			switch {
			case errors.As(err, &e):
				c.ErrCode = &e.Code
				c.ErrMsg = e.Msg
				return e
			default:
				c.ErrCode = &domains.StatusFail
				c.ErrMsg = e.Error()
				return c.FailErr(err)
			}
		}
		return c.OkStruct(r)
	}
}

func CtxHandlerFunc[T any, R any](fn func(ctx *Ctx, t *T) (*R, error)) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		c := Trans(ctx)
		params := new(T)
		err := c.StructBind(params)
		if err != nil {
			c.ErrCode = &domains.StatusFail
			c.ErrMsg = err.Error()
			return c.FailErr(err)
		}
		r, err := fn(c, params)
		if err != nil {
			var e res.Response
			switch {
			case errors.As(err, &e):
				c.ErrCode = &e.Code
				c.ErrMsg = e.Msg
				return e
			default:
				c.ErrCode = &domains.StatusFail
				c.ErrMsg = e.Error()
				return c.FailErr(err)
			}
		}
		return c.OkStruct(r)
	}
}

func CtxHandlerNilFunc[T any](fn func(ctx *Ctx, t *T) error) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		c := Trans(ctx)
		params := new(T)
		err := c.StructBind(params)
		if err != nil {
			c.ErrCode = &domains.StatusFail
			c.ErrMsg = err.Error()
			return c.FailErr(err)
		}
		err = fn(c, params)
		if err != nil {
			var e res.Response
			switch {
			case errors.As(err, &e):
				c.ErrCode = &e.Code
				c.ErrMsg = e.Msg
				return e
			default:
				c.ErrCode = &domains.StatusFail
				c.ErrMsg = e.Error()
				return c.FailErr(err)
			}
		}
		return c.Ok0()
	}
}

type Ctx struct {
	fiber.DefaultCtx
	Config       *config.Config
	TraceId      string
	PrivateKey   string
	Language     string
	SessionInfo  *auth.SessionInfo
	SessionRaw   *auth.Session
	logger       *zap.Logger
	LogResFields []zap.Field
	ErrCode      *int
	ErrMsg       string
	IsPanic      bool
}

func (ctx *Ctx) SetLogger(logger *zap.Logger) {
	if logger != nil {
		ctx.logger = logger
	}
}

func (ctx *Ctx) AddLogFields(field ...zap.Field) {
	ctx.logger = ctx.logger.With(field...)
}

func (ctx *Ctx) AddResLogFields(field zap.Field) {
	ctx.LogResFields = append(ctx.LogResFields, field)
}

func (ctx *Ctx) L() *zap.Logger {
	return ctx.logger
}

func (ctx *Ctx) S() *zap.SugaredLogger {
	return ctx.logger.Sugar()
}

func (ctx *Ctx) FailMsg(msg string) error {
	return ctx.JSON(res.FailMsg(msg))
}

func (ctx *Ctx) FailMsgMust(msg string) {
	_ = ctx.FailMsg(msg)
}

func (ctx *Ctx) FailErr(msg error) error {
	return ctx.JSON(res.FailMsg(msg.Error()))
}

func (ctx *Ctx) FailErrMust(msg error) {
	_ = ctx.FailErr(msg)
}

func (ctx *Ctx) Ok0() error {
	return ctx.JSON(domains.JsonOk)
}

func (ctx *Ctx) Ok(data any) error {
	if data == nil {
		return ctx.JSON(domains.JsonEmpty)
	}
	return ctx.JSON(res.OkRes(data))
}

func (ctx *Ctx) OkStruct(data any) error {
	if data == nil {
		return ctx.JSON(domains.JsonEmptyStruct)
	}

	return ctx.JSON(res.OkRes(data))
}

func (ctx *Ctx) OkSlice(data any) error {
	if data == nil {
		return ctx.JSON(domains.JsonEmptySlice)
	}
	return ctx.JSON(res.OkRes(data))
}

func (ctx *Ctx) Ok0Must() {
	_ = ctx.Ok0()
}

func (ctx *Ctx) OkMust(data any) {
	_ = ctx.Ok(data)
}

func (ctx *Ctx) OkStructMust(data any) {
	_ = ctx.OkStruct(data)
}

func (ctx *Ctx) OkSliceMust(data any) {
	_ = ctx.OkSlice(data)
}

func (ctx *Ctx) bindFn(obj any, fn func() error) error {
	if !types.IsPointer(obj) {
		return errors.New("obj is not a pointer")
	}
	err := fn()
	if err != nil {
		return err
	}
	if err = validator.Struct(obj); err != nil {
		return err
	}
	switch v := obj.(type) {
	case validator.Validator:
		return v.Validate()
	default:
		return nil
	}
}

func (ctx *Ctx) StructBind(obj any) error {
	return ctx.bindFn(obj, func() error {
		if _, ok := obj.(protoreflect.ProtoMessage); ok {
			body := ctx.Request().Body()
			err := protojson.Unmarshal(body, obj.(protoreflect.ProtoMessage))
			return err
		}
		return ctx.DefaultCtx.Bind().All(obj)
	})
}

func (ctx *Ctx) BodyBind(obj any) error {
	return ctx.bindFn(obj, func() error {
		return ctx.DefaultCtx.Bind().Body(obj)
	})
}

func (ctx *Ctx) URIBind(obj any) error {
	return ctx.bindFn(obj, func() error {
		return ctx.DefaultCtx.Bind().URI(obj)
	})
}

func (ctx *Ctx) QueryBind(obj any) error {
	return ctx.bindFn(obj, func() error {
		return ctx.DefaultCtx.Bind().Query(obj)
	})
}

func (ctx *Ctx) CookieBind(obj any) error {
	return ctx.bindFn(obj, func() error {
		return ctx.DefaultCtx.Bind().Cookie(obj)
	})
}

func (ctx *Ctx) HeaderBind(obj any) error {
	return ctx.bindFn(obj, func() error {
		return ctx.DefaultCtx.Bind().Header(obj)
	})
}

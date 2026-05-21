package middleware

import (
	"admin/internal/fiberc/handler"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/mileusna/useragent"
	"go-common/utils/coroutine"
	"go-common/utils/function"
	"go-common/utils/ip_util"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const maxApiLogFieldLength = 64 * 1024

var (
	apiLogSensitiveTextPattern = regexp.MustCompile(`(?i)("?(authorization|cookie|set-cookie|token|access_token|refresh_token|password|pwd|secret)"?\s*[:=]\s*)("[^"]*"|[^,\s&}]+)`)
	apiLogSensitiveKeys        = map[string]struct{}{
		"authorization": {},
		"cookie":        {},
		"set-cookie":    {},
		"token":         {},
		"access_token":  {},
		"refresh_token": {},
		"password":      {},
		"pwd":           {},
		"secret":        {},
	}
)

func sanitizeApiLogPayload(raw string) string {
	if raw == "" {
		return ""
	}

	var payload any
	if err := sonic.UnmarshalString(raw, &payload); err == nil {
		redactApiLogJSONValue(payload)
		if sanitized, err := sonic.MarshalString(payload); err == nil {
			return truncateApiLogField(sanitized)
		}
	}

	return truncateApiLogField(apiLogSensitiveTextPattern.ReplaceAllString(raw, `${1}"***"`))
}

func redactApiLogJSONValue(value any) {
	switch typedValue := value.(type) {
	case map[string]any:
		for key, item := range typedValue {
			if isApiLogSensitiveKey(key) {
				typedValue[key] = "***"
				continue
			}
			redactApiLogJSONValue(item)
		}
	case []any:
		for _, item := range typedValue {
			redactApiLogJSONValue(item)
		}
	}
}

func isApiLogSensitiveKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	if _, ok := apiLogSensitiveKeys[key]; ok {
		return true
	}
	return strings.Contains(key, "token") || strings.Contains(key, "password")
}

func truncateApiLogField(value string) string {
	if len(value) <= maxApiLogFieldLength {
		return value
	}
	return value[:maxApiLogFieldLength] + "...[truncated]"
}

func DiffChange(before, after any) string {
	if after == nil {
		return ""
	}

	// 反射查找 ChangeString 方法
	afterVal := reflect.ValueOf(after)
	method := afterVal.MethodByName("ChangeString")
	if method.IsValid() {
		var beforeVal reflect.Value
		if before == nil {
			// before 为 nil，传零值
			beforeVal = reflect.Zero(afterVal.Type())
		} else {
			beforeVal = reflect.ValueOf(before)
		}
		results := method.Call([]reflect.Value{beforeVal})
		return results[0].String()
	}

	if before == nil {
		return diffNewByTag(after)
	}
	return diffByTag(before, after)
}

func diffByTag(before, after any) string {
	bVal := reflect.ValueOf(before)
	aVal := reflect.ValueOf(after)
	// 解指针
	for bVal.Kind() == reflect.Pointer {
		bVal = bVal.Elem()
	}
	for aVal.Kind() == reflect.Pointer {
		aVal = aVal.Elem()
	}
	if bVal.Kind() != reflect.Struct || aVal.Kind() != reflect.Struct {
		return ""
	}
	t := bVal.Type()
	var parts []string
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("change")
		if tag == "" {
			continue
		}
		bField := bVal.Field(i)
		aField := aVal.Field(i)
		if !bField.CanInterface() || !aField.CanInterface() {
			continue
		}
		if !reflect.DeepEqual(bField.Interface(), aField.Interface()) {
			parts = append(parts, fmt.Sprintf("%s：%s->%s", tag, formatChangeFieldValue(bField), formatChangeFieldValue(aField)))
		}
	}
	return strings.Join(parts, ", ")
}

func diffNewByTag(after any) string {
	aVal := reflect.ValueOf(after)
	for aVal.Kind() == reflect.Pointer {
		aVal = aVal.Elem()
	}
	if aVal.Kind() != reflect.Struct {
		return ""
	}
	t := aVal.Type()
	var parts []string
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("change")
		if tag == "" {
			continue
		}
		aField := aVal.Field(i)
		if !aField.CanInterface() {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s：%s", tag, formatChangeFieldValue(aField)))
	}
	return strings.Join(parts, ", ")
}

func formatChangeFieldValue(value reflect.Value) string {
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return "空"
		}
		value = value.Elem()
	}
	if !value.IsValid() || !value.CanInterface() {
		return "空"
	}
	return fmt.Sprint(value.Interface())
}

type Changeable[T any] interface {
	ChangeString(before T) string
}

type ChangeQueryHandler func(ctx *handler.Ctx) (any, error)

func ChangeQueryNilFn[R any](fn func(ctx *handler.Ctx) (*R, error)) ChangeQueryHandler {
	return func(ctx *handler.Ctx) (any, error) {
		r, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

}

func ChangeQueryParamsFn[T any, R any](fn func(ctx *handler.Ctx, t *T) (*R, error)) ChangeQueryHandler {
	return func(ctx *handler.Ctx) (any, error) {
		t := new(T)
		err := ctx.DefaultCtx.Bind().All(t)
		if err != nil {
			return nil, err
		}
		r, err := fn(ctx, t)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}

type ApiLogConfig struct {
	BeforeChangeQuery ChangeQueryHandler
	AfterChangeQuery  ChangeQueryHandler
	Module            string
}

type Option func(*ApiLogConfig)

func WithChangeQuery(fn ChangeQueryHandler) Option {
	return func(config *ApiLogConfig) {
		config.BeforeChangeQuery = fn
		config.AfterChangeQuery = fn
	}
}

func WithModule(module string) Option {
	return func(config *ApiLogConfig) {
		config.Module = module
	}
}

func isForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "foreign key") || strings.Contains(msg, "violates foreign key constraint")
}

func createSysApiLog(logger *zap.Logger, m *models.SysApiLog) {
	err := query.SysApiLog.Create(m)
	if err == nil {
		return
	}
	if m.SysUserID != nil && isForeignKeyViolation(err) {
		invalidSysUserID := *m.SysUserID
		m.SysUserID = nil
		retryErr := query.SysApiLog.Create(m)
		if retryErr == nil {
			logger.Warn("SysApiLog.Create retry without sys_user_id", zap.Uint64("sysUserID", invalidSysUserID), zap.Error(err))
			return
		}
		logger.Error("SysApiLog.Create retry without sys_user_id fail", zap.Uint64("sysUserID", invalidSysUserID), zap.Error(retryErr), zap.NamedError("firstError", err))
		return
	}
	logger.Error("SysApiLog.Create fail", zap.Error(err))
}

func ApiLogMiddleware(options ...Option) fiber.Handler {
	conf := &ApiLogConfig{}
	for _, option := range options {
		option(conf)
	}
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		now := time.Now()
		var beforeChange string
		var beforeChangeData any
		if conf.BeforeChangeQuery != nil {
			function.RecFn(func() {
				var err error
				beforeChangeData, err = conf.BeforeChangeQuery(ctx)
				if err != nil {
					ctx.L().Error("BeforeChangeQuery fail", zap.Error(err))
					return
				}
				beforeChange, err = sonic.MarshalString(beforeChangeData)
				if err != nil {
					ctx.L().Error("BeforeChangeQuery marshal fail", zap.Error(err))
				}
			})
		}
		defer func() {
			nextRecoverErr := recover()
			logger := ctx.L()
			var afterChange string
			var afterChangeData any
			if conf.AfterChangeQuery != nil {
				function.RecFn(func() {
					var err error
					afterChangeData, err = conf.AfterChangeQuery(ctx)
					if err != nil {
						logger.Error("AfterChangeQuery fail", zap.Error(err))
						return
					}
					afterChange, err = sonic.MarshalString(afterChangeData)
					if err != nil {
						logger.Error("AfterChangeQuery marshal fail", zap.Error(err))
					}
				})
			}

			ip := ctx.IP()
			method := ctx.Method()
			path := ctx.Path()
			requestID := ctx.TraceId
			requestBody := sanitizeApiLogPayload(string(ctx.Request().Body()))
			responseBody := sanitizeApiLogPayload(string(ctx.Response().Body()))
			statusCode := ctx.Response().StatusCode()
			headers, _ := sonic.MarshalString(ctx.GetReqHeaders())
			headers = sanitizeApiLogPayload(headers)
			referer := ctx.Referer()
			requestURI := ctx.Request().URI().String()
			userAgent := ctx.UserAgent()
			userId := ctx.SessionInfo.Id
			var sysUserID *uint64
			if userId > 0 {
				sysUserID = &userId
			}
			errMsg := ctx.ErrMsg
			errCode := ctx.ErrCode
			costTime := time.Since(now).Milliseconds()
			var formatChange string
			function.RecFn(func() {
				formatChange = DiffChange(beforeChangeData, afterChangeData)
			})
			coroutine.Launch(func() {
				result, err := ip_util.Client.Query(ip)
				if err != nil {
					logger.Error("Query Ip fail", zap.Error(err))
				}
				ua := useragent.Parse(userAgent)
				var clientName string
				if ua.Device != "" {
					clientName = ua.Device
				} else {
					if ua.Desktop {
						clientName = "PC"
					}
				}
				m := &models.SysApiLog{
					RequestID:      requestID,
					Method:         method,
					Module:         conf.Module,
					Path:           path,
					Referer:        referer,
					BeforeChange:   beforeChange,
					AfterChange:    afterChange,
					FormatChange:   formatChange,
					RequestURI:     requestURI,
					RequestBody:    requestBody,
					RequestHeader:  headers,
					Response:       responseBody,
					CostTime:       costTime,
					SysUserID:      sysUserID,
					ClientIP:       ip,
					StatusCode:     statusCode,
					Reason:         errMsg,
					Success:        errCode == nil,
					Location:       result.String(),
					UserAgent:      userAgent,
					BrowserName:    ua.Name,
					BrowserVersion: ua.Version,
					ClientID:       "",
					ClientName:     clientName,
					OSName:         ua.OS,
					OSVersion:      ua.OSVersion,
				}
				createSysApiLog(logger, m)
			})

			if nextRecoverErr != nil {
				panic(nextRecoverErr)
			}
		}()
		return ctx.Next()
	})
}

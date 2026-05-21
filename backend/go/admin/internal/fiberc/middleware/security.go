package middleware

import (
	domains2 "admin/internal/domains"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	redisc2 "admin/internal/services/redisc"
	"go-common/utils/encrypt"
	"go-common/utils/encrypt/aes_util"
	"go-common/utils/encrypt/rsa_util"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func TimestampMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		timestamp := fiber.GetReqHeader[int64](ctx, domains2.HeaderXRequestTimestamp)
		if timestamp != 0 {
			now := time.Now().UnixMilli()
			if abs(now-timestamp) > encrypt.REQUEST_EXPIRE_TIME.Milliseconds() {
				return res.FailRequestExpired
			}
		}
		return ctx.Next()
	})
}

func NonceMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		if redisc2.Client != nil {
			nonce := ctx.RequestID()
			if nonce != "" {
				nonceKey := encrypt.NONCE_REDIS_KEY_PREFIX.Sprintf(nonce)

				err := redisc2.Client.Do(ctx.Context(), redisc2.Client.B().Set().Key(nonceKey).Value(nonce).Ex(encrypt.NONCE_EXPIRE_TIME).Build()).Error()
				if err != nil {
					ctx.L().Error("redis error", zap.Error(err))
					return res.FailDefault
				}
			}
		}
		return ctx.Next()
	})
}

func EncryptMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		ctx.Set(domains2.HeaderXResponseIsEncrypt, "false")
		aesKey, decrypt, err := DecryptRequest(ctx)
		if err != nil {
			return err
		}
		if len(ctx.Request().Body()) > 0 {
			ctx.Request().SetBody(decrypt)
		}
		err = ctx.Next()
		if err != nil {
			return err
		}
		resBody := ctx.Response().Body()
		resBodyStr := string(resBody)
		ctx.AddResLogFields(zap.String(domains2.LogFieldDecryptResBody, resBodyStr))
		encryptedBody, err := EncryptText(resBodyStr, aesKey)
		if err != nil {
			ctx.L().Error("encrypt error", zap.Error(err))
			return res.FailDefault
		}
		ctx.Set(domains2.HeaderXResponseIsEncrypt, "true")
		ctx.Response().SetBody([]byte(encryptedBody))
		return err
	})
}

func DecryptRequest(ctx *handler.Ctx) (string, []byte, error) {
	aesKey, err := DecryptRequestAESKey(ctx)
	if err != nil {
		return "", nil, err
	}
	decrypt, err := DecryptRequestBody(ctx, aesKey)
	if err != nil {
		return "", nil, err
	}
	return aesKey, decrypt, nil
}

func DecryptRequestAESKey(ctx *handler.Ctx) (string, error) {
	privateKey := ctx.PrivateKey
	if privateKey == "" {
		ctx.L().Error("privateKey is empty")
		return "", res.FailRequest
	}
	encryptedKey := fiber.GetReqHeader[string](ctx, domains2.HeaderXRequestEncryptedKey)
	if encryptedKey == "" {
		ctx.L().Error("encrypted key is empty")
		return "", res.FailRequest
	}
	aesKey, err := rsa_util.Decrypt(encryptedKey, privateKey)
	if err != nil {
		ctx.L().Error("encrypt error", zap.Error(err))
		return "", res.FailDefault
	}
	return aesKey, nil
}

func DecryptRequestBody(ctx *handler.Ctx, aesKey string) ([]byte, error) {
	sign := fiber.GetReqHeader[string](ctx, domains2.HeaderXRequestSignature)
	decrypt, err := aes_util.DecryptCiphertextAndTag(string(ctx.Request().Body()), sign, aesKey, RequestAAD(ctx))
	if err != nil {
		ctx.L().Error("decrypt error", zap.Error(err))
		return nil, res.FailDefault
	}
	return decrypt, nil
}

func RequestAAD(ctx *handler.Ctx) string {
	timestamp := fiber.GetReqHeader[int64](ctx, domains2.HeaderXRequestTimestamp)
	nonce := ctx.RequestID()
	params := map[string]any{
		domains2.HeaderXRequestID:        nonce,
		domains2.HeaderXRequestTimestamp: timestamp,
	}

	for k, v := range ctx.Queries() {
		params[k] = v
	}
	return encrypt.UriSort(params, func(key string) bool {
		return true
	})
}

func EncryptText(plainText string, aesKey string) (string, error) {
	result, err := aes_util.Encrypt(plainText, aesKey, "")
	if err != nil {
		return "", err
	}
	return result.Combined, nil
}

const SigData = "signData"

func SignMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		privateKey := ctx.PrivateKey
		if privateKey == "" {
			ctx.L().Error("privateKey is empty")
			return res.FailRequest
		}
		encryptedKey := fiber.GetReqHeader[string](ctx, domains2.HeaderXRequestEncryptedKey)
		if encryptedKey == "" {
			ctx.L().Error("encrypted key is empty")
			return res.FailRequest
		}
		aesKey, err := rsa_util.Decrypt(encryptedKey, privateKey)
		if err != nil {
			ctx.L().Error("encrypt error", zap.Error(err))
			return res.FailDefault
		}
		timestamp := fiber.GetReqHeader[int64](ctx, domains2.HeaderXRequestTimestamp)
		nonce := ctx.RequestID()
		params := map[string]any{
			domains2.HeaderXRequestID:        nonce,
			domains2.HeaderXRequestTimestamp: timestamp,
		}

		for k, v := range ctx.Queries() {
			params[k] = v
		}
		reqBody := ctx.Request().Body()

		params[SigData] = reqBody
		add := encrypt.UriSort(params, func(key string) bool {
			return true
		})
		sign := fiber.GetReqHeader[string](ctx, domains2.HeaderXRequestSignature)
		_, err = aes_util.DecryptCiphertextAndTag("", sign, aesKey, add)
		if err != nil {
			ctx.L().Error("decrypt error", zap.Error(err))
			return res.FailDefault
		}
		return ctx.Next()
	})
}

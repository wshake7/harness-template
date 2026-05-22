package logic

import (
	"admin/internal/appsvc"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"errors"

	"go.uber.org/zap"
)

type EncryptHandler struct {
	Cache appsvc.RedisCache
}

func NewEncryptHandler(cache appsvc.RedisCache) *EncryptHandler {
	return &EncryptHandler{Cache: cache}
}

type ResPublicKey struct {
	PublicKey string `json:"publicKey"`
}

// @Summary 获取加密公钥
// @Remark 获取当前全局 RSA 公钥，不存在时自动生成并缓存
// @Tags Encrypt
// @Produce json
// @Success 200 {object} res.Response{data=ResPublicKey} "成功"
// @Router /api/encrypt/public/key [get]
func (h *EncryptHandler) PublicKey(ctx *handler.Ctx) (*ResPublicKey, error) {
	publicKey, _, err := h.Cache.GetEncryptKeyPair(ctx)
	if err == nil {
		return &ResPublicKey{PublicKey: publicKey}, nil
	}
	if !errors.Is(err, appsvc.ErrCacheMiss) {
		ctx.L().Error("获取全局Key错误", zap.Error(err))
		return nil, res.FailDefault
	}

	publicKey, _, err = appsvc.GenerateAndCacheKeyPair(ctx, h.Cache)
	if err != nil {
		ctx.L().Error("生成或保存rsaKey错误", zap.Error(err))
		return nil, res.FailDefault
	}
	return &ResPublicKey{PublicKey: publicKey}, nil
}

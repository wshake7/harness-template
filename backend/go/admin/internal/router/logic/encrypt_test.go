package logic

import (
	"admin/internal/fiberc/res"
	"admin/internal/mock"
	"admin/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEncryptHandler_PublicKey_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCache := mock.NewMockRedisCache(ctrl)
	h := NewEncryptHandler(mockCache)

	mockCache.EXPECT().GetEncryptKeyPair(gomock.Any()).Return("pubkey123", "privkey456", nil)

	ctx := newTestCtx(t)
	result, err := h.PublicKey(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "pubkey123", result.PublicKey)
}

func TestEncryptHandler_PublicKey_CacheMiss(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCache := mock.NewMockRedisCache(ctrl)
	h := NewEncryptHandler(mockCache)

	mockCache.EXPECT().GetEncryptKeyPair(gomock.Any()).Return("", "", service.ErrCacheMiss)
	mockCache.EXPECT().SetEncryptKeyPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	ctx := newTestCtx(t)
	result, err := h.PublicKey(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.PublicKey)
}

func TestEncryptHandler_PublicKey_RedisError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCache := mock.NewMockRedisCache(ctrl)
	h := NewEncryptHandler(mockCache)

	mockCache.EXPECT().GetEncryptKeyPair(gomock.Any()).Return("", "", assert.AnError)

	ctx := newTestCtx(t)
	_, err := h.PublicKey(ctx)
	assert.Equal(t, res.FailDefault, err)
}

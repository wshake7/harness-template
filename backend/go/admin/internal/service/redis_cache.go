package service

import (
	"admin/internal/services/redisc"
	"context"
	"errors"

	"github.com/redis/rueidis"
	"go-common/utils/encrypt/rsa_util"
)

//go:generate mockgen -source=redis_cache.go -destination=../mock/mock_redis_cache.go -package=mock -typed

type RedisCache interface {
	GetEncryptKeyPair(ctx context.Context) (publicKey, privateKey string, err error)
	SetEncryptKeyPair(ctx context.Context, publicKey, privateKey string) error
}

type redisCacheImpl struct{}

func NewRedisCache() RedisCache {
	return &redisCacheImpl{}
}

var ErrCacheMiss = errors.New("cache miss")

func (r *redisCacheImpl) GetEncryptKeyPair(ctx context.Context) (string, string, error) {
	var keyPair redisc.DtoKeyPair
	err := redisc.Client.GetJson(ctx, redisc.KeyGlobalEncryptPublicKey, &keyPair)
	if err != nil {
		if errors.Is(err, rueidis.Nil) {
			return "", "", ErrCacheMiss
		}
		return "", "", err
	}
	return keyPair.PublicKey, keyPair.PrivateKey, nil
}

func (r *redisCacheImpl) SetEncryptKeyPair(ctx context.Context, publicKey, privateKey string) error {
	keyPair := redisc.DtoKeyPair{PublicKey: publicKey, PrivateKey: privateKey}
	return redisc.Client.Do(ctx, redisc.Client.B().Set().Key(redisc.KeyGlobalEncryptPublicKey).Value(rueidis.JSON(keyPair)).Build()).Error()
}

// GenerateAndCacheKeyPair is a helper used by EncryptHandler on cache miss.
func GenerateAndCacheKeyPair(ctx context.Context, cache RedisCache) (string, string, error) {
	privateKey, publicKey, err := rsa_util.GenerateKeyPair()
	if err != nil {
		return "", "", err
	}
	if err = cache.SetEncryptKeyPair(ctx, publicKey, privateKey); err != nil {
		return "", "", err
	}
	return publicKey, privateKey, nil
}

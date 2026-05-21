package services

import (
	"admin/internal/config"
	"admin/internal/services/redisc"
	"context"
	"sa-token/rueidis"

	"github.com/click33/sa-token-go/core"
	"github.com/click33/sa-token-go/core/pool"
	"github.com/click33/sa-token-go/stputil"
)

type Auth struct {
	authConf config.AuthConfig
	client   *redisc.RedisClient
}

func NewAuth(conf config.AuthConfig, client *redisc.RedisClient) *Auth {
	return &Auth{
		authConf: conf,
		client:   client,
	}
}

func (a *Auth) Start(ctx context.Context) error {
	manager := core.NewManager(rueidis.NewStorageFromClient(a.client), &core.Config{
		TokenName:              a.authConf.TokenName,
		Timeout:                a.authConf.Timeout,
		MaxRefresh:             a.authConf.MaxRefresh,
		RenewInterval:          a.authConf.RenewInterval,
		ActiveTimeout:          a.authConf.ActiveTimeout,
		IsConcurrent:           a.authConf.IsConcurrent,
		IsShare:                a.authConf.IsShare,
		MaxLoginCount:          a.authConf.MaxLoginCount,
		IsReadBody:             a.authConf.IsReadBody,
		IsReadHeader:           a.authConf.IsReadHeader,
		IsReadCookie:           a.authConf.IsReadCookie,
		TokenStyle:             a.authConf.TokenStyle,
		DataRefreshPeriod:      a.authConf.DataRefreshPeriod,
		TokenSessionCheckLogin: a.authConf.TokenSessionCheckLogin,
		AutoRenew:              a.authConf.AutoRenew,
		JwtSecretKey:           a.authConf.JwtSecretKey,
		IsLog:                  a.authConf.IsLog,
		IsPrintBanner:          a.authConf.IsPrintBanner,
		KeyPrefix:              a.authConf.KeyPrefix,
		CookieConfig: &core.CookieConfig{
			Domain:   a.authConf.CookieConfig.Domain,
			Path:     a.authConf.CookieConfig.Path,
			Secure:   a.authConf.CookieConfig.Secure,
			HttpOnly: a.authConf.CookieConfig.HttpOnly,
			SameSite: a.authConf.CookieConfig.SameSite,
			MaxAge:   a.authConf.CookieConfig.MaxAge,
		},
		RenewPoolConfig: &pool.RenewPoolConfig{
			MinSize:             a.authConf.RenewPoolConfig.MinSize,
			MaxSize:             a.authConf.RenewPoolConfig.MaxSize,
			ScaleUpRate:         a.authConf.RenewPoolConfig.ScaleUpRate,
			ScaleDownRate:       a.authConf.RenewPoolConfig.ScaleDownRate,
			CheckInterval:       a.authConf.RenewPoolConfig.CheckInterval,
			Expiry:              a.authConf.RenewPoolConfig.Expiry,
			PrintStatusInterval: a.authConf.RenewPoolConfig.PrintStatusInterval,
			PreAlloc:            a.authConf.RenewPoolConfig.PreAlloc,
			NonBlocking:         a.authConf.RenewPoolConfig.NonBlocking,
		},
	})
	stputil.SetManager(manager)
	return nil
}

func (a *Auth) String() string {
	return "auth"
}

func (a *Auth) State(ctx context.Context) (string, error) {
	return "HEALTHY", nil
}

func (a *Auth) Terminate(ctx context.Context) error {
	stputil.CloseManager()
	return nil
}

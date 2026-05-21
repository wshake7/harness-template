package services

import (
	"admin/internal/config"
	"admin/internal/services/orm"
	"admin/internal/services/redisc"
)

func New(conf *config.Config) {
	ormService := NewOrm(conf.Orm)
	redisService := NewRedis(conf.Redis)
	conf.Fiber.Services = append(conf.Fiber.Services, NewHttpc(), ormService, redisService, NewAuth(conf.Auth, redisc.Client), NewGeo(), NewAsynq(conf.Redis), NewTemporal(conf.Temporal), NewCasbin(orm.Client.DB))
}

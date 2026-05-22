package services

import (
	"admin/internal/config"
	"admin/internal/services/orm"
	"admin/internal/services/redisc"
	"admin/internal/services/temporaljob"
)

func New(conf *config.Config, registerWorkflows func(temporaljob.WorkerRegistry)) {
	ormService := NewOrm(conf.Orm)
	redisService := NewRedis(conf.Redis)
	conf.Fiber.Services = append(conf.Fiber.Services, NewHttpc(), ormService, redisService, NewAuth(conf.Auth, redisc.Client), NewGeo(), NewAsynq(conf.Redis), NewMilvus(conf.Milvus), NewTemporal(conf.Temporal, registerWorkflows), NewCasbin(orm.Client.DB))
}

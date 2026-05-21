package main

import (
	"admin/internal/config"
	"admin/internal/fiberc"
	"admin/internal/router"
	"admin/internal/services"
	"flag"
	"go-common/log"
	"go-common/viperc"
)

var configFile = flag.String("f", "./etc/config.yaml", "the config file")

// @title Admin API
// @version 1.0
// @description Admin 后端服务 API 文档。
// @description 所有接口均返回 HTTP 200，通过响应体中的 code 区分业务状态：
// @description | Code | Msg | 说明 |
// @description | :--- | :--- | :--- |
// @description | 1 | success | 成功 |
// @description | 2 | 服务繁忙 | 通用失败 |
// @description | 3 | 请求超时 | 请求过期 |
// @description | 4 | 请求重放 | Nonce 校验失败 |
// @description | 5 | 请求错误 | 客户端请求错误 |
// @description | 100 | - | 登录/权限相关失败 |
// @description | 200 | - | 授权相关失败|
// @BasePath /
func main() {
	flag.Parse()
	conf := config.Conf
	_, err := viperc.ParseFile(*configFile, conf)
	if err != nil {
		panic(err)
	}
	log.Init(log.DevConfig())
	services.New(conf)
	app := fiberc.NewFiber(conf)
	group := app.Group(conf.RestPrefix)
	r := router.Router{}
	r.RegisterRouters(group)
	app.Start()
}

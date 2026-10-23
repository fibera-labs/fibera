package main

import (
	"fibera/internal/app"
)

// @title						Fibera API
// @version						v0.1.0
// @description					基于 Go 与 Fiber v3 的后端服务脚手架。
// @BasePath					/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					除登录接口外，其余接口均需在请求头携带 Bearer <token>
func main() {
	app.Start()
}

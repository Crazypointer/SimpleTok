package main

import (
	"github.com/RaymondCode/simple-demo/core"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 读取配置文件
	core.InitConfig()
	// 初始化数据库
	global.DB = core.InitGorm()

	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

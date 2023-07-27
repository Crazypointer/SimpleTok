package main

import (
	"log"

	"github.com/Crazypointer/simple-tok/core"
	"github.com/Crazypointer/simple-tok/global"
	"github.com/Crazypointer/simple-tok/models"
	"github.com/Crazypointer/simple-tok/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 读取配置文件
	core.InitConfig()
	// 初始化数据库
	global.DB = core.InitGorm()
	// 生成数据库表结构
	// Makemigrations()

	go service.RunMessageServer()
	r := gin.Default()
	initRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 生成数据表
func Makemigrations() {
	err := global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.User{},
		&models.Video{},
		&models.Message{},
		&models.Comment{},
		&models.MessagePushEvent{},
		&models.MessageSendEvent{},
	)
	if err != nil {
		log.Fatalln("生成数据库表结构失败")
		return
	}
	log.Fatalln("生成数据库表结构成功！")
}

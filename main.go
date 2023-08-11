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

	// 初始化redis
	global.Redis = core.InitRedis()

	// 生成数据库表结构 第一次运行项目 或 修改表结构 时解开注释
	// Makemigrations()

	go service.RunMessageServer()
	r := gin.Default()

	initRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 生成数据表
func Makemigrations() {
	err := global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.User{},               // 用户表
		&models.Video{},              // 视频表
		&models.UserFavoriteVideo{},  // 用户点赞视频表
		&models.Message{},            // 消息表
		&models.Comment{},            // 评论表
		&models.UserFollowRelation{}, // 用户关注关系表
		&models.MessagePushEvent{},   // 消息推送事件表
		&models.MessageSendEvent{},   // 消息发送事件表
	)
	if err != nil {
		log.Fatalln("生成数据库表结构失败")
		return
	}
	log.Fatalln("生成数据库表结构成功！")
}

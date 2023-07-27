package core

import (
	"fmt"
	"log"
	"time"

	"github.com/RaymondCode/simple-tok/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm() *gorm.DB {
	if global.Config.Mysql.Host == "" {
		log.Println("未配grom，取消mysql链接！")
		return nil
	}
	dsn := global.Config.Mysql.Dsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(fmt.Sprintf("[%s] mysql链接失败！", dsn))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               //最大空闲链接数
	sqlDB.SetMaxOpenConns(100)              //最多可容纳
	sqlDB.SetConnMaxLifetime(time.Hour * 4) //链接最大复用时间，不能超过mysql的wait_timeout
	log.Println("mysql链接成功！")
	return db
}

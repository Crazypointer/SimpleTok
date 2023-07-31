package core

import (
	"context"
	"strconv"
	"time"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// InitRedis 初始化redis
func InitRedis() *redis.Client {
	addr := global.Config.Redis.Host + ":" + strconv.Itoa(global.Config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.DB,       // use default DB
		PoolSize: global.Config.Redis.PoolSize, // 连接池大小
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		logrus.Errorf("redis 链接失败: %s, %s", addr, err)
		return nil
	}
	logrus.Info("redis on" + addr + " 链接成功!")
	return rdb
}

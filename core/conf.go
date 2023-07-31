package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/Crazypointer/simple-tok/config"
	"github.com/Crazypointer/simple-tok/global"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const ConfigFile = "settings.yaml"

func InitConfig() {
	c := &config.Config{}
	yamlConfig, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yamlConf erro: %s", err))
	}
	err = yaml.Unmarshal(yamlConfig, c)
	if err != nil {
		panic(fmt.Errorf("unmarshal yamlConf erro: %s", err))
	}
	log.Println("config init success!")
	global.Config = c
}

// InitRedis 初始化redis
func InitRedis() *redis.Client {
	addr := global.Config.Redis.IP + ":" + strconv.Itoa(global.Config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: global.Config.Redis.Password, // no password set
		DB:       global.Config.Redis.DB,       // use default DB
		PoolSize: global.Config.Redis.PoolSize, // 连接池大小
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		logrus.Errorf("redis 链接失败: %s", addr)
	}
	logrus.Info("redis 链接成功: ", addr)
	return rdb
}

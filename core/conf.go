package core

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/RaymondCode/simple-tok/config"
	"github.com/RaymondCode/simple-tok/global"
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

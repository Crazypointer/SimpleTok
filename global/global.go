package global

import (
	"github.com/RaymondCode/simple-demo/config"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
)

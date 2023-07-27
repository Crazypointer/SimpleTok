package global

import (
	"github.com/RaymondCode/simple-tok/config"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
)

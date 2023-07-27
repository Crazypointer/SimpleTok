package global

import (
	"github.com/Crazypointer/simple-tok/config"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
)

package globle

import (
	"IM/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 全局变量
var (
	Config   config.Config
	DataBase *gorm.DB
	Logger   *logrus.Logger
)

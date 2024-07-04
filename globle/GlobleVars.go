package globle

import (
	"IM/config"
	"IM/db/query"
	"github.com/sirupsen/logrus"
)

// 全局变量
var (
	Config config.Config
	//DataBase *gorm.DB
	Logger *logrus.Logger
	Query  *query.Query
)

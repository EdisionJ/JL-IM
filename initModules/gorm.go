package initModules

import (
	"IM/db/query"
	"IM/globle"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func initGorm() {
	if viper.GetString("mysql.host") == "" {
		err := errors.New("Host Must Be Provided!")
		log.Fatalf(fmt.Sprintf("Config Load ERROR:%v", err))
	}

	// new logger
	var l logger2.LogLevel
	if viper.GetString("mysql.log_level") == "info" {
		l = logger2.Info
	} else {
		l = logger2.Error
	}
	logger := logger2.New(
		log.New(os.Stdout, "[DB]", log.LstdFlags),
		logger2.Config{LogLevel: l,
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		})

	// build connect
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/im?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("mysql.user"), viper.GetString("mysql.passwd"), viper.GetString("mysql.host"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		log.Fatalf(fmt.Sprintf("DataBase Open ERROR : %v", err))
	}
	globle.Db = query.Use(db)
	log.Println("DataBase Init Success.")
}

package utils

import (
	"IM/config"
	"IM/globle"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// 配置初始化
func InitSystem(ConfigPath string) {
	conf := getConfig(ConfigPath)
	initGorm(conf)
	initLogger(conf)
}

// 根据配置文件内容获取配置
func getConfig(ConfigPath string) *config.Config {
	yamlFile, err := os.ReadFile(ConfigPath)
	if err != nil {
		panic(fmt.Errorf("Config File Read ERROR:%v", err))
	}
	err = yaml.Unmarshal(yamlFile, &globle.Config)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Config Load ERROR:%v", err))
	}
	log.Println("Config Load Success.")
	return &globle.Config
}

// 初始化Gorm
func initGorm(c *config.Config) {
	if c.Mysql.Host == "" {
		err := errors.New("Host Must Be Provided!")
		log.Fatalf(fmt.Sprintf("Config Load ERROR:%v", err))
	}

	// new logger
	var l logger2.LogLevel
	if c.Mysql.LogLevel == "info" {
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
	dsn := c.Mysql.Dsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		panic(fmt.Sprintf("DataBase Open ERROR : %v", err))
	}
	globle.DataBase = db
	log.Println("DataBase Connect Success.")
}

// 初始化系统Logger
func initLogger(c *config.Config) {
	globle.Logger = InitLogger(c.Logger)
	log.Println("Logger Init Success:", globle.Logger)
}

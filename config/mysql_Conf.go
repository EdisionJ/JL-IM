package config

import (
	"fmt"
)

// Mysql配置
type DBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
	DB       string `yaml:"database"`
	LogLevel string `yaml:"log_level"`
}

// 获取用于Gorm连接的Dsn
func (db DBConf) Dsn() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		db.User, db.Passwd, db.Host, db.Port, db.DB)
}

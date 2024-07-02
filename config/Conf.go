package config

// 系统配置结构体
type Config struct {
	System SysConf `yaml:"system"`
	Mysql  DBConf  `yaml:"mysql"`
	Logger LogConf `yaml:"logger"`
}

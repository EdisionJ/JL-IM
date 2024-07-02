package test

import (
	"IM.com/config"
	"IM.com/utils"
	"testing"
)

func TestLogger(t *testing.T) {
	c := config.LogConf{
		Level:        "info",
		Prefix:       "[Demo]",
		Director:     "log",
		ShowLine:     true,
		LogInConsole: true,
	}
	logger := utils.InitLogger(c)
	logger.Warning("123")
}

package main

import (
	"IM/routers"
	_ "IM/service/eventProcessor"
	"github.com/spf13/viper"
)

func main() {
	r := routers.Router()
	r.Run(viper.GetString("system.host"))
}

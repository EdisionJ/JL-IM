package main

import (
	"IM/globle"
	"IM/routers"
	"IM/utils"
	"fmt"
)

func main() {
	ConfigPath := "config.yaml"
	utils.InitSystem(ConfigPath)
	config := globle.Config
	r := routers.Router()
	r.Run(fmt.Sprintf("%v:%v", config.System.Host, config.System.Port))
}

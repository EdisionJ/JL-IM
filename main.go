package main

import (
	"IM/globle"
	"IM/routers"
	_ "IM/service/listener"
	"IM/websocketSereve"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	go startWebsocketServe()
	r := routers.Router()
	r.Run(viper.GetString("system.host"))
}

func startWebsocketServe() {
	http.HandleFunc("/websocket", websocketSereve.Connect)
	err := http.ListenAndServe(viper.GetString("system.websocket"), nil)
	if err != nil {
		globle.Logger.Fatalln("websocket服务启动失败：", err)
	}
}

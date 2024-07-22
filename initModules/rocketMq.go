package initModules

import (
	"IM/globle"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/spf13/viper"
	"log"
)

func initRocketMq() {
	host := viper.GetString("rocketmq.host")
	group := viper.GetString("rocketmq.group")
	rlog.SetLogLevel("warn")
	globle.RocketProducer, _ = rocketmq.NewProducer(
		producer.WithNameServer([]string{host}),
		producer.WithRetry(3),
		producer.WithGroupName(group),
	)
	// 开始连接
	err := globle.RocketProducer.Start()
	if err != nil {
		log.Fatalf("start producer error: %s", err.Error())
		log.Fatalf("RocketProducer Init Fail! ")
	} else {
		log.Println("RocketProducer Init Success.")
	}
}

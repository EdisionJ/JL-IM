package initModules

import (
	"IM/globle"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

func initRedis() {
	globle.Rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.passwd"), // 没有密码，默认值
		DB:       viper.GetInt("redis.db"),        // 默认DB 0
	})
	if globle.Rdb != nil {
		globle.Rdb.FlushAll(context.Background()).Result()
		log.Println("Rdb Init Success.")
	} else {
		log.Fatalf("Rdb Init Fail! ")
	}
}

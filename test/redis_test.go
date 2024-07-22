package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

type person struct {
	Name string
	Age  int
}

func TestRedis(t *testing.T) {
	data := person{
		"li hua",
		15,
	}
	marshal, err := json.Marshal(data)

	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.175.136:6379",
		Password: "123456", // 没有密码，默认值
		DB:       0,        // 默认DB 0
	})
	ctx := context.Background()
	err = client.Set(ctx, "demo", marshal, 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "demo").Result()
	if err != nil {
		panic(err)
	}
	var p person
	json.Unmarshal([]byte(val), &p)
	fmt.Println("demo:", p)

}

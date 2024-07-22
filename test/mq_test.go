package test

import (
	"IM/globle"
	"IM/service/enum"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"testing"
	"time"
)

func TestMq(t *testing.T) {
	// 创建一个生产者的实例
	rlog.SetLogLevel("warn")
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.175.136:9876"}),
	)
	if err != nil {
		fmt.Printf("Failed to create producer, error: %v\n", err)
		return
	}

	// 启动生产者
	if err = p.Start(); err != nil {
		fmt.Printf("Failed to start producer, error: %v\n", err)
		return
	}
	defer p.Shutdown()

	// 创建一个上下文
	ctx := context.Background()

	// 构建消息体
	msg := &primitive.Message{
		Topic: enum.UserLogInOrOut,
		Body:  []byte("Hello RocketMQ"),
	}

	// 发送消息
	result, err := p.SendSync(ctx, msg)
	if err != nil {
		fmt.Printf("Failed to send message, error: %v\n", err)
		return
	}
	fmt.Printf("\nMessage sent successfully, result: %s\n\n", result.String())
	func1()
}

func func1() {
	host := "192.168.175.136:9876"
	pc, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(enum.UserLogInOrOutGroup),
		consumer.WithNameServer([]string{host}),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)

	err := pc.Subscribe(enum.UserLogInOrOut, consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, messageExt := range ext {
				fmt.Println("*************", string(messageExt.Body), "*************")
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		globle.Logger.Fatal("UserLogin subscribe error: ", err)
	}
	err = pc.Start()
	if err != nil {
		globle.Logger.Fatal("UserLogin start error: ", err)
	}
	time.Sleep(5 * time.Second)
}

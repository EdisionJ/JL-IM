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

	host := "192.168.175.136:9876"
	pc, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(enum.UserLogInOrOutGroup),
		consumer.WithNameServer([]string{host}),
	)
	func1(pc)

	//ctx := context.Background()
	//i := 5
	//for i != 0 {
	//	msg := &primitive.Message{
	//		Topic: enum.UserLogInOrOut,
	//		Body:  []byte("msg:" + strconv.Itoa(i) + "   " + time.Now().String()),
	//	}
	//
	//	// 发送消息
	//	_, err := p.SendSync(ctx, msg)
	//	if err != nil {
	//		fmt.Printf("Failed to send message, error: %v\n", err)
	//		return
	//	}
	//	i--
	//	time.Sleep(time.Second * 1)
	//}
	fmt.Println("SEND OVER ! ")
	time.Sleep(time.Second * 2)
	defer pc.Shutdown()

}

func func1(pc rocketmq.PushConsumer) {

	pc.Subscribe(enum.UserSignUp, consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, e := range ext {

				fmt.Println(string(e.Body))

				fmt.Println("*************")
			}
			return consumer.ConsumeSuccess, nil
		})

	pc.Subscribe(enum.UserLogInOrOut, consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, e := range ext {

				fmt.Println(string(e.Body))

				fmt.Println("*************")
			}
			return consumer.ConsumeSuccess, nil
		})
	pc.Subscribe(enum.UserFriendReq, consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, e := range ext {

				fmt.Println(string(e.Body))

				fmt.Println("*************")
			}
			return consumer.ConsumeSuccess, nil
		})
	pc.Subscribe(enum.UserFriendAdd, consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, e := range ext {

				fmt.Println(string(e.Body))

				fmt.Println("*************")
			}
			return consumer.ConsumeSuccess, nil
		})
	err := pc.Start()
	if err != nil {
		globle.Logger.Fatal("UserLogin start error: ", err)
	}

}

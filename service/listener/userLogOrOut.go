package listener

import (
	"IM/db/model"
	"IM/globle"
	_ "IM/initModules"
	"IM/service/enum"
	"IM/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var UserQ = globle.Db.User

func init() {
	host := viper.GetString("rocketmq.host")
	pc, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(enum.UserLogInOrOutGroup),
		consumer.WithNameServer([]string{host}),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	err := pc.Subscribe(enum.UserLogInOrOut, consumer.MessageSelector{}, login)
	if err != nil {
		globle.Logger.Fatal("登录事件订阅失败: ", err)
	}
	err = pc.Start()
	if err != nil {
		globle.Logger.Fatal("登录服务启动失败: ", err)
	}

	go func() {
		// 监听信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan // 等待信号

		// 收到信号后，关闭消费者
		if err := pc.Shutdown(); err != nil {
			log.Printf("关闭消费者失败: %v", err)
		}
		os.Exit(0)
	}()
}

func login(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range ext {
		var userInfo globle.UserInfo
		err := json.Unmarshal(msg.Body, &userInfo)
		if err != nil {
			globle.Logger.Errorf("json.Unmarshal发生错误: %v", err)
			return consumer.ConsumeRetryLater, err
		}
		if msg.Flag == enum.LogOut {
			//设置离线状态
			_, err = UserQ.WithContext(ctx).
				Where(UserQ.ID.Eq(userInfo.ID)).
				Select(UserQ.LastOnlineTime, UserQ.IsOnline).
				Updates(model.User{LastOnlineTime: time.Now(), IsOnline: enum.UserStatusOffline})
		} else {
			keyId := fmt.Sprintf(enum.UserCacheByID, userInfo.ID)
			keyName := fmt.Sprintf(enum.UserCacheByName, userInfo.Name)
			err = utils.SetToCache(keyId, userInfo)
			err = utils.SetToCache(keyName, userInfo)
			if err != nil {
				globle.Logger.Warnf("添加缓存时发生错误: %v", err)
				return consumer.ConsumeRetryLater, err
			}
			//设置在线状态
			_, err = UserQ.WithContext(ctx).
				Where(UserQ.ID.Eq(userInfo.ID)).
				Update(UserQ.IsOnline, enum.UserStatusOnline)
		}
		if err != nil {
			globle.Logger.Error("在线状态切换失败： ", err)
			return consumer.ConsumeRetryLater, err
		}
	}
	return consumer.ConsumeSuccess, nil
}

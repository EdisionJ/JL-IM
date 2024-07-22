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
}

func login(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range ext {
		var userInfo globle.UserInfo
		err := json.Unmarshal(msg.Body, &userInfo)
		if err != nil {
			globle.Logger.Errorf("json.Unmarshal发生错误: %v", err)
			return consumer.Rollback, err
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
				return consumer.Rollback, err
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

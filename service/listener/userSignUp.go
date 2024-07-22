package listener

import (
	"IM/db/model"
	"IM/globle"
	"IM/service/RR"
	"IM/service/enum"
	"IM/utils"
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/spf13/viper"
)

func init() {
	host := viper.GetString("rocketmq.host")
	pc, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(enum.UserLogSignUpGroup),
		consumer.WithNameServer([]string{host}),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	err := pc.Subscribe(enum.UserLogSignUp, consumer.MessageSelector{}, signUp)
	if err != nil {
		globle.Logger.Fatal("订阅注册事件失败: ", err)
	}
	err = pc.Start()
	if err != nil {
		globle.Logger.Fatal("注册服务启动失败: ", err)
	}
}

func signUp(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range ext {
		var userInfo RR.UserSingUp
		err := json.Unmarshal(msg.Body, &userInfo)
		if err != nil {
			globle.Logger.Warnf("json.Unmarshal发生错误 ", err)
			return consumer.ConsumeRetryLater, err
		}
		var NewUser model.User
		NewUser.ID = utils.GenID()
		NewUser.Name = userInfo.Name
		NewUser.PassWd = utils.EncodeWithSHA256(userInfo.PassWD)
		NewUser.PhoneNumber = userInfo.PhoneNumber
		NewUser.Email = userInfo.Email
		err = UserQ.WithContext(ctx).
			Select(UserQ.ID, UserQ.Name, UserQ.PassWd, UserQ.PhoneNumber, UserQ.Email).
			Create(&NewUser)
		if err != nil {
			globle.Logger.Error("添加用户时发生错误！", err)
			return consumer.ConsumeRetryLater, err
		}
	}
	return consumer.ConsumeSuccess, nil
}

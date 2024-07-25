package listener

import (
	"IM/db/model"
	"IM/globle"
	"IM/service/enum"
	"IM/service/requestModels"
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
)

var FriendReqQ = globle.Db.FriendReq
var FriendQ = globle.Db.Friend
var GroupQ = globle.Db.Group

func init() {
	host := viper.GetString("rocketmq.host")
	pc, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(enum.UserFriendServiceGroup),
		consumer.WithNameServer([]string{host}),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)

	//好友请求
	err := pc.Subscribe(enum.UserFriendReq, consumer.MessageSelector{}, friendReq)
	if err != nil {
		globle.Logger.Fatal("请求好友事件订阅失败: ", err)
	}

	//添加好友
	err = pc.Subscribe(enum.UserFriendAdd, consumer.MessageSelector{}, friendAdd)
	if err != nil {
		globle.Logger.Fatal("请求好友事件订阅失败: ", err)
	}

	err = pc.Start()
	if err != nil {
		globle.Logger.Fatal("好友服务启动失败: ", err)
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

func friendReq(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range ext {
		var reqInfo requestModels.AddFriendInfo
		err := json.Unmarshal(msg.Body, &reqInfo)
		if err != nil {
			globle.Logger.Println("json.Unmarshal发生错误 ", err)
			return consumer.Rollback, err
		}
		var fr model.FriendReq
		fr.ID = reqInfo.Uid
		fr.FriendID = reqInfo.FriendId
		fr.Msg = reqInfo.ReqMsg

		//写入数据库
		err = FriendReqQ.WithContext(ctx).
			Select(FriendReqQ.ID, FriendReqQ.FriendID, FriendReqQ.Msg).
			Create(&fr)
		if err != nil {
			globle.Logger.Warnf("数据库错误: %v", err)
			return consumer.ConsumeRetryLater, err
		}

		//写入缓存
		key := fmt.Sprintf(enum.UserApplyCacheByUidAndFriendUid, fr.ID, fr.FriendID)
		err = utils.SetToCache(key, fr)
		key = fmt.Sprintf(enum.UserApplyCacheByFriendUid, fr.FriendID)
		err = utils.SetToCache(key, fr)
		if err != nil {
			globle.Logger.Warnf("设置缓存时遇到错误: %v", err)
			return consumer.ConsumeRetryLater, err
		}
	}
	return consumer.ConsumeSuccess, nil
}

func friendAdd(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range ext {
		var info requestModels.AddFriendInfo
		err := json.Unmarshal(msg.Body, &info)
		if err != nil {
			globle.Logger.Warnf("json.Unmarshal发生错误 ", err)
			return consumer.ConsumeRetryLater, err
		}

		switch info.Flag {
		//同意好友请求
		case enum.FriendReqAgree:
			//删除缓存
			key := fmt.Sprintf(enum.UserApplyCacheByUidAndFriendUid, info.FriendId, info.Uid)
			err := utils.RemoveCacheData(key)
			if err != nil {
				globle.Logger.Warnf("删除缓存时遇到错误: %v", err)
				return consumer.ConsumeRetryLater, err
			}
			//开始事务
			tx := globle.Db.Begin()
			//设置好友请求状态
			_, err = FriendReqQ.WithContext(ctx).
				Where(FriendReqQ.ID.Eq(info.FriendId), FriendReqQ.FriendID.Eq(info.Uid)).
				Update(FriendReqQ.IsAgree, enum.FriendReqAgree)
			if err != nil {
				globle.Logger.Error("数据库数据更新失败： ", err)
				err := tx.Rollback()
				if err != nil {
					globle.Logger.Error("事务回滚失败： ", err)
				}
				return consumer.ConsumeRetryLater, err
			}

			//获取用户信息
			var u1 = model.User{}
			key1 := fmt.Sprintf(enum.UserCacheByID, info.Uid)
			var u2 = model.User{}
			key2 := fmt.Sprintf(enum.UserCacheByID, info.FriendId)
			err = utils.Get(key1, &u1, func() (any, error) {
				return UserQ.WithContext(ctx).Where(UserQ.ID.Eq(info.Uid)).First()
			})
			err = utils.Get(key2, &u2, func() (any, error) {
				return UserQ.WithContext(ctx).Where(UserQ.ID.Eq(info.FriendId)).First()
			})
			if err != nil {
				globle.Logger.Error("查询数据出错： ", err)
				return consumer.ConsumeRetryLater, err
			}

			//新建聊天室
			roomID := utils.GenID()
			room := model.Group{
				ID:   roomID,
				Type: enum.RoomTypePrivate,
			}
			err = GroupQ.WithContext(ctx).
				Select(GroupQ.ID, GroupQ.Type).
				Create(&room)
			if err != nil {
				globle.Logger.Error("数据库插入数据失败： ", err)
				err := tx.Rollback()
				if err != nil {
					globle.Logger.Error("事务回滚失败： ", err)
				}
				return consumer.ConsumeRetryLater, err
			}
			//更新好友列表
			var friendShip1 model.Friend
			friendShip1.ID = info.Uid
			friendShip1.FriendID = info.FriendId
			friendShip1.NickName = u1.Name
			friendShip1.Room = roomID
			var friendShip2 model.Friend
			friendShip2.ID = info.FriendId
			friendShip2.FriendID = info.Uid
			friendShip2.NickName = u2.Name
			friendShip2.Room = roomID
			err = FriendQ.WithContext(ctx).
				Select(FriendQ.ID, FriendQ.FriendID, FriendQ.NickName, FriendQ.Room).
				Create(&friendShip1, &friendShip2)
			if err != nil {
				globle.Logger.Error("数据库插入数据失败： ", err)
				err := tx.Rollback()
				if err != nil {
					globle.Logger.Error("事务回滚失败： ", err)
				}
				return consumer.ConsumeRetryLater, err
			}

			if err != nil {
				globle.Logger.Error("数据库插入数据失败： ", err)
				err := tx.Rollback()
				if err != nil {
					globle.Logger.Error("事务回滚失败： ", err)
				}
				return consumer.ConsumeRetryLater, err
			}
			//提交事务
			tx.Commit()
			//
			key = fmt.Sprintf(enum.UserFriendCacheByUidAndFriendUid, info.Uid, info.FriendId)
			err = utils.SetToCache(key, &friendShip1)
			key = fmt.Sprintf(enum.UserFriendCacheByUidAndFriendUid, info.FriendId, info.Uid)
			err = utils.SetToCache(key, &friendShip2)
			if err != nil {
				globle.Logger.Warnln("设置好友缓存失败： ", err)
			}
			return consumer.ConsumeSuccess, nil
		//拒绝好友请求
		case enum.FriendReqRefuse:
			//删除缓存
			key := fmt.Sprintf(enum.UserApplyCacheByUidAndFriendUid, info.FriendId, info.Uid)
			err := utils.RemoveCacheData(key)
			if err != nil {
				globle.Logger.Warnf("删除缓存时遇到错误: %v", err)
				return consumer.ConsumeRetryLater, err
			}
			//设置好友请求状态
			_, err = FriendReqQ.WithContext(ctx).
				Where(FriendReqQ.ID.Eq(info.FriendId), FriendReqQ.FriendID.Eq(info.Uid)).
				Update(FriendReqQ.IsAgree, enum.FriendReqRefuse)
			if err != nil {
				globle.Logger.Error("数据库数据更新失败： ", err)
				return consumer.ConsumeRetryLater, err
			}
			return consumer.ConsumeSuccess, nil
		}
	}
	return consumer.ConsumeSuccess, nil
}

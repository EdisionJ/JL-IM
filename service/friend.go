package service

import (
	"IM/db/model"
	"IM/globle"
	"IM/service/RR"
	"IM/service/enum"
	"IM/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var FriendReqQ = globle.Db.FriendReq
var FriendQ = globle.Db.Friend

func FriendReq(c *gin.Context) {
	ctx := context.Background()
	var reqInfo RR.AddFriend
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		utils.RspWithMsg(c, http.StatusBadRequest, false, "非法请求！")
		return
	}

	from := reqInfo.Uid
	to := reqInfo.Uid

	//用户是否存在
	key := fmt.Sprintf(enum.UserCacheByID, to)
	u := model.User{}
	err = utils.Get(key, &u, func() (any, error) {
		return UserQ.WithContext(ctx).Where(UserQ.ID.Eq(to)).First()
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RspWithMsg(c, http.StatusBadRequest, false, "用户不存在！")
			return
		}
		globle.Logger.Error("查找用户时发生错误！", err)
		utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
		return
	}

	//是否已经是好友
	isFriend, err := IsFriend(from, to)
	if err != nil {
		globle.Logger.Error("查找好友关系时发生错误！", err)
		utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
		return
	}
	if isFriend {
		utils.RspWithMsg(c, http.StatusOK, true, "已经是好友了！")
		return
	}

	//不要重复发送请求
	key = fmt.Sprintf(enum.UserApplyCacheByUidAndFriendUid, from, to)
	fr := model.FriendReq{}
	err = utils.Get(key, &fr, func() (any, error) {
		return FriendReqQ.WithContext(ctx).Where(FriendReqQ.UID.Eq(from), FriendReqQ.FriendID.Eq(to)).First()
	})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			globle.Logger.Error("查找好友请求时发生错误！", err)
			utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
			return
		}
	} else {
		utils.RspWithMsg(c, http.StatusOK, false, "已发送过好友请求，请等待对方同意！")
		return
	}

	//复用该接口，用以处理对方发送的好友请求
	if reqInfo.Flag != enum.FriendReqNotYetAgreed {
		addFriend(ctx, c, reqInfo)
		utils.RspWithMsg(c, http.StatusOK, true, "操作成功")
		return
	}

	//如果对方也发送了好友请求，那么直接加为好友
	key = fmt.Sprintf(enum.UserApplyCacheByUidAndFriendUid, to, from)
	fr = model.FriendReq{}
	err = utils.Get(key, &fr, func() (any, error) {
		return FriendReqQ.WithContext(ctx).Where(FriendReqQ.UID.Eq(to), FriendReqQ.FriendID.Eq(from)).First()
	})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			globle.Logger.Error("查找好友请求时发生错误！", err)
			utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
			return
		}
	} else {
		//对方发送了好友请求
		reqInfo.Flag = enum.FriendReqAgree
		addFriend(ctx, c, reqInfo)
		utils.RspWithMsg(c, http.StatusOK, true, "好友添加成功")
		return
	}

	//没有查到记录，那么发送好友请求
	reqBytes, err := json.Marshal(reqInfo)
	if err != nil {
		globle.Logger.Error("json.Marshal发生错误！", err)
		utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
		return
	}
	msg := &primitive.Message{
		Topic: enum.UserFriendReq,
		Body:  reqBytes,
	}

	r, err := globle.RocketProducer.SendSync(ctx, msg)
	if err != nil {
		globle.Logger.Error("好友请求消息发送失败！", err, r)
		utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
		return
	}
	utils.RspWithMsg(c, http.StatusOK, true, "好友请求发送成功")
	return
}

func IsFriend(uid1, uid2 int64) (bool, error) {
	key := fmt.Sprintf(enum.UserFriendCacheByUidAndFriendUid, uid1, uid2)
	fr := model.Friend{}
	err := utils.Get(key, &fr, func() (any, error) {
		return FriendQ.WithContext(context.Background()).Where(FriendQ.ID.Eq(uid1), FriendQ.FriendID.Eq(uid2), FriendQ.Status.Eq(enum.FriendShipNormal)).First()
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func addFriend(ctx context.Context, c *gin.Context, req RR.AddFriend) {
	reqByte, err := json.Marshal(req)
	if err != nil {
		globle.Logger.Error("json.Marshal发生错误！", err)
		utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
		return
	}
	msg := &primitive.Message{
		Topic: enum.UserFriendAdd,
		Body:  reqByte,
	}

	r, err := globle.RocketProducer.SendSync(ctx, msg)
	if err != nil {
		globle.Logger.Error("添加好友消息发送失败！", err, r)
		utils.RspWithMsg(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
		return
	}
}

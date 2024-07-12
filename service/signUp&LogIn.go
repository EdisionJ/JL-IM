package service

import (
	"IM/db/model"
	"IM/globle"
	"IM/globle/enum/event"
	"IM/service/RR"
	"IM/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// @Summary 注册用户
// @Param user_info body RR.UserSingUp true "用户注册所需信息，email, phone_number、passwd、re_passwd是必须提供的"
// @Accept json
// @Produce json
// @Success 200 {object} utils.RepData
// @Router /user/signUp [post]
func UserSignUp(c *gin.Context) {
	ctx := context.Background()
	UserQ := globle.Db.User
	var userInfo RR.UserSingUp
	err := c.ShouldBindJSON(&userInfo)
	isValid, _ := govalidator.ValidateStruct(userInfo)
	switch {
	case !isValid:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "请输入正确格式的Email或电话！")
		return
	case err != nil:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "非法请求！")
		return
	case userInfo.PassWD != userInfo.RePassWD:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "两次输入的密码不一致！")
		return
	case userInfo.Email == "" || userInfo.PhoneNumber == "":
		utils.RspWithMsg(c, http.StatusBadRequest, false, "Email和电话不能为空！")
		return
	default:
		userIsExist, err := UserQ.WithContext(ctx).Where(UserQ.PhoneNumber.Eq(userInfo.PhoneNumber)).Or(UserQ.Email.Eq(userInfo.Email)).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			globle.Logger.Error("查找用户时发生错误！", err)
			utils.RspWithMsg(c, http.StatusBadRequest, false, "系统错误，请稍后再试！")
			return
		}
		if userIsExist != nil {
			utils.RspWithMsg(c, http.StatusOK, false, "当前电话号码或Email已被注册！")
			return
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
			utils.RspWithMsg(c, http.StatusBadRequest, false, "由于系统错误注册失败！")
			return
		} else {
			utils.RspWithMsg(c, http.StatusOK, true, "注册成功！")
			return
		}
	}

}

// @Summary 用户登录
// @Param user_info body RR.UserLogIn true "用户登录所需信息，{name/phone_number/email}、passwd是必须提供的"
// @Accept json
// @Produce json
// @Success 200 {object} utils.RepData
// @Router /user/login [post]
func UserLogIn(c *gin.Context) {
	ctx := context.Background()
	UserQ := globle.Db.User
	var loginer RR.UserLogIn
	err := c.ShouldBindJSON(&loginer)
	isValid, _ := govalidator.ValidateStruct(loginer)
	switch {
	case err != nil:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "非法请求！")
		return
	case !isValid:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "请输入正确格式的Email或电话！")
		return
	case loginer.Email == "" && loginer.PhoneNumber == "":
		utils.RspWithMsg(c, http.StatusBadRequest, false, "请输入电话或Email以登录！")
		return
	default:
		var user *model.User
		var err error
		if loginer.Email != "" {
			user, err = UserQ.WithContext(ctx).
				Where(UserQ.Email.Eq(loginer.Email)).First()
		} else {
			user, err = UserQ.WithContext(ctx).
				Where(UserQ.PhoneNumber.Eq(loginer.PhoneNumber)).First()
		}
		if err != nil {
			if err.Error() == "record not found" {
				utils.RspWithMsg(c, http.StatusOK, false, "用户未注册！")
				return
			}
			globle.Logger.Error("查找用户时发生错误！", err)
			utils.RspWithMsg(c, http.StatusBadRequest, false, "系统错误，请稍后再试！")
			return
		}
		if user == nil {
			utils.RspWithMsg(c, http.StatusOK, false, "用户不存在！")
			return
		}
		if utils.EncodeWithSHA256(loginer.PassWD) != user.PassWd {
			utils.RspWithMsg(c, http.StatusOK, false, "密码错误！")
			return
		} else {
			token, err := utils.GenToken(user.ID)
			if err != nil {
				globle.Logger.Error("生成Token时发生错误！", err)
				utils.RspWithMsg(c, http.StatusBadRequest, false, "系统错误，请稍后再试！")
				return
			}
			data := RR.UserLogOk{
				ID:           user.ID,
				Name:         user.Name,
				Token:        token,
				SelfDescribe: user.SelfDescribe,
				PhoneNumber:  user.PhoneNumber,
				Email:        user.Email,
			}
			userByte, err := json.Marshal(data)
			if err != nil {
				globle.Logger.Error("json 序列化失败！", err)
				utils.RspWithMsg(c, http.StatusBadRequest, false, "系统错误，请稍后再试！")
				return
			}
			msg := &primitive.Message{
				Topic: event.UserLogIn,
				Body:  userByte,
			}
			r, err := globle.RocketProducer.SendSync(ctx, msg)
			if err != nil {
				globle.Logger.Error("登录消息发送失败！", err, r)
				utils.RspWithMsg(c, http.StatusBadRequest, false, "系统错误，请稍后再试！")
				return
			}
			utils.RspWithData(c, http.StatusOK, true, data)
			return
		}
	}
}

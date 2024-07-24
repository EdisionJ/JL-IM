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
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var UserQ = globle.Db.User

func UserSignUp(c *gin.Context) {
	ctx := context.Background()
	var userInfo RR.UserSingUp
	err := c.ShouldBindJSON(&userInfo)
	isValid, _ := govalidator.ValidateStruct(userInfo)
	switch {
	case !isValid:
		utils.DefaultRsp(c, http.StatusBadRequest, false, "请输入正确格式的Email或电话！")
		return
	case err != nil:
		utils.DefaultRsp(c, http.StatusBadRequest, false, "非法请求！")
		return
	case userInfo.PassWD != userInfo.RePassWD:
		utils.DefaultRsp(c, http.StatusBadRequest, false, "两次输入的密码不一致！")
		return
	case userInfo.Email == "" || userInfo.PhoneNumber == "":
		utils.DefaultRsp(c, http.StatusBadRequest, false, "Email和电话不能为空！")
		return
	default:
		//查看Email或者电话号码是否已注册
		userIsExist, err := UserQ.WithContext(ctx).Where(UserQ.PhoneNumber.Eq(userInfo.PhoneNumber)).Or(UserQ.Email.Eq(userInfo.Email)).Count()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			globle.Logger.Error("查找用户时发生错误！", err)
			utils.DefaultRsp(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
			return
		}
		//能找到用户
		if userIsExist != 0 {
			utils.DefaultRsp(c, http.StatusOK, false, "当前电话号码或Email已被注册！")
			return
		}
		//发送用户注册消息
		infoBytes, err := json.Marshal(userInfo)
		if err != nil {
			globle.Logger.Error("json.Marshal发生错误！", err)
			utils.DefaultRsp(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
			return
		}
		msg := &primitive.Message{
			Topic: enum.UserSignUp,
			Body:  infoBytes,
		}
		r, err := globle.RocketProducer.SendSync(ctx, msg)
		if err != nil {
			globle.Logger.Error("注册消息发送失败！", err, r)
			utils.DefaultRsp(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
			return
		}

		utils.DefaultRsp(c, http.StatusOK, true, "注册成功！")
		return
	}

}

func UserLogIn(c *gin.Context) {
	ctx := context.Background()
	var loginer RR.UserLogIn
	err := c.ShouldBindJSON(&loginer)
	isValid, _ := govalidator.ValidateStruct(loginer)
	switch {
	case err != nil:
		utils.DefaultRsp(c, http.StatusBadRequest, false, "非法请求！")
		return
	case !isValid:
		utils.DefaultRsp(c, http.StatusBadRequest, false, "请输入正确格式的Email或电话！")
		return
	case loginer.Email == "" && loginer.PhoneNumber == "":
		utils.DefaultRsp(c, http.StatusBadRequest, false, "请输入电话或Email以登录！")
		return
	default:
		var user *model.User
		var err error
		//查看Email或者电话号码是否已注册
		if loginer.Email != "" {
			user, err = UserQ.WithContext(ctx).
				Where(UserQ.Email.Eq(loginer.Email)).First()
		} else {
			user, err = UserQ.WithContext(ctx).
				Where(UserQ.PhoneNumber.Eq(loginer.PhoneNumber)).First()
		}
		if err != nil {
			//找不到用户信息，未注册
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.DefaultRsp(c, http.StatusOK, false, "用户不存在！")
				return
			}
			globle.Logger.Error("查找用户时发生错误！", err)
			utils.DefaultRsp(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
			return
		}
		//密码验证
		if loginer.PassWD != user.PassWd {
			utils.DefaultRsp(c, http.StatusOK, false, "密码错误！")
			return
		} else {
			token, err := utils.GenToken(user.ID)
			c.Writer.Header().Set("Token", token)
			if err != nil {
				globle.Logger.Error("生成Token时发生错误！", err)
				utils.DefaultRsp(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
				return
			}
			data := globle.UserInfo{
				ID:           user.ID,
				Name:         user.Name,
				SelfDescribe: user.SelfDescribe,
				PhoneNumber:  user.PhoneNumber,
				Email:        user.Email,
			}

			//发送登陆消息
			userBytes, err := json.Marshal(data)
			if err != nil {
				globle.Logger.Error("json.Marshal发生错误！", err)
				utils.DefaultRsp(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
				return
			}
			msg := &primitive.Message{
				Topic: enum.UserLogInOrOut,
				Body:  userBytes,
				Flag:  enum.LogIn,
			}
			r, err := globle.RocketProducer.SendSync(ctx, msg)
			if err != nil {
				globle.Logger.Error("登录消息发送失败！", err, r)
				utils.DefaultRsp(c, http.StatusInternalServerError, false, "系统错误，请稍后再试！")
				return
			}
			utils.RspWithData(c, http.StatusOK, true, "登录成功！", data)
			return
		}
	}
}

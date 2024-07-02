package service

import (
	"IM/db/model"
	"IM/db/query"
	"IM/globle"
	"IM/utils"
	"IM/utils/RR"
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 注册用户
// @Param user_info body RR.UserSingUp true "用户注册所需信息，name、phone_number、passwd、re_passwd是必须提供的"
// @Accept json
// @Produce json
// @Success 200 {object} utils.RepData
// @Router /user/signUp [post]
func UserSignUp(c *gin.Context) {
	ctx := context.Background()
	var newUser RR.UserSingUp
	err := c.ShouldBindJSON(&newUser)
	isValid, _ := govalidator.ValidateStruct(newUser)
	switch {
	case !isValid:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "请输入正确格式的Email或电话！")
		return
	case err != nil:
		globle.Logger.Error("ShouldBindJSON发生错误！", err)
		utils.RspWithMsg(c, http.StatusBadRequest, false, "非法请求！")
		return
	case newUser.PassWD != newUser.RePassWD:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "两次输入的密码不一致！")
		return
	case newUser.Name == "" || newUser.PhoneNumber == "":
		utils.RspWithMsg(c, http.StatusBadRequest, false, "姓名或电话不能为空！")
		return
	default:
		user, err := UserQ.WithContext(ctx).Where(query.User.PhoneNumber.Eq(newUser.PhoneNumber)).Or(query.User.Name.Eq(newUser.Name)).First()
		if err != nil && err.Error() != "record not found" {
			globle.Logger.Error("查找用户时发生错误！", err)
			utils.RspWithMsg(c, http.StatusBadRequest, false, "系统错误，请稍后再试！")
			return
		}
		if user != nil {
			utils.RspWithMsg(c, http.StatusOK, false, "当前电话号码或用户名已被注册！")
			return
		}
		var NewUser model.User
		NewUser.Name = newUser.Name
		NewUser.PassWd = utils.EncodeWithSHA256(newUser.PassWD)
		NewUser.PhoneNumber = newUser.PhoneNumber
		if err := UserQ.WithContext(ctx).Create(&NewUser); err != nil {
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
	var loginer RR.UserLogIn
	err := c.ShouldBindJSON(&loginer)
	isValid, _ := govalidator.ValidateStruct(loginer)
	switch {
	case err != nil:
		globle.Logger.Error("ShouldBindJSON发生错误！", err)
		utils.RspWithMsg(c, http.StatusBadRequest, false, "非法请求！")
		return
	case !isValid:
		utils.RspWithMsg(c, http.StatusBadRequest, false, "请输入正确格式的Email或电话！")
		return
	case loginer.Name == "" || loginer.Email == "" || loginer.PhoneNumber == "":
		utils.RspWithMsg(c, http.StatusBadRequest, false, "请输入用户名、电话或Email以登录！")
		return
	default:
		user, err := UserQ.WithContext(ctx).Where(
			query.User.PhoneNumber.Eq(loginer.PhoneNumber)).
			Or(query.User.Name.Eq(loginer.Name)).
			Or(query.User.Email.Eq(loginer.Email)).First()
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
			data := RR.UserLogOk{
				ID:           user.ID,
				Name:         user.Name,
				SelfDescribe: user.SelfDescribe,
				PhoneNumber:  user.PhoneNumber,
				Email:        user.Email,
			}
			utils.RspWithData(c, http.StatusOK, true, data)
			return
		}
	}
}

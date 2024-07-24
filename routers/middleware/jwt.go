package middleware

import (
	"IM/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		strToken := c.Request.Header.Get("Authorization")
		if strToken == "" {
			utils.DefaultRsp(c, 1, false, "无权限访问")
			c.Abort()
			return
		}

		token, err := utils.ParseToken(strToken)
		if err != nil {
			utils.DefaultRsp(c, 1, false, "无权限访问")
			c.Abort()
			return
		}

		if token.ExpiresAt.Before(time.Now()) {
			utils.DefaultRsp(c, 1, false, "登陆状态失效，请重新登录")
			c.Abort()
			return
		}

		if token.ExpiresAt.Sub(time.Now()) <= time.Hour*time.Duration(viper.GetInt("jwt.reIssueToken_time")) {
			jwtToken, _ := utils.GenToken(token.UID)
			c.Set("Token", jwtToken)
		}
		c.Set("uid", token.UID)
		c.Next()

		if Token, isExist := c.Get("Token"); isExist {
			c.Writer.Header().Set("Token", Token.(string))
		}
	}
}

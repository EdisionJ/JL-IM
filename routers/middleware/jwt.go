package middleware

import (
	"IM/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			utils.RspWithMsg(c, 1, false, "无权限访问")
			c.Abort()
			return
		}
		//ID@eyJhbGciOiJIUzI1N...
		parts := strings.SplitN(authHeader, "@", 2)
		if !(len(parts) == 2 && parts[0] == "ID") {
			utils.RspWithMsg(c, 1, false, "无权限访问")
			c.Abort()
			return
		}

		token, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.RspWithMsg(c, 1, false, "无权限访问")
			c.Abort()
			return
		}

		c.Set("uid", token.UID)
		c.Set("expire", token.ExpiresAt)
		c.Next()
	}
}

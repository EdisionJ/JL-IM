package routers

import (
	"IM/docs"
	"IM/routers/middleware"
	"IM/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiPublic := r.Group("/api/public")
	{
		apiPublic.POST("/signUp", service.UserSignUp)
		apiPublic.POST("/login", service.UserLogIn)
	}

	user := r.Group("/api/user")
	user.Use(middleware.Jwt())
	{
		user.POST("add")
		user.GET("getFriendReqList")
	}

	chatRoom := r.Group("/api/group")
	chatRoom.Use(middleware.Jwt())
	{
		chatRoom.POST("newRoom")
		chatRoom.POST("joinRoom")

	}

	chat := r.Group("/contact")
	chat.Use(middleware.Jwt())
	{
		chat.POST("sendMsg")
		chat.GET("getMsg")
	}

	return r
}

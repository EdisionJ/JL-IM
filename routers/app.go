package routers

import (
	"IM/docs"
	"IM/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiPublic := r.Group("/user")
	{
		apiPublic.POST("/signUp", service.UserSignUp)
		apiPublic.POST("/login", service.UserLogIn)
	}

	return r
}

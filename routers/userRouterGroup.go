package routers

import (
	"devops/api"
	"github.com/gin-gonic/gin"
)

func UserRouterGroupInit(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.GET("/userInfo", api.User{}.GetUserInfo)

	}
}

package initialize

import (
	"devops/routers"
	"github.com/gin-gonic/gin"
)

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.Default()

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}

	routers.UserRouterGroupInit(Router)
	routers.KubernetesGroupInit(Router)
	return Router
}

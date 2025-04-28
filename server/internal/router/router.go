package router

import (
	"context"
	"devops/internal/controller"
	"devops/internal/middleware"
	"devops/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var R = new(Router)

type Router struct{}

func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	group.Middleware(service.Middleware().MiddlewareCORS)
	group.Middleware(middleware.MiddlewareResponseHandler)
	group.Middleware(middleware.Auth)
	group.Group("/api", func(group *ghttp.RouterGroup) {
		group.Bind(
			//登录
			controller.Login,
			controller.User,
		)

	})
}

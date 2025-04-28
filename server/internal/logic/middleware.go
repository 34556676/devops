package logic

import (
	"devops/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

func init() {
	service.RegisterMiddleware(NewMiddleware())
}

func NewMiddleware() service.IMiddleware {
	return &sMiddleware{}
}

type sMiddleware struct{}

func (s sMiddleware) MiddlewareCORS(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	// you can set options
	//corsOptions.AllowDomain = []string{"goframe.org", "baidu.com"}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

package cmd

import (
	"context"
	"devops/internal/mounter"
	"devops/internal/router"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			//调用注册已挂载相关组件
			mounter.DoMount(ctx, s)

			s.Group("/", func(group *ghttp.RouterGroup) {
				
				router.R.BindController(ctx, group)
			})
			s.Run()
			return nil
		},
	}
)

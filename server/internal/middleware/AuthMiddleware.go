package middleware

import (
	libUtils "devops/library"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"strings"
)

// Auth 权限判断处理中间件
func Auth(r *ghttp.Request) {
	ctx := gctx.New()
	excludePaths, err := g.Cfg().Get(ctx, "gfToken.excludePaths")
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 2. 检查是否是无需认证的路由
	for _, route := range excludePaths.Strings() {
		if strings.HasPrefix(r.URL.Path, route) {
			r.Middleware.Next()
			return
		}
	}

	// 3. 获取 Authorization 头
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Response.WriteJsonExit(g.Map{
			"code": 401,
			"msg":  "Authorization header is required",
			"data": nil,
		})
		return
	}

	// 4. 检查 Bearer token 格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "Invalid Authorization header format",
			"data":    nil,
		})
		return
	}

	token := parts[1]
	// 5. 解析和验证 token
	claims, err := libUtils.ParseToken(token, "devops")
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "Invalid token: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 6. 将用户信息存入上下文
	userData := claims.Data.(map[string]interface{})

	r.SetCtxVar("userId", userData["id"])
	r.SetCtxVar("userName", userData["userName"])

	// 如果验证通过，则调用 Next 方法继续执行下一个中间件或最终的处理器
	r.Middleware.Next()
}

package api

import (
	"devops/internal/model"
	"github.com/gogf/gf/v2/frame/g"
)

type UserLoginReq struct {
	g.Meta   `path:"/login" tags:"系统后台/登录" method:"post" summary:"用户登录"`
	Username string `p:"username" v:"required#用户名不能为空"`
	Password string `p:"password" v:"required#密码不能为空"`
}

type UserLoginRes struct {
	g.Meta   `mime:"application/json"`
	UserInfo *model.LoginUserRes `json:"userInfo"`
	Token    string              `json:"accessToken"`
}

type HelloRes struct {
}

type UserLoginOutReq struct {
	g.Meta `path:"/logout" tags:"系统后台/登录" method:"get" summary:"退出登录"`
}

type UserLoginOutRes struct {
}

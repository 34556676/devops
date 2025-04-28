package api

import "github.com/gogf/gf/v2/frame/g"



type UserInfoReq struct {
	g.Meta `path:"/user/info" tags:"系统后台/用户信息" method:"get" summary:"用户信息"`
}

type UserInfoRes struct {
}


type UserListReq struct {
	g.Meta `path:"/user/list" tags:"系统后台/登录" method:"get" summary:"用户列表"`
}

type UserListRes struct {
}

type UserRestPasswordReq struct {
	g.Meta      `path:"/user/resetPassword" tags:"系统后台/用户中心" method:"put" summary:"重置密码"`
	NewPassword string `json:"newPassword"`
}

type UserRestPasswordRes struct {
}

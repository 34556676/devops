package api

import "github.com/gogf/gf/v2/frame/g"

type RoleListReq struct {
	g.Meta `path:"/list" tags:"系统后台/角色" method:"post" summary:"角色列表"`
}

type RoleListRes struct{}

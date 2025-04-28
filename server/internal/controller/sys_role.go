package controller

import (
	"context"
	"devops/api"
)

var (
	Role = roleController{}
)

type roleController struct {
	BaseController
}

func (c *roleController) List(ctx context.Context, req *api.RoleListReq) (res *api.RoleListRes, err error) {
	res = new(api.RoleListRes)

	return res, nil
}

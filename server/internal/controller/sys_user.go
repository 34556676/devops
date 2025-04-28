package controller

import (
	"context"
	"devops/api"
	"devops/internal/service"
)

var (
	User = userController{}
)

type userController struct {
	BaseController
}

func (c *userController) Info(ctx context.Context, req *api.UserInfoReq) (res *api.UserInfoRes, err error) {
	res = new(api.UserInfoRes)

	return res, nil
}


func (c *userController) List(ctx context.Context, req *api.UserListReq) (res *api.UserListRes, err error) {
	res = new(api.UserListRes)

	return res, nil
}

func (c *userController) RestPassword(ctx context.Context, req *api.UserRestPasswordReq) (res *api.UserRestPasswordRes, err error) {
	res = new(api.UserRestPasswordRes)
	err = service.SysUser().UpdateUserPasswordByUsername(ctx, req.NewPassword)
	return res, err
}

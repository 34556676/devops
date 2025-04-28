package service

import (
	"context"
	"devops/api"
	"devops/internal/model"
)

type (
	ISysUser interface {
		GetAdminUserByUsernamePassword(ctx context.Context, req *api.UserLoginReq) (user *model.LoginUserRes, err error)
		GetUserByUsername(ctx context.Context, userName string) (user *model.LoginUserRes, err error)
		UpdateUserPasswordByUsername(ctx context.Context, newPassword string) (err error)
	}
)

var (
	localSysUser ISysUser
)

func SysUser() ISysUser {
	if localSysUser == nil {
		panic("implement not found for interface ISysUser, forgot register?")
	}
	return localSysUser
}

func RegisterSysUser(i ISysUser) {
	localSysUser = i
}

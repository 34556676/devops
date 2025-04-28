package controller

import (
	"context"
	"devops/api"
	"devops/internal/model"
	"devops/internal/service"
	libUtils "devops/library"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	Login = loginController{}
)

type loginController struct {
	BaseController
}

func (c *loginController) Login(ctx context.Context, req *api.UserLoginReq) (res *api.UserLoginRes, err error) {
	var (
		user  *model.LoginUserRes
		token string
	)

	//ip := libUtils.GetClientIp(ctx)
	//userAgent := libUtils.GetUserAgent(ctx)

	user, err = service.SysUser().GetAdminUserByUsernamePassword(ctx, req)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.NewCode(gcode.New(10001, "用户名不存在", nil))
	}

	//判断密码是否正确
	if !libUtils.ComparePasswords(user.UserPassword, req.Password, user.UserSalt) {
		return nil, gerror.NewCode(gcode.New(10002, "密码错误", nil))
	}

	//key := gconv.String(user.Id) + "-" + gmd5.MustEncryptString(user.UserName) + gmd5.MustEncryptString(user.UserPassword)
	//if g.Cfg().MustGet(ctx, "gfToken.multiLogin").Bool() {
	//	key = gconv.String(user.Id) + "-" + gmd5.MustEncryptString(user.UserName) + gmd5.MustEncryptString(user.UserPassword+ip+userAgent)
	//}
	token, err = libUtils.GenerateToken(ctx, "devops", user)
	if err != nil {
		g.Log().Error(ctx, err)
		err = gerror.New("登录失败，后端服务出现错误")
		return nil, gerror.NewCode(gcode.New(10003, "登录失败，后端服务出现错误", nil))
	}

	res = &api.UserLoginRes{
		UserInfo: user,
		Token:    token,
	}
	fmt.Println(*res)
	return res, nil
}

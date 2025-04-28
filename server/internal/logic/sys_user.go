package logic

import (
	"context"
	"devops/api"
	"devops/internal/dao"
	"devops/internal/model"
	"devops/internal/service"
	libUtils "devops/library"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterSysUser(NewSysUser())
}

type sSysUser struct {
}

func (s sSysUser) UpdateUserPasswordByUsername(ctx context.Context, newPassword string) (err error) {
	salt, err := libUtils.GenerateSecureSalt() // 使用上面定义的generateSecureSalt函数。
	if err != nil {
		fmt.Println("Error generating salt:", err)
		return // 或者其他错误处理。
	}

	saltPassword, err := libUtils.EncryptPassword(newPassword, salt)
	if err != nil {
		return gerror.NewCode(gcode.New(10008, "修改密码失败", nil))
	}
	_, err = dao.SysUser.SysUserDao.Ctx(ctx).Data(g.Map{
		"password":  saltPassword,
		"user_salt": salt,
	}).Where("username", ctx.Value("userName")).Update()

	if err != nil {
		return gerror.NewCode(gcode.New(10009, "修改密码失败", nil))
	}

	return
}

func NewSysUser() service.ISysUser {
	return &sSysUser{}
}

func (s sSysUser) GetUserByUsername(ctx context.Context, userName string) (user *model.LoginUserRes, err error) {
	user = &model.LoginUserRes{}
	err = dao.SysUser.Ctx(ctx).Fields(user).Where(dao.SysUser.Columns().Username, userName).Scan(user)
	return user, err
}

func (s sSysUser) GetAdminUserByUsernamePassword(ctx context.Context, req *api.UserLoginReq) (user *model.LoginUserRes, err error) {
	user, err = s.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("用户名不存在")
	}
	//账号状态
	if user.UserStatus == 0 {
		return nil, errors.New("账号已被冻结")
	}

	return user, nil
}

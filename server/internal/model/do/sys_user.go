// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure of table sys_user for DAO operations like Where/Data.
type SysUser struct {
	g.Meta      `orm:"table:sys_user, do:true"`
	Id          interface{} // 主键ID
	Username    interface{} // 用户名
	Password    interface{} // 密码
	UserSalt    interface{} //
	NickName    interface{} // 姓名
	Email       interface{} // 邮箱
	Phone       interface{} // 手机号
	AuthorityId interface{} // 用户角色ID
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}

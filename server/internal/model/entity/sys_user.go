// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure for table sys_user.
type SysUser struct {
	Id          int         `json:"id"          orm:"id"           description:"主键ID"`   // 主键ID
	Username    string      `json:"username"    orm:"username"     description:"用户名"`    // 用户名
	Password    string      `json:"password"    orm:"password"     description:"密码"`     // 密码
	UserSalt    string      `json:"userSalt"    orm:"user_salt"    description:""`       //
	NickName    string      `json:"nickName"    orm:"nick_name"    description:"姓名"`     // 姓名
	Email       string      `json:"email"       orm:"email"        description:"邮箱"`     // 邮箱
	Phone       string      `json:"phone"       orm:"phone"        description:"手机号"`    // 手机号
	AuthorityId int64       `json:"authorityId" orm:"authority_id" description:"用户角色ID"` // 用户角色ID
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`   // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"`   // 更新时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:"删除时间"`   // 删除时间
}

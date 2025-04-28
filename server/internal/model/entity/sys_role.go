// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRole is the golang structure for table sys_role.
type SysRole struct {
	RoleId    int64       `json:"roleId"    orm:"role_id"    description:"角色ID"`  // 角色ID
	RoleName  string      `json:"roleName"  orm:"role_name"  description:"角色名"`   // 角色名
	ParentId  int64       `json:"parentId"  orm:"parent_id"  description:"父角色ID"` // 父角色ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`      //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`      //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""`      //
}

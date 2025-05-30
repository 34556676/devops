// =================================================================================
// This file is auto-generated by the GoFrame CLI tool. You may modify it as needed.
// =================================================================================

package dao

import (
	"devops/internal/dao/internal"
)

// sysRoleDao is the data access object for the table sys_role.
// You can define custom methods on it to extend its functionality as needed.
type sysRoleDao struct {
	*internal.SysRoleDao
}

var (
	// SysRole is a globally accessible object for table sys_role operations.
	SysRole = sysRoleDao{internal.NewSysRoleDao()}
)

// Add your custom methods and functionality below.

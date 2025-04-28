// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysRoleDao is the data access object for the table sys_role.
type SysRoleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysRoleColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysRoleColumns defines and stores column names for the table sys_role.
type SysRoleColumns struct {
	RoleId    string // 角色ID
	RoleName  string // 角色名
	ParentId  string // 父角色ID
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// sysRoleColumns holds the columns for the table sys_role.
var sysRoleColumns = SysRoleColumns{
	RoleId:    "role_id",
	RoleName:  "role_name",
	ParentId:  "parent_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewSysRoleDao creates and returns a new DAO object for table data access.
func NewSysRoleDao(handlers ...gdb.ModelHandler) *SysRoleDao {
	return &SysRoleDao{
		group:    "default",
		table:    "sys_role",
		columns:  sysRoleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysRoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysRoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysRoleDao) Columns() SysRoleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysRoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysRoleDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SysRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeviceGroupsDao is the data access object for the table device_groups.
type DeviceGroupsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  DeviceGroupsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// DeviceGroupsColumns defines and stores column names for the table device_groups.
type DeviceGroupsColumns struct {
	Id        string //
	Name      string //
	SortOrder string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// deviceGroupsColumns holds the columns for the table device_groups.
var deviceGroupsColumns = DeviceGroupsColumns{
	Id:        "id",
	Name:      "name",
	SortOrder: "sort_order",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewDeviceGroupsDao creates and returns a new DAO object for table data access.
func NewDeviceGroupsDao(handlers ...gdb.ModelHandler) *DeviceGroupsDao {
	return &DeviceGroupsDao{
		group:    "default",
		table:    "device_groups",
		columns:  deviceGroupsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DeviceGroupsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DeviceGroupsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DeviceGroupsDao) Columns() DeviceGroupsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DeviceGroupsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DeviceGroupsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DeviceGroupsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

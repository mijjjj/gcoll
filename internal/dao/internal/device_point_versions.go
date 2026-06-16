// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DevicePointVersionsDao is the data access object for the table device_point_versions.
type DevicePointVersionsDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  DevicePointVersionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// DevicePointVersionsColumns defines and stores column names for the table device_point_versions.
type DevicePointVersionsColumns struct {
	Id           string //
	PointId      string //
	DeviceId     string //
	PluginId     string //
	Version      string //
	SnapshotJson string //
	ChangeNote   string //
	CreatedAt    string //
}

// devicePointVersionsColumns holds the columns for the table device_point_versions.
var devicePointVersionsColumns = DevicePointVersionsColumns{
	Id:           "id",
	PointId:      "point_id",
	DeviceId:     "device_id",
	PluginId:     "plugin_id",
	Version:      "version",
	SnapshotJson: "snapshot_json",
	ChangeNote:   "change_note",
	CreatedAt:    "created_at",
}

// NewDevicePointVersionsDao creates and returns a new DAO object for table data access.
func NewDevicePointVersionsDao(handlers ...gdb.ModelHandler) *DevicePointVersionsDao {
	return &DevicePointVersionsDao{
		group:    "default",
		table:    "device_point_versions",
		columns:  devicePointVersionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DevicePointVersionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DevicePointVersionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DevicePointVersionsDao) Columns() DevicePointVersionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DevicePointVersionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DevicePointVersionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DevicePointVersionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

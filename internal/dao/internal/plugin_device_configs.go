// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PluginDeviceConfigsDao is the data access object for the table plugin_device_configs.
type PluginDeviceConfigsDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  PluginDeviceConfigsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// PluginDeviceConfigsColumns defines and stores column names for the table plugin_device_configs.
type PluginDeviceConfigsColumns struct {
	Id         string //
	DeviceId   string //
	PluginId   string //
	Version    string //
	ConfigJson string //
	ReportMode string //
	Enabled    string //
	Active     string //
	CreatedAt  string //
	UpdatedAt  string //
}

// pluginDeviceConfigsColumns holds the columns for the table plugin_device_configs.
var pluginDeviceConfigsColumns = PluginDeviceConfigsColumns{
	Id:         "id",
	DeviceId:   "device_id",
	PluginId:   "plugin_id",
	Version:    "version",
	ConfigJson: "config_json",
	ReportMode: "report_mode",
	Enabled:    "enabled",
	Active:     "active",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewPluginDeviceConfigsDao creates and returns a new DAO object for table data access.
func NewPluginDeviceConfigsDao(handlers ...gdb.ModelHandler) *PluginDeviceConfigsDao {
	return &PluginDeviceConfigsDao{
		group:    "default",
		table:    "plugin_device_configs",
		columns:  pluginDeviceConfigsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PluginDeviceConfigsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PluginDeviceConfigsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PluginDeviceConfigsDao) Columns() PluginDeviceConfigsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PluginDeviceConfigsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PluginDeviceConfigsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PluginDeviceConfigsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

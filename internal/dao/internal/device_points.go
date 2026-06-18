// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DevicePointsDao is the data access object for the table device_points.
type DevicePointsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  DevicePointsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// DevicePointsColumns defines and stores column names for the table device_points.
type DevicePointsColumns struct {
	Id           string //
	DeviceId     string //
	PluginId     string //
	Name         string //
	Description  string //
	Address      string //
	ValueType    string //
	Unit         string //
	Enabled      string //
	TagsJson     string //
	MetadataJson string //
	CreatedAt    string //
	UpdatedAt    string //
}

// devicePointsColumns holds the columns for the table device_points.
var devicePointsColumns = DevicePointsColumns{
	Id:           "id",
	DeviceId:     "device_id",
	PluginId:     "plugin_id",
	Name:         "name",
	Description:  "description",
	Address:      "address",
	ValueType:    "value_type",
	Unit:         "unit",
	Enabled:      "enabled",
	TagsJson:     "tags_json",
	MetadataJson: "metadata_json",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewDevicePointsDao creates and returns a new DAO object for table data access.
func NewDevicePointsDao(handlers ...gdb.ModelHandler) *DevicePointsDao {
	return &DevicePointsDao{
		group:    "default",
		table:    "device_points",
		columns:  devicePointsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DevicePointsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DevicePointsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DevicePointsDao) Columns() DevicePointsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DevicePointsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DevicePointsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DevicePointsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

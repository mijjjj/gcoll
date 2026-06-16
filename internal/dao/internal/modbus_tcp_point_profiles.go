// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ModbusTcpPointProfilesDao is the data access object for the table modbus_tcp_point_profiles.
type ModbusTcpPointProfilesDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  ModbusTcpPointProfilesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// ModbusTcpPointProfilesColumns defines and stores column names for the table modbus_tcp_point_profiles.
type ModbusTcpPointProfilesColumns struct {
	Id         string //
	DeviceId   string //
	PointId    string //
	PluginId   string //
	Version    string //
	Area       string //
	Address    string //
	Quantity   string //
	Mode       string //
	ValueType  string //
	ByteOrder  string //
	WordOrder  string //
	Scale      string //
	Offset     string //
	ReportMode string //
	Enabled    string //
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
}

// modbusTcpPointProfilesColumns holds the columns for the table modbus_tcp_point_profiles.
var modbusTcpPointProfilesColumns = ModbusTcpPointProfilesColumns{
	Id:         "id",
	DeviceId:   "device_id",
	PointId:    "point_id",
	PluginId:   "plugin_id",
	Version:    "version",
	Area:       "area",
	Address:    "address",
	Quantity:   "quantity",
	Mode:       "mode",
	ValueType:  "value_type",
	ByteOrder:  "byte_order",
	WordOrder:  "word_order",
	Scale:      "scale",
	Offset:     "offset",
	ReportMode: "report_mode",
	Enabled:    "enabled",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewModbusTcpPointProfilesDao creates and returns a new DAO object for table data access.
func NewModbusTcpPointProfilesDao(handlers ...gdb.ModelHandler) *ModbusTcpPointProfilesDao {
	return &ModbusTcpPointProfilesDao{
		group:    "default",
		table:    "modbus_tcp_point_profiles",
		columns:  modbusTcpPointProfilesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ModbusTcpPointProfilesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ModbusTcpPointProfilesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ModbusTcpPointProfilesDao) Columns() ModbusTcpPointProfilesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ModbusTcpPointProfilesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ModbusTcpPointProfilesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ModbusTcpPointProfilesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

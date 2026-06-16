// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ModbusTcpDeviceProfilesDao is the data access object for the table modbus_tcp_device_profiles.
type ModbusTcpDeviceProfilesDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  ModbusTcpDeviceProfilesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// ModbusTcpDeviceProfilesColumns defines and stores column names for the table modbus_tcp_device_profiles.
type ModbusTcpDeviceProfilesColumns struct {
	Id               string //
	DeviceId         string //
	PluginId         string //
	Version          string //
	Host             string //
	Port             string //
	UnitId           string //
	TimeoutMs        string //
	PollIntervalMs   string //
	ReportMode       string //
	MaxCoilBatch     string //
	MaxRegisterBatch string //
	LowLatencyMs     string //
	HighLatencyMs    string //
	DebugEnabled     string //
	Enabled          string //
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// modbusTcpDeviceProfilesColumns holds the columns for the table modbus_tcp_device_profiles.
var modbusTcpDeviceProfilesColumns = ModbusTcpDeviceProfilesColumns{
	Id:               "id",
	DeviceId:         "device_id",
	PluginId:         "plugin_id",
	Version:          "version",
	Host:             "host",
	Port:             "port",
	UnitId:           "unit_id",
	TimeoutMs:        "timeout_ms",
	PollIntervalMs:   "poll_interval_ms",
	ReportMode:       "report_mode",
	MaxCoilBatch:     "max_coil_batch",
	MaxRegisterBatch: "max_register_batch",
	LowLatencyMs:     "low_latency_ms",
	HighLatencyMs:    "high_latency_ms",
	DebugEnabled:     "debug_enabled",
	Enabled:          "enabled",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewModbusTcpDeviceProfilesDao creates and returns a new DAO object for table data access.
func NewModbusTcpDeviceProfilesDao(handlers ...gdb.ModelHandler) *ModbusTcpDeviceProfilesDao {
	return &ModbusTcpDeviceProfilesDao{
		group:    "default",
		table:    "modbus_tcp_device_profiles",
		columns:  modbusTcpDeviceProfilesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ModbusTcpDeviceProfilesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ModbusTcpDeviceProfilesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ModbusTcpDeviceProfilesDao) Columns() ModbusTcpDeviceProfilesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ModbusTcpDeviceProfilesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ModbusTcpDeviceProfilesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ModbusTcpDeviceProfilesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

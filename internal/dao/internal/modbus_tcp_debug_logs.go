// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ModbusTcpDebugLogsDao is the data access object for the table modbus_tcp_debug_logs.
type ModbusTcpDebugLogsDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  ModbusTcpDebugLogsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// ModbusTcpDebugLogsColumns defines and stores column names for the table modbus_tcp_debug_logs.
type ModbusTcpDebugLogsColumns struct {
	Id         string //
	DeviceId   string //
	TaskId     string //
	PointId    string //
	TraceId    string //
	Level      string //
	Message    string //
	Area       string //
	Address    string //
	LatencyMs  string //
	RawHex     string //
	FieldsJson string //
	CreatedAt  string //
}

// modbusTcpDebugLogsColumns holds the columns for the table modbus_tcp_debug_logs.
var modbusTcpDebugLogsColumns = ModbusTcpDebugLogsColumns{
	Id:         "id",
	DeviceId:   "device_id",
	TaskId:     "task_id",
	PointId:    "point_id",
	TraceId:    "trace_id",
	Level:      "level",
	Message:    "message",
	Area:       "area",
	Address:    "address",
	LatencyMs:  "latency_ms",
	RawHex:     "raw_hex",
	FieldsJson: "fields_json",
	CreatedAt:  "created_at",
}

// NewModbusTcpDebugLogsDao creates and returns a new DAO object for table data access.
func NewModbusTcpDebugLogsDao(handlers ...gdb.ModelHandler) *ModbusTcpDebugLogsDao {
	return &ModbusTcpDebugLogsDao{
		group:    "default",
		table:    "modbus_tcp_debug_logs",
		columns:  modbusTcpDebugLogsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ModbusTcpDebugLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ModbusTcpDebugLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ModbusTcpDebugLogsDao) Columns() ModbusTcpDebugLogsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ModbusTcpDebugLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ModbusTcpDebugLogsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ModbusTcpDebugLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

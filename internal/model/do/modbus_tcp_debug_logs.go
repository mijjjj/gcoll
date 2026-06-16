// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ModbusTcpDebugLogs is the golang structure of table modbus_tcp_debug_logs for DAO operations like Where/Data.
type ModbusTcpDebugLogs struct {
	g.Meta     `orm:"table:modbus_tcp_debug_logs, do:true"`
	Id         any //
	DeviceId   any //
	TaskId     any //
	PointId    any //
	TraceId    any //
	Level      any //
	Message    any //
	Area       any //
	Address    any //
	LatencyMs  any //
	RawHex     any //
	FieldsJson any //
	CreatedAt  any //
}

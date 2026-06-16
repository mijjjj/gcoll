// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ModbusTcpDeviceProfiles is the golang structure of table modbus_tcp_device_profiles for DAO operations like Where/Data.
type ModbusTcpDeviceProfiles struct {
	g.Meta           `orm:"table:modbus_tcp_device_profiles, do:true"`
	Id               any //
	DeviceId         any //
	PluginId         any //
	Version          any //
	Host             any //
	Port             any //
	UnitId           any //
	TimeoutMs        any //
	PollIntervalMs   any //
	ReportMode       any //
	MaxCoilBatch     any //
	MaxRegisterBatch any //
	LowLatencyMs     any //
	HighLatencyMs    any //
	DebugEnabled     any //
	Enabled          any //
	CreatedAt        any //
	UpdatedAt        any //
	DeletedAt        any //
}

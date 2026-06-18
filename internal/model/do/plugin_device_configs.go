// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PluginDeviceConfigs is the golang structure of table plugin_device_configs for DAO operations like Where/Data.
type PluginDeviceConfigs struct {
	g.Meta     `orm:"table:plugin_device_configs, do:true"`
	Id         any //
	DeviceId   any //
	PluginId   any //
	Version    any //
	ConfigJson any //
	ReportMode any //
	Enabled    any //
	Active     any //
	CreatedAt  any //
	UpdatedAt  any //
}

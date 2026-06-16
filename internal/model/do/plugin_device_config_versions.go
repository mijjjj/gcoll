// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PluginDeviceConfigVersions is the golang structure of table plugin_device_config_versions for DAO operations like Where/Data.
type PluginDeviceConfigVersions struct {
	g.Meta     `orm:"table:plugin_device_config_versions, do:true"`
	Id         any //
	ConfigId   any //
	DeviceId   any //
	PluginId   any //
	Version    any //
	ConfigJson any //
	ChangeNote any //
	CreatedAt  any //
}

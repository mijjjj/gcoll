// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// DevicePointVersions is the golang structure of table device_point_versions for DAO operations like Where/Data.
type DevicePointVersions struct {
	g.Meta       `orm:"table:device_point_versions, do:true"`
	Id           any //
	PointId      any //
	DeviceId     any //
	PluginId     any //
	Version      any //
	SnapshotJson any //
	ChangeNote   any //
	CreatedAt    any //
}

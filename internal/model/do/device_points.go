// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// DevicePoints is the golang structure of table device_points for DAO operations like Where/Data.
type DevicePoints struct {
	g.Meta       `orm:"table:device_points, do:true"`
	Id           any //
	DeviceId     any //
	PluginId     any //
	Name         any //
	Description  any //
	Address      any //
	ValueType    any //
	Unit         any //
	Enabled      any //
	TagsJson     any //
	MetadataJson any //
	CreatedAt    any //
	UpdatedAt    any //
	DeletedAt    any //
}

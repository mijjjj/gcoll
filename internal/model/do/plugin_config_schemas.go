// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PluginConfigSchemas is the golang structure of table plugin_config_schemas for DAO operations like Where/Data.
type PluginConfigSchemas struct {
	g.Meta          `orm:"table:plugin_config_schemas, do:true"`
	Id              any //
	PluginId        any //
	PluginVersionId any //
	SchemaVersion   any //
	SchemaJson      any //
	CreatedAt       any //
	DeletedAt       any //
}

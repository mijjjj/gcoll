// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PluginVersions is the golang structure of table plugin_versions for DAO operations like Where/Data.
type PluginVersions struct {
	g.Meta           `orm:"table:plugin_versions, do:true"`
	Id               any //
	PluginId         any //
	Version          any //
	PackagePath      any //
	ManifestJson     any //
	PermissionsJson  any //
	CapabilitiesJson any //
	Checksum         any //
	Active           any //
	InstalledAt      any //
	CreatedAt        any //
	UpdatedAt        any //
}

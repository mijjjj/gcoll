// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PluginVersions is the golang structure for table plugin_versions.
type PluginVersions struct {
	Id               string `json:"id"               orm:"id"                ` //
	PluginId         string `json:"pluginId"         orm:"plugin_id"         ` //
	Version          string `json:"version"          orm:"version"           ` //
	PackagePath      string `json:"packagePath"      orm:"package_path"      ` //
	ManifestJson     string `json:"manifestJson"     orm:"manifest_json"     ` //
	PermissionsJson  string `json:"permissionsJson"  orm:"permissions_json"  ` //
	CapabilitiesJson string `json:"capabilitiesJson" orm:"capabilities_json" ` //
	Checksum         string `json:"checksum"         orm:"checksum"          ` //
	Active           int    `json:"active"           orm:"active"            ` //
	InstalledAt      string `json:"installedAt"      orm:"installed_at"      ` //
	CreatedAt        string `json:"createdAt"        orm:"created_at"        ` //
	UpdatedAt        string `json:"updatedAt"        orm:"updated_at"        ` //
	DeletedAt        string `json:"deletedAt"        orm:"deleted_at"        ` //
}

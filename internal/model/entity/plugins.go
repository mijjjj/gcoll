// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Plugins is the golang structure for table plugins.
type Plugins struct {
	Id            string `json:"id"            orm:"id"             ` //
	Name          string `json:"name"          orm:"name"           ` //
	Type          string `json:"type"          orm:"type"           ` //
	Runtime       string `json:"runtime"       orm:"runtime"        ` //
	Protocol      string `json:"protocol"      orm:"protocol"       ` //
	Status        string `json:"status"        orm:"status"         ` //
	ActiveVersion string `json:"activeVersion" orm:"active_version" ` //
	Source        string `json:"source"        orm:"source"         ` //
	Description   string `json:"description"   orm:"description"    ` //
	Enabled       int    `json:"enabled"       orm:"enabled"        ` //
	InstalledAt   string `json:"installedAt"   orm:"installed_at"   ` //
	UpdatedAt     string `json:"updatedAt"     orm:"updated_at"     ` //
	DeletedAt     string `json:"deletedAt"     orm:"deleted_at"     ` //
}

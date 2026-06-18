// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Devices is the golang structure for table devices.
type Devices struct {
	Id           string `json:"id"           orm:"id"            ` //
	Name         string `json:"name"         orm:"name"          ` //
	Code         string `json:"code"         orm:"code"          ` //
	GroupId      string `json:"groupId"      orm:"group_id"      ` //
	PluginId     string `json:"pluginId"     orm:"plugin_id"     ` //
	Status       string `json:"status"       orm:"status"        ` //
	Enabled      int    `json:"enabled"      orm:"enabled"       ` //
	ReportMode   string `json:"reportMode"   orm:"report_mode"   ` //
	LastSeenAt   string `json:"lastSeenAt"   orm:"last_seen_at"  ` //
	Description  string `json:"description"  orm:"description"   ` //
	MetadataJson string `json:"metadataJson" orm:"metadata_json" ` //
	CreatedAt    string `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    string `json:"updatedAt"    orm:"updated_at"    ` //
}

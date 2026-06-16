// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// DevicePoints is the golang structure for table device_points.
type DevicePoints struct {
	Id           string `json:"id"           orm:"id"            ` //
	DeviceId     string `json:"deviceId"     orm:"device_id"     ` //
	PluginId     string `json:"pluginId"     orm:"plugin_id"     ` //
	Name         string `json:"name"         orm:"name"          ` //
	Description  string `json:"description"  orm:"description"   ` //
	Address      string `json:"address"      orm:"address"       ` //
	ValueType    string `json:"valueType"    orm:"value_type"    ` //
	Unit         string `json:"unit"         orm:"unit"          ` //
	Enabled      int    `json:"enabled"      orm:"enabled"       ` //
	TagsJson     string `json:"tagsJson"     orm:"tags_json"     ` //
	MetadataJson string `json:"metadataJson" orm:"metadata_json" ` //
	CreatedAt    string `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    string `json:"updatedAt"    orm:"updated_at"    ` //
	DeletedAt    string `json:"deletedAt"    orm:"deleted_at"    ` //
}

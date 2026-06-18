// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PluginConfigSchemas is the golang structure for table plugin_config_schemas.
type PluginConfigSchemas struct {
	Id              string `json:"id"              orm:"id"                ` //
	PluginId        string `json:"pluginId"        orm:"plugin_id"         ` //
	PluginVersionId string `json:"pluginVersionId" orm:"plugin_version_id" ` //
	SchemaVersion   int    `json:"schemaVersion"   orm:"schema_version"    ` //
	SchemaJson      string `json:"schemaJson"      orm:"schema_json"       ` //
	CreatedAt       string `json:"createdAt"       orm:"created_at"        ` //
}

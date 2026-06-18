// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CollectionTasks is the golang structure for table collection_tasks.
type CollectionTasks struct {
	Id              string `json:"id"              orm:"id"                ` //
	Name            string `json:"name"            orm:"name"              ` //
	DeviceId        string `json:"deviceId"        orm:"device_id"         ` //
	SouthPluginId   string `json:"southPluginId"   orm:"south_plugin_id"   ` //
	ReportMode      string `json:"reportMode"      orm:"report_mode"       ` //
	Status          string `json:"status"          orm:"status"            ` //
	RuleHitRate     string `json:"ruleHitRate"     orm:"rule_hit_rate"     ` //
	Rate            string `json:"rate"            orm:"rate"              ` //
	LastCollectedAt string `json:"lastCollectedAt" orm:"last_collected_at" ` //
	CreatedAt       string `json:"createdAt"       orm:"created_at"        ` //
	UpdatedAt       string `json:"updatedAt"       orm:"updated_at"        ` //
}

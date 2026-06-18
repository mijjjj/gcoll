// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CollectionTasks is the golang structure of table collection_tasks for DAO operations like Where/Data.
type CollectionTasks struct {
	g.Meta          `orm:"table:collection_tasks, do:true"`
	Id              any //
	Name            any //
	DeviceId        any //
	SouthPluginId   any //
	ReportMode      any //
	Status          any //
	RuleHitRate     any //
	Rate            any //
	LastCollectedAt any //
	CreatedAt       any //
	UpdatedAt       any //
}

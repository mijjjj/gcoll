package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// OverviewReq 描述工作台总览请求。
type OverviewReq struct {
	g.Meta `path:"/runtime/overview" method:"get" tags:"Runtime" summary:"获取运行时工作台总览"`
}

// OverviewRes 描述工作台总览响应。
type OverviewRes struct {
	Metrics       []commonv1.MetricItem      `json:"metrics"`
	Runtime       commonv1.RuntimeStatus     `json:"runtime"`
	DataPlane     commonv1.DataPlaneStatus   `json:"dataPlane"`
	Tasks         []commonv1.TaskSummary     `json:"tasks"`
	RecentEvents  []commonv1.RuntimeEvent    `json:"recentEvents"`
	PluginSummary commonv1.PluginSummary     `json:"pluginSummary"`
	Network       commonv1.RuntimeDependency `json:"network"`
}

package v1

import "github.com/gogf/gf/v2/frame/g"

// OverviewReq 描述工作台总览请求。
type OverviewReq struct {
	g.Meta `path:"/runtime/overview" method:"get" tags:"Runtime" summary:"获取运行时工作台总览"`
}

// OverviewRes 描述工作台总览响应。
type OverviewRes struct {
	Metrics       []MetricItem      `json:"metrics"`
	Runtime       RuntimeStatus     `json:"runtime"`
	DataPlane     DataPlaneStatus   `json:"dataPlane"`
	Tasks         []TaskSummary     `json:"tasks"`
	RecentEvents  []RuntimeEvent    `json:"recentEvents"`
	PluginSummary PluginSummary     `json:"pluginSummary"`
	Network       RuntimeDependency `json:"network"`
}

// PluginsReq 描述插件列表请求。
type PluginsReq struct {
	g.Meta `path:"/plugins" method:"get" tags:"Plugins" summary:"获取插件列表"`
}

// PluginsRes 描述插件列表响应。
type PluginsRes struct {
	Items []PluginItem `json:"items"`
}

// DevicesReq 描述设备列表请求。
type DevicesReq struct {
	g.Meta `path:"/devices" method:"get" tags:"Devices" summary:"获取设备列表"`
}

// DevicesRes 描述设备列表响应。
type DevicesRes struct {
	Groups []DeviceGroup `json:"groups"`
	Items  []DeviceItem  `json:"items"`
}

// DevicePointsReq 描述设备点位列表请求。
type DevicePointsReq struct {
	g.Meta   `path:"/devices/{deviceId}/points" method:"get" tags:"Devices" summary:"获取设备点位列表"`
	DeviceId string `json:"deviceId" in:"path"`
}

// DevicePointsRes 描述设备点位列表响应。
type DevicePointsRes struct {
	Items []PointItem `json:"items"`
}

// TasksReq 描述采集任务列表请求。
type TasksReq struct {
	g.Meta `path:"/tasks" method:"get" tags:"Tasks" summary:"获取采集任务列表"`
}

// TasksRes 描述采集任务列表响应。
type TasksRes struct {
	Items []TaskSummary `json:"items"`
}

// PointCacheReq 描述最新点位缓存请求。
type PointCacheReq struct {
	g.Meta `path:"/point-cache" method:"get" tags:"PointCache" summary:"获取最新点位缓存"`
}

// PointCacheRes 描述最新点位缓存响应。
type PointCacheRes struct {
	Items []PointCacheItem `json:"items"`
}

// PipelineRulesReq 描述规则过滤列表请求。
type PipelineRulesReq struct {
	g.Meta `path:"/pipeline/rules" method:"get" tags:"Pipeline" summary:"获取规则过滤列表"`
}

// PipelineRulesRes 描述规则过滤列表响应。
type PipelineRulesRes struct {
	Items []PipelineRuleItem `json:"items"`
}

// TargetsReq 描述北向转发目标列表请求。
type TargetsReq struct {
	g.Meta `path:"/targets" method:"get" tags:"Targets" summary:"获取北向转发目标列表"`
}

// TargetsRes 描述北向转发目标列表响应。
type TargetsRes struct {
	Items []ForwardTargetItem `json:"items"`
}

// LogsReq 描述运行日志列表请求。
type LogsReq struct {
	g.Meta `path:"/logs" method:"get" tags:"Logs" summary:"获取运行日志列表"`
}

// LogsRes 描述运行日志列表响应。
type LogsRes struct {
	Items []RuntimeEvent `json:"items"`
}

// MetricItem 描述工作台指标。
type MetricItem struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
	Tone  string `json:"tone"`
}

// RuntimeStatus 描述核心运行时状态。
type RuntimeStatus struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Mode      string `json:"mode"`
	CheckedAt string `json:"checkedAt"`
	ApiBase   string `json:"apiBase"`
	Database  string `json:"database"`
}

// RuntimeDependency 描述运行时依赖状态。
type RuntimeDependency struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

// DataPlaneStatus 描述数据面状态。
type DataPlaneStatus struct {
	QueueUsagePercent int    `json:"queueUsagePercent"`
	RuleHitPercent    int    `json:"ruleHitPercent"`
	ForwardPercent    int    `json:"forwardPercent"`
	Throughput        string `json:"throughput"`
	Latency           string `json:"latency"`
	Backpressure      string `json:"backpressure"`
}

// PluginSummary 描述插件运行汇总。
type PluginSummary struct {
	Running int `json:"running"`
	Total   int `json:"total"`
}

// PluginItem 描述插件列表项。
type PluginItem struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Version     string   `json:"version"`
	Runtime     string   `json:"runtime"`
	Protocol    string   `json:"protocol"`
	Status      string   `json:"status"`
	Permissions []string `json:"permissions"`
	UpdatedAt   string   `json:"updatedAt"`
}

// DeviceGroup 描述设备分组。
type DeviceGroup struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// DeviceItem 描述设备列表项。
type DeviceItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	GroupId     string `json:"groupId"`
	PluginId    string `json:"pluginId"`
	PluginName  string `json:"pluginName"`
	Status      string `json:"status"`
	Enabled     bool   `json:"enabled"`
	PointCount  int    `json:"pointCount"`
	ReportMode  string `json:"reportMode"`
	LastSeenAt  string `json:"lastSeenAt"`
	Description string `json:"description"`
}

// PointItem 描述通用点位表列表项。
type PointItem struct {
	Id          string            `json:"id"`
	DeviceId    string            `json:"deviceId"`
	PluginId    string            `json:"pluginId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Address     string            `json:"address"`
	ValueType   string            `json:"valueType"`
	Unit        string            `json:"unit"`
	Enabled     bool              `json:"enabled"`
	Tags        map[string]string `json:"tags"`
}

// TaskSummary 描述采集任务摘要。
type TaskSummary struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	DeviceId        string `json:"deviceId"`
	DeviceName      string `json:"deviceName"`
	SouthPluginName string `json:"southPluginName"`
	PointCount      int    `json:"pointCount"`
	ReportMode      string `json:"reportMode"`
	Status          string `json:"status"`
	Rate            string `json:"rate"`
	RuleHitRate     string `json:"ruleHitRate"`
	LastCollectedAt string `json:"lastCollectedAt"`
}

// PointCacheItem 描述最新点位缓存项。
type PointCacheItem struct {
	PointId   string `json:"pointId"`
	DeviceId  string `json:"deviceId"`
	PointName string `json:"pointName"`
	Value     string `json:"value"`
	Quality   string `json:"quality"`
	Changed   bool   `json:"changed"`
	UpdatedAt string `json:"updatedAt"`
}

// PipelineRuleItem 描述规则过滤项。
type PipelineRuleItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Expression  string `json:"expression"`
	Matched     int    `json:"matched"`
	TargetCount int    `json:"targetCount"`
	UpdatedAt   string `json:"updatedAt"`
}

// ForwardTargetItem 描述北向转发目标。
type ForwardTargetItem struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PluginName string `json:"pluginName"`
	Status     string `json:"status"`
	Endpoint   string `json:"endpoint"`
	LastError  string `json:"lastError"`
	UpdatedAt  string `json:"updatedAt"`
}

// RuntimeEvent 描述运行时事件或日志。
type RuntimeEvent struct {
	Id       string `json:"id"`
	Time     string `json:"time"`
	Level    string `json:"level"`
	Source   string `json:"source"`
	PluginId string `json:"pluginId"`
	DeviceId string `json:"deviceId"`
	TaskId   string `json:"taskId"`
	Message  string `json:"message"`
	TraceId  string `json:"traceId"`
}

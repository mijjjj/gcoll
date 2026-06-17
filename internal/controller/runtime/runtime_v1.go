package runtime

import (
	"context"

	runtimev1 "github.com/mijjjj/gcoll/api/runtime/v1"

	runtimesvc "github.com/mijjjj/gcoll/internal/service/runtime"
)

// ControllerV1 提供运行时 API V1 控制器。
type ControllerV1 struct {
	runtimeSvc *runtimesvc.Service
}

// NewV1 创建运行时 API V1 控制器。
func NewV1() *ControllerV1 {
	return &ControllerV1{
		runtimeSvc: runtimesvc.New(),
	}
}

// Health 返回服务端运行时健康状态。
func (c *ControllerV1) Health(ctx context.Context, req *runtimev1.HealthReq) (res *runtimev1.HealthRes, err error) {
	_ = req

	return c.runtimeSvc.GetHealth(ctx), nil
}

// Overview 返回运行时工作台总览。
func (c *ControllerV1) Overview(ctx context.Context, req *runtimev1.OverviewReq) (res *runtimev1.OverviewRes, err error) {
	_ = req

	return c.runtimeSvc.GetOverview(ctx)
}

// Plugins 返回本地插件列表。
func (c *ControllerV1) Plugins(ctx context.Context, req *runtimev1.PluginsReq) (res *runtimev1.PluginsRes, err error) {
	_ = req

	return c.runtimeSvc.GetPlugins(ctx)
}

// ImportPlugin 导入插件清单。
func (c *ControllerV1) ImportPlugin(ctx context.Context, req *runtimev1.ImportPluginReq) (res *runtimev1.ImportPluginRes, err error) {
	return c.runtimeSvc.ImportPlugin(ctx, req.PackagePath)
}

// DevicePluginConfigPage 返回当前设备的插件配置页数据。
func (c *ControllerV1) DevicePluginConfigPage(ctx context.Context, req *runtimev1.DevicePluginConfigPageReq) (res *runtimev1.DevicePluginConfigPageRes, err error) {
	return c.runtimeSvc.GetDevicePluginConfigPage(ctx, req.DeviceId)
}

// UpdateDevicePluginConfig 保存当前设备的插件配置。
func (c *ControllerV1) UpdateDevicePluginConfig(ctx context.Context, req *runtimev1.UpdateDevicePluginConfigReq) (res *runtimev1.UpdateDevicePluginConfigRes, err error) {
	return c.runtimeSvc.UpdateDevicePluginConfig(ctx, req)
}

// TestDevicePluginConnection 测试当前设备的插件连接。
func (c *ControllerV1) TestDevicePluginConnection(ctx context.Context, req *runtimev1.TestDevicePluginConnectionReq) (res *runtimev1.TestDevicePluginConnectionRes, err error) {
	return c.runtimeSvc.TestDevicePluginConnection(ctx, req.DeviceId)
}

// Devices 返回设备列表。
func (c *ControllerV1) Devices(ctx context.Context, req *runtimev1.DevicesReq) (res *runtimev1.DevicesRes, err error) {
	_ = req

	return c.runtimeSvc.GetDevices(ctx)
}

// CreateDeviceGroup 新增设备分组。
func (c *ControllerV1) CreateDeviceGroup(ctx context.Context, req *runtimev1.CreateDeviceGroupReq) (res *runtimev1.CreateDeviceGroupRes, err error) {
	return c.runtimeSvc.CreateDeviceGroup(ctx, req)
}

// CreateDevice 新增设备。
func (c *ControllerV1) CreateDevice(ctx context.Context, req *runtimev1.CreateDeviceReq) (res *runtimev1.CreateDeviceRes, err error) {
	return c.runtimeSvc.CreateDevice(ctx, req)
}

// MoveDeviceToGroup 移动设备所属分组。
func (c *ControllerV1) MoveDeviceToGroup(ctx context.Context, req *runtimev1.MoveDeviceToGroupReq) (res *runtimev1.MoveDeviceToGroupRes, err error) {
	return c.runtimeSvc.MoveDeviceToGroup(ctx, req)
}

// DeleteDevice 删除设备。
func (c *ControllerV1) DeleteDevice(ctx context.Context, req *runtimev1.DeleteDeviceReq) (res *runtimev1.DeleteDeviceRes, err error) {
	return c.runtimeSvc.DeleteDevice(ctx, req.DeviceId)
}

// DevicePoints 返回指定设备的点位列表。
func (c *ControllerV1) DevicePoints(ctx context.Context, req *runtimev1.DevicePointsReq) (res *runtimev1.DevicePointsRes, err error) {
	return c.runtimeSvc.GetDevicePoints(ctx, req.DeviceId)
}

// CreateDevicePoint 新增指定设备的点位。
func (c *ControllerV1) CreateDevicePoint(ctx context.Context, req *runtimev1.CreateDevicePointReq) (res *runtimev1.CreateDevicePointRes, err error) {
	return c.runtimeSvc.CreateDevicePoint(ctx, req)
}

// Tasks 返回采集任务列表。
func (c *ControllerV1) Tasks(ctx context.Context, req *runtimev1.TasksReq) (res *runtimev1.TasksRes, err error) {
	_ = req

	return c.runtimeSvc.GetTasks(ctx)
}

// StartDeviceCollectionTask 启动设备默认采集任务。
func (c *ControllerV1) StartDeviceCollectionTask(ctx context.Context, req *runtimev1.StartDeviceCollectionTaskReq) (res *runtimev1.CollectionTaskActionRes, err error) {
	return c.runtimeSvc.StartDeviceCollectionTask(ctx, req.DeviceId)
}

// StartCollectionTask 启动采集任务。
func (c *ControllerV1) StartCollectionTask(ctx context.Context, req *runtimev1.StartCollectionTaskReq) (res *runtimev1.CollectionTaskActionRes, err error) {
	return c.runtimeSvc.StartCollectionTask(ctx, req.TaskId)
}

// StopCollectionTask 停止采集任务。
func (c *ControllerV1) StopCollectionTask(ctx context.Context, req *runtimev1.StopCollectionTaskReq) (res *runtimev1.CollectionTaskActionRes, err error) {
	return c.runtimeSvc.StopCollectionTask(ctx, req.TaskId)
}

// PointCache 返回最新点位缓存。
func (c *ControllerV1) PointCache(ctx context.Context, req *runtimev1.PointCacheReq) (res *runtimev1.PointCacheRes, err error) {
	_ = req

	return c.runtimeSvc.GetPointCache(ctx), nil
}

// PipelineRules 返回规则过滤列表。
func (c *ControllerV1) PipelineRules(ctx context.Context, req *runtimev1.PipelineRulesReq) (res *runtimev1.PipelineRulesRes, err error) {
	_ = req

	return c.runtimeSvc.GetPipelineRules(ctx), nil
}

// Targets 返回北向转发目标列表。
func (c *ControllerV1) Targets(ctx context.Context, req *runtimev1.TargetsReq) (res *runtimev1.TargetsRes, err error) {
	_ = req

	return c.runtimeSvc.GetTargets(ctx), nil
}

// Logs 返回运行日志列表。
func (c *ControllerV1) Logs(ctx context.Context, req *runtimev1.LogsReq) (res *runtimev1.LogsRes, err error) {
	_ = req

	return c.runtimeSvc.GetLogs(ctx)
}

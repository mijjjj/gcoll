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

	return c.runtimeSvc.GetOverview(ctx), nil
}

// Plugins 返回本地插件列表。
func (c *ControllerV1) Plugins(ctx context.Context, req *runtimev1.PluginsReq) (res *runtimev1.PluginsRes, err error) {
	_ = req

	return c.runtimeSvc.GetPlugins(ctx), nil
}

// Devices 返回设备列表。
func (c *ControllerV1) Devices(ctx context.Context, req *runtimev1.DevicesReq) (res *runtimev1.DevicesRes, err error) {
	_ = req

	return c.runtimeSvc.GetDevices(ctx), nil
}

// DevicePoints 返回指定设备的点位列表。
func (c *ControllerV1) DevicePoints(ctx context.Context, req *runtimev1.DevicePointsReq) (res *runtimev1.DevicePointsRes, err error) {
	return c.runtimeSvc.GetDevicePoints(ctx, req.DeviceId), nil
}

// Tasks 返回采集任务列表。
func (c *ControllerV1) Tasks(ctx context.Context, req *runtimev1.TasksReq) (res *runtimev1.TasksRes, err error) {
	_ = req

	return c.runtimeSvc.GetTasks(ctx), nil
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

	return c.runtimeSvc.GetLogs(ctx), nil
}

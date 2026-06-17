package runtime

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
	devicev1 "github.com/mijjjj/gcoll/api/device/v1"
	logv1 "github.com/mijjjj/gcoll/api/log/v1"
	pipelinev1 "github.com/mijjjj/gcoll/api/pipeline/v1"
	pluginv1 "github.com/mijjjj/gcoll/api/plugin/v1"
	pointcachev1 "github.com/mijjjj/gcoll/api/pointcache/v1"
	runtimev1 "github.com/mijjjj/gcoll/api/runtime/v1"
	targetv1 "github.com/mijjjj/gcoll/api/target/v1"
	taskv1 "github.com/mijjjj/gcoll/api/task/v1"
	"github.com/mijjjj/gcoll/internal/consts"
	"github.com/mijjjj/gcoll/internal/dao"
	"github.com/mijjjj/gcoll/internal/model/do"
	"github.com/mijjjj/gcoll/internal/model/entity"
	devicesvc "github.com/mijjjj/gcoll/internal/service/device"
	pluginhostsvc "github.com/mijjjj/gcoll/internal/service/pluginhost"
	pluginmgmtsvc "github.com/mijjjj/gcoll/internal/service/pluginmgmt"
	pointsvc "github.com/mijjjj/gcoll/internal/service/point"
)

// Service 提供运行时相关服务。
type Service struct {
	pluginSvc     *pluginmgmtsvc.Service
	deviceSvc     *devicesvc.Service
	pointSvc      *pointsvc.Service
	pluginHostSvc *pluginhostsvc.Service
}

// New 创建运行时服务。
func New() *Service {
	return &Service{
		pluginSvc:     pluginmgmtsvc.New(),
		deviceSvc:     devicesvc.New(),
		pointSvc:      pointsvc.New(),
		pluginHostSvc: pluginhostsvc.Instance(),
	}
}

// GetHealth 返回服务端运行时健康状态。
func (s *Service) GetHealth(ctx context.Context) *runtimev1.HealthRes {
	_ = ctx

	return &runtimev1.HealthRes{
		Status:    "ok",
		Service:   consts.ServiceName,
		Version:   consts.Version,
		Mode:      consts.RuntimeModeServer,
		CheckedAt: gtime.Now().Format("Y-m-d H:i:s"),
	}
}

// GetOverview 返回控制台工作台总览。
func (s *Service) GetOverview(ctx context.Context) (*runtimev1.OverviewRes, error) {
	devices, err := s.GetDevices(ctx)
	if err != nil {
		return nil, err
	}
	points := s.GetPointCache(ctx)
	tasks, err := s.GetTasks(ctx)
	if err != nil {
		return nil, err
	}
	plugins, err := s.GetPlugins(ctx)
	if err != nil {
		return nil, err
	}
	logs, err := s.GetLogs(ctx)
	if err != nil {
		return nil, err
	}

	var (
		runningCount int
		onlineCount  int
	)
	for _, plugin := range plugins.Items {
		if plugin.Status == "running" {
			runningCount++
		}
	}
	for _, device := range devices.Items {
		if device.Status == "online" {
			onlineCount++
		}
	}

	return &runtimev1.OverviewRes{
		Metrics: []commonv1.MetricItem{
			{Key: "runtime", Label: "运行时", Value: "运行中", Hint: "本机服务", Tone: "primary"},
			{Key: "devices", Label: "运行设备", Value: gconv.String(onlineCount), Hint: "共 " + gconv.String(len(devices.Items)) + " 台设备", Tone: "success"},
			{Key: "points", Label: "缓存点位", Value: gconv.String(len(points.Items)), Hint: "最新点位缓存", Tone: "primary"},
			{Key: "plugins", Label: "插件进程", Value: gconv.String(runningCount) + "/" + gconv.String(len(plugins.Items)), Hint: "本地插件模式", Tone: "warning"},
		},
		Runtime: commonv1.RuntimeStatus{
			Status:    "running",
			Service:   consts.ServiceName,
			Version:   consts.Version,
			Mode:      consts.RuntimeModeServer,
			CheckedAt: gtime.Now().Format("Y-m-d H:i:s"),
			ApiBase:   "http://127.0.0.1:8260",
			Database:  "SQLite 开发模式",
		},
		DataPlane: commonv1.DataPlaneStatus{
			QueueUsagePercent: 0,
			RuleHitPercent:    0,
			ForwardPercent:    0,
			Throughput:        "等待插件上报",
			Latency:           "未采集",
			Backpressure:      "正常",
		},
		Tasks:        tasks.Items,
		RecentEvents: logs.Items,
		PluginSummary: commonv1.PluginSummary{
			Running: runningCount,
			Total:   len(plugins.Items),
		},
		Network: commonv1.RuntimeDependency{
			Name:   "网络状态",
			Status: "offline",
			Detail: "离线模式",
		},
	}, nil
}

// GetPlugins 返回内存插件列表。
func (s *Service) GetPlugins(ctx context.Context) (*pluginv1.PluginsRes, error) {
	items, err := s.pluginSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return &pluginv1.PluginsRes{Items: items}, nil
}

// ImportPlugin 导入插件清单并返回插件列表项。
func (s *Service) ImportPlugin(ctx context.Context, packagePath string) (*pluginv1.ImportPluginRes, error) {
	item, err := s.pluginSvc.Import(ctx, packagePath)
	if err != nil {
		return nil, err
	}
	return &pluginv1.ImportPluginRes{Plugin: *item}, nil
}

// GetDevices 返回设备列表和分组。
func (s *Service) GetDevices(ctx context.Context) (*devicev1.DevicesRes, error) {
	return s.deviceSvc.List(ctx)
}

// CreateDeviceGroup 新增设备分组。
func (s *Service) CreateDeviceGroup(ctx context.Context, req *devicev1.CreateDeviceGroupReq) (*devicev1.CreateDeviceGroupRes, error) {
	group, err := s.deviceSvc.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	return &devicev1.CreateDeviceGroupRes{Group: *group}, nil
}

// DeleteDeviceGroup 删除空设备分组。
func (s *Service) DeleteDeviceGroup(ctx context.Context, groupId string) (*devicev1.DeleteDeviceGroupRes, error) {
	if err := s.deviceSvc.DeleteGroup(ctx, groupId); err != nil {
		return nil, err
	}
	return &devicev1.DeleteDeviceGroupRes{GroupId: groupId}, nil
}

// CreateDevice 新增设备和设备插件配置。
func (s *Service) CreateDevice(ctx context.Context, req *devicev1.CreateDeviceReq) (*devicev1.CreateDeviceRes, error) {
	item, err := s.deviceSvc.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &devicev1.CreateDeviceRes{Device: *item}, nil
}

// MoveDeviceToGroup 移动设备所属分组。
func (s *Service) MoveDeviceToGroup(ctx context.Context, req *devicev1.MoveDeviceToGroupReq) (*devicev1.MoveDeviceToGroupRes, error) {
	item, err := s.deviceSvc.MoveToGroup(ctx, req.DeviceId, req.GroupId)
	if err != nil {
		return nil, err
	}
	return &devicev1.MoveDeviceToGroupRes{Device: *item}, nil
}

// DeleteDevice 删除设备及其控制面关联数据。
func (s *Service) DeleteDevice(ctx context.Context, deviceId string) (*devicev1.DeleteDeviceRes, error) {
	if err := s.deviceSvc.Delete(ctx, deviceId); err != nil {
		return nil, err
	}
	return &devicev1.DeleteDeviceRes{DeviceId: deviceId}, nil
}

// GetDevicePluginConfigPage 返回指定设备的通用插件配置页数据。
func (s *Service) GetDevicePluginConfigPage(ctx context.Context, deviceId string) (*devicev1.DevicePluginConfigPageRes, error) {
	device, err := s.deviceSvc.Get(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	plugin, err := s.pluginHostSvc.Plugin(ctx, device.PluginId)
	if err != nil {
		return nil, err
	}
	points, err := s.GetDevicePoints(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	config, configured, err := s.deviceConfig(ctx, deviceId, device.PluginId)
	if err != nil {
		return nil, err
	}
	events, err := s.recentDeviceEvents(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	configPage, configHtml, configJs, err := s.pluginHostSvc.CustomConfigPageContent(ctx, device.PluginId)
	if err != nil {
		return nil, err
	}
	pointPage, pointHtml, pointJs, err := s.pluginHostSvc.CustomPointPageContent(ctx, device.PluginId)
	if err != nil {
		return nil, err
	}

	return &devicev1.DevicePluginConfigPageRes{
		Plugin: commonv1.PluginItem{
			Id:          plugin.Manifest.Id,
			Name:        plugin.Manifest.Name,
			Type:        plugin.Manifest.Type,
			Version:     plugin.Manifest.Version,
			Runtime:     plugin.Manifest.Runtime,
			Protocol:    plugin.Manifest.Protocol,
			Status:      plugin.Status,
			Permissions: plugin.Manifest.Permissions,
			UpdatedAt:   plugin.UpdatedAt,
		},
		Device: commonv1.DeviceItem{
			Id:          device.Id,
			Name:        device.Name,
			Code:        device.Code,
			GroupId:     device.GroupId,
			PluginId:    device.PluginId,
			PluginName:  plugin.Manifest.Name,
			Status:      device.Status,
			Enabled:     device.Enabled == 1,
			PointCount:  len(points.Items),
			ReportMode:  device.ReportMode,
			LastSeenAt:  displayTime(device.LastSeenAt, "尚未连接"),
			Description: device.Description,
		},
		Config:       config,
		ConfigSchema: emptyAnyMap(plugin.Manifest.ConfigSchema),
		CustomConfigPage: commonv1.PluginCustomConfigPage{
			Enabled: configPage.Enabled,
			Entry:   configPage.Entry,
			Script:  configPage.Script,
			Html:    configHtml,
			Js:      configJs,
		},
		CustomPointPage: commonv1.PluginCustomPointPage{
			Enabled: pointPage.Enabled,
			Entry:   pointPage.Entry,
			Script:  pointPage.Script,
			Html:    pointHtml,
			Js:      pointJs,
		},
		Configured:   configured,
		Points:       points.Items,
		RecentEvents: events,
		Operations: []commonv1.PluginOperation{
			{Key: "saveConfig", Label: "保存配置", Description: "保存当前设备的插件运行配置。", Enabled: true},
			{Key: "test", Label: "测试连接", Description: "通过南向插件 gRPC 服务测试连接。", Enabled: configured},
			{Key: "startTask", Label: "启动采集", Description: "启动设备采集任务并等待插件上报。", Enabled: configured && len(points.Items) > 0},
		},
		Warnings: []string{
			"设备配置和点位扩展由插件解释，宿主只保存通用 JSON 结构。",
			"采集明细不落库，插件只能通过标准 gRPC 服务向宿主提交记录。",
		},
	}, nil
}

// UpdateDevicePluginConfig 保存指定设备的通用插件配置。
func (s *Service) UpdateDevicePluginConfig(ctx context.Context, req *devicev1.UpdateDevicePluginConfigReq) (*devicev1.UpdateDevicePluginConfigRes, error) {
	config, err := s.deviceSvc.SavePluginConfig(ctx, req.DeviceId, req.Config)
	if err != nil {
		return nil, err
	}
	return &devicev1.UpdateDevicePluginConfigRes{Config: config}, nil
}

// TestDevicePluginConnection 测试指定设备的插件连接。
func (s *Service) TestDevicePluginConnection(ctx context.Context, deviceId string) (*devicev1.TestDevicePluginConnectionRes, error) {
	result, err := s.pluginHostSvc.TestConnection(ctx, deviceId)
	if result == nil {
		return nil, err
	}
	res := &devicev1.TestDevicePluginConnectionRes{
		Success:   result.Success,
		Message:   result.Message,
		LatencyMs: result.LatencyMs,
		TraceId:   result.TraceID,
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetDevicePoints 返回指定设备的通用点位表。
func (s *Service) GetDevicePoints(ctx context.Context, deviceId string) (*devicev1.DevicePointsRes, error) {
	return s.pointSvc.ListByDevice(ctx, deviceId)
}

// CreateDevicePoint 新增指定设备的通用点位。
func (s *Service) CreateDevicePoint(ctx context.Context, req *devicev1.CreateDevicePointReq) (*devicev1.CreateDevicePointRes, error) {
	item, err := s.pointSvc.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &devicev1.CreateDevicePointRes{Point: *item}, nil
}

// UpdateDevicePoints 保存指定设备的完整点位表。
func (s *Service) UpdateDevicePoints(ctx context.Context, req *devicev1.UpdateDevicePointsReq) (*devicev1.UpdateDevicePointsRes, error) {
	result, err := s.pointSvc.ReplaceByDevice(ctx, req.DeviceId, req.Items)
	if err != nil {
		return nil, err
	}
	return &devicev1.UpdateDevicePointsRes{Items: result.Items}, nil
}

// GetTasks 返回采集任务列表。
func (s *Service) GetTasks(ctx context.Context) (*taskv1.TasksRes, error) {
	var (
		tasks   []entity.CollectionTasks
		devices []entity.Devices
		points  []entity.DevicePoints
	)
	if err := dao.CollectionTasks.Ctx(ctx).OrderDesc(dao.CollectionTasks.Columns().UpdatedAt).Scan(&tasks); err != nil {
		return nil, gerror.Wrap(err, "读取采集任务失败")
	}
	if err := dao.Devices.Ctx(ctx).Scan(&devices); err != nil {
		return nil, gerror.Wrap(err, "读取任务设备失败")
	}
	if err := dao.DevicePoints.Ctx(ctx).Scan(&points); err != nil {
		return nil, gerror.Wrap(err, "读取任务点位数量失败")
	}

	deviceNames := map[string]string{}
	for _, device := range devices {
		deviceNames[device.Id] = device.Name
	}
	pluginNames, err := s.pluginHostSvc.PluginNameMap(ctx)
	if err != nil {
		return nil, err
	}
	pointCounts := map[string]int{}
	for _, point := range points {
		pointCounts[point.DeviceId]++
	}

	items := make([]commonv1.TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, commonv1.TaskSummary{
			Id:              task.Id,
			Name:            task.Name,
			DeviceId:        task.DeviceId,
			DeviceName:      deviceNames[task.DeviceId],
			SouthPluginName: pluginNames[task.SouthPluginId],
			PointCount:      pointCounts[task.DeviceId],
			ReportMode:      task.ReportMode,
			Status:          task.Status,
			Rate:            task.Rate,
			RuleHitRate:     task.RuleHitRate,
			LastCollectedAt: displayTime(task.LastCollectedAt, "尚未采集"),
		})
	}
	return &taskv1.TasksRes{Items: items}, nil
}

// StartDeviceCollectionTask 启动设备默认采集任务，缺失任务时创建控制面任务。
func (s *Service) StartDeviceCollectionTask(ctx context.Context, deviceId string) (*taskv1.CollectionTaskActionRes, error) {
	task, err := s.ensureDeviceTask(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	return s.StartCollectionTask(ctx, task.Id)
}

// StartCollectionTask 启动指定采集任务。
func (s *Service) StartCollectionTask(ctx context.Context, taskId string) (*taskv1.CollectionTaskActionRes, error) {
	if err := s.pluginHostSvc.StartTask(ctx, taskId); err != nil {
		return nil, err
	}
	return s.taskActionResult(ctx, taskId)
}

// StopCollectionTask 停止指定采集任务。
func (s *Service) StopCollectionTask(ctx context.Context, taskId string) (*taskv1.CollectionTaskActionRes, error) {
	if err := s.pluginHostSvc.StopTask(ctx, taskId); err != nil {
		return nil, err
	}
	return s.taskActionResult(ctx, taskId)
}

// GetPointCache 返回最新点位缓存。
func (s *Service) GetPointCache(ctx context.Context) *pointcachev1.PointCacheRes {
	_ = ctx
	latest := s.pluginHostSvc.LatestValues()
	items := make([]commonv1.PointCacheItem, 0, len(latest))
	for _, item := range latest {
		items = append(items, commonv1.PointCacheItem{
			PointId:   item.PointID,
			DeviceId:  item.DeviceID,
			PointName: item.PointName,
			Value:     item.Value,
			Quality:   item.Quality,
			Changed:   item.Changed,
			UpdatedAt: item.UpdatedAt,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].DeviceId != items[j].DeviceId {
			return items[i].DeviceId < items[j].DeviceId
		}
		return items[i].PointId < items[j].PointId
	})
	return &pointcachev1.PointCacheRes{Items: items}
}

// GetPipelineRules 返回规则过滤列表。
func (s *Service) GetPipelineRules(ctx context.Context) *pipelinev1.PipelineRulesRes {
	_ = ctx

	return &pipelinev1.PipelineRulesRes{
		Items: []commonv1.PipelineRuleItem{
			{Id: "rule-good-quality", Name: "仅转发良好质量数据", Enabled: true, Expression: "quality == \"good\"", Matched: 0, TargetCount: 0, UpdatedAt: "尚未启用"},
			{Id: "rule-change-only", Name: "变化值进入北向链路", Enabled: true, Expression: "report_mode == \"change\" && changed == true", Matched: 0, TargetCount: 0, UpdatedAt: "尚未启用"},
		},
	}
}

// GetTargets 返回北向转发目标列表。
func (s *Service) GetTargets(ctx context.Context) *targetv1.TargetsRes {
	_ = ctx

	return &targetv1.TargetsRes{
		Items: []commonv1.ForwardTargetItem{},
	}
}

// GetLogs 返回运行日志和事件列表。
func (s *Service) GetLogs(ctx context.Context) (*logv1.LogsRes, error) {
	var events []entity.RuntimeEvents
	if err := dao.RuntimeEvents.Ctx(ctx).OrderDesc(dao.RuntimeEvents.Columns().Time).Limit(50).Scan(&events); err != nil {
		return nil, gerror.Wrap(err, "读取运行事件失败")
	}
	return &logv1.LogsRes{Items: runtimeEvents(events)}, nil
}

func (s *Service) ensureDeviceTask(ctx context.Context, deviceId string) (*entity.CollectionTasks, error) {
	device, err := s.deviceSvc.Get(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	var task entity.CollectionTasks
	if err := dao.CollectionTasks.Ctx(ctx).Where(do.CollectionTasks{
		DeviceId:      deviceId,
		SouthPluginId: device.PluginId,
	}).OrderAsc(dao.CollectionTasks.Columns().CreatedAt).Scan(&task); err != nil {
		return nil, gerror.Wrapf(err, "读取设备采集任务失败: %s", deviceId)
	}
	if task.Id != "" {
		return &task, nil
	}
	plugin, err := s.pluginHostSvc.Plugin(ctx, device.PluginId)
	if err != nil {
		return nil, err
	}
	task = entity.CollectionTasks{
		Id:            "task-" + guid.S(),
		Name:          device.Name + " 采集任务",
		DeviceId:      device.Id,
		SouthPluginId: device.PluginId,
		ReportMode:    device.ReportMode,
		Status:        "stopped",
		RuleHitRate:   "0%",
		Rate:          "0 条/秒",
	}
	if _, err := dao.CollectionTasks.Ctx(ctx).Data(do.CollectionTasks{
		Id:            task.Id,
		Name:          task.Name,
		DeviceId:      task.DeviceId,
		SouthPluginId: task.SouthPluginId,
		ReportMode:    task.ReportMode,
		Status:        task.Status,
		RuleHitRate:   task.RuleHitRate,
		Rate:          task.Rate,
	}).Insert(); err != nil {
		return nil, gerror.Wrapf(err, "创建设备采集任务失败: %s", deviceId)
	}
	_ = plugin
	return &task, nil
}

func (s *Service) taskActionResult(ctx context.Context, taskId string) (*taskv1.CollectionTaskActionRes, error) {
	tasks, err := s.GetTasks(ctx)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks.Items {
		if task.Id == taskId {
			return &taskv1.CollectionTaskActionRes{Task: task}, nil
		}
	}
	return nil, gerror.Newf("采集任务不存在: %s", taskId)
}

func (s *Service) deviceConfig(ctx context.Context, deviceId string, pluginId string) (map[string]any, bool, error) {
	var config entity.PluginDeviceConfigs
	if err := dao.PluginDeviceConfigs.Ctx(ctx).Where(do.PluginDeviceConfigs{
		DeviceId: deviceId,
		PluginId: pluginId,
		Active:   1,
	}).Scan(&config); err != nil {
		return nil, false, gerror.Wrapf(err, "读取设备插件配置失败: %s", deviceId)
	}
	if config.Id == "" || strings.TrimSpace(config.ConfigJson) == "" {
		return map[string]any{}, false, nil
	}
	values := map[string]any{}
	if err := json.Unmarshal([]byte(config.ConfigJson), &values); err != nil {
		return nil, false, gerror.Wrapf(err, "解析设备插件配置失败: %s", deviceId)
	}
	return values, true, nil
}

func (s *Service) recentDeviceEvents(ctx context.Context, deviceId string) ([]commonv1.RuntimeEvent, error) {
	var events []entity.RuntimeEvents
	if err := dao.RuntimeEvents.Ctx(ctx).
		Where(do.RuntimeEvents{DeviceId: deviceId}).
		OrderDesc(dao.RuntimeEvents.Columns().Time).
		Limit(20).
		Scan(&events); err != nil {
		return nil, gerror.Wrapf(err, "读取设备运行事件失败: %s", deviceId)
	}
	return runtimeEvents(events), nil
}

func runtimeEvents(events []entity.RuntimeEvents) []commonv1.RuntimeEvent {
	items := make([]commonv1.RuntimeEvent, 0, len(events))
	for _, event := range events {
		items = append(items, commonv1.RuntimeEvent{
			Id:       event.Id,
			Time:     event.Time,
			Level:    event.Level,
			Source:   event.Source,
			PluginId: event.PluginId,
			DeviceId: event.DeviceId,
			TaskId:   event.TaskId,
			Message:  event.Message,
			TraceId:  event.TraceId,
		})
	}
	return items
}

func displayTime(value string, empty string) string {
	if value == "" {
		return empty
	}
	return value
}

func emptyAnyMap(values map[string]any) map[string]any {
	if values == nil {
		return map[string]any{}
	}
	return values
}

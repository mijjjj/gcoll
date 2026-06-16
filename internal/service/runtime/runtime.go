package runtime

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	runtimev1 "github.com/mijjjj/gcoll/api/runtime/v1"
	"github.com/mijjjj/gcoll/internal/consts"
	"github.com/mijjjj/gcoll/internal/dao"
	"github.com/mijjjj/gcoll/internal/model/do"
	"github.com/mijjjj/gcoll/internal/model/entity"
	devicesvc "github.com/mijjjj/gcoll/internal/service/device"
	pluginmgmtsvc "github.com/mijjjj/gcoll/internal/service/pluginmgmt"
	pointsvc "github.com/mijjjj/gcoll/internal/service/point"
)

const modbusPluginId = "com.gcoll.modbus-tcp"

// Service 提供运行时相关服务。
type Service struct {
	pluginSvc *pluginmgmtsvc.Service
	deviceSvc *devicesvc.Service
	pointSvc  *pointsvc.Service
}

// New 创建运行时服务。
func New() *Service {
	return &Service{
		pluginSvc: pluginmgmtsvc.New(),
		deviceSvc: devicesvc.New(),
		pointSvc:  pointsvc.New(),
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
		Metrics: []runtimev1.MetricItem{
			{Key: "runtime", Label: "运行时", Value: "运行中", Hint: "本机服务", Tone: "primary"},
			{Key: "devices", Label: "运行设备", Value: gconv.String(onlineCount), Hint: "共 " + gconv.String(len(devices.Items)) + " 台设备", Tone: "success"},
			{Key: "points", Label: "启用点位", Value: gconv.String(len(points.Items)), Hint: "最新缓存 " + gconv.String(len(points.Items)) + " 条", Tone: "primary"},
			{Key: "plugins", Label: "插件进程", Value: gconv.String(runningCount) + "/" + gconv.String(len(plugins.Items)), Hint: "本地插件模式", Tone: "warning"},
		},
		Runtime: runtimev1.RuntimeStatus{
			Status:    "running",
			Service:   consts.ServiceName,
			Version:   consts.Version,
			Mode:      consts.RuntimeModeServer,
			CheckedAt: gtime.Now().Format("Y-m-d H:i:s"),
			ApiBase:   "http://127.0.0.1:8260",
			Database:  "SQLite 开发模式",
		},
		DataPlane: runtimev1.DataPlaneStatus{
			QueueUsagePercent: 18,
			RuleHitPercent:    72,
			ForwardPercent:    64,
			Throughput:        "128 条/秒",
			Latency:           "42 ms",
			Backpressure:      "正常",
		},
		Tasks:        tasks.Items,
		RecentEvents: logs.Items,
		PluginSummary: runtimev1.PluginSummary{
			Running: runningCount,
			Total:   len(plugins.Items),
		},
		Network: runtimev1.RuntimeDependency{
			Name:   "网络状态",
			Status: "offline",
			Detail: "离线模式",
		},
	}, nil
}

// GetPlugins 返回数据库中的插件列表。
func (s *Service) GetPlugins(ctx context.Context) (*runtimev1.PluginsRes, error) {
	items, err := s.pluginSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return &runtimev1.PluginsRes{Items: items}, nil
}

// ImportPlugin 导入插件清单并返回插件列表项。
func (s *Service) ImportPlugin(ctx context.Context, packagePath string) (*runtimev1.ImportPluginRes, error) {
	item, err := s.pluginSvc.Import(ctx, packagePath)
	if err != nil {
		return nil, err
	}
	return &runtimev1.ImportPluginRes{Plugin: *item}, nil
}

// GetDevices 返回设备列表和分组。
func (s *Service) GetDevices(ctx context.Context) (*runtimev1.DevicesRes, error) {
	return s.deviceSvc.List(ctx)
}

// CreateDevice 新增设备和设备插件配置。
func (s *Service) CreateDevice(ctx context.Context, req *runtimev1.CreateDeviceReq) (*runtimev1.CreateDeviceRes, error) {
	item, err := s.deviceSvc.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &runtimev1.CreateDeviceRes{Device: *item}, nil
}

// GetDevicePoints 返回指定设备的通用点位表。
func (s *Service) GetDevicePoints(ctx context.Context, deviceId string) (*runtimev1.DevicePointsRes, error) {
	return s.pointSvc.ListByDevice(ctx, deviceId)
}

// CreateDevicePoint 新增指定设备的通用点位。
func (s *Service) CreateDevicePoint(ctx context.Context, req *runtimev1.CreateDevicePointReq) (*runtimev1.CreateDevicePointRes, error) {
	item, err := s.pointSvc.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &runtimev1.CreateDevicePointRes{Point: *item}, nil
}

// GetTasks 返回采集任务列表。
func (s *Service) GetTasks(ctx context.Context) (*runtimev1.TasksRes, error) {
	var (
		tasks   []entity.CollectionTasks
		devices []entity.Devices
		plugins []entity.Plugins
		points  []entity.DevicePoints
	)
	if err := dao.CollectionTasks.Ctx(ctx).OrderDesc(dao.CollectionTasks.Columns().UpdatedAt).Scan(&tasks); err != nil {
		return nil, gerror.Wrap(err, "读取采集任务失败")
	}
	if err := dao.Devices.Ctx(ctx).Scan(&devices); err != nil {
		return nil, gerror.Wrap(err, "读取任务设备失败")
	}
	if err := dao.Plugins.Ctx(ctx).Scan(&plugins); err != nil {
		return nil, gerror.Wrap(err, "读取任务插件失败")
	}
	if err := dao.DevicePoints.Ctx(ctx).Scan(&points); err != nil {
		return nil, gerror.Wrap(err, "读取任务点位数量失败")
	}

	deviceNames := map[string]string{}
	for _, device := range devices {
		deviceNames[device.Id] = device.Name
	}
	pluginNames := map[string]string{}
	for _, plugin := range plugins {
		pluginNames[plugin.Id] = plugin.Name
	}
	pointCounts := map[string]int{}
	for _, point := range points {
		pointCounts[point.DeviceId]++
	}

	items := make([]runtimev1.TaskSummary, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, runtimev1.TaskSummary{
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
	return &runtimev1.TasksRes{Items: items}, nil
}

// GetPointCache 返回最新点位缓存。
func (s *Service) GetPointCache(ctx context.Context) *runtimev1.PointCacheRes {
	_ = ctx

	return &runtimev1.PointCacheRes{
		Items: []runtimev1.PointCacheItem{
			{PointId: "pt-temperature", DeviceId: "dev-edge-gw-a01", PointName: "TEMP_01", Value: "25.6 ℃", Quality: "good", Changed: true, UpdatedAt: "2026-06-16 10:30:18.123"},
			{PointId: "pt-pressure", DeviceId: "dev-edge-gw-a01", PointName: "PRESS_01", Value: "101.3 kPa", Quality: "good", Changed: false, UpdatedAt: "2026-06-16 10:30:18.120"},
			{PointId: "pt-motor-state", DeviceId: "dev-edge-gw-a01", PointName: "MOTOR_RUN", Value: "ON", Quality: "good", Changed: true, UpdatedAt: "2026-06-16 10:30:18.118"},
			{PointId: "pt-energy", DeviceId: "dev-edge-gw-a01", PointName: "ENERGY_TOTAL", Value: "1288.4 kWh", Quality: "uncertain", Changed: false, UpdatedAt: "2026-06-16 10:30:18.110"},
		},
	}
}

// GetPipelineRules 返回规则过滤列表。
func (s *Service) GetPipelineRules(ctx context.Context) *runtimev1.PipelineRulesRes {
	_ = ctx

	return &runtimev1.PipelineRulesRes{
		Items: []runtimev1.PipelineRuleItem{
			{Id: "rule-good-quality", Name: "仅转发良好质量数据", Enabled: true, Expression: "quality == \"good\"", Matched: 1860, TargetCount: 1, UpdatedAt: "2026-06-15 09:20:00"},
			{Id: "rule-change-only", Name: "变化值进入北向链路", Enabled: true, Expression: "report_mode == \"change\" && changed == true", Matched: 940, TargetCount: 1, UpdatedAt: "2026-06-15 09:22:00"},
		},
	}
}

// GetTargets 返回北向转发目标列表。
func (s *Service) GetTargets(ctx context.Context) *runtimev1.TargetsRes {
	_ = ctx

	return &runtimev1.TargetsRes{
		Items: []runtimev1.ForwardTargetItem{
			{Id: "target-http-demo", Name: "演示 HTTP 接收端", PluginName: "HTTP 北向转发", Status: "running", Endpoint: "https://example.internal/ingest", UpdatedAt: "2026-06-15 09:28:00"},
			{Id: "target-http-standby", Name: "备用 HTTP 接收端", PluginName: "HTTP 北向转发", Status: "stopped", Endpoint: "https://backup.example.internal/ingest", LastError: "未启用", UpdatedAt: "2026-06-15 08:50:00"},
		},
	}
}

// GetLogs 返回运行日志和事件列表。
func (s *Service) GetLogs(ctx context.Context) (*runtimev1.LogsRes, error) {
	var events []entity.RuntimeEvents
	if err := dao.RuntimeEvents.Ctx(ctx).OrderDesc(dao.RuntimeEvents.Columns().Time).Limit(50).Scan(&events); err != nil {
		return nil, gerror.Wrap(err, "读取运行事件失败")
	}
	items := make([]runtimev1.RuntimeEvent, 0, len(events))
	for _, event := range events {
		items = append(items, runtimev1.RuntimeEvent{
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
	return &runtimev1.LogsRes{Items: items}, nil
}

// GetModbusTcpDeviceConfigPage 返回指定设备的 Modbus TCP 协议配置页数据。
func (s *Service) GetModbusTcpDeviceConfigPage(ctx context.Context, deviceId string) (*runtimev1.ModbusTcpDeviceConfigPageRes, error) {
	device, err := s.deviceSvc.Get(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	if device.PluginId != modbusPluginId {
		return nil, gerror.Newf("设备未使用 Modbus TCP 插件: %s", deviceId)
	}

	plugin, err := s.pluginById(ctx, modbusPluginId)
	if err != nil {
		return nil, err
	}
	modbusConfig, err := s.modbusConfig(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	points, err := s.modbusPoints(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	readPlan := buildReadPlan(points)
	debugLogs, err := s.modbusDebugLogs(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	return &runtimev1.ModbusTcpDeviceConfigPageRes{
		Plugin: runtimev1.PluginItem{
			Id:          plugin.Id,
			Name:        plugin.Name,
			Type:        plugin.Type,
			Version:     plugin.ActiveVersion,
			Runtime:     plugin.Runtime,
			Protocol:    plugin.Protocol,
			Status:      plugin.Status,
			Permissions: []string{"network.outbound", "config.read", "runtime.events"},
			UpdatedAt:   plugin.UpdatedAt,
		},
		Device: runtimev1.DeviceItem{
			Id:          device.Id,
			Name:        device.Name,
			Code:        device.Code,
			GroupId:     device.GroupId,
			PluginId:    device.PluginId,
			PluginName:  plugin.Name,
			Status:      device.Status,
			Enabled:     device.Enabled == 1,
			PointCount:  len(points),
			ReportMode:  device.ReportMode,
			LastSeenAt:  displayTime(device.LastSeenAt, "尚未连接"),
			Description: device.Description,
		},
		Config:    *modbusConfig,
		ReadPlan:  readPlan,
		Points:    points,
		DebugLogs: debugLogs,
		Operations: []runtimev1.ModbusTcpOperation{
			{Key: "test", Label: "测试连接", Description: "建立 TCP 连接并执行一次轻量读取。", Enabled: true},
			{Key: "readOnce", Label: "读取一次", Description: "按当前点位表执行一次读取并输出调试日志。", Enabled: true},
			{Key: "toggleDebug", Label: "调试模式", Description: "采集请求、响应耗时和原始报文摘要。", Enabled: true},
			{Key: "writePoint", Label: "写入点位", Description: "仅允许写入线圈和保持寄存器点位。", Enabled: true},
		},
		Warnings: []string{
			"采集明细不落库，调试日志只保存最近窗口并由宿主统一收集。",
			"离散输入和输入寄存器为只读区域，不能配置写入模式。",
		},
	}, nil
}

func (s *Service) pluginById(ctx context.Context, pluginId string) (*entity.Plugins, error) {
	var plugin entity.Plugins
	if err := dao.Plugins.Ctx(ctx).Where(do.Plugins{Id: pluginId}).Scan(&plugin); err != nil {
		return nil, gerror.Wrapf(err, "读取插件失败: %s", pluginId)
	}
	if plugin.Id == "" {
		return nil, gerror.Newf("插件不存在: %s", pluginId)
	}
	return &plugin, nil
}

func (s *Service) modbusConfig(ctx context.Context, deviceId string) (*runtimev1.ModbusTcpDeviceConfig, error) {
	var config entity.ModbusTcpDeviceProfiles
	if err := dao.ModbusTcpDeviceProfiles.Ctx(ctx).Where(do.ModbusTcpDeviceProfiles{
		DeviceId: deviceId,
		PluginId: modbusPluginId,
	}).Scan(&config); err != nil {
		return nil, gerror.Wrapf(err, "读取 Modbus TCP 设备配置失败: %s", deviceId)
	}
	if config.Id == "" {
		return nil, gerror.Newf("设备缺少 Modbus TCP 配置: %s", deviceId)
	}
	return &runtimev1.ModbusTcpDeviceConfig{
		Host:             config.Host,
		Port:             config.Port,
		UnitId:           config.UnitId,
		TimeoutMs:        config.TimeoutMs,
		PollIntervalMs:   config.PollIntervalMs,
		ReportMode:       config.ReportMode,
		DebugEnabled:     config.DebugEnabled == 1,
		MaxCoilBatch:     config.MaxCoilBatch,
		MaxRegisterBatch: config.MaxRegisterBatch,
		LowLatencyMs:     config.LowLatencyMs,
		HighLatencyMs:    config.HighLatencyMs,
	}, nil
}

func (s *Service) modbusPoints(ctx context.Context, deviceId string) ([]runtimev1.ModbusTcpPoint, error) {
	var (
		points   []entity.DevicePoints
		profiles []entity.ModbusTcpPointProfiles
	)
	if err := dao.DevicePoints.Ctx(ctx).Where(do.DevicePoints{DeviceId: deviceId, PluginId: modbusPluginId}).Scan(&points); err != nil {
		return nil, gerror.Wrapf(err, "读取设备点位失败: %s", deviceId)
	}
	if err := dao.ModbusTcpPointProfiles.Ctx(ctx).Where(do.ModbusTcpPointProfiles{DeviceId: deviceId, PluginId: modbusPluginId}).Scan(&profiles); err != nil {
		return nil, gerror.Wrapf(err, "读取 Modbus TCP 点位配置失败: %s", deviceId)
	}
	pointMap := make(map[string]entity.DevicePoints, len(points))
	for _, point := range points {
		pointMap[point.Id] = point
	}
	items := make([]runtimev1.ModbusTcpPoint, 0, len(profiles))
	for _, profile := range profiles {
		point, ok := pointMap[profile.PointId]
		if !ok {
			return nil, gerror.Newf("Modbus TCP 点位缺少通用点位记录: %s", profile.PointId)
		}
		items = append(items, runtimev1.ModbusTcpPoint{
			Id:          point.Id,
			Name:        point.Name,
			Area:        profile.Area,
			Address:     profile.Address,
			Quantity:    profile.Quantity,
			ValueType:   profile.ValueType,
			Mode:        profile.Mode,
			ReportMode:  profile.ReportMode,
			Enabled:     point.Enabled == 1 && profile.Enabled == 1,
			ByteOrder:   profile.ByteOrder,
			WordOrder:   profile.WordOrder,
			Scale:       gconv.String(profile.Scale),
			Current:     currentValue(point.Id),
			Quality:     currentQuality(point.Id),
			LastReadAt:  currentTime(point.Id),
			Description: point.Description,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if areaOrder(items[i].Area) != areaOrder(items[j].Area) {
			return areaOrder(items[i].Area) < areaOrder(items[j].Area)
		}
		if items[i].Address != items[j].Address {
			return items[i].Address < items[j].Address
		}
		return items[i].Id < items[j].Id
	})
	return items, nil
}

func buildReadPlan(points []runtimev1.ModbusTcpPoint) []runtimev1.ModbusTcpReadBlock {
	blocks := make([]runtimev1.ModbusTcpReadBlock, 0, len(points))
	for _, point := range points {
		if !point.Enabled || point.Mode != "read" {
			continue
		}
		blocks = append(blocks, runtimev1.ModbusTcpReadBlock{
			Area:     point.Area,
			Start:    point.Address,
			Quantity: point.Quantity,
			PointIds: []string{point.Id},
		})
	}
	return blocks
}

func (s *Service) modbusDebugLogs(ctx context.Context, deviceId string) ([]runtimev1.ModbusTcpDebugLog, error) {
	var logs []entity.ModbusTcpDebugLogs
	if err := dao.ModbusTcpDebugLogs.Ctx(ctx).
		Where(do.ModbusTcpDebugLogs{DeviceId: deviceId}).
		OrderDesc(dao.ModbusTcpDebugLogs.Columns().CreatedAt).
		Limit(20).
		Scan(&logs); err != nil {
		return nil, gerror.Wrapf(err, "读取 Modbus TCP 调试日志失败: %s", deviceId)
	}
	items := make([]runtimev1.ModbusTcpDebugLog, 0, len(logs))
	for _, log := range logs {
		items = append(items, runtimev1.ModbusTcpDebugLog{
			Time:    log.CreatedAt,
			Level:   log.Level,
			Message: log.Message,
			TraceId: log.TraceId,
			Area:    log.Area,
			Address: gconv.String(log.Address),
			CostMs:  log.LatencyMs,
			RawHex:  log.RawHex,
		})
	}
	return items, nil
}

func areaOrder(area string) int {
	switch area {
	case "coil":
		return 1
	case "discrete_input":
		return 2
	case "holding_register":
		return 3
	case "input_register":
		return 4
	default:
		return 99
	}
}

func displayTime(value string, empty string) string {
	if value == "" {
		return empty
	}
	return value
}

func currentValue(pointId string) string {
	values := map[string]string{
		"pt-temperature": "25.6 ℃",
		"pt-pressure":    "101.3 kPa",
		"pt-motor-state": "ON",
		"pt-energy":      "1288.4 kWh",
		"pt-speed-set":   "960 rpm",
		"pt-emergency":   "OFF",
	}
	return values[pointId]
}

func currentQuality(pointId string) string {
	if pointId == "pt-energy" {
		return "uncertain"
	}
	if currentValue(pointId) == "" {
		return ""
	}
	return "good"
}

func currentTime(pointId string) string {
	times := map[string]string{
		"pt-temperature": "2026-06-16 10:30:18.123",
		"pt-pressure":    "2026-06-16 10:30:18.120",
		"pt-motor-state": "2026-06-16 10:30:18.118",
		"pt-energy":      "2026-06-16 10:30:18.110",
		"pt-speed-set":   "2026-06-16 10:29:55.010",
		"pt-emergency":   "2026-06-16 10:30:18.106",
	}
	return times[pointId]
}

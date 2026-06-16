package runtime

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	runtimev1 "github.com/mijjjj/gcoll/api/runtime/v1"
	"github.com/mijjjj/gcoll/internal/consts"
)

// Service 提供运行时相关服务。
type Service struct{}

// New 创建运行时服务。
func New() *Service {
	return &Service{}
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
func (s *Service) GetOverview(ctx context.Context) *runtimev1.OverviewRes {
	_ = ctx

	var (
		devices      = s.GetDevices(ctx)
		points       = s.GetPointCache(ctx)
		tasks        = s.GetTasks(ctx)
		plugins      = s.GetPlugins(ctx)
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
		RecentEvents: s.GetLogs(ctx).Items,
		PluginSummary: runtimev1.PluginSummary{
			Running: runningCount,
			Total:   len(plugins.Items),
		},
		Network: runtimev1.RuntimeDependency{
			Name:   "网络状态",
			Status: "offline",
			Detail: "离线模式",
		},
	}
}

// GetPlugins 返回本地插件列表。
func (s *Service) GetPlugins(ctx context.Context) *runtimev1.PluginsRes {
	_ = ctx

	return &runtimev1.PluginsRes{
		Items: []runtimev1.PluginItem{
			{
				Id:          "com.gcoll.http-collector",
				Name:        "HTTP 轮询采集",
				Type:        "southbound",
				Version:     "0.1.0",
				Runtime:     "process",
				Protocol:    "grpc",
				Status:      "running",
				Permissions: []string{"network.outbound", "config.read"},
				UpdatedAt:   "2026-06-15 09:10:00",
			},
			{
				Id:          "com.gcoll.http-forwarder",
				Name:        "HTTP 北向转发",
				Type:        "northbound",
				Version:     "0.1.0",
				Runtime:     "process",
				Protocol:    "grpc",
				Status:      "running",
				Permissions: []string{"network.outbound", "secret.read"},
				UpdatedAt:   "2026-06-15 09:12:00",
			},
		},
	}
}

// GetDevices 返回设备列表和分组。
func (s *Service) GetDevices(ctx context.Context) *runtimev1.DevicesRes {
	_ = ctx

	return &runtimev1.DevicesRes{
		Groups: []runtimev1.DeviceGroup{
			{Id: "edge", Name: "边缘现场", Count: 2},
			{Id: "test", Name: "测试分组", Count: 0},
		},
		Items: []runtimev1.DeviceItem{
			{
				Id:          "dev-edge-gw-a01",
				Name:        "边缘网关 A01",
				Code:        "DEV-EDGE-A01",
				GroupId:     "edge",
				PluginId:    "com.gcoll.http-collector",
				PluginName:  "HTTP 轮询采集",
				Status:      "online",
				Enabled:     true,
				PointCount:  4,
				ReportMode:  "change",
				LastSeenAt:  "2026-06-15 09:30:18",
				Description: "用于验证 MVP 采集、过滤、转发闭环的本地网关。",
			},
			{
				Id:          "dev-sim-line-b02",
				Name:        "模拟产线 B02",
				Code:        "DEV-SIM-B02",
				GroupId:     "edge",
				PluginId:    "com.gcoll.http-collector",
				PluginName:  "HTTP 轮询采集",
				Status:      "offline",
				Enabled:     false,
				PointCount:  0,
				ReportMode:  "change",
				LastSeenAt:  "尚未连接",
				Description: "保留给现场调试的模拟设备。",
			},
		},
	}
}

// GetDevicePoints 返回指定设备的通用点位表。
func (s *Service) GetDevicePoints(ctx context.Context, deviceId string) *runtimev1.DevicePointsRes {
	_ = ctx

	points := []runtimev1.PointItem{
		{
			Id:          "pt-temperature",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.http-collector",
			Name:        "TEMP_01",
			Description: "环境温度",
			Address:     "$.temperature",
			ValueType:   "float",
			Unit:        "℃",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "environment"},
		},
		{
			Id:          "pt-pressure",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.http-collector",
			Name:        "PRESS_01",
			Description: "管线压力",
			Address:     "$.pressure",
			ValueType:   "float",
			Unit:        "kPa",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "process"},
		},
		{
			Id:          "pt-motor-state",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.http-collector",
			Name:        "MOTOR_RUN",
			Description: "电机运行状态",
			Address:     "$.motor.running",
			ValueType:   "bool",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "status"},
		},
		{
			Id:          "pt-energy",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.http-collector",
			Name:        "ENERGY_TOTAL",
			Description: "累计能耗",
			Address:     "$.energy.total",
			ValueType:   "float",
			Unit:        "kWh",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "meter"},
		},
	}

	filtered := make([]runtimev1.PointItem, 0, len(points))
	for _, point := range points {
		if point.DeviceId == deviceId {
			filtered = append(filtered, point)
		}
	}

	return &runtimev1.DevicePointsRes{Items: filtered}
}

// GetTasks 返回采集任务列表。
func (s *Service) GetTasks(ctx context.Context) *runtimev1.TasksRes {
	_ = ctx

	return &runtimev1.TasksRes{
		Items: []runtimev1.TaskSummary{
			{
				Id:              "task-http-a01",
				Name:            "样例 HTTP 采集链路",
				DeviceId:        "dev-edge-gw-a01",
				DeviceName:      "边缘网关 A01",
				SouthPluginName: "HTTP 轮询采集",
				PointCount:      4,
				ReportMode:      "change",
				Status:          "running",
				Rate:            "128 条/秒",
				RuleHitRate:     "72%",
				LastCollectedAt: "2026-06-15 09:30:18",
			},
		},
	}
}

// GetPointCache 返回最新点位缓存。
func (s *Service) GetPointCache(ctx context.Context) *runtimev1.PointCacheRes {
	_ = ctx

	return &runtimev1.PointCacheRes{
		Items: []runtimev1.PointCacheItem{
			{PointId: "pt-temperature", DeviceId: "dev-edge-gw-a01", PointName: "TEMP_01", Value: "25.6 ℃", Quality: "good", Changed: true, UpdatedAt: "2026-06-15 09:30:18.123"},
			{PointId: "pt-pressure", DeviceId: "dev-edge-gw-a01", PointName: "PRESS_01", Value: "101.3 kPa", Quality: "good", Changed: false, UpdatedAt: "2026-06-15 09:30:18.120"},
			{PointId: "pt-motor-state", DeviceId: "dev-edge-gw-a01", PointName: "MOTOR_RUN", Value: "ON", Quality: "good", Changed: true, UpdatedAt: "2026-06-15 09:30:18.118"},
			{PointId: "pt-energy", DeviceId: "dev-edge-gw-a01", PointName: "ENERGY_TOTAL", Value: "1288.4 kWh", Quality: "uncertain", Changed: false, UpdatedAt: "2026-06-15 09:30:18.110"},
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
func (s *Service) GetLogs(ctx context.Context) *runtimev1.LogsRes {
	_ = ctx

	return &runtimev1.LogsRes{
		Items: []runtimev1.RuntimeEvent{
			{Id: "evt-001", Time: "2026-06-15 09:30:18", Level: "INFO", Source: "collector", PluginId: "com.gcoll.http-collector", DeviceId: "dev-edge-gw-a01", TaskId: "task-http-a01", Message: "已接收 128 条采集记录并写入内存缓冲。", TraceId: "trace-demo-001"},
			{Id: "evt-002", Time: "2026-06-15 09:30:18", Level: "INFO", Source: "pipeline", DeviceId: "dev-edge-gw-a01", TaskId: "task-http-a01", Message: "规则过滤命中 92 条记录，准备交给北向转发。", TraceId: "trace-demo-001"},
			{Id: "evt-003", Time: "2026-06-15 09:30:19", Level: "WARN", Source: "delivery", PluginId: "com.gcoll.http-forwarder", TaskId: "task-http-a01", Message: "备用转发目标未启用，已跳过。", TraceId: "trace-demo-002"},
		},
	}
}

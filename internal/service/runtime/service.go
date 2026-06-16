package runtime

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
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
				Id:          "com.gcoll.modbus-tcp",
				Name:        "Modbus TCP 采集",
				Type:        "southbound",
				Version:     "0.1.0",
				Runtime:     "process",
				Protocol:    "grpc",
				Status:      "running",
				Permissions: []string{"network.outbound", "config.read", "runtime.events"},
				UpdatedAt:   "2026-06-16 10:15:00",
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
				PluginId:    "com.gcoll.modbus-tcp",
				PluginName:  "Modbus TCP 采集",
				Status:      "online",
				Enabled:     true,
				PointCount:  6,
				ReportMode:  "change",
				LastSeenAt:  "2026-06-16 10:30:18",
				Description: "用于验证 Modbus TCP 采集、过滤、转发闭环的本地网关。",
			},
			{
				Id:          "dev-sim-line-b02",
				Name:        "模拟产线 B02",
				Code:        "DEV-SIM-B02",
				GroupId:     "edge",
				PluginId:    "com.gcoll.modbus-tcp",
				PluginName:  "Modbus TCP 采集",
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
			PluginId:    "com.gcoll.modbus-tcp",
			Name:        "TEMP_01",
			Description: "环境温度",
			Address:     "holding_register:40001",
			ValueType:   "float",
			Unit:        "℃",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "environment"},
		},
		{
			Id:          "pt-pressure",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.modbus-tcp",
			Name:        "PRESS_01",
			Description: "管线压力",
			Address:     "holding_register:40003",
			ValueType:   "float",
			Unit:        "kPa",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "process"},
		},
		{
			Id:          "pt-motor-state",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.modbus-tcp",
			Name:        "MOTOR_RUN",
			Description: "电机运行状态",
			Address:     "coil:00001",
			ValueType:   "bool",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "status"},
		},
		{
			Id:          "pt-energy",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.modbus-tcp",
			Name:        "ENERGY_TOTAL",
			Description: "累计能耗",
			Address:     "input_register:30001",
			ValueType:   "float",
			Unit:        "kWh",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "meter"},
		},
		{
			Id:          "pt-speed-set",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.modbus-tcp",
			Name:        "SPEED_SET",
			Description: "速度设定值",
			Address:     "holding_register:40110",
			ValueType:   "int",
			Unit:        "rpm",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "write"},
		},
		{
			Id:          "pt-emergency",
			DeviceId:    "dev-edge-gw-a01",
			PluginId:    "com.gcoll.modbus-tcp",
			Name:        "EMERGENCY_STOP",
			Description: "急停输入状态",
			Address:     "discrete_input:10001",
			ValueType:   "bool",
			Enabled:     true,
			Tags:        map[string]string{"area": "A", "kind": "safety"},
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
				Name:            "样例 Modbus TCP 采集链路",
				DeviceId:        "dev-edge-gw-a01",
				DeviceName:      "边缘网关 A01",
				SouthPluginName: "Modbus TCP 采集",
				PointCount:      6,
				ReportMode:      "change",
				Status:          "running",
				Rate:            "128 条/秒",
				RuleHitRate:     "72%",
				LastCollectedAt: "2026-06-16 10:30:18",
			},
		},
	}
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
func (s *Service) GetLogs(ctx context.Context) *runtimev1.LogsRes {
	_ = ctx

	return &runtimev1.LogsRes{
		Items: []runtimev1.RuntimeEvent{
			{Id: "evt-001", Time: "2026-06-16 10:30:18", Level: "INFO", Source: "collector", PluginId: "com.gcoll.modbus-tcp", DeviceId: "dev-edge-gw-a01", TaskId: "task-modbus-a01", Message: "已接收 128 条采集记录并写入内存缓冲。", TraceId: "trace-demo-001"},
			{Id: "evt-002", Time: "2026-06-16 10:30:18", Level: "INFO", Source: "pipeline", DeviceId: "dev-edge-gw-a01", TaskId: "task-modbus-a01", Message: "规则过滤命中 92 条记录，准备交给北向转发。", TraceId: "trace-demo-001"},
			{Id: "evt-003", Time: "2026-06-16 10:30:19", Level: "WARN", Source: "delivery", PluginId: "com.gcoll.http-forwarder", TaskId: "task-modbus-a01", Message: "备用转发目标未启用，已跳过。", TraceId: "trace-demo-002"},
		},
	}
}

// GetModbusTcpDeviceConfigPage 返回指定设备的 Modbus TCP 协议配置页示例数据。
func (s *Service) GetModbusTcpDeviceConfigPage(ctx context.Context, deviceId string) (*runtimev1.ModbusTcpDeviceConfigPageRes, error) {
	devices := s.GetDevices(ctx)
	var device *runtimev1.DeviceItem
	for index := range devices.Items {
		if devices.Items[index].Id == deviceId {
			device = &devices.Items[index]
			break
		}
	}
	if device == nil {
		return nil, gerror.Newf("设备不存在: %s", deviceId)
	}
	if device.PluginId != "com.gcoll.modbus-tcp" {
		return nil, gerror.Newf("设备未使用 Modbus TCP 插件: %s", deviceId)
	}

	configs := map[string]runtimev1.ModbusTcpDeviceConfig{
		"dev-edge-gw-a01": {
			Host:             "192.168.10.25",
			Port:             502,
			UnitId:           1,
			TimeoutMs:        2000,
			PollIntervalMs:   1000,
			ReportMode:       "change",
			DebugEnabled:     true,
			MaxCoilBatch:     512,
			MaxRegisterBatch: 64,
			LowLatencyMs:     80,
			HighLatencyMs:    1000,
		},
		"dev-sim-line-b02": {
			Host:             "192.168.10.88",
			Port:             502,
			UnitId:           2,
			TimeoutMs:        3000,
			PollIntervalMs:   1500,
			ReportMode:       "change",
			DebugEnabled:     false,
			MaxCoilBatch:     256,
			MaxRegisterBatch: 32,
			LowLatencyMs:     100,
			HighLatencyMs:    1200,
		},
	}
	config, ok := configs[deviceId]
	if !ok {
		return nil, gerror.Newf("设备缺少 Modbus TCP 配置: %s", deviceId)
	}

	var (
		readPlan  []runtimev1.ModbusTcpReadBlock
		points    []runtimev1.ModbusTcpPoint
		debugLogs []runtimev1.ModbusTcpDebugLog
	)
	if deviceId == "dev-edge-gw-a01" {
		readPlan = []runtimev1.ModbusTcpReadBlock{
			{Area: "coil", Start: 0, Quantity: 1, PointIds: []string{"pt-motor-state"}, LatencyMs: 16},
			{Area: "discrete_input", Start: 0, Quantity: 1, PointIds: []string{"pt-emergency"}, LatencyMs: 14},
			{Area: "holding_register", Start: 0, Quantity: 4, PointIds: []string{"pt-temperature", "pt-pressure"}, LatencyMs: 22},
			{Area: "holding_register", Start: 109, Quantity: 1, PointIds: []string{"pt-speed-set"}, LatencyMs: 18},
			{Area: "input_register", Start: 0, Quantity: 2, PointIds: []string{"pt-energy"}, LatencyMs: 20},
		}
		points = []runtimev1.ModbusTcpPoint{
			{Id: "pt-temperature", Name: "TEMP_01", Area: "holding_register", Address: 0, Quantity: 2, ValueType: "float32", Mode: "read", ReportMode: "change", Enabled: true, ByteOrder: "big", WordOrder: "big", Scale: "1", Current: "25.6 ℃", Quality: "good", LastReadAt: "2026-06-16 10:30:18.123", Description: "环境温度"},
			{Id: "pt-pressure", Name: "PRESS_01", Area: "holding_register", Address: 2, Quantity: 2, ValueType: "float32", Mode: "read", ReportMode: "change", Enabled: true, ByteOrder: "big", WordOrder: "big", Scale: "1", Current: "101.3 kPa", Quality: "good", LastReadAt: "2026-06-16 10:30:18.120", Description: "管线压力"},
			{Id: "pt-motor-state", Name: "MOTOR_RUN", Area: "coil", Address: 0, Quantity: 1, ValueType: "bool", Mode: "read", ReportMode: "change", Enabled: true, ByteOrder: "big", WordOrder: "big", Scale: "1", Current: "ON", Quality: "good", LastReadAt: "2026-06-16 10:30:18.118", Description: "电机运行状态"},
			{Id: "pt-energy", Name: "ENERGY_TOTAL", Area: "input_register", Address: 0, Quantity: 2, ValueType: "float32", Mode: "read", ReportMode: "change", Enabled: true, ByteOrder: "big", WordOrder: "little", Scale: "1", Current: "1288.4 kWh", Quality: "uncertain", LastReadAt: "2026-06-16 10:30:18.110", Description: "累计能耗"},
			{Id: "pt-speed-set", Name: "SPEED_SET", Area: "holding_register", Address: 109, Quantity: 1, ValueType: "uint16", Mode: "write", ReportMode: "change", Enabled: true, ByteOrder: "big", WordOrder: "big", Scale: "1", Current: "960 rpm", Quality: "good", LastReadAt: "2026-06-16 10:29:55.010", Description: "速度设定值"},
			{Id: "pt-emergency", Name: "EMERGENCY_STOP", Area: "discrete_input", Address: 0, Quantity: 1, ValueType: "bool", Mode: "read", ReportMode: "change", Enabled: true, ByteOrder: "big", WordOrder: "big", Scale: "1", Current: "OFF", Quality: "good", LastReadAt: "2026-06-16 10:30:18.106", Description: "急停输入状态"},
		}
		debugLogs = []runtimev1.ModbusTcpDebugLog{
			{Time: "2026-06-16 10:30:18.101", Level: "DEBUG", Message: "批量读取成功", TraceId: "trace-demo-001", Area: "holding_register", Address: "0", CostMs: 22, RawHex: "41CC000042CA999A"},
			{Time: "2026-06-16 10:30:18.118", Level: "DEBUG", Message: "线圈读取成功", TraceId: "trace-demo-001", Area: "coil", Address: "0", CostMs: 16, RawHex: "01"},
			{Time: "2026-06-16 10:30:19.002", Level: "INFO", Message: "自适应读取上限保持稳定", TraceId: "trace-demo-001", Area: "holding_register", Address: "0", CostMs: 0, RawHex: ""},
		}
	}

	return &runtimev1.ModbusTcpDeviceConfigPageRes{
		Plugin: runtimev1.PluginItem{
			Id:          "com.gcoll.modbus-tcp",
			Name:        "Modbus TCP 采集",
			Type:        "southbound",
			Version:     "0.1.0",
			Runtime:     "process",
			Protocol:    "grpc",
			Status:      "running",
			Permissions: []string{"network.outbound", "config.read", "runtime.events"},
			UpdatedAt:   "2026-06-16 10:15:00",
		},
		Device:    *device,
		Config:    config,
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

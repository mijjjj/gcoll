package runtime

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	runtimev1 "github.com/mijjjj/gcoll/api/runtime/v1"
	"github.com/mijjjj/gcoll/internal/consts"
	storagesvc "github.com/mijjjj/gcoll/internal/service/storage"
)

func TestServiceGetHealth(t *testing.T) {
	service := New()
	resp := service.GetHealth(context.Background())

	if resp.Status != "ok" {
		t.Fatalf("健康状态不符合预期: %s", resp.Status)
	}
	if resp.Service != consts.ServiceName {
		t.Fatalf("服务名称不符合预期: %s", resp.Service)
	}
	if resp.Version == "" {
		t.Fatal("版本号不能为空")
	}
	if resp.CheckedAt == "" {
		t.Fatal("检查时间不能为空")
	}
}

func TestServiceGetOverview(t *testing.T) {
	prepareTestDatabase(t)
	service := New()
	resp, err := service.GetOverview(context.Background())
	if err != nil {
		t.Fatalf("获取工作台总览失败: %v", err)
	}

	if len(resp.Metrics) != 4 {
		t.Fatalf("工作台指标数量不符合预期: %d", len(resp.Metrics))
	}
	if resp.PluginSummary.Total != 1 {
		t.Fatalf("插件总数不符合预期: %d", resp.PluginSummary.Total)
	}
	if len(resp.Tasks) == 0 {
		t.Fatal("采集任务不能为空")
	}
	if resp.DataPlane.Backpressure != "正常" {
		t.Fatalf("背压状态不符合预期: %s", resp.DataPlane.Backpressure)
	}
}

func TestServiceGetDevicePoints(t *testing.T) {
	prepareTestDatabase(t)
	service := New()

	resp, err := service.GetDevicePoints(context.Background(), "dev-edge-gw-a01")
	if err != nil {
		t.Fatalf("获取设备点位失败: %v", err)
	}
	if len(resp.Items) != 6 {
		t.Fatalf("设备点位数量不符合预期: %d", len(resp.Items))
	}

	_, err = service.GetDevicePoints(context.Background(), "unknown")
	if err == nil {
		t.Fatal("未知设备不应返回点位")
	}
}

func TestServiceGetModbusTcpDeviceConfigPage(t *testing.T) {
	prepareTestDatabase(t)
	service := New()

	resp, err := service.GetModbusTcpDeviceConfigPage(context.Background(), "dev-edge-gw-a01")
	if err != nil {
		t.Fatalf("获取设备协议配置失败: %v", err)
	}
	if resp.Plugin.Id != "com.gcoll.modbus-tcp" {
		t.Fatalf("插件 ID 不符合预期: %s", resp.Plugin.Id)
	}
	if resp.Device.Id != "dev-edge-gw-a01" {
		t.Fatalf("设备 ID 不符合预期: %s", resp.Device.Id)
	}
	if len(resp.ReadPlan) == 0 {
		t.Fatal("读取计划不能为空")
	}
	if len(resp.DebugLogs) == 0 {
		t.Fatal("调试日志不能为空")
	}

	_, err = service.GetModbusTcpDeviceConfigPage(context.Background(), "unknown")
	if err == nil {
		t.Fatal("未知设备不应返回协议配置")
	}
}

func TestServiceCreateDeviceAndPoint(t *testing.T) {
	prepareTestDatabase(t)
	service := New()
	ctx := context.Background()

	deviceResp, err := service.CreateDevice(ctx, &runtimev1.CreateDeviceReq{
		Id:          "dev-test-create",
		Name:        "测试设备",
		Code:        "DEV-TEST-CREATE",
		GroupId:     "test",
		PluginId:    "com.gcoll.modbus-tcp",
		Enabled:     true,
		ReportMode:  "change",
		Description: "用于验证数据库设备写入。",
		Config: map[string]any{
			"host":             "192.168.1.10",
			"port":             502,
			"unitId":           1,
			"timeoutMs":        1000,
			"pollIntervalMs":   1000,
			"debugEnabled":     false,
			"maxCoilBatch":     128,
			"maxRegisterBatch": 32,
			"lowLatencyMs":     80,
			"highLatencyMs":    1000,
		},
	})
	if err != nil {
		t.Fatalf("新增设备失败: %v", err)
	}
	if deviceResp.Device.Id != "dev-test-create" {
		t.Fatalf("设备 ID 不符合预期: %s", deviceResp.Device.Id)
	}

	pointResp, err := service.CreateDevicePoint(ctx, &runtimev1.CreateDevicePointReq{
		DeviceId:  "dev-test-create",
		Id:        "pt-test-create",
		PluginId:  "com.gcoll.modbus-tcp",
		Name:      "TEST_REGISTER",
		Address:   "holding_register:40011",
		ValueType: "int",
		Enabled:   true,
		Metadata: map[string]any{
			"area":       "holding_register",
			"address":    10,
			"quantity":   1,
			"mode":       "read",
			"valueType":  "uint16",
			"byteOrder":  "big",
			"wordOrder":  "big",
			"scale":      1,
			"offset":     0,
			"reportMode": "change",
		},
	})
	if err != nil {
		t.Fatalf("新增点位失败: %v", err)
	}
	if pointResp.Point.Id != "pt-test-create" {
		t.Fatalf("点位 ID 不符合预期: %s", pointResp.Point.Id)
	}

	points, err := service.GetDevicePoints(ctx, "dev-test-create")
	if err != nil {
		t.Fatalf("读取新增设备点位失败: %v", err)
	}
	if len(points.Items) != 1 {
		t.Fatalf("新增设备点位数量不符合预期: %d", len(points.Items))
	}

	_, err = service.CreateDevicePoint(ctx, &runtimev1.CreateDevicePointReq{
		DeviceId:  "dev-test-create",
		Id:        "pt-invalid-write",
		PluginId:  "com.gcoll.modbus-tcp",
		Name:      "INVALID_WRITE",
		Address:   "input_register:30011",
		ValueType: "int",
		Enabled:   true,
		Metadata: map[string]any{
			"area":      "input_register",
			"address":   10,
			"quantity":  1,
			"mode":      "write",
			"valueType": "uint16",
			"byteOrder": "big",
			"wordOrder": "big",
		},
	})
	if err == nil {
		t.Fatal("只读区写入点位不应创建成功")
	}
}

func prepareTestDatabase(t *testing.T) {
	t.Helper()

	err := gdb.SetConfigGroup("default", gdb.ConfigGroup{{
		Type:             "sqlite",
		Link:             "sqlite::@file(:memory:)",
		MaxOpenConnCount: 1,
	}})
	if err != nil {
		t.Fatalf("设置测试数据库失败: %v", err)
	}
	if err := storagesvc.Init(context.Background()); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
}

package runtime

import (
	"context"
	"testing"

	"github.com/mijjjj/gcoll/internal/consts"
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
	service := New()
	resp := service.GetOverview(context.Background())

	if len(resp.Metrics) != 4 {
		t.Fatalf("工作台指标数量不符合预期: %d", len(resp.Metrics))
	}
	if resp.PluginSummary.Total != 2 {
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
	service := New()

	resp := service.GetDevicePoints(context.Background(), "dev-edge-gw-a01")
	if len(resp.Items) != 6 {
		t.Fatalf("设备点位数量不符合预期: %d", len(resp.Items))
	}

	emptyResp := service.GetDevicePoints(context.Background(), "unknown")
	if len(emptyResp.Items) != 0 {
		t.Fatalf("未知设备不应返回点位: %d", len(emptyResp.Items))
	}
}

func TestServiceGetModbusTcpDeviceConfigPage(t *testing.T) {
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

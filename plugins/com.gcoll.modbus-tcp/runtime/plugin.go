package modbustcp

import (
	"fmt"
	"time"
)

// Runtime 提供宿主调用 Modbus TCP 插件的入口。
type Runtime struct {
	host HostClient
}

// NewRuntime 创建插件运行时。
func NewRuntime(host HostClient) *Runtime {
	return &Runtime{host: host}
}

// Manifest 返回插件清单核心信息。
func (r *Runtime) Manifest() map[string]any {
	return map[string]any{
		"id":       PluginID,
		"name":     "Modbus TCP 采集",
		"type":     "southbound",
		"version":  "0.1.0",
		"runtime":  "process",
		"protocol": "grpc",
		"capabilities": []string{
			"southbound.collector.modbus_tcp",
			"southbound.writer.modbus_tcp",
			"runtime.debug_logs",
		},
	}
}

// ValidateConfig 校验宿主下发的设备配置。
func (r *Runtime) ValidateConfig(config ConnectionConfig, points []Point) error {
	if err := config.Normalize().Validate(); err != nil {
		return err
	}
	for _, point := range points {
		if err := point.Normalize().Validate(); err != nil {
			return err
		}
	}
	return nil
}

// TestConnection 从宿主读取配置后执行一次 TCP 连接测试。
func (r *Runtime) TestConnection(deviceID string) (time.Duration, error) {
	if r.host == nil {
		return 0, fmt.Errorf("宿主客户端未配置")
	}
	config, err := r.host.LoadDeviceConfig(deviceID)
	if err != nil {
		return 0, err
	}
	client, err := NewTCPClient(config)
	if err != nil {
		return 0, err
	}
	started := time.Now()
	if err := client.Connect(); err != nil {
		return time.Since(started), err
	}
	defer client.Close()
	return time.Since(started), nil
}

// ReadOnce 从宿主读取配置和点位后执行一次采集。
func (r *Runtime) ReadOnce(deviceID string, taskID string, traceID string) ([]ReportRecord, error) {
	if r.host == nil {
		return nil, fmt.Errorf("宿主客户端未配置")
	}
	config, err := r.host.LoadDeviceConfig(deviceID)
	if err != nil {
		return nil, err
	}
	points, err := r.host.LoadDevicePoints(deviceID)
	if err != nil {
		return nil, err
	}
	collector, err := NewCollector(config, points, nil)
	if err != nil {
		return nil, err
	}
	records, err := collector.ReadOnce(deviceID, taskID, traceID)
	if reportErr := r.host.ReportDebugLogs(collector.DebugLogs()); reportErr != nil && err == nil {
		err = reportErr
	}
	if err != nil {
		return records, err
	}
	return records, r.host.SubmitRecords(records)
}

package modbustcp

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

const (
	// PluginID 是内置 Modbus TCP 南向插件 ID。
	PluginID = "com.gcoll.modbus-tcp"

	modeRead  = "read"
	modeWrite = "write"

	areaCoil            = "coil"
	areaDiscreteInput   = "discrete_input"
	areaHoldingRegister = "holding_register"
	areaInputRegister   = "input_register"

	reportModeAll    = "all"
	reportModeChange = "change"

	defaultPort             = 502
	defaultTimeoutMs        = 2000
	defaultPollIntervalMs   = 1000
	defaultMaxCoilBatch     = 512
	defaultMaxRegisterBatch = 64
	defaultDebugLogLimit    = 500

	modbusMaxReadBits       = 2000
	modbusMaxReadRegisters  = 125
	modbusMaxWriteCoils     = 1968
	modbusMaxWriteRegisters = 123
)

// ConnectionConfig 描述设备维度的 Modbus TCP 连接配置。
type ConnectionConfig struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	UnitID             byte   `json:"unitId"`
	TimeoutMs          int    `json:"timeoutMs"`
	PollIntervalMs     int    `json:"pollIntervalMs"`
	MaxCoilBatch       int    `json:"maxCoilBatch"`
	MaxRegisterBatch   int    `json:"maxRegisterBatch"`
	ReportMode         string `json:"reportMode"`
	DebugEnabled       bool   `json:"debugEnabled"`
	DebugLogLimit      int    `json:"debugLogLimit"`
	LowLatencyMs       int    `json:"lowLatencyMs"`
	HighLatencyMs      int    `json:"highLatencyMs"`
	MaxRetryPerRequest int    `json:"maxRetryPerRequest"`
}

// Normalize 归一化本地数值默认项，不选择任何业务资源。
func (c ConnectionConfig) Normalize() ConnectionConfig {
	if c.Port == 0 {
		c.Port = defaultPort
	}
	if c.TimeoutMs == 0 {
		c.TimeoutMs = defaultTimeoutMs
	}
	if c.PollIntervalMs == 0 {
		c.PollIntervalMs = defaultPollIntervalMs
	}
	if c.MaxCoilBatch == 0 {
		c.MaxCoilBatch = defaultMaxCoilBatch
	}
	if c.MaxRegisterBatch == 0 {
		c.MaxRegisterBatch = defaultMaxRegisterBatch
	}
	if c.ReportMode == "" {
		c.ReportMode = reportModeChange
	}
	if c.DebugLogLimit == 0 {
		c.DebugLogLimit = defaultDebugLogLimit
	}
	if c.LowLatencyMs == 0 {
		c.LowLatencyMs = 80
	}
	if c.HighLatencyMs == 0 {
		c.HighLatencyMs = c.TimeoutMs / 2
	}
	if c.HighLatencyMs <= c.LowLatencyMs {
		c.HighLatencyMs = c.LowLatencyMs * 2
	}
	if c.MaxRetryPerRequest == 0 {
		c.MaxRetryPerRequest = 1
	}
	return c
}

// Validate 校验连接配置，缺失必需连接参数时直接返回错误。
func (c ConnectionConfig) Validate() error {
	c = c.Normalize()
	if strings.TrimSpace(c.Host) == "" {
		return fmt.Errorf("Modbus TCP 主机地址不能为空")
	}
	if net.ParseIP(c.Host) == nil && !isLikelyHostName(c.Host) {
		return fmt.Errorf("Modbus TCP 主机地址格式无效: %s", c.Host)
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("Modbus TCP 端口必须在 1 到 65535 之间")
	}
	if c.UnitID > 247 {
		return fmt.Errorf("Modbus Unit ID 必须在 0 到 247 之间")
	}
	if c.TimeoutMs < 100 {
		return fmt.Errorf("超时时间不能小于 100ms")
	}
	if c.PollIntervalMs < 100 {
		return fmt.Errorf("轮询间隔不能小于 100ms")
	}
	if c.MaxCoilBatch < 1 || c.MaxCoilBatch > modbusMaxReadBits {
		return fmt.Errorf("线圈批量读取长度必须在 1 到 %d 之间", modbusMaxReadBits)
	}
	if c.MaxRegisterBatch < 1 || c.MaxRegisterBatch > modbusMaxReadRegisters {
		return fmt.Errorf("寄存器批量读取长度必须在 1 到 %d 之间", modbusMaxReadRegisters)
	}
	if c.ReportMode != reportModeAll && c.ReportMode != reportModeChange {
		return fmt.Errorf("上报模式只支持 all 或 change")
	}
	return nil
}

// Address 返回 host:port 形式的网络地址。
func (c ConnectionConfig) Address() string {
	c = c.Normalize()
	return fmt.Sprintf("%s:%d", strings.TrimSpace(c.Host), c.Port)
}

// Timeout 返回请求超时时间。
func (c ConnectionConfig) Timeout() time.Duration {
	c = c.Normalize()
	return time.Duration(c.TimeoutMs) * time.Millisecond
}

func isLikelyHostName(value string) bool {
	if strings.ContainsAny(value, " /\\") {
		return false
	}
	return strings.Contains(value, ".") || strings.EqualFold(value, "localhost")
}

// Point 描述宿主下发的通用点位和 Modbus 扩展元数据。
type Point struct {
	ID          string            `json:"id"`
	DeviceID    string            `json:"deviceId"`
	PluginID    string            `json:"pluginId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Address     uint16            `json:"address"`
	Area        string            `json:"area"`
	ValueType   string            `json:"valueType"`
	Mode        string            `json:"mode"`
	Enabled     bool              `json:"enabled"`
	Quantity    uint16            `json:"quantity"`
	Scale       float64           `json:"scale"`
	Offset      float64           `json:"offset"`
	ByteOrder   string            `json:"byteOrder"`
	WordOrder   string            `json:"wordOrder"`
	Tags        map[string]string `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
}

// Normalize 补全点位的本地解析默认项。
func (p Point) Normalize() Point {
	if p.Mode == "" {
		p.Mode = modeRead
	}
	if p.Quantity == 0 {
		p.Quantity = quantityForValueType(p.ValueType)
	}
	if p.Scale == 0 {
		p.Scale = 1
	}
	if p.ByteOrder == "" {
		p.ByteOrder = "big"
	}
	if p.WordOrder == "" {
		p.WordOrder = "big"
	}
	return p
}

// Validate 校验点位元数据。
func (p Point) Validate() error {
	p = p.Normalize()
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("点位 ID 不能为空")
	}
	if p.PluginID != "" && p.PluginID != PluginID {
		return fmt.Errorf("点位 %s 的插件 ID 与 Modbus TCP 插件不匹配", p.ID)
	}
	if !isValidArea(p.Area) {
		return fmt.Errorf("点位 %s 的 Modbus 区域无效: %s", p.ID, p.Area)
	}
	if p.Mode != modeRead && p.Mode != modeWrite {
		return fmt.Errorf("点位 %s 的模式只支持 read 或 write", p.ID)
	}
	if p.Quantity == 0 {
		return fmt.Errorf("点位 %s 的读取长度必须大于 0", p.ID)
	}
	if p.Area == areaDiscreteInput || p.Area == areaInputRegister {
		if p.Mode == modeWrite {
			return fmt.Errorf("点位 %s 使用只读区域，不能配置写入模式", p.ID)
		}
	}
	return nil
}

// EndAddress 返回点位占用区间的结束地址。
func (p Point) EndAddress() uint32 {
	p = p.Normalize()
	return uint32(p.Address) + uint32(p.Quantity)
}

// ReportRecord 描述插件提交给宿主的数据记录。
type ReportRecord struct {
	PointID       string            `json:"pointId"`
	DeviceID      string            `json:"deviceId"`
	Quality       string            `json:"quality"`
	Value         any               `json:"value"`
	RawHex        string            `json:"rawHex"`
	Changed       bool              `json:"changed"`
	CollectedAt   time.Time         `json:"collectedAt"`
	TraceID       string            `json:"traceId"`
	Tags          map[string]string `json:"tags"`
	SourceAddress uint16            `json:"sourceAddress"`
}

// WriteRequest 描述写入点位请求。
type WriteRequest struct {
	PointID string `json:"pointId"`
	Value   any    `json:"value"`
	TraceID string `json:"traceId"`
}

// HostClient 描述插件与宿主交互所需能力。
type HostClient interface {
	LoadDeviceConfig(deviceID string) (ConnectionConfig, error)
	LoadDevicePoints(deviceID string) ([]Point, error)
	SubmitRecords(records []ReportRecord) error
	ReportDebugLogs(logs []DebugLog) error
}

// SortPointsForRead 将读取点位按 Unit、区域和地址排序，降低网络往返次数。
func SortPointsForRead(points []Point) ([]Point, error) {
	enabled := make([]Point, 0, len(points))
	for _, point := range points {
		point = point.Normalize()
		if !point.Enabled || point.Mode != modeRead {
			continue
		}
		if err := point.Validate(); err != nil {
			return nil, err
		}
		enabled = append(enabled, point)
	}
	sort.SliceStable(enabled, func(i, j int) bool {
		if enabled[i].Area != enabled[j].Area {
			return areaWeight(enabled[i].Area) < areaWeight(enabled[j].Area)
		}
		if enabled[i].Address != enabled[j].Address {
			return enabled[i].Address < enabled[j].Address
		}
		return enabled[i].ID < enabled[j].ID
	})
	return enabled, nil
}

func areaWeight(area string) int {
	switch area {
	case areaCoil:
		return 1
	case areaDiscreteInput:
		return 2
	case areaHoldingRegister:
		return 3
	case areaInputRegister:
		return 4
	default:
		return 99
	}
}

func isValidArea(area string) bool {
	switch area {
	case areaCoil, areaDiscreteInput, areaHoldingRegister, areaInputRegister:
		return true
	default:
		return false
	}
}

func quantityForValueType(valueType string) uint16 {
	switch valueType {
	case "bool", "int16", "uint16":
		return 1
	case "int32", "uint32", "float32":
		return 2
	case "int64", "uint64", "float64":
		return 4
	default:
		return 1
	}
}

// ParseAddressSpec 解析地址字段中的“读取功能:地址”格式。
func ParseAddressSpec(value string) (area string, address uint16, err error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return "", 0, fmt.Errorf("点位地址不能为空")
	}
	parts := strings.SplitN(text, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("点位地址格式无效，必须使用 读取功能:地址，例如 holding_register:0")
	}
	area, ok := canonicalArea(parts[0])
	if !ok {
		return "", 0, fmt.Errorf("点位地址中的读取功能无效: %s", strings.TrimSpace(parts[0]))
	}
	number, parseErr := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 16)
	if parseErr != nil {
		return "", 0, fmt.Errorf("点位地址中的寄存器地址无效: %s", strings.TrimSpace(parts[1]))
	}
	return area, uint16(number), nil
}

// PointFromCommonItem 将宿主通用点位转换为 Modbus 点位。
func PointFromCommonItem(item commonv1.PointItem) (Point, error) {
	metadata := emptyStringMap(item.Metadata)
	area, address, err := ParseAddressSpec(item.Address)
	if err != nil {
		return Point{}, err
	}
	point := Point{
		ID:          item.Id,
		DeviceID:    item.DeviceId,
		PluginID:    item.PluginId,
		Name:        item.Name,
		Description: item.Description,
		Address:     address,
		Area:        area,
		ValueType:   metadataString(metadata, "valueType", item.ValueType),
		Mode:        metadataString(metadata, "mode", ""),
		Enabled:     item.Enabled,
		Quantity:    metadataUint16(metadata, "quantity"),
		Scale:       metadataFloat64(metadata, "scale"),
		Offset:      metadataFloat64(metadata, "offset"),
		ByteOrder:   metadataString(metadata, "byteOrder", ""),
		WordOrder:   metadataString(metadata, "wordOrder", ""),
		Tags:        item.Tags,
		Metadata:    metadata,
	}
	return point.Normalize(), nil
}

func canonicalArea(value string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	switch normalized {
	case "coil", "coils", "01", "fc01", "read_coils":
		return areaCoil, true
	case "discrete_input", "discrete_inputs", "02", "fc02", "read_discrete_inputs":
		return areaDiscreteInput, true
	case "holding_register", "holding_registers", "03", "fc03", "read_holding_registers", "hr":
		return areaHoldingRegister, true
	case "input_register", "input_registers", "04", "fc04", "read_input_registers", "ir":
		return areaInputRegister, true
	default:
		return "", false
	}
}

func metadataString(metadata map[string]string, key string, fallback string) string {
	value := strings.TrimSpace(metadata[key])
	if value == "" {
		return fallback
	}
	return value
}

func metadataUint16(metadata map[string]string, key string) uint16 {
	value := strings.TrimSpace(metadata[key])
	if value == "" {
		return 0
	}
	number, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(number)
}

func metadataFloat64(metadata map[string]string, key string) float64 {
	value := strings.TrimSpace(metadata[key])
	if value == "" {
		return 0
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return number
}

func emptyStringMap(values map[string]any) map[string]string {
	if len(values) == 0 {
		return map[string]string{}
	}
	result := make(map[string]string, len(values))
	for key, value := range values {
		if value == nil {
			continue
		}
		result[key] = strings.TrimSpace(fmt.Sprint(value))
	}
	return result
}

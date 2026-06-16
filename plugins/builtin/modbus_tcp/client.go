package modbustcp

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/goburrow/modbus"
)

// ModbusClient 描述插件使用的 Modbus TCP 操作。
type ModbusClient interface {
	Connect() error
	Close() error
	ReadCoils(address uint16, quantity uint16) ([]byte, error)
	ReadDiscreteInputs(address uint16, quantity uint16) ([]byte, error)
	ReadHoldingRegisters(address uint16, quantity uint16) ([]byte, error)
	ReadInputRegisters(address uint16, quantity uint16) ([]byte, error)
	WriteSingleCoil(address uint16, value uint16) ([]byte, error)
	WriteMultipleCoils(address uint16, quantity uint16, value []byte) ([]byte, error)
	WriteSingleRegister(address uint16, value uint16) ([]byte, error)
	WriteMultipleRegisters(address uint16, quantity uint16, value []byte) ([]byte, error)
}

type tcpClient struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
}

// NewTCPClient 创建 Modbus TCP 客户端。
func NewTCPClient(config ConnectionConfig) (ModbusClient, error) {
	config = config.Normalize()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	handler := modbus.NewTCPClientHandler(config.Address())
	handler.Timeout = config.Timeout()
	handler.SlaveId = config.UnitID
	return &tcpClient{
		handler: handler,
		client:  modbus.NewClient(handler),
	}, nil
}

func (c *tcpClient) Connect() error {
	return c.handler.Connect()
}

func (c *tcpClient) Close() error {
	return c.handler.Close()
}

func (c *tcpClient) ReadCoils(address uint16, quantity uint16) ([]byte, error) {
	return c.client.ReadCoils(address, quantity)
}

func (c *tcpClient) ReadDiscreteInputs(address uint16, quantity uint16) ([]byte, error) {
	return c.client.ReadDiscreteInputs(address, quantity)
}

func (c *tcpClient) ReadHoldingRegisters(address uint16, quantity uint16) ([]byte, error) {
	return c.client.ReadHoldingRegisters(address, quantity)
}

func (c *tcpClient) ReadInputRegisters(address uint16, quantity uint16) ([]byte, error) {
	return c.client.ReadInputRegisters(address, quantity)
}

func (c *tcpClient) WriteSingleCoil(address uint16, value uint16) ([]byte, error) {
	return c.client.WriteSingleCoil(address, value)
}

func (c *tcpClient) WriteMultipleCoils(address uint16, quantity uint16, value []byte) ([]byte, error) {
	return c.client.WriteMultipleCoils(address, quantity, value)
}

func (c *tcpClient) WriteSingleRegister(address uint16, value uint16) ([]byte, error) {
	return c.client.WriteSingleRegister(address, value)
}

func (c *tcpClient) WriteMultipleRegisters(address uint16, quantity uint16, value []byte) ([]byte, error) {
	return c.client.WriteMultipleRegisters(address, quantity, value)
}

// Collector 执行 Modbus TCP 采集和写入。
type Collector struct {
	config  ConnectionConfig
	points  []Point
	client  ModbusClient
	planner *AdaptivePlanner
	debug   *DebugBuffer
	cache   map[string]any
}

// NewCollector 创建采集器。
func NewCollector(config ConnectionConfig, points []Point, client ModbusClient) (*Collector, error) {
	config = config.Normalize()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	normalized := make([]Point, 0, len(points))
	for _, point := range points {
		point = point.Normalize()
		if err := point.Validate(); err != nil {
			return nil, err
		}
		normalized = append(normalized, point)
	}
	if client == nil {
		created, err := NewTCPClient(config)
		if err != nil {
			return nil, err
		}
		client = created
	}
	return &Collector{
		config:  config,
		points:  normalized,
		client:  client,
		planner: NewAdaptivePlanner(config),
		debug:   NewDebugBuffer(config.DebugLogLimit, config.DebugEnabled),
		cache:   map[string]any{},
	}, nil
}

// ReadOnce 执行一次读取，按配置选择全量或变化上报。
func (c *Collector) ReadOnce(deviceID string, taskID string, traceID string) ([]ReportRecord, error) {
	blocks, err := c.planner.BuildReadPlan(c.points)
	if err != nil {
		return nil, err
	}
	records := make([]ReportRecord, 0, len(c.points))
	for _, block := range blocks {
		started := time.Now()
		raw, err := c.readBlock(block)
		c.planner.Observe(block.Area, block.Quantity, time.Since(started), err)
		if err != nil {
			c.debug.Append(DebugLog{
				Level:    "ERROR",
				DeviceID: deviceID,
				TaskID:   taskID,
				Message:  "批量读取失败",
				TraceID:  traceID,
				Fields: map[string]string{
					"area":     block.Area,
					"start":    fmt.Sprint(block.Start),
					"quantity": fmt.Sprint(block.Quantity),
					"error":    err.Error(),
				},
			})
			return records, err
		}
		c.debug.Append(DebugLog{
			Level:    "DEBUG",
			DeviceID: deviceID,
			TaskID:   taskID,
			Message:  "批量读取成功",
			TraceID:  traceID,
			Fields: map[string]string{
				"area":       block.Area,
				"start":      fmt.Sprint(block.Start),
				"quantity":   fmt.Sprint(block.Quantity),
				"latency_ms": fmt.Sprint(time.Since(started).Milliseconds()),
			},
		})
		for _, point := range block.Points {
			pointRaw, err := sliceRawForPoint(block, point, raw)
			if err != nil {
				return records, err
			}
			value, err := DecodeValue(point, pointRaw)
			if err != nil {
				return records, err
			}
			previous, existed := c.cache[point.ID]
			changed := !existed || fmt.Sprint(previous) != fmt.Sprint(value)
			c.cache[point.ID] = value
			if c.config.ReportMode == reportModeChange && !changed {
				continue
			}
			records = append(records, ReportRecord{
				PointID:       point.ID,
				DeviceID:      deviceID,
				Quality:       "good",
				Value:         value,
				RawHex:        strings.ToUpper(hex.EncodeToString(pointRaw)),
				Changed:       changed,
				CollectedAt:   time.Now(),
				TraceID:       traceID,
				Tags:          point.Tags,
				SourceAddress: point.Address,
			})
		}
	}
	return records, nil
}

// WritePoint 向可写点位写入值。
func (c *Collector) WritePoint(request WriteRequest) error {
	point, ok := c.findPoint(request.PointID)
	if !ok {
		return fmt.Errorf("点位不存在: %s", request.PointID)
	}
	if point.Mode != modeWrite {
		return fmt.Errorf("点位 %s 未配置为写入模式", request.PointID)
	}
	raw, err := EncodeWriteValue(point, request.Value)
	if err != nil {
		return err
	}
	switch point.Area {
	case areaCoil:
		value := uint16(0)
		if len(raw) >= 2 && raw[0] == 0xff {
			value = 0xff00
		}
		_, err = c.client.WriteSingleCoil(point.Address, value)
	case areaHoldingRegister:
		if point.Quantity == 1 {
			_, err = c.client.WriteSingleRegister(point.Address, uint16(raw[0])<<8|uint16(raw[1]))
		} else {
			if point.Quantity > modbusMaxWriteRegisters {
				return fmt.Errorf("点位 %s 写入寄存器长度超过 %d", point.ID, modbusMaxWriteRegisters)
			}
			_, err = c.client.WriteMultipleRegisters(point.Address, point.Quantity, raw)
		}
	default:
		err = fmt.Errorf("点位 %s 的区域不支持写入: %s", point.ID, point.Area)
	}
	if err != nil {
		return err
	}
	c.debug.Append(DebugLog{
		Level:   "INFO",
		PointID: point.ID,
		Message: "点位写入成功",
		TraceID: request.TraceID,
		Fields: map[string]string{
			"area":    point.Area,
			"address": fmt.Sprint(point.Address),
		},
	})
	return nil
}

// DebugLogs 返回调试日志快照。
func (c *Collector) DebugLogs() []DebugLog {
	return c.debug.Snapshot()
}

func (c *Collector) readBlock(block ReadBlock) ([]byte, error) {
	switch block.Area {
	case areaCoil:
		return c.client.ReadCoils(block.Start, block.Quantity)
	case areaDiscreteInput:
		return c.client.ReadDiscreteInputs(block.Start, block.Quantity)
	case areaHoldingRegister:
		return c.client.ReadHoldingRegisters(block.Start, block.Quantity)
	case areaInputRegister:
		return c.client.ReadInputRegisters(block.Start, block.Quantity)
	default:
		return nil, fmt.Errorf("不支持的 Modbus 区域: %s", block.Area)
	}
}

func (c *Collector) findPoint(pointID string) (Point, bool) {
	for _, point := range c.points {
		if point.ID == pointID {
			return point, true
		}
	}
	return Point{}, false
}

func sliceRawForPoint(block ReadBlock, point Point, raw []byte) ([]byte, error) {
	offset := int(point.Address - block.Start)
	if block.Area == areaCoil || block.Area == areaDiscreteInput {
		byteOffset := offset / 8
		if byteOffset >= len(raw) {
			return nil, fmt.Errorf("点位 %s 的线圈响应越界", point.ID)
		}
		bitOffset := offset % 8
		if raw[byteOffset]&(1<<bitOffset) != 0 {
			return []byte{1}, nil
		}
		return []byte{0}, nil
	}
	byteOffset := offset * 2
	byteLength := int(point.Quantity) * 2
	if byteOffset+byteLength > len(raw) {
		return nil, fmt.Errorf("点位 %s 的寄存器响应越界", point.ID)
	}
	return raw[byteOffset : byteOffset+byteLength], nil
}

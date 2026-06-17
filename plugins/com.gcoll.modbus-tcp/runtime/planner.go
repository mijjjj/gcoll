package modbustcp

import "time"

// ReadBlock 描述一次可合并的 Modbus 读取请求。
type ReadBlock struct {
	Area     string  `json:"area"`
	Start    uint16  `json:"start"`
	Quantity uint16  `json:"quantity"`
	Points   []Point `json:"points"`
}

// AdaptivePlanner 根据点位表和网络反馈调整批量读取长度。
type AdaptivePlanner struct {
	config        ConnectionConfig
	coilLimit     uint16
	registerLimit uint16
}

// NewAdaptivePlanner 创建自适应读取规划器。
func NewAdaptivePlanner(config ConnectionConfig) *AdaptivePlanner {
	config = config.Normalize()
	return &AdaptivePlanner{
		config:        config,
		coilLimit:     uint16(min(min(config.MaxCoilBatch, defaultMaxCoilBatch), modbusMaxReadBits)),
		registerLimit: uint16(min(min(config.MaxRegisterBatch, defaultMaxRegisterBatch), modbusMaxReadRegisters)),
	}
}

// BuildReadPlan 按区域和连续地址合并点位，生成批量读取计划。
func (p *AdaptivePlanner) BuildReadPlan(points []Point) ([]ReadBlock, error) {
	sorted, err := SortPointsForRead(points)
	if err != nil {
		return nil, err
	}
	blocks := make([]ReadBlock, 0, len(sorted))
	for _, point := range sorted {
		limit := p.limitForArea(point.Area)
		if len(blocks) == 0 || !canMerge(blocks[len(blocks)-1], point, limit) {
			blocks = append(blocks, ReadBlock{
				Area:     point.Area,
				Start:    point.Address,
				Quantity: point.Quantity,
				Points:   []Point{point},
			})
			continue
		}
		block := &blocks[len(blocks)-1]
		block.Points = append(block.Points, point)
		block.Quantity = uint16(point.EndAddress() - uint32(block.Start))
	}
	return blocks, nil
}

// Observe 根据最近一次读取结果调整批量读取限制。
func (p *AdaptivePlanner) Observe(area string, quantity uint16, latency time.Duration, err error) {
	limit := p.limitForArea(area)
	if err != nil || latency >= time.Duration(p.config.HighLatencyMs)*time.Millisecond {
		limit = maxUint16(1, uint16(float64(limit)*0.65))
	} else if latency <= time.Duration(p.config.LowLatencyMs)*time.Millisecond && quantity >= limit {
		limit = minUint16(p.maxForArea(area), uint16(float64(limit)*1.25)+1)
	}
	p.setLimitForArea(area, limit)
}

// Limits 返回当前自适应读取上限。
func (p *AdaptivePlanner) Limits() (coilLimit uint16, registerLimit uint16) {
	return p.coilLimit, p.registerLimit
}

func (p *AdaptivePlanner) limitForArea(area string) uint16 {
	if area == areaCoil || area == areaDiscreteInput {
		return p.coilLimit
	}
	return p.registerLimit
}

func (p *AdaptivePlanner) setLimitForArea(area string, limit uint16) {
	if area == areaCoil || area == areaDiscreteInput {
		p.coilLimit = minUint16(limit, modbusMaxReadBits)
		return
	}
	p.registerLimit = minUint16(limit, modbusMaxReadRegisters)
}

func (p *AdaptivePlanner) maxForArea(area string) uint16 {
	if area == areaCoil || area == areaDiscreteInput {
		return minUint16(uint16(p.config.MaxCoilBatch), modbusMaxReadBits)
	}
	return minUint16(uint16(p.config.MaxRegisterBatch), modbusMaxReadRegisters)
}

func canMerge(block ReadBlock, point Point, limit uint16) bool {
	if block.Area != point.Area {
		return false
	}
	if uint32(point.Address) < uint32(block.Start)+uint32(block.Quantity) {
		return true
	}
	required := point.EndAddress() - uint32(block.Start)
	return required <= uint32(limit)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func minUint16(a uint16, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func maxUint16(a uint16, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

package modbustcp

import (
	"testing"
	"time"
)

func TestBuildReadPlanSortsAndMergesPoints(t *testing.T) {
	planner := NewAdaptivePlanner(ConnectionConfig{
		Host:             "127.0.0.1",
		MaxRegisterBatch: 10,
		MaxCoilBatch:     10,
	})
	blocks, err := planner.BuildReadPlan([]Point{
		{ID: "hr-12", PluginID: PluginID, Area: areaHoldingRegister, Address: 12, ValueType: "uint16", Enabled: true},
		{ID: "coil-1", PluginID: PluginID, Area: areaCoil, Address: 1, ValueType: "bool", Enabled: true},
		{ID: "hr-10", PluginID: PluginID, Area: areaHoldingRegister, Address: 10, ValueType: "float32", Enabled: true},
	})
	if err != nil {
		t.Fatalf("生成读取计划失败: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("读取块数量 = %d, want 2", len(blocks))
	}
	if blocks[0].Area != areaCoil || blocks[0].Start != 1 {
		t.Fatalf("第一个读取块未按区域和地址排序: %#v", blocks[0])
	}
	if blocks[1].Area != areaHoldingRegister || blocks[1].Start != 10 || blocks[1].Quantity != 3 {
		t.Fatalf("寄存器读取块未正确合并: %#v", blocks[1])
	}
}

func TestAdaptivePlannerAdjustsLimits(t *testing.T) {
	planner := NewAdaptivePlanner(ConnectionConfig{
		Host:             "127.0.0.1",
		MaxRegisterBatch: 100,
		LowLatencyMs:     10,
		HighLatencyMs:    100,
	})
	_, before := planner.Limits()
	planner.Observe(areaHoldingRegister, before, 5*time.Millisecond, nil)
	_, increased := planner.Limits()
	if increased <= before {
		t.Fatalf("低延迟成功读取后未提高批量长度: before=%d after=%d", before, increased)
	}
	planner.Observe(areaHoldingRegister, increased, 150*time.Millisecond, nil)
	_, decreased := planner.Limits()
	if decreased >= increased {
		t.Fatalf("高延迟读取后未降低批量长度: before=%d after=%d", increased, decreased)
	}
}

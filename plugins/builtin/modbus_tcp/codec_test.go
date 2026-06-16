package modbustcp

import "testing"

func TestDecodeValueWithScale(t *testing.T) {
	value, err := DecodeValue(Point{
		ID:        "temperature",
		PluginID:  PluginID,
		Area:      areaHoldingRegister,
		ValueType: "int16",
		Quantity:  1,
		Scale:     0.1,
		Enabled:   true,
	}, []byte{0x00, 0xfa})
	if err != nil {
		t.Fatalf("解码失败: %v", err)
	}
	if value != 25.0 {
		t.Fatalf("解码值 = %v, want 25", value)
	}
}

func TestDecodeLittleWordOrderFloat32(t *testing.T) {
	value, err := DecodeValue(Point{
		ID:        "flow",
		PluginID:  PluginID,
		Area:      areaHoldingRegister,
		ValueType: "float32",
		Quantity:  2,
		WordOrder: "little",
		Enabled:   true,
	}, []byte{0x00, 0x00, 0x42, 0x48})
	if err != nil {
		t.Fatalf("解码失败: %v", err)
	}
	if value != 50.0 {
		t.Fatalf("解码值 = %v, want 50", value)
	}
}

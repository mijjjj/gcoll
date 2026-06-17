package modbustcp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

// DecodeValue 将 Modbus 响应字节解码为点位值。
func DecodeValue(point Point, raw []byte) (any, error) {
	point = point.Normalize()
	if point.Area == areaCoil || point.Area == areaDiscreteInput {
		if len(raw) == 0 {
			return nil, fmt.Errorf("点位 %s 的线圈响应为空", point.ID)
		}
		value := raw[0]&0x01 == 0x01
		return value, nil
	}
	if len(raw) < int(point.Quantity)*2 {
		return nil, fmt.Errorf("点位 %s 的寄存器响应长度不足", point.ID)
	}
	ordered := reorderRegisterBytes(point, raw[:int(point.Quantity)*2])
	switch point.ValueType {
	case "bool":
		return binary.BigEndian.Uint16(ordered[:2]) != 0, nil
	case "int16":
		return applyScale(float64(int16(binary.BigEndian.Uint16(ordered[:2]))), point), nil
	case "uint16", "int":
		return applyScale(float64(binary.BigEndian.Uint16(ordered[:2])), point), nil
	case "int32":
		return applyScale(float64(int32(binary.BigEndian.Uint32(ordered[:4]))), point), nil
	case "uint32":
		return applyScale(float64(binary.BigEndian.Uint32(ordered[:4])), point), nil
	case "float32", "float":
		return applyScale(float64(math.Float32frombits(binary.BigEndian.Uint32(ordered[:4]))), point), nil
	case "int64":
		return applyScale(float64(int64(binary.BigEndian.Uint64(ordered[:8]))), point), nil
	case "uint64":
		return applyScale(float64(binary.BigEndian.Uint64(ordered[:8])), point), nil
	case "float64":
		return applyScale(math.Float64frombits(binary.BigEndian.Uint64(ordered[:8])), point), nil
	default:
		return hex.EncodeToString(raw[:int(point.Quantity)*2]), nil
	}
}

// EncodeWriteValue 将写入值编码为 Modbus 写入字节。
func EncodeWriteValue(point Point, value any) ([]byte, error) {
	point = point.Normalize()
	if point.Area == areaCoil {
		boolValue, err := asBool(value)
		if err != nil {
			return nil, err
		}
		if boolValue {
			return []byte{0xff, 0x00}, nil
		}
		return []byte{0x00, 0x00}, nil
	}
	registers := make([]byte, int(point.Quantity)*2)
	number, err := asFloat64(value)
	if err != nil {
		return nil, err
	}
	number = (number - point.Offset) / point.Scale
	switch point.ValueType {
	case "bool":
		if value == true {
			binary.BigEndian.PutUint16(registers[:2], 1)
		}
	case "int16":
		binary.BigEndian.PutUint16(registers[:2], uint16(int16(number)))
	case "uint16", "int":
		binary.BigEndian.PutUint16(registers[:2], uint16(number))
	case "int32":
		binary.BigEndian.PutUint32(registers[:4], uint32(int32(number)))
	case "uint32":
		binary.BigEndian.PutUint32(registers[:4], uint32(number))
	case "float32", "float":
		binary.BigEndian.PutUint32(registers[:4], math.Float32bits(float32(number)))
	case "int64":
		binary.BigEndian.PutUint64(registers[:8], uint64(int64(number)))
	case "uint64":
		binary.BigEndian.PutUint64(registers[:8], uint64(number))
	case "float64":
		binary.BigEndian.PutUint64(registers[:8], math.Float64bits(number))
	default:
		return nil, fmt.Errorf("点位 %s 不支持写入值类型: %s", point.ID, point.ValueType)
	}
	return restoreRegisterBytes(point, registers), nil
}

func applyScale(value float64, point Point) float64 {
	return value*point.Scale + point.Offset
}

func reorderRegisterBytes(point Point, raw []byte) []byte {
	result := append([]byte(nil), raw...)
	if point.ByteOrder == "little" {
		for i := 0; i+1 < len(result); i += 2 {
			result[i], result[i+1] = result[i+1], result[i]
		}
	}
	if point.WordOrder == "little" && len(result) > 2 {
		for left, right := 0, len(result)-2; left < right; left, right = left+2, right-2 {
			result[left], result[right] = result[right], result[left]
			result[left+1], result[right+1] = result[right+1], result[left+1]
		}
	}
	return result
}

func restoreRegisterBytes(point Point, raw []byte) []byte {
	return reorderRegisterBytes(point, raw)
}

func asBool(value any) (bool, error) {
	switch typed := value.(type) {
	case bool:
		return typed, nil
	case int:
		return typed != 0, nil
	case float64:
		return typed != 0, nil
	default:
		return false, fmt.Errorf("写入线圈需要 bool 或数值")
	}
}

func asFloat64(value any) (float64, error) {
	switch typed := value.(type) {
	case int:
		return float64(typed), nil
	case int16:
		return float64(typed), nil
	case int32:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case uint16:
		return float64(typed), nil
	case uint32:
		return float64(typed), nil
	case uint64:
		return float64(typed), nil
	case float32:
		return float64(typed), nil
	case float64:
		return typed, nil
	default:
		return 0, fmt.Errorf("写入寄存器需要数值")
	}
}

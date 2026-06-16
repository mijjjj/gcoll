package point

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"

	runtimev1 "github.com/mijjjj/gcoll/api/runtime/v1"
	"github.com/mijjjj/gcoll/internal/dao"
	"github.com/mijjjj/gcoll/internal/model/do"
	"github.com/mijjjj/gcoll/internal/model/entity"
)

const modbusPluginId = "com.gcoll.modbus-tcp"

// Service 提供点位表服务。
type Service struct{}

// New 创建点位表服务。
func New() *Service {
	return &Service{}
}

// ListByDevice 返回指定设备的通用点位表。
func (s *Service) ListByDevice(ctx context.Context, deviceId string) (*runtimev1.DevicePointsRes, error) {
	device, err := s.getDevice(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	var points []entity.DevicePoints
	if err := dao.DevicePoints.Ctx(ctx).
		Where(do.DevicePoints{DeviceId: device.Id}).
		OrderAsc(dao.DevicePoints.Columns().CreatedAt).
		Scan(&points); err != nil {
		return nil, gerror.Wrapf(err, "读取设备点位失败: %s", deviceId)
	}

	items := make([]runtimev1.PointItem, 0, len(points))
	for _, point := range points {
		items = append(items, toPointItem(point))
	}
	return &runtimev1.DevicePointsRes{Items: items}, nil
}

// Create 为设备新增点位。
func (s *Service) Create(ctx context.Context, req *runtimev1.CreateDevicePointReq) (*runtimev1.PointItem, error) {
	device, err := s.getDevice(ctx, req.DeviceId)
	if err != nil {
		return nil, err
	}
	if device.PluginId != req.PluginId {
		return nil, gerror.Newf("点位插件与设备插件不一致: %s", req.PluginId)
	}
	if err := validatePointMetadata(req.PluginId, req.Metadata); err != nil {
		return nil, err
	}

	pointId := strings.TrimSpace(req.Id)
	if pointId == "" {
		pointId = "pt-" + guid.S()
	}
	tagsJSON, err := json.Marshal(emptyStringMap(req.Tags))
	if err != nil {
		return nil, gerror.Wrap(err, "序列化点位标签失败")
	}
	metadataJSON, err := json.Marshal(emptyAnyMap(req.Metadata))
	if err != nil {
		return nil, gerror.Wrap(err, "序列化点位扩展信息失败")
	}
	snapshotJSON, err := json.Marshal(map[string]any{
		"name":      req.Name,
		"address":   req.Address,
		"valueType": req.ValueType,
		"metadata":  emptyAnyMap(req.Metadata),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "序列化点位版本失败")
	}

	if err := dao.DevicePoints.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_ = ctx
		if _, err := tx.Insert(dao.DevicePoints.Table(), do.DevicePoints{
			Id:           pointId,
			DeviceId:     req.DeviceId,
			PluginId:     req.PluginId,
			Name:         req.Name,
			Description:  req.Description,
			Address:      req.Address,
			ValueType:    req.ValueType,
			Unit:         req.Unit,
			Enabled:      boolInt(req.Enabled),
			TagsJson:     string(tagsJSON),
			MetadataJson: string(metadataJSON),
		}); err != nil {
			return gerror.Wrap(err, "新增设备点位失败")
		}
		if _, err := tx.Insert(dao.DevicePointVersions.Table(), do.DevicePointVersions{
			Id:           "dpv-" + pointId + "-1",
			PointId:      pointId,
			DeviceId:     req.DeviceId,
			PluginId:     req.PluginId,
			Version:      1,
			SnapshotJson: string(snapshotJSON),
			ChangeNote:   "新增设备点位",
		}); err != nil {
			return gerror.Wrap(err, "新增设备点位版本失败")
		}
		if req.PluginId == modbusPluginId {
			if _, err := tx.Insert(dao.ModbusTcpPointProfiles.Table(), do.ModbusTcpPointProfiles{
				Id:        "mpp-" + pointId,
				DeviceId:  req.DeviceId,
				PointId:   pointId,
				PluginId:  req.PluginId,
				Version:   1,
				Area:      gconv.String(req.Metadata["area"]),
				Address:   gconv.Int(req.Metadata["address"]),
				Quantity:  gconv.Int(req.Metadata["quantity"]),
				Mode:      gconv.String(req.Metadata["mode"]),
				ValueType: gconv.String(req.Metadata["valueType"]),
				ByteOrder: gconv.String(req.Metadata["byteOrder"]),
				WordOrder: gconv.String(req.Metadata["wordOrder"]),
				Scale:     gconv.Float64(req.Metadata["scale"]),
				Offset:    gconv.Float64(req.Metadata["offset"]),
				ReportMode: func() string {
					if value := gconv.String(req.Metadata["reportMode"]); value != "" {
						return value
					}
					return device.ReportMode
				}(),
				Enabled: boolInt(req.Enabled),
			}); err != nil {
				return gerror.Wrap(err, "新增 Modbus TCP 点位配置失败")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	item := runtimev1.PointItem{
		Id:          pointId,
		DeviceId:    req.DeviceId,
		PluginId:    req.PluginId,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		ValueType:   req.ValueType,
		Unit:        req.Unit,
		Enabled:     req.Enabled,
		Tags:        emptyStringMap(req.Tags),
		Metadata:    emptyAnyMap(req.Metadata),
	}
	return &item, nil
}

func (s *Service) getDevice(ctx context.Context, deviceId string) (*entity.Devices, error) {
	var device entity.Devices
	if err := dao.Devices.Ctx(ctx).Where(do.Devices{Id: deviceId}).Scan(&device); err != nil {
		return nil, gerror.Wrapf(err, "读取设备失败: %s", deviceId)
	}
	if device.Id == "" {
		return nil, gerror.Newf("设备不存在: %s", deviceId)
	}
	return &device, nil
}

func toPointItem(point entity.DevicePoints) runtimev1.PointItem {
	return runtimev1.PointItem{
		Id:          point.Id,
		DeviceId:    point.DeviceId,
		PluginId:    point.PluginId,
		Name:        point.Name,
		Description: point.Description,
		Address:     point.Address,
		ValueType:   point.ValueType,
		Unit:        point.Unit,
		Enabled:     point.Enabled == 1,
		Tags:        stringMapFromJSON(point.TagsJson),
		Metadata:    anyMapFromJSON(point.MetadataJson),
	}
}

func validatePointMetadata(pluginId string, metadata map[string]any) error {
	if pluginId != modbusPluginId {
		return nil
	}
	if metadata == nil {
		return gerror.New("Modbus TCP 点位 metadata 不能为空")
	}
	area := gconv.String(metadata["area"])
	switch area {
	case "coil", "discrete_input", "holding_register", "input_register":
	default:
		return gerror.New("Modbus TCP 点位 area 不支持")
	}
	mode := gconv.String(metadata["mode"])
	if mode != "read" && mode != "write" {
		return gerror.New("Modbus TCP 点位 mode 必须是 read 或 write")
	}
	if mode == "write" && (area == "discrete_input" || area == "input_register") {
		return gerror.New("Modbus TCP 只读区域不能配置写入点位")
	}
	if gconv.Int(metadata["quantity"]) < 1 {
		return gerror.New("Modbus TCP 点位 quantity 必须大于 0")
	}
	if strings.TrimSpace(gconv.String(metadata["valueType"])) == "" {
		return gerror.New("Modbus TCP 点位 metadata 缺少 valueType")
	}
	byteOrder := gconv.String(metadata["byteOrder"])
	if byteOrder != "big" && byteOrder != "little" {
		return gerror.New("Modbus TCP 点位 byteOrder 必须是 big 或 little")
	}
	wordOrder := gconv.String(metadata["wordOrder"])
	if wordOrder != "big" && wordOrder != "little" {
		return gerror.New("Modbus TCP 点位 wordOrder 必须是 big 或 little")
	}
	return nil
}

func stringMapFromJSON(raw string) map[string]string {
	values := map[string]string{}
	if raw == "" {
		return values
	}
	_ = json.Unmarshal([]byte(raw), &values)
	return values
}

func anyMapFromJSON(raw string) map[string]any {
	values := map[string]any{}
	if raw == "" {
		return values
	}
	_ = json.Unmarshal([]byte(raw), &values)
	return values
}

func emptyStringMap(values map[string]string) map[string]string {
	if values == nil {
		return map[string]string{}
	}
	return values
}

func emptyAnyMap(values map[string]any) map[string]any {
	if values == nil {
		return map[string]any{}
	}
	return values
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

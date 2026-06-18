package point

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/guid"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
	devicev1 "github.com/mijjjj/gcoll/api/device/v1"
	"github.com/mijjjj/gcoll/internal/dao"
	"github.com/mijjjj/gcoll/internal/model/do"
	"github.com/mijjjj/gcoll/internal/model/entity"
	pluginhostsvc "github.com/mijjjj/gcoll/internal/service/pluginhost"
)

// Service 提供点位表服务。
type Service struct {
	pluginHost *pluginhostsvc.Service
}

// New 创建点位表服务。
func New() *Service {
	return &Service{pluginHost: pluginhostsvc.Instance()}
}

// ListByDevice 返回指定设备的通用点位表。
func (s *Service) ListByDevice(ctx context.Context, deviceId string) (*devicev1.DevicePointsRes, error) {
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

	items := make([]commonv1.PointItem, 0, len(points))
	for _, point := range points {
		items = append(items, toPointItem(point))
	}
	return &devicev1.DevicePointsRes{Items: items}, nil
}

// Create 为设备新增点位。
func (s *Service) Create(ctx context.Context, req *devicev1.CreateDevicePointReq) (*commonv1.PointItem, error) {
	device, err := s.getDevice(ctx, req.DeviceId)
	if err != nil {
		return nil, err
	}
	if device.PluginId != req.PluginId {
		return nil, gerror.Newf("点位插件与设备插件不一致: %s", req.PluginId)
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

	item := commonv1.PointItem{
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
	currentPoints, err := s.currentPoints(ctx, req.DeviceId)
	if err != nil {
		return nil, err
	}
	config, err := s.currentConfig(ctx, req.DeviceId, device.PluginId)
	if err != nil {
		return nil, err
	}
	if err := s.pluginHost.ValidateRuntimeConfig(ctx, device.PluginId, config, append(currentPoints, item)); err != nil {
		return nil, err
	}

	if err := dao.DevicePoints.Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		if _, err := dao.DevicePoints.Ctx(ctx).Data(do.DevicePoints{
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
		}).Insert(); err != nil {
			return gerror.Wrap(err, "新增设备点位失败")
		}
		if _, err := dao.DevicePointVersions.Ctx(ctx).Data(do.DevicePointVersions{
			Id:           "dpv-" + pointId + "-1",
			PointId:      pointId,
			DeviceId:     req.DeviceId,
			PluginId:     req.PluginId,
			Version:      1,
			SnapshotJson: string(snapshotJSON),
			ChangeNote:   "新增设备点位",
		}).Insert(); err != nil {
			return gerror.Wrap(err, "新增设备点位版本失败")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &item, nil
}

// ReplaceByDevice 保存设备完整点位表。
func (s *Service) ReplaceByDevice(ctx context.Context, deviceId string, items []commonv1.PointItem) (*devicev1.DevicePointsRes, error) {
	device, err := s.getDevice(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	normalized := make([]commonv1.PointItem, 0, len(items))
	for _, item := range items {
		item.DeviceId = device.Id
		item.PluginId = device.PluginId
		if strings.TrimSpace(item.Id) == "" {
			item.Id = "pt-" + guid.S()
		}
		item.Tags = emptyStringMap(item.Tags)
		item.Metadata = emptyAnyMap(item.Metadata)
		normalized = append(normalized, item)
	}
	config, err := s.currentConfig(ctx, device.Id, device.PluginId)
	if err != nil {
		return nil, err
	}
	if err := s.pluginHost.ValidateRuntimeConfig(ctx, device.PluginId, config, normalized); err != nil {
		return nil, err
	}

	var existing []entity.DevicePoints
	if err := dao.DevicePoints.Ctx(ctx).Where(do.DevicePoints{DeviceId: device.Id}).Scan(&existing); err != nil {
		return nil, gerror.Wrapf(err, "读取设备现有点位失败: %s", device.Id)
	}
	existingMap := map[string]entity.DevicePoints{}
	for _, point := range existing {
		existingMap[point.Id] = point
	}
	keep := map[string]bool{}

	if err := dao.DevicePoints.Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		for _, item := range normalized {
			keep[item.Id] = true
			tagsJSON, err := json.Marshal(emptyStringMap(item.Tags))
			if err != nil {
				return gerror.Wrap(err, "序列化点位标签失败")
			}
			metadataJSON, err := json.Marshal(emptyAnyMap(item.Metadata))
			if err != nil {
				return gerror.Wrap(err, "序列化点位扩展信息失败")
			}
			data := do.DevicePoints{
				Id:           item.Id,
				DeviceId:     item.DeviceId,
				PluginId:     item.PluginId,
				Name:         item.Name,
				Description:  item.Description,
				Address:      item.Address,
				ValueType:    item.ValueType,
				Unit:         item.Unit,
				Enabled:      boolInt(item.Enabled),
				TagsJson:     string(tagsJSON),
				MetadataJson: string(metadataJSON),
			}
			if _, exists := existingMap[item.Id]; exists {
				if _, err := dao.DevicePoints.Ctx(ctx).
					Where(do.DevicePoints{Id: item.Id}).
					Data(data).
					Update(); err != nil {
					return gerror.Wrap(err, "更新设备点位失败")
				}
			} else {
				if _, err := dao.DevicePoints.Ctx(ctx).Data(data).Insert(); err != nil {
					return gerror.Wrap(err, "新增设备点位失败")
				}
			}
			if err := insertPointVersion(ctx, item, "保存设备点位表"); err != nil {
				return err
			}
		}
		for _, point := range existing {
			if keep[point.Id] {
				continue
			}
			if err := ensurePointID(point.Id); err != nil {
				return err
			}
			if _, err := dao.DevicePoints.Ctx(ctx).
				Delete(do.DevicePoints{Id: point.Id}); err != nil {
				return gerror.Wrap(err, "删除设备点位失败")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.ListByDevice(ctx, device.Id)
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

func (s *Service) currentPoints(ctx context.Context, deviceId string) ([]commonv1.PointItem, error) {
	result, err := s.ListByDevice(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (s *Service) currentConfig(ctx context.Context, deviceId string, pluginId string) (map[string]any, error) {
	var config entity.PluginDeviceConfigs
	if err := dao.PluginDeviceConfigs.Ctx(ctx).Where(do.PluginDeviceConfigs{
		DeviceId: deviceId,
		PluginId: pluginId,
		Active:   1,
	}).Scan(&config); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, gerror.Wrapf(err, "读取设备插件配置失败: %s", deviceId)
		}
	}
	if config.Id == "" || strings.TrimSpace(config.ConfigJson) == "" {
		return nil, gerror.Newf("设备缺少插件运行配置: %s", deviceId)
	}
	values := map[string]any{}
	if err := json.Unmarshal([]byte(config.ConfigJson), &values); err != nil {
		return nil, gerror.Wrapf(err, "解析设备插件配置失败: %s", deviceId)
	}
	return values, nil
}

func insertPointVersion(ctx context.Context, point commonv1.PointItem, note string) error {
	snapshotJSON, err := json.Marshal(map[string]any{
		"name":      point.Name,
		"address":   point.Address,
		"valueType": point.ValueType,
		"metadata":  emptyAnyMap(point.Metadata),
	})
	if err != nil {
		return gerror.Wrap(err, "序列化点位版本失败")
	}
	if _, err := dao.DevicePointVersions.Ctx(ctx).Data(do.DevicePointVersions{
		Id:           "dpv-" + guid.S(),
		PointId:      point.Id,
		DeviceId:     point.DeviceId,
		PluginId:     point.PluginId,
		Version:      1,
		SnapshotJson: string(snapshotJSON),
		ChangeNote:   note,
	}).Insert(); err != nil {
		return gerror.Wrap(err, "新增设备点位版本失败")
	}
	return nil
}

func toPointItem(point entity.DevicePoints) commonv1.PointItem {
	return commonv1.PointItem{
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

func ensurePointID(pointID string) error {
	if strings.TrimSpace(pointID) == "" {
		return gerror.New("点位删除条件不能为空: pointId")
	}
	return nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

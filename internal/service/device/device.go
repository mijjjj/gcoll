package device

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

// Service 提供设备控制面服务。
type Service struct {
	pluginHost *pluginhostsvc.Service
}

// New 创建设备服务。
func New() *Service {
	return &Service{pluginHost: pluginhostsvc.Instance()}
}

// List 返回设备分组和设备列表。
func (s *Service) List(ctx context.Context) (*devicev1.DevicesRes, error) {
	var (
		groups  []entity.DeviceGroups
		devices []entity.Devices
		points  []entity.DevicePoints
	)
	if err := dao.DeviceGroups.Ctx(ctx).OrderAsc(dao.DeviceGroups.Columns().SortOrder).Scan(&groups); err != nil {
		return nil, gerror.Wrap(err, "读取设备分组失败")
	}
	if err := dao.Devices.Ctx(ctx).OrderAsc(dao.Devices.Columns().CreatedAt).Scan(&devices); err != nil {
		return nil, gerror.Wrap(err, "读取设备列表失败")
	}
	if err := dao.DevicePoints.Ctx(ctx).Scan(&points); err != nil {
		return nil, gerror.Wrap(err, "读取设备点位数量失败")
	}

	pluginNames, err := s.pluginHost.PluginNameMap(ctx)
	if err != nil {
		return nil, err
	}
	pointCounts := make(map[string]int)
	for _, point := range points {
		pointCounts[point.DeviceId]++
	}
	groupCounts := make(map[string]int)
	for _, device := range devices {
		groupCounts[device.GroupId]++
	}

	res := &devicev1.DevicesRes{
		Groups: make([]commonv1.DeviceGroup, 0, len(groups)),
		Items:  make([]commonv1.DeviceItem, 0, len(devices)),
	}
	for _, group := range groups {
		res.Groups = append(res.Groups, commonv1.DeviceGroup{
			Id:    group.Id,
			Name:  group.Name,
			Count: groupCounts[group.Id],
		})
	}
	for _, device := range devices {
		res.Items = append(res.Items, s.toDeviceItem(device, pluginNames[device.PluginId], pointCounts[device.Id]))
	}
	return res, nil
}

// CreateGroup 新增设备分组。
func (s *Service) CreateGroup(ctx context.Context, req *devicev1.CreateDeviceGroupReq) (*commonv1.DeviceGroup, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, gerror.New("设备分组名称不能为空")
	}

	groupID := strings.TrimSpace(req.Id)
	if groupID == "" {
		groupID = "grp-" + guid.S()
	}

	var lastGroups []entity.DeviceGroups
	if err := dao.DeviceGroups.Ctx(ctx).
		OrderDesc(dao.DeviceGroups.Columns().SortOrder).
		OrderDesc(dao.DeviceGroups.Columns().CreatedAt).
		Limit(1).
		Scan(&lastGroups); err != nil {
		return nil, gerror.Wrap(err, "读取设备分组排序失败")
	}

	sortOrder := 10
	if len(lastGroups) > 0 {
		sortOrder = lastGroups[0].SortOrder + 10
	}
	if _, err := dao.DeviceGroups.Ctx(ctx).Data(do.DeviceGroups{
		Id:        groupID,
		Name:      name,
		SortOrder: sortOrder,
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "新增设备分组失败")
	}

	return &commonv1.DeviceGroup{
		Id:    groupID,
		Name:  name,
		Count: 0,
	}, nil
}

// DeleteGroup 删除空设备分组。
func (s *Service) DeleteGroup(ctx context.Context, groupID string) error {
	if err := s.ensureGroup(ctx, groupID); err != nil {
		return err
	}
	deviceCount, err := dao.Devices.Ctx(ctx).Where(do.Devices{GroupId: groupID}).Count()
	if err != nil {
		return gerror.Wrapf(err, "读取设备分组占用失败: %s", groupID)
	}
	if deviceCount > 0 {
		return gerror.Newf("设备分组下仍有设备，不能删除: %s", groupID)
	}
	if _, err := dao.DeviceGroups.Ctx(ctx).Delete(do.DeviceGroups{Id: groupID}); err != nil {
		return gerror.Wrapf(err, "删除设备分组失败: %s", groupID)
	}
	return nil
}

// Create 新增设备并按需保存设备插件配置。
func (s *Service) Create(ctx context.Context, req *devicev1.CreateDeviceReq) (*commonv1.DeviceItem, error) {
	if err := s.ensureGroup(ctx, req.GroupId); err != nil {
		return nil, err
	}
	plugin, err := s.ensurePlugin(ctx, req.PluginId)
	if err != nil {
		return nil, err
	}
	hasConfig := len(req.Config) > 0
	if hasConfig {
		if err := s.pluginHost.ValidateRuntimeConfig(ctx, req.PluginId, req.Config, nil); err != nil {
			return nil, err
		}
	}
	if !hasConfig && req.Enabled {
		return nil, gerror.New("设备启用前必须先保存有效插件配置")
	}

	deviceId := strings.TrimSpace(req.Id)
	if deviceId == "" {
		deviceId = "dev-" + guid.S()
	}
	configId := "pdc-" + deviceId + "-" + strings.ReplaceAll(req.PluginId, ".", "-")
	configVersionId := configId + "-1"
	configJSON, err := marshalConfig(req.Config)
	if err != nil {
		return nil, err
	}

	if err := dao.Devices.Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		if _, err := dao.Devices.Ctx(ctx).Data(do.Devices{
			Id:          deviceId,
			Name:        req.Name,
			GroupId:     req.GroupId,
			PluginId:    req.PluginId,
			Status:      "offline",
			Enabled:     boolInt(req.Enabled),
			ReportMode:  req.ReportMode,
			Description: req.Description,
		}).Insert(); err != nil {
			return gerror.Wrap(err, "新增设备失败")
		}
		if _, err := dao.PluginDeviceConfigs.Ctx(ctx).Data(do.PluginDeviceConfigs{
			Id:         configId,
			DeviceId:   deviceId,
			PluginId:   req.PluginId,
			Version:    initialConfigVersion(hasConfig),
			ConfigJson: initialConfigJSON(hasConfig, configJSON),
			ReportMode: req.ReportMode,
			Enabled:    boolInt(req.Enabled),
			Active:     1,
		}).Insert(); err != nil {
			return gerror.Wrap(err, "新增设备插件配置失败")
		}
		if !hasConfig {
			return nil
		}
		if _, err := dao.PluginDeviceConfigVersions.Ctx(ctx).Data(do.PluginDeviceConfigVersions{
			Id:         configVersionId,
			ConfigId:   configId,
			DeviceId:   deviceId,
			PluginId:   req.PluginId,
			Version:    1,
			ConfigJson: configJSON,
			ChangeNote: "初始化设备插件配置",
		}).Insert(); err != nil {
			return gerror.Wrap(err, "新增设备配置版本失败")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &commonv1.DeviceItem{
		Id:          deviceId,
		Name:        req.Name,
		GroupId:     req.GroupId,
		PluginId:    req.PluginId,
		PluginName:  plugin.Manifest.Name,
		Status:      "offline",
		Enabled:     req.Enabled,
		PointCount:  0,
		ReportMode:  req.ReportMode,
		LastSeenAt:  "尚未连接",
		Description: req.Description,
	}, nil
}

// MoveToGroup 移动设备所属分组。
func (s *Service) MoveToGroup(ctx context.Context, deviceID string, groupID string) (*commonv1.DeviceItem, error) {
	device, err := s.Get(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureGroup(ctx, groupID); err != nil {
		return nil, err
	}
	if device.GroupId == groupID {
		return s.buildDeviceItem(ctx, *device)
	}
	if _, err := dao.Devices.Ctx(ctx).
		Where(do.Devices{Id: deviceID}).
		Data(do.Devices{GroupId: groupID}).
		Update(); err != nil {
		return nil, gerror.Wrapf(err, "移动设备分组失败: %s", deviceID)
	}
	device.GroupId = groupID
	return s.buildDeviceItem(ctx, *device)
}

// Delete 删除设备及其控制面关联数据。
func (s *Service) Delete(ctx context.Context, deviceID string) error {
	if err := ensureDeleteConditionValue("deviceID", deviceID); err != nil {
		return err
	}
	if _, err := s.Get(ctx, deviceID); err != nil {
		return err
	}

	var tasks []entity.CollectionTasks
	if err := dao.CollectionTasks.Ctx(ctx).
		Where(do.CollectionTasks{DeviceId: deviceID}).
		Scan(&tasks); err != nil {
		return gerror.Wrapf(err, "读取设备采集任务失败: %s", deviceID)
	}
	for _, task := range tasks {
		if err := s.pluginHost.StopTask(ctx, task.Id); err != nil {
			return gerror.Wrapf(err, "停止设备采集任务失败: %s", task.Id)
		}
	}

	if err := dao.Devices.Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		if _, err := dao.PluginDeviceConfigVersions.Ctx(ctx).Delete(do.PluginDeviceConfigVersions{DeviceId: deviceID}); err != nil {
			return gerror.Wrap(err, "删除设备插件配置版本失败")
		}
		if _, err := dao.PluginDeviceConfigs.Ctx(ctx).Delete(do.PluginDeviceConfigs{DeviceId: deviceID}); err != nil {
			return gerror.Wrap(err, "删除设备插件配置失败")
		}
		if _, err := dao.DevicePointVersions.Ctx(ctx).Delete(do.DevicePointVersions{DeviceId: deviceID}); err != nil {
			return gerror.Wrap(err, "删除设备点位版本失败")
		}
		if _, err := dao.DevicePoints.Ctx(ctx).Delete(do.DevicePoints{DeviceId: deviceID}); err != nil {
			return gerror.Wrap(err, "删除设备点位失败")
		}
		if _, err := dao.CollectionTasks.Ctx(ctx).Delete(do.CollectionTasks{DeviceId: deviceID}); err != nil {
			return gerror.Wrap(err, "删除设备采集任务失败")
		}
		if _, err := dao.Devices.Ctx(ctx).Delete(do.Devices{Id: deviceID}); err != nil {
			return gerror.Wrap(err, "删除设备失败")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// SavePluginConfig 保存设备维度的通用插件运行配置。
func (s *Service) SavePluginConfig(ctx context.Context, deviceId string, config map[string]any) (map[string]any, error) {
	device, err := s.Get(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	plugin, err := s.ensurePlugin(ctx, device.PluginId)
	if err != nil {
		return nil, err
	}
	points, err := s.devicePoints(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	if err := s.pluginHost.ValidateRuntimeConfig(ctx, plugin.Manifest.Id, config, points); err != nil {
		return nil, err
	}
	configJSON, err := marshalConfig(config)
	if err != nil {
		return nil, err
	}

	var current entity.PluginDeviceConfigs
	if err := dao.PluginDeviceConfigs.Ctx(ctx).Where(do.PluginDeviceConfigs{
		DeviceId: deviceId,
		PluginId: device.PluginId,
		Active:   1,
	}).Scan(&current); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, gerror.Wrapf(err, "读取设备插件配置失败: %s", deviceId)
		}
	}
	nextVersion := current.Version + 1
	currentExists := current.Id != ""
	if current.Id == "" {
		current.Id = "pdc-" + deviceId + "-" + strings.ReplaceAll(device.PluginId, ".", "-")
		nextVersion = 1
	}

	if err := dao.Devices.Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		configData := do.PluginDeviceConfigs{
			Id:         current.Id,
			DeviceId:   deviceId,
			PluginId:   device.PluginId,
			Version:    nextVersion,
			ConfigJson: configJSON,
			ReportMode: device.ReportMode,
			Enabled:    boolInt(device.Enabled == 1),
			Active:     1,
		}
		if currentExists {
			if _, err := dao.PluginDeviceConfigs.Ctx(ctx).
				Where(do.PluginDeviceConfigs{Id: current.Id}).
				Data(configData).
				Update(); err != nil {
				return gerror.Wrap(err, "保存设备插件配置失败")
			}
		} else {
			if _, err := dao.PluginDeviceConfigs.Ctx(ctx).Data(configData).Insert(); err != nil {
				return gerror.Wrap(err, "保存设备插件配置失败")
			}
		}
		if _, err := dao.PluginDeviceConfigVersions.Ctx(ctx).Data(do.PluginDeviceConfigVersions{
			Id:         "pdcv-" + guid.S(),
			ConfigId:   current.Id,
			DeviceId:   deviceId,
			PluginId:   device.PluginId,
			Version:    nextVersion,
			ConfigJson: configJSON,
			ChangeNote: "保存设备插件配置",
		}).Insert(); err != nil {
			return gerror.Wrap(err, "新增设备配置版本失败")
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return emptyAnyMap(config), nil
}

func (s *Service) devicePoints(ctx context.Context, deviceId string) ([]commonv1.PointItem, error) {
	var points []entity.DevicePoints
	if err := dao.DevicePoints.Ctx(ctx).
		Where(do.DevicePoints{DeviceId: deviceId}).
		OrderAsc(dao.DevicePoints.Columns().CreatedAt).
		Scan(&points); err != nil {
		return nil, gerror.Wrapf(err, "读取设备点位失败: %s", deviceId)
	}
	items := make([]commonv1.PointItem, 0, len(points))
	for _, point := range points {
		items = append(items, commonv1.PointItem{
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
		})
	}
	return items, nil
}

// Get 返回指定设备。
func (s *Service) Get(ctx context.Context, deviceId string) (*entity.Devices, error) {
	var item entity.Devices
	if err := dao.Devices.Ctx(ctx).Where(do.Devices{Id: deviceId}).Scan(&item); err != nil {
		return nil, gerror.Wrapf(err, "读取设备失败: %s", deviceId)
	}
	if item.Id == "" {
		return nil, gerror.Newf("设备不存在: %s", deviceId)
	}
	return &item, nil
}

func (s *Service) ensureGroup(ctx context.Context, groupId string) error {
	var group entity.DeviceGroups
	if err := dao.DeviceGroups.Ctx(ctx).Where(do.DeviceGroups{Id: groupId}).Scan(&group); err != nil {
		return gerror.Wrapf(err, "读取设备分组失败: %s", groupId)
	}
	if group.Id == "" {
		return gerror.Newf("设备分组不存在: %s", groupId)
	}
	return nil
}

func (s *Service) ensurePlugin(ctx context.Context, pluginId string) (*pluginhostsvc.RuntimePlugin, error) {
	return s.pluginHost.Plugin(ctx, pluginId)
}

func (s *Service) buildDeviceItem(ctx context.Context, device entity.Devices) (*commonv1.DeviceItem, error) {
	plugin, err := s.ensurePlugin(ctx, device.PluginId)
	if err != nil {
		return nil, err
	}
	pointCount, err := dao.DevicePoints.Ctx(ctx).Where(do.DevicePoints{DeviceId: device.Id}).Count()
	if err != nil {
		return nil, gerror.Wrapf(err, "读取设备点位数量失败: %s", device.Id)
	}
	item := s.toDeviceItem(device, plugin.Manifest.Name, pointCount)
	return &item, nil
}

func (s *Service) toDeviceItem(device entity.Devices, pluginName string, pointCount int) commonv1.DeviceItem {
	return commonv1.DeviceItem{
		Id:          device.Id,
		Name:        device.Name,
		GroupId:     device.GroupId,
		PluginId:    device.PluginId,
		PluginName:  pluginName,
		Status:      device.Status,
		Enabled:     device.Enabled == 1,
		PointCount:  pointCount,
		ReportMode:  device.ReportMode,
		LastSeenAt:  displayTime(device.LastSeenAt, "尚未连接"),
		Description: device.Description,
	}
}

func marshalConfig(config map[string]any) (string, error) {
	content, err := json.Marshal(emptyAnyMap(config))
	if err != nil {
		return "", gerror.Wrap(err, "序列化设备插件配置失败")
	}
	return string(content), nil
}

func initialConfigVersion(hasConfig bool) int {
	if hasConfig {
		return 1
	}
	return 0
}

func initialConfigJSON(hasConfig bool, configJSON string) string {
	if hasConfig {
		return configJSON
	}
	return ""
}

func ensureDeleteConditionValue(name string, value string) error {
	if strings.TrimSpace(value) == "" {
		return gerror.Newf("删除条件不能为空: %s", name)
	}
	return nil
}

func displayTime(value string, empty string) string {
	if strings.TrimSpace(value) == "" {
		return empty
	}
	return value
}

func emptyAnyMap(values map[string]any) map[string]any {
	if values == nil {
		return map[string]any{}
	}
	return values
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

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

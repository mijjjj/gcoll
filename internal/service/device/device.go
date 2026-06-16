package device

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

// Service 提供设备控制面服务。
type Service struct{}

// New 创建设备服务。
func New() *Service {
	return &Service{}
}

// List 返回设备分组和设备列表。
func (s *Service) List(ctx context.Context) (*runtimev1.DevicesRes, error) {
	var (
		groups  []entity.DeviceGroups
		devices []entity.Devices
		plugins []entity.Plugins
		points  []entity.DevicePoints
	)
	if err := dao.DeviceGroups.Ctx(ctx).OrderAsc(dao.DeviceGroups.Columns().SortOrder).Scan(&groups); err != nil {
		return nil, gerror.Wrap(err, "读取设备分组失败")
	}
	if err := dao.Devices.Ctx(ctx).OrderAsc(dao.Devices.Columns().CreatedAt).Scan(&devices); err != nil {
		return nil, gerror.Wrap(err, "读取设备列表失败")
	}
	if err := dao.Plugins.Ctx(ctx).Scan(&plugins); err != nil {
		return nil, gerror.Wrap(err, "读取插件列表失败")
	}
	if err := dao.DevicePoints.Ctx(ctx).Scan(&points); err != nil {
		return nil, gerror.Wrap(err, "读取设备点位数量失败")
	}

	pluginNames := make(map[string]string, len(plugins))
	for _, plugin := range plugins {
		pluginNames[plugin.Id] = plugin.Name
	}
	pointCounts := make(map[string]int)
	for _, point := range points {
		pointCounts[point.DeviceId]++
	}
	groupCounts := make(map[string]int)
	for _, device := range devices {
		groupCounts[device.GroupId]++
	}

	res := &runtimev1.DevicesRes{
		Groups: make([]runtimev1.DeviceGroup, 0, len(groups)),
		Items:  make([]runtimev1.DeviceItem, 0, len(devices)),
	}
	for _, group := range groups {
		res.Groups = append(res.Groups, runtimev1.DeviceGroup{
			Id:    group.Id,
			Name:  group.Name,
			Count: groupCounts[group.Id],
		})
	}
	for _, device := range devices {
		res.Items = append(res.Items, runtimev1.DeviceItem{
			Id:          device.Id,
			Name:        device.Name,
			Code:        device.Code,
			GroupId:     device.GroupId,
			PluginId:    device.PluginId,
			PluginName:  pluginNames[device.PluginId],
			Status:      device.Status,
			Enabled:     device.Enabled == 1,
			PointCount:  pointCounts[device.Id],
			ReportMode:  device.ReportMode,
			LastSeenAt:  displayTime(device.LastSeenAt, "尚未连接"),
			Description: device.Description,
		})
	}
	return res, nil
}

// Create 新增设备并保存设备插件配置。
func (s *Service) Create(ctx context.Context, req *runtimev1.CreateDeviceReq) (*runtimev1.DeviceItem, error) {
	if err := s.ensureGroup(ctx, req.GroupId); err != nil {
		return nil, err
	}
	plugin, err := s.ensurePlugin(ctx, req.PluginId)
	if err != nil {
		return nil, err
	}
	if err := validateDeviceConfig(req.PluginId, req.Config); err != nil {
		return nil, err
	}

	deviceId := strings.TrimSpace(req.Id)
	if deviceId == "" {
		deviceId = "dev-" + guid.S()
	}
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return nil, gerror.Wrap(err, "序列化设备插件配置失败")
	}
	configId := "pdc-" + deviceId + "-" + strings.ReplaceAll(req.PluginId, ".", "-")
	configVersionId := configId + "-1"

	if err := dao.Devices.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_ = ctx
		if _, err := tx.Insert(dao.Devices.Table(), do.Devices{
			Id:          deviceId,
			Name:        req.Name,
			Code:        req.Code,
			GroupId:     req.GroupId,
			PluginId:    req.PluginId,
			Status:      "offline",
			Enabled:     boolInt(req.Enabled),
			ReportMode:  req.ReportMode,
			Description: req.Description,
		}); err != nil {
			return gerror.Wrap(err, "新增设备失败")
		}
		if _, err := tx.Insert(dao.PluginDeviceConfigs.Table(), do.PluginDeviceConfigs{
			Id:         configId,
			DeviceId:   deviceId,
			PluginId:   req.PluginId,
			Version:    1,
			ConfigJson: string(configJSON),
			ReportMode: req.ReportMode,
			Enabled:    boolInt(req.Enabled),
			Active:     1,
		}); err != nil {
			return gerror.Wrap(err, "新增设备插件配置失败")
		}
		if _, err := tx.Insert(dao.PluginDeviceConfigVersions.Table(), do.PluginDeviceConfigVersions{
			Id:         configVersionId,
			ConfigId:   configId,
			DeviceId:   deviceId,
			PluginId:   req.PluginId,
			Version:    1,
			ConfigJson: string(configJSON),
			ChangeNote: "初始化设备插件配置",
		}); err != nil {
			return gerror.Wrap(err, "新增设备配置版本失败")
		}
		if req.PluginId == modbusPluginId {
			if _, err := tx.Insert(dao.ModbusTcpDeviceProfiles.Table(), do.ModbusTcpDeviceProfiles{
				Id:               "mdp-" + deviceId,
				DeviceId:         deviceId,
				PluginId:         req.PluginId,
				Version:          1,
				Host:             gconv.String(req.Config["host"]),
				Port:             gconv.Int(req.Config["port"]),
				UnitId:           gconv.Int(req.Config["unitId"]),
				TimeoutMs:        gconv.Int(req.Config["timeoutMs"]),
				PollIntervalMs:   gconv.Int(req.Config["pollIntervalMs"]),
				ReportMode:       req.ReportMode,
				MaxCoilBatch:     gconv.Int(req.Config["maxCoilBatch"]),
				MaxRegisterBatch: gconv.Int(req.Config["maxRegisterBatch"]),
				LowLatencyMs:     gconv.Int(req.Config["lowLatencyMs"]),
				HighLatencyMs:    gconv.Int(req.Config["highLatencyMs"]),
				DebugEnabled:     boolInt(gconv.Bool(req.Config["debugEnabled"])),
				Enabled:          boolInt(req.Enabled),
			}); err != nil {
				return gerror.Wrap(err, "新增 Modbus TCP 设备配置失败")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &runtimev1.DeviceItem{
		Id:          deviceId,
		Name:        req.Name,
		Code:        req.Code,
		GroupId:     req.GroupId,
		PluginId:    req.PluginId,
		PluginName:  plugin.Name,
		Status:      "offline",
		Enabled:     req.Enabled,
		PointCount:  0,
		ReportMode:  req.ReportMode,
		LastSeenAt:  "尚未连接",
		Description: req.Description,
	}, nil
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

func (s *Service) ensurePlugin(ctx context.Context, pluginId string) (*entity.Plugins, error) {
	var plugin entity.Plugins
	if err := dao.Plugins.Ctx(ctx).Where(do.Plugins{Id: pluginId}).Scan(&plugin); err != nil {
		return nil, gerror.Wrapf(err, "读取插件失败: %s", pluginId)
	}
	if plugin.Id == "" {
		return nil, gerror.Newf("插件不存在: %s", pluginId)
	}
	if plugin.Enabled != 1 {
		return nil, gerror.Newf("插件未启用: %s", pluginId)
	}
	return &plugin, nil
}

func validateDeviceConfig(pluginId string, config map[string]any) error {
	if config == nil {
		return gerror.New("设备插件配置不能为空")
	}
	if pluginId != modbusPluginId {
		return nil
	}
	if strings.TrimSpace(gconv.String(config["host"])) == "" {
		return gerror.New("Modbus TCP 配置缺少 host")
	}
	port := gconv.Int(config["port"])
	if port < 1 || port > 65535 {
		return gerror.New("Modbus TCP 端口必须在 1 到 65535 之间")
	}
	unitId := gconv.Int(config["unitId"])
	if unitId < 0 || unitId > 247 {
		return gerror.New("Modbus TCP unitId 必须在 0 到 247 之间")
	}
	if gconv.Int(config["timeoutMs"]) < 100 {
		return gerror.New("Modbus TCP timeoutMs 不能小于 100")
	}
	if gconv.Int(config["pollIntervalMs"]) < 100 {
		return gerror.New("Modbus TCP pollIntervalMs 不能小于 100")
	}
	if gconv.Int(config["maxCoilBatch"]) < 1 || gconv.Int(config["maxCoilBatch"]) > 2000 {
		return gerror.New("Modbus TCP maxCoilBatch 必须在 1 到 2000 之间")
	}
	if gconv.Int(config["maxRegisterBatch"]) < 1 || gconv.Int(config["maxRegisterBatch"]) > 125 {
		return gerror.New("Modbus TCP maxRegisterBatch 必须在 1 到 125 之间")
	}
	return nil
}

func displayTime(value string, empty string) string {
	if strings.TrimSpace(value) == "" {
		return empty
	}
	return value
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

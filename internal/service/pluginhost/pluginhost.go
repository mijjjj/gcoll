package pluginhost

import (
	"archive/zip"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"gopkg.in/yaml.v3"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
	"github.com/mijjjj/gcoll/internal/dao"
	"github.com/mijjjj/gcoll/internal/model/do"
	"github.com/mijjjj/gcoll/internal/model/entity"
	modbustcp "github.com/mijjjj/gcoll/plugins/com.gcoll.modbus-tcp/runtime"
)

const (
	pluginRootDir           = "plugins"
	pluginModeDevelopment   = "development"
	pluginModeProduction    = "production"
	configPluginMode        = "pluginHost.mode"
	configPluginAutoCompile = "pluginHost.autoCompile"
)

var defaultService = New()

// Manifest 描述插件目录中的 plugin.yaml。
type Manifest struct {
	SchemaVersion    int               `yaml:"schemaVersion" json:"schemaVersion"`
	Id               string            `yaml:"id"            json:"id"`
	Name             string            `yaml:"name"          json:"name"`
	Type             string            `yaml:"type"          json:"type"`
	Version          string            `yaml:"version"       json:"version"`
	Runtime          string            `yaml:"runtime"       json:"runtime"`
	Protocol         string            `yaml:"protocol"      json:"protocol"`
	Entry            map[string]string `yaml:"entry" json:"entry"`
	Capabilities     []string          `yaml:"capabilities"  json:"capabilities"`
	Permissions      []string          `yaml:"permissions"   json:"permissions"`
	ConfigSchema     map[string]any    `yaml:"configSchema"  json:"configSchema"`
	CustomConfigPage CustomAssetPage   `yaml:"customConfigPage" json:"customConfigPage"`
	CustomPointPage  CustomAssetPage   `yaml:"customPointPage"  json:"customPointPage"`
}

// CustomAssetPage 描述插件自带配置页面资源。
type CustomAssetPage struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Entry   string `yaml:"entry"   json:"entry"`
	Script  string `yaml:"script"  json:"script"`
}

// RuntimePlugin 描述内存中的插件运行时记录。
type RuntimePlugin struct {
	Manifest  Manifest
	Directory string
	Source    string
	EntryPath string
	Status    string
	UpdatedAt string
	Cmd       *exec.Cmd
}

// LatestPointValue 描述内存中的最新点位值。
type LatestPointValue struct {
	PointID   string
	DeviceID  string
	PointName string
	Value     string
	Quality   string
	Changed   bool
	UpdatedAt string
}

// TestConnectionResult 描述插件连接测试结果。
type TestConnectionResult struct {
	Success   bool
	Message   string
	LatencyMs int
	TraceID   string
}

// CollectionRecord 描述南向插件提交给宿主的通用采集记录。
type CollectionRecord struct {
	PointID     string
	DeviceID    string
	TaskID      string
	Quality     string
	Value       any
	Changed     bool
	CollectedAt time.Time
	TraceID     string
}

// Service 管理插件目录、插件进程、任务状态和宿主侧数据交换。
type Service struct {
	mu          sync.Mutex
	loaded      bool
	registry    map[string]*RuntimePlugin
	taskCancels map[string]context.CancelFunc
	latest      map[string]LatestPointValue
}

// New 创建插件宿主服务。
func New() *Service {
	return &Service{
		registry:    map[string]*RuntimePlugin{},
		taskCancels: map[string]context.CancelFunc{},
		latest:      map[string]LatestPointValue{},
	}
}

// Instance 返回进程内插件宿主单例。
func Instance() *Service {
	return defaultService
}

// ResetForTest 重置测试用插件宿主内存注册表。
func ResetForTest() {
	defaultService.mu.Lock()
	defer defaultService.mu.Unlock()
	defaultService.loaded = false
	defaultService.registry = map[string]*RuntimePlugin{}
	defaultService.taskCancels = map[string]context.CancelFunc{}
	defaultService.latest = map[string]LatestPointValue{}
}

// Init 扫描插件目录并启动插件进程。
func Init(ctx context.Context) error {
	return defaultService.LoadAndStart(ctx)
}

// LoadAndStart 从插件目录加载插件到内存并启动进程。
func (s *Service) LoadAndStart(ctx context.Context) error {
	settings, err := loadSettings(ctx)
	if err != nil {
		return err
	}

	loaded := map[string]*RuntimePlugin{}
	plugins, err := s.loadPluginDirectories(ctx, settings)
	if err != nil {
		return err
	}
	for _, plugin := range plugins {
		loaded[plugin.Manifest.Id] = plugin
	}

	s.mu.Lock()
	s.registry = loaded
	s.loaded = true
	s.mu.Unlock()

	for _, plugin := range loaded {
		if err := s.startRuntimeProcess(ctx, plugin); err != nil {
			plugin.Status = "failed"
			_ = s.appendRuntimeEvent(ctx, "ERROR", "pluginhost", plugin.Manifest.Id, "", "", err.Error(), "")
			return err
		}
	}
	return nil
}

// List 返回内存插件注册表。
func (s *Service) List(ctx context.Context) ([]commonv1.PluginItem, error) {
	if err := s.ensureLoaded(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	items := make([]commonv1.PluginItem, 0, len(s.registry))
	for _, plugin := range s.registry {
		items = append(items, toPluginItem(plugin))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Id < items[j].Id
	})
	return items, nil
}

// Plugin 返回指定插件的内存记录。
func (s *Service) Plugin(ctx context.Context, pluginID string) (*RuntimePlugin, error) {
	if err := s.ensureLoaded(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	plugin := s.registry[pluginID]
	if plugin == nil {
		return nil, gerror.Newf("插件未加载: %s", pluginID)
	}
	return plugin, nil
}

// PluginNameMap 返回插件 ID 到名称的映射。
func (s *Service) PluginNameMap(ctx context.Context) (map[string]string, error) {
	items, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	names := make(map[string]string, len(items))
	for _, item := range items {
		names[item.Id] = item.Name
	}
	return names, nil
}

// CustomConfigPageContent 读取插件自带设备配置页内容。
func (s *Service) CustomConfigPageContent(ctx context.Context, pluginID string) (CustomAssetPage, string, string, error) {
	plugin, err := s.Plugin(ctx, pluginID)
	if err != nil {
		return CustomAssetPage{}, "", "", err
	}
	page := plugin.Manifest.CustomConfigPage
	if !page.Enabled {
		return page, "", "", nil
	}
	html, err := readPluginAsset(plugin.Directory, page.Entry)
	if err != nil {
		return page, "", "", err
	}
	js := ""
	if strings.TrimSpace(page.Script) != "" {
		js, err = readPluginAsset(plugin.Directory, page.Script)
		if err != nil {
			return page, "", "", err
		}
	}
	return page, html, js, nil
}

// CustomPointPageContent 读取插件自带点位配置页内容。
func (s *Service) CustomPointPageContent(ctx context.Context, pluginID string) (CustomAssetPage, string, string, error) {
	plugin, err := s.Plugin(ctx, pluginID)
	if err != nil {
		return CustomAssetPage{}, "", "", err
	}
	page := plugin.Manifest.CustomPointPage
	if !page.Enabled {
		return page, "", "", nil
	}
	html, err := readPluginAsset(plugin.Directory, page.Entry)
	if err != nil {
		return page, "", "", err
	}
	js := ""
	if strings.TrimSpace(page.Script) != "" {
		js, err = readPluginAsset(plugin.Directory, page.Script)
		if err != nil {
			return page, "", "", err
		}
	}
	return page, html, js, nil
}

// ValidateRuntimeConfig 执行保存前插件配置校验。
func (s *Service) ValidateRuntimeConfig(ctx context.Context, pluginID string, config map[string]any, points []commonv1.PointItem) error {
	plugin, err := s.Plugin(ctx, pluginID)
	if err != nil {
		return err
	}
	if plugin.Status != "running" {
		return gerror.Newf("插件校验服务未运行: %s", pluginID)
	}
	if err := validateConfigSchema(plugin.Manifest.ConfigSchema, config); err != nil {
		return err
	}
	for _, point := range points {
		if strings.TrimSpace(point.Id) == "" {
			return gerror.New("点位 ID 不能为空")
		}
		if point.PluginId != pluginID {
			return gerror.Newf("点位插件与设备插件不一致: %s", point.Id)
		}
		if strings.TrimSpace(point.Name) == "" {
			return gerror.Newf("点位名称不能为空: %s", point.Id)
		}
		if strings.TrimSpace(point.Address) == "" {
			return gerror.Newf("点位地址不能为空: %s", point.Id)
		}
	}
	return nil
}

// ImportPackage 将插件包解压或复制到插件目录，加载到内存并启动插件进程。
func (s *Service) ImportPackage(ctx context.Context, packagePath string) (*commonv1.PluginItem, error) {
	if strings.TrimSpace(packagePath) == "" {
		return nil, gerror.New("插件包路径不能为空")
	}
	settings, err := loadSettings(ctx)
	if err != nil {
		return nil, err
	}
	manifest, err := readManifest(packagePath)
	if err != nil {
		return nil, err
	}
	if err := validateManifest(manifest); err != nil {
		return nil, err
	}
	targetDir := filepath.Join(settings.RootDir, manifest.Id)
	if err := replacePluginDirectory(packagePath, targetDir, settings.RootDir); err != nil {
		return nil, err
	}
	plugin, err := s.loadPluginDirectory(ctx, targetDir, settings)
	if err != nil {
		return nil, err
	}
	if err := s.startRuntimeProcess(ctx, plugin); err != nil {
		plugin.Status = "failed"
		return nil, err
	}
	s.mu.Lock()
	s.registry[plugin.Manifest.Id] = plugin
	s.loaded = true
	s.mu.Unlock()
	item := toPluginItem(plugin)
	return &item, nil
}

// TestConnection 通过通用南向插件连接测试入口发起测试。
func (s *Service) TestConnection(ctx context.Context, deviceID string, config map[string]any) (*TestConnectionResult, error) {
	start := time.Now()
	traceID := "trace-" + guid.S()
	device, storedConfig, err := s.deviceAndConfig(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if err := s.ensurePluginStarted(ctx, device.PluginId); err != nil {
		return nil, err
	}

	testConfig, err := s.testConnectionConfig(ctx, device, storedConfig, config)
	if err != nil {
		return nil, err
	}

	if device.PluginId == modbustcp.PluginID {
		result, testErr := s.testModbusConnection(ctx, deviceID, device.PluginId, testConfig, traceID)
		if testErr != nil {
			return result, testErr
		}
		return result, nil
	}

	message := "插件 gRPC TestConnection 尚未接入，宿主已完成设备配置和插件进程检查。"
	result := &TestConnectionResult{
		Success:   false,
		Message:   message,
		LatencyMs: int(time.Since(start).Milliseconds()),
		TraceID:   traceID,
	}
	_ = s.appendRuntimeEvent(ctx, "WARN", "pluginhost", device.PluginId, deviceID, "", message, traceID)
	return result, gerror.New(message)
}

func (s *Service) testConnectionConfig(ctx context.Context, device *entity.Devices, storedConfig *entity.PluginDeviceConfigs, draftConfig map[string]any) (map[string]any, error) {
	if len(draftConfig) > 0 {
		if err := s.ValidateRuntimeConfig(ctx, device.PluginId, draftConfig, nil); err != nil {
			return nil, err
		}
		return draftConfig, nil
	}
	if storedConfig == nil || strings.TrimSpace(storedConfig.ConfigJson) == "" {
		return nil, gerror.Newf("设备缺少插件运行配置: %s", device.Id)
	}
	values := map[string]any{}
	if err := json.Unmarshal([]byte(storedConfig.ConfigJson), &values); err != nil {
		return nil, gerror.Wrapf(err, "解析设备插件配置失败: %s", device.Id)
	}
	if err := s.ValidateRuntimeConfig(ctx, device.PluginId, values, nil); err != nil {
		return nil, err
	}
	return values, nil
}

func (s *Service) testModbusConnection(ctx context.Context, deviceID string, pluginID string, config map[string]any, traceID string) (*TestConnectionResult, error) {
	connectionConfig, err := modbusConnectionConfig(config)
	if err != nil {
		return nil, err
	}
	client, err := modbustcp.NewTCPClient(connectionConfig)
	if err != nil {
		return nil, err
	}
	connectStartedAt := time.Now()
	if err := client.Connect(); err != nil {
		latency := int(time.Since(connectStartedAt).Milliseconds())
		message := "Modbus TCP 连接测试失败: " + err.Error()
		result := &TestConnectionResult{
			Success:   false,
			Message:   message,
			LatencyMs: latency,
			TraceID:   traceID,
		}
		_ = s.appendRuntimeEvent(ctx, "WARN", "pluginhost", pluginID, deviceID, "", message, traceID)
		return result, gerror.New(message)
	}
	defer client.Close()

	latency := int(time.Since(connectStartedAt).Milliseconds())
	message := "Modbus TCP 连接测试成功。"
	result := &TestConnectionResult{
		Success:   true,
		Message:   message,
		LatencyMs: latency,
		TraceID:   traceID,
	}
	_ = s.appendRuntimeEvent(ctx, "INFO", "pluginhost", pluginID, deviceID, "", message, traceID)
	return result, nil
}

func modbusConnectionConfig(config map[string]any) (modbustcp.ConnectionConfig, error) {
	content, err := json.Marshal(config)
	if err != nil {
		return modbustcp.ConnectionConfig{}, gerror.Wrap(err, "序列化 Modbus TCP 测试配置失败")
	}
	var connectionConfig modbustcp.ConnectionConfig
	if err := json.Unmarshal(content, &connectionConfig); err != nil {
		return modbustcp.ConnectionConfig{}, gerror.Wrap(err, "解析 Modbus TCP 测试配置失败")
	}
	if err := connectionConfig.Validate(); err != nil {
		return modbustcp.ConnectionConfig{}, gerror.Wrap(err, "校验 Modbus TCP 测试配置失败")
	}
	return connectionConfig.Normalize(), nil
}

// StartTask 启动采集任务控制面状态，实际采集必须由南向插件通过 gRPC 完成。
func (s *Service) StartTask(ctx context.Context, taskID string) error {
	task, err := s.task(ctx, taskID)
	if err != nil {
		return err
	}
	device, _, err := s.deviceRuntimeConfig(ctx, task.DeviceId)
	if err != nil {
		return err
	}
	if device.PluginId != task.SouthPluginId {
		return gerror.Newf("采集任务插件与设备插件不一致: %s", taskID)
	}
	if err := s.ensurePluginStarted(ctx, task.SouthPluginId); err != nil {
		return err
	}
	if err := s.ensureDeviceRunnable(ctx, task.DeviceId); err != nil {
		return err
	}

	s.mu.Lock()
	if _, exists := s.taskCancels[taskID]; exists {
		s.mu.Unlock()
		return nil
	}
	runCtx, cancel := context.WithCancel(context.Background())
	s.taskCancels[taskID] = cancel
	s.mu.Unlock()

	if err := s.markTaskRunning(ctx, task); err != nil {
		cancel()
		s.mu.Lock()
		delete(s.taskCancels, taskID)
		s.mu.Unlock()
		return err
	}
	go s.waitTaskCancel(runCtx, taskID)
	return nil
}

// StopTask 停止采集任务循环。
func (s *Service) StopTask(ctx context.Context, taskID string) error {
	task, err := s.task(ctx, taskID)
	if err != nil {
		return err
	}
	s.mu.Lock()
	cancel := s.taskCancels[taskID]
	delete(s.taskCancels, taskID)
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if _, err := dao.CollectionTasks.Ctx(ctx).
		Where(do.CollectionTasks{Id: taskID}).
		Data(do.CollectionTasks{Status: "stopped", Rate: "0 条/秒"}).
		Update(); err != nil {
		return gerror.Wrapf(err, "停止采集任务失败: %s", taskID)
	}
	return s.appendRuntimeEvent(ctx, "INFO", "collector", task.SouthPluginId, task.DeviceId, taskID, "采集任务已停止。", "")
}

// LatestValues 返回内存中的最新点位值。
func (s *Service) LatestValues() []LatestPointValue {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := make([]LatestPointValue, 0, len(s.latest))
	for _, item := range s.latest {
		items = append(items, item)
	}
	return items
}

// SubmitRecords 接收插件采集结果，只更新内存最新值和低频运行状态。
func (s *Service) SubmitRecords(ctx context.Context, records []CollectionRecord) error {
	if len(records) == 0 {
		return nil
	}
	pointNames, err := pointNames(ctx, records)
	if err != nil {
		return err
	}
	now := gtime.Now().Format("Y-m-d H:i:s.u")
	s.mu.Lock()
	for _, record := range records {
		s.latest[record.PointID] = LatestPointValue{
			PointID:   record.PointID,
			DeviceID:  record.DeviceID,
			PointName: pointNames[record.PointID],
			Value:     fmt.Sprint(record.Value),
			Quality:   record.Quality,
			Changed:   record.Changed,
			UpdatedAt: record.CollectedAt.Format("2006-01-02 15:04:05.000"),
		}
	}
	s.mu.Unlock()

	first := records[0]
	if _, err := dao.Devices.Ctx(ctx).
		Where(do.Devices{Id: first.DeviceID}).
		Data(do.Devices{Status: "online", LastSeenAt: now}).
		Update(); err != nil {
		return gerror.Wrapf(err, "更新设备采集状态失败: %s", first.DeviceID)
	}
	return s.appendRuntimeEvent(ctx, "INFO", "collector", "", first.DeviceID, first.TaskID, fmt.Sprintf("已接收 %d 条采集记录并写入内存缓存。", len(records)), first.TraceID)
}

func (s *Service) ensureLoaded(ctx context.Context) error {
	s.mu.Lock()
	loaded := s.loaded
	s.mu.Unlock()
	if loaded {
		return nil
	}
	settings, err := loadSettings(ctx)
	if err != nil {
		return err
	}
	plugins := map[string]*RuntimePlugin{}
	items, err := s.loadPluginDirectories(ctx, settings.withoutCompile())
	if err != nil {
		return err
	}
	for _, plugin := range items {
		plugins[plugin.Manifest.Id] = plugin
	}
	s.mu.Lock()
	if !s.loaded {
		s.registry = plugins
		s.loaded = true
	}
	s.mu.Unlock()
	return nil
}

func (s *Service) ensurePluginStarted(ctx context.Context, pluginID string) error {
	plugin, err := s.Plugin(ctx, pluginID)
	if err != nil {
		return err
	}
	if plugin.Status == "running" && plugin.Cmd != nil && plugin.Cmd.Process != nil {
		return nil
	}
	return s.startRuntimeProcess(ctx, plugin)
}

func (s *Service) loadPluginDirectories(ctx context.Context, settings hostSettings) ([]*RuntimePlugin, error) {
	if strings.TrimSpace(settings.RootDir) == "" {
		return nil, gerror.New("插件目录不能为空")
	}
	if !gfile.IsDir(settings.RootDir) {
		return nil, gerror.Newf("插件目录不存在: %s", settings.RootDir)
	}
	entries, err := os.ReadDir(settings.RootDir)
	if err != nil {
		return nil, gerror.Wrapf(err, "读取插件目录失败: %s", settings.RootDir)
	}
	plugins := make([]*RuntimePlugin, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(settings.RootDir, entry.Name())
		if !gfile.IsFile(filepath.Join(dir, "plugin.yaml")) {
			continue
		}
		plugin, err := s.loadPluginDirectory(ctx, dir, settings)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

func (s *Service) loadPluginDirectory(ctx context.Context, dir string, settings hostSettings) (*RuntimePlugin, error) {
	manifest, err := readManifest(filepath.Join(dir, "plugin.yaml"))
	if err != nil {
		return nil, err
	}
	if err := validateManifest(manifest); err != nil {
		return nil, err
	}
	entryPath := ""
	if settings.Mode == pluginModeDevelopment {
		if settings.AutoCompile {
			entryPath, err = buildDevelopmentPlugin(ctx, dir, manifest)
			if err != nil {
				return nil, err
			}
		} else {
			entryPath, err = installedEntryPath(dir, manifest)
			if err != nil {
				return nil, err
			}
		}
	} else {
		entryPath, err = installedEntryPath(dir, manifest)
		if err != nil {
			return nil, err
		}
	}
	return &RuntimePlugin{
		Manifest:  *manifest,
		Directory: dir,
		Source:    settings.Mode,
		EntryPath: entryPath,
		Status:    "loaded",
		UpdatedAt: gtime.Now().Format("Y-m-d H:i:s"),
	}, nil
}

func (s *Service) startRuntimeProcess(ctx context.Context, plugin *RuntimePlugin) error {
	if plugin.EntryPath == "" && plugin.Source == pluginModeDevelopment {
		entryPath, err := buildDevelopmentPlugin(ctx, plugin.Directory, &plugin.Manifest)
		if err != nil {
			return err
		}
		plugin.EntryPath = entryPath
	}
	if strings.TrimSpace(plugin.EntryPath) == "" {
		return gerror.Newf("插件缺少启动文件: %s", plugin.Manifest.Id)
	}
	if !gfile.IsFile(plugin.EntryPath) {
		return gerror.Newf("插件启动文件不存在: %s", plugin.EntryPath)
	}

	s.mu.Lock()
	if plugin.Cmd != nil && plugin.Cmd.Process != nil && plugin.Status == "running" {
		s.mu.Unlock()
		return nil
	}
	cmd := exec.CommandContext(context.Background(), plugin.EntryPath)
	cmd.Dir = plugin.Directory
	plugin.Cmd = cmd
	plugin.Status = "starting"
	s.mu.Unlock()

	if err := cmd.Start(); err != nil {
		plugin.Status = "failed"
		return gerror.Wrapf(err, "启动插件进程失败: %s", plugin.EntryPath)
	}
	plugin.Status = "running"
	go func() {
		err := cmd.Wait()
		s.mu.Lock()
		defer s.mu.Unlock()
		if current := s.registry[plugin.Manifest.Id]; current != nil && current.Cmd == cmd {
			if err != nil {
				current.Status = "failed"
			} else {
				current.Status = "stopped"
			}
		}
	}()
	return s.appendRuntimeEvent(ctx, "INFO", "pluginhost", plugin.Manifest.Id, "", "", "插件进程已启动。", "")
}

func buildDevelopmentPlugin(ctx context.Context, dir string, manifest *Manifest) (string, error) {
	mainPath, err := findPluginMain(dir)
	if err != nil {
		return "", err
	}
	entry := strings.TrimSpace(entryForPlatform(manifest))
	if entry == "" {
		name := strings.TrimSuffix(filepath.Base(dir), filepath.Ext(filepath.Base(dir)))
		if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
			name += ".exe"
		}
		entry = filepath.Join("bin", platformKey(), name)
	}
	outputPath := filepath.Clean(filepath.Join(dir, entry))
	if !isPathInside(outputPath, dir) {
		return "", gerror.Newf("开发插件构建输出路径越界: %s", entry)
	}
	if err := gfile.Mkdir(filepath.Dir(outputPath)); err != nil {
		return "", gerror.Wrapf(err, "创建插件构建输出目录失败: %s", filepath.Dir(outputPath))
	}
	cmd := exec.CommandContext(ctx, "go", "build", "-trimpath", "-o", outputPath, mainPath)
	cmd.Dir = dir
	content, err := cmd.CombinedOutput()
	if err != nil {
		return "", gerror.Wrapf(err, "编译开发插件失败: %s\n%s", manifest.Id, strings.TrimSpace(string(content)))
	}
	return outputPath, nil
}

func findPluginMain(dir string) (string, error) {
	mainPath := filepath.Join(dir, "main.go")
	if !gfile.IsFile(mainPath) {
		return "", gerror.Newf("开发插件缺少根目录 main.go: %s", dir)
	}
	return mainPath, nil
}

func installedEntryPath(dir string, manifest *Manifest) (string, error) {
	entry := entryForPlatform(manifest)
	if strings.TrimSpace(entry) == "" {
		return "", gerror.Newf("插件清单缺少当前平台启动入口: %s", manifest.Id)
	}
	path := filepath.Clean(filepath.Join(dir, entry))
	if !isPathInside(path, dir) {
		return "", gerror.Newf("插件启动入口越界: %s", entry)
	}
	return path, nil
}

func entryForPlatform(manifest *Manifest) string {
	if manifest.Entry == nil {
		return ""
	}
	key := platformKey()
	if value := manifest.Entry[key]; value != "" {
		return value
	}
	return ""
}

func platformKey() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}

func replacePluginDirectory(packagePath string, targetDir string, rootDir string) error {
	if !isPathInside(targetDir, rootDir) || filepath.Clean(targetDir) == filepath.Clean(rootDir) {
		return gerror.Newf("插件安装目录不在允许范围内: %s", targetDir)
	}
	if err := os.RemoveAll(targetDir); err != nil {
		return gerror.Wrapf(err, "清理插件安装目录失败: %s", targetDir)
	}
	if err := gfile.Mkdir(targetDir); err != nil {
		return gerror.Wrapf(err, "创建插件安装目录失败: %s", targetDir)
	}
	if strings.EqualFold(filepath.Ext(packagePath), ".gcpkg") || strings.EqualFold(filepath.Ext(packagePath), ".zip") {
		return unzipPluginPackage(packagePath, targetDir)
	}
	stat, err := os.Stat(packagePath)
	if err != nil {
		return gerror.Wrapf(err, "读取插件源失败: %s", packagePath)
	}
	if !stat.IsDir() {
		return gerror.Newf("插件导入源必须是 .gcpkg 或目录: %s", packagePath)
	}
	return copyDirectory(packagePath, targetDir)
}

func unzipPluginPackage(packagePath string, targetDir string) error {
	reader, err := zip.OpenReader(packagePath)
	if err != nil {
		return gerror.Wrapf(err, "打开插件包失败: %s", packagePath)
	}
	defer reader.Close()
	for _, file := range reader.File {
		target := filepath.Clean(filepath.Join(targetDir, file.Name))
		if !isPathInside(target, targetDir) {
			return gerror.Newf("插件包包含越界路径: %s", file.Name)
		}
		if file.FileInfo().IsDir() {
			if err := gfile.Mkdir(target); err != nil {
				return err
			}
			continue
		}
		if err := gfile.Mkdir(filepath.Dir(target)); err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return gerror.Wrapf(err, "读取插件包文件失败: %s", file.Name)
		}
		if err := writeFileFromReader(target, src); err != nil {
			_ = src.Close()
			return err
		}
		_ = src.Close()
	}
	return nil
}

func copyDirectory(source string, target string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dst := filepath.Join(target, rel)
		if entry.IsDir() {
			return gfile.Mkdir(dst)
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		return writeFileFromReader(dst, src)
	})
}

func writeFileFromReader(path string, reader io.Reader) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return gerror.Wrapf(err, "写入插件文件失败: %s", path)
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	return gerror.Wrapf(err, "写入插件文件失败: %s", path)
}

func isPathInside(path string, root string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && !filepath.IsAbs(rel))
}

func readManifest(packagePath string) (*Manifest, error) {
	if stat, err := os.Stat(packagePath); err == nil && stat.IsDir() {
		return readManifestFile(filepath.Join(packagePath, "plugin.yaml"))
	}
	if strings.EqualFold(filepath.Ext(packagePath), ".gcpkg") || strings.EqualFold(filepath.Ext(packagePath), ".zip") {
		return readManifestZip(packagePath)
	}
	return readManifestFile(packagePath)
}

func readManifestFile(path string) (*Manifest, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, gerror.Wrapf(err, "读取插件清单失败: %s", path)
	}
	var manifest Manifest
	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, gerror.Wrapf(err, "解析插件清单失败: %s", path)
	}
	return &manifest, nil
}

func readManifestZip(path string) (*Manifest, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, gerror.Wrapf(err, "打开插件包失败: %s", path)
	}
	defer reader.Close()
	for _, file := range reader.File {
		if filepath.Clean(file.Name) != "plugin.yaml" {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, gerror.Wrap(err, "读取插件包清单失败")
		}
		defer rc.Close()
		content, err := io.ReadAll(rc)
		if err != nil {
			return nil, gerror.Wrap(err, "读取插件包清单失败")
		}
		var manifest Manifest
		if err := yaml.Unmarshal(content, &manifest); err != nil {
			return nil, gerror.Wrap(err, "解析插件包清单失败")
		}
		return &manifest, nil
	}
	return nil, gerror.New("插件包缺少 plugin.yaml")
}

func readPluginAsset(pluginDir string, assetPath string) (string, error) {
	assetPath = filepath.Clean(strings.TrimSpace(assetPath))
	if assetPath == "." || assetPath == "" {
		return "", gerror.New("插件页面资源路径不能为空")
	}
	path := filepath.Clean(filepath.Join(pluginDir, assetPath))
	if !isPathInside(path, pluginDir) {
		return "", gerror.Newf("插件页面资源路径越界: %s", assetPath)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", gerror.Wrapf(err, "读取插件页面资源失败: %s", assetPath)
	}
	return string(content), nil
}

func validateConfigSchema(schema map[string]any, config map[string]any) error {
	if config == nil {
		return gerror.New("设备插件配置不能为空")
	}
	required, _ := schema["required"].([]any)
	for _, field := range required {
		name, ok := field.(string)
		if !ok || strings.TrimSpace(name) == "" {
			continue
		}
		value, exists := config[name]
		if !exists || isEmptyConfigValue(value) {
			return gerror.Newf("设备插件配置缺少必填字段: %s", name)
		}
	}
	properties, _ := schema["properties"].(map[string]any)
	for name, rawProperty := range properties {
		value, exists := config[name]
		if !exists || isEmptyConfigValue(value) {
			continue
		}
		property, _ := rawProperty.(map[string]any)
		if err := validateConfigProperty(name, value, property); err != nil {
			return err
		}
	}
	return nil
}

func validateConfigProperty(name string, value any, property map[string]any) error {
	rules, _ := property["rules"].(map[string]any)
	if len(rules) == 0 {
		rules = property
	}
	if maximum, exists := numberRule(rules, "maximum"); exists {
		if number, ok := numericValue(value); ok && number > maximum {
			return gerror.Newf("设备插件配置字段 %s 不能大于 %v", name, maximum)
		}
	}
	if minimum, exists := numberRule(rules, "minimum"); exists {
		if number, ok := numericValue(value); ok && number < minimum {
			return gerror.Newf("设备插件配置字段 %s 不能小于 %v", name, minimum)
		}
	}
	if maxLength, exists := numberRule(rules, "maxLength"); exists {
		if text, ok := value.(string); ok && float64(len([]rune(text))) > maxLength {
			return gerror.Newf("设备插件配置字段 %s 长度不能超过 %v", name, maxLength)
		}
	}
	enumValues, _ := property["enum"].([]any)
	if len(enumValues) > 0 {
		text := fmt.Sprint(value)
		for _, enumValue := range enumValues {
			if text == fmt.Sprint(enumValue) {
				return nil
			}
		}
		return gerror.Newf("设备插件配置字段 %s 不在允许范围内", name)
	}
	return nil
}

func numberRule(values map[string]any, key string) (float64, bool) {
	value, exists := values[key]
	if !exists {
		return 0, false
	}
	number, ok := numericValue(value)
	return number, ok
}

func numericValue(value any) (float64, bool) {
	switch item := value.(type) {
	case int:
		return float64(item), true
	case int64:
		return float64(item), true
	case float64:
		return item, true
	case float32:
		return float64(item), true
	case json.Number:
		number, err := item.Float64()
		return number, err == nil
	default:
		text := strings.TrimSpace(gconv.String(value))
		if text == "" {
			return 0, false
		}
		number, err := strconv.ParseFloat(text, 64)
		return number, err == nil
	}
}

func isEmptyConfigValue(value any) bool {
	switch item := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(item) == ""
	default:
		return false
	}
}

func validateManifest(manifest *Manifest) error {
	if manifest.Id == "" {
		return gerror.New("插件清单缺少 id")
	}
	if manifest.Name == "" {
		return gerror.New("插件清单缺少 name")
	}
	if manifest.Type == "" {
		return gerror.New("插件清单缺少 type")
	}
	if manifest.Version == "" {
		return gerror.New("插件清单缺少 version")
	}
	if manifest.Runtime != "process" {
		return gerror.Newf("插件 runtime 只支持 process: %s", manifest.Id)
	}
	if manifest.Protocol != "grpc" {
		return gerror.Newf("插件 protocol 只支持 grpc: %s", manifest.Id)
	}
	return nil
}

func toPluginItem(plugin *RuntimePlugin) commonv1.PluginItem {
	return commonv1.PluginItem{
		Id:          plugin.Manifest.Id,
		Name:        plugin.Manifest.Name,
		Type:        plugin.Manifest.Type,
		Version:     plugin.Manifest.Version,
		Runtime:     plugin.Manifest.Runtime,
		Protocol:    plugin.Manifest.Protocol,
		Status:      plugin.Status,
		Permissions: append([]string(nil), plugin.Manifest.Permissions...),
		UpdatedAt:   plugin.UpdatedAt,
	}
}

type hostSettings struct {
	RootDir     string
	Mode        string
	AutoCompile bool
}

func (s hostSettings) withoutCompile() hostSettings {
	s.AutoCompile = false
	return s
}

func loadSettings(ctx context.Context) (hostSettings, error) {
	root, err := projectRoot()
	if err != nil {
		return hostSettings{}, err
	}
	if err := ensureConfigPath(root); err != nil {
		return hostSettings{}, err
	}
	mode, err := configString(ctx, configPluginMode)
	if err != nil {
		return hostSettings{}, err
	}
	if mode == "" {
		mode = pluginModeProduction
	}
	if mode != pluginModeDevelopment && mode != pluginModeProduction {
		return hostSettings{}, gerror.Newf("不支持的插件宿主模式: %s", mode)
	}
	autoCompile, err := configBool(ctx, configPluginAutoCompile, true)
	if err != nil {
		return hostSettings{}, err
	}
	return hostSettings{
		RootDir:     filepath.Join(root, pluginRootDir),
		Mode:        mode,
		AutoCompile: mode == pluginModeDevelopment && autoCompile,
	}, nil
}

func configString(ctx context.Context, key string) (string, error) {
	value, err := g.Cfg().GetEffective(ctx, key)
	if err != nil {
		return "", gerror.Wrapf(err, "读取配置失败: %s", key)
	}
	if value == nil {
		return "", nil
	}
	return strings.TrimSpace(value.String()), nil
}

func configBool(ctx context.Context, key string, defaultValue bool) (bool, error) {
	value, err := g.Cfg().GetEffective(ctx, key, defaultValue)
	if err != nil {
		return false, gerror.Wrapf(err, "读取配置失败: %s", key)
	}
	if value == nil {
		return defaultValue, nil
	}
	raw := strings.TrimSpace(value.String())
	if raw == "" {
		return defaultValue, nil
	}
	return gconv.Bool(raw), nil
}

func ensureConfigPath(root string) error {
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		return nil
	}
	configPath := filepath.Join(root, "manifest", "config")
	if err := adapter.SetPath(configPath); err != nil {
		return gerror.Wrapf(err, "设置配置目录失败: %s", configPath)
	}
	return nil
}

func projectRoot() (string, error) {
	dir := gfile.Pwd()
	for {
		if gfile.IsFile(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", gerror.Wrap(err, "读取当前工作目录失败")
	}
	return "", gerror.Newf("无法定位项目根目录: %s", wd)
}

func (s *Service) waitTaskCancel(ctx context.Context, taskID string) {
	<-ctx.Done()
	s.mu.Lock()
	delete(s.taskCancels, taskID)
	s.mu.Unlock()
}

func (s *Service) deviceAndConfig(ctx context.Context, deviceID string) (*entity.Devices, *entity.PluginDeviceConfigs, error) {
	var device entity.Devices
	if err := dao.Devices.Ctx(ctx).Where(do.Devices{Id: deviceID}).Scan(&device); err != nil {
		return nil, nil, gerror.Wrapf(err, "读取设备失败: %s", deviceID)
	}
	if device.Id == "" {
		return nil, nil, gerror.Newf("设备不存在: %s", deviceID)
	}
	var config entity.PluginDeviceConfigs
	if err := dao.PluginDeviceConfigs.Ctx(ctx).Where(do.PluginDeviceConfigs{
		DeviceId: deviceID,
		PluginId: device.PluginId,
		Active:   1,
	}).Scan(&config); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, gerror.Wrapf(err, "读取设备插件配置失败: %s", deviceID)
		}
	}
	if config.Id == "" {
		return &device, nil, nil
	}
	return &device, &config, nil
}

func (s *Service) task(ctx context.Context, taskID string) (*entity.CollectionTasks, error) {
	var task entity.CollectionTasks
	if err := dao.CollectionTasks.Ctx(ctx).Where(do.CollectionTasks{Id: taskID}).Scan(&task); err != nil {
		return nil, gerror.Wrapf(err, "读取采集任务失败: %s", taskID)
	}
	if task.Id == "" {
		return nil, gerror.Newf("采集任务不存在: %s", taskID)
	}
	return &task, nil
}

func (s *Service) deviceRuntimeConfig(ctx context.Context, deviceID string) (*entity.Devices, *entity.PluginDeviceConfigs, error) {
	device, config, err := s.deviceAndConfig(ctx, deviceID)
	if err != nil {
		return nil, nil, err
	}
	if config == nil || strings.TrimSpace(config.ConfigJson) == "" {
		return nil, nil, gerror.Newf("设备缺少插件运行配置: %s", deviceID)
	}
	return device, config, nil
}

func (s *Service) ensureDeviceRunnable(ctx context.Context, deviceID string) error {
	device, config, err := s.deviceRuntimeConfig(ctx, deviceID)
	if err != nil {
		return err
	}
	if device.Enabled != 1 {
		return gerror.Newf("设备未启用: %s", deviceID)
	}
	if config.Enabled != 1 {
		return gerror.Newf("设备插件配置未启用: %s", deviceID)
	}
	pointCount, err := dao.DevicePoints.Ctx(ctx).Where(do.DevicePoints{
		DeviceId: deviceID,
		PluginId: device.PluginId,
		Enabled:  1,
	}).Count()
	if err != nil {
		return gerror.Wrapf(err, "读取设备启用点位失败: %s", deviceID)
	}
	if pointCount == 0 {
		return gerror.Newf("设备没有启用点位: %s", deviceID)
	}
	return nil
}

func (s *Service) markTaskRunning(ctx context.Context, task *entity.CollectionTasks) error {
	now := gtime.Now().Format("Y-m-d H:i:s.u")
	if _, err := dao.CollectionTasks.Ctx(ctx).
		Where(do.CollectionTasks{Id: task.Id}).
		Data(do.CollectionTasks{Status: "running", Rate: "等待插件上报", LastCollectedAt: now}).
		Update(); err != nil {
		return gerror.Wrapf(err, "启动采集任务失败: %s", task.Id)
	}
	if _, err := dao.Devices.Ctx(ctx).
		Where(do.Devices{Id: task.DeviceId}).
		Data(do.Devices{Status: "online", LastSeenAt: now}).
		Update(); err != nil {
		return gerror.Wrapf(err, "更新设备在线状态失败: %s", task.DeviceId)
	}
	return s.appendRuntimeEvent(ctx, "INFO", "collector", task.SouthPluginId, task.DeviceId, task.Id, "采集任务已启动，等待南向插件通过 gRPC 提交记录。", "")
}

func (s *Service) appendRuntimeEvent(ctx context.Context, level string, source string, pluginID string, deviceID string, taskID string, message string, traceID string) error {
	if traceID == "" {
		traceID = "trace-" + guid.S()
	}
	_, err := dao.RuntimeEvents.Ctx(ctx).Data(do.RuntimeEvents{
		Id:       "evt-" + guid.S(),
		Time:     gtime.Now().Format("Y-m-d H:i:s.u"),
		Level:    level,
		Source:   source,
		PluginId: pluginID,
		DeviceId: deviceID,
		TaskId:   taskID,
		Message:  message,
		TraceId:  traceID,
	}).Insert()
	return gerror.Wrap(err, "写入运行事件失败")
}

func pointNames(ctx context.Context, records []CollectionRecord) (map[string]string, error) {
	ids := make([]string, 0, len(records))
	seen := map[string]bool{}
	for _, record := range records {
		if !seen[record.PointID] {
			ids = append(ids, record.PointID)
			seen[record.PointID] = true
		}
	}
	var points []entity.DevicePoints
	if err := dao.DevicePoints.Ctx(ctx).WhereIn(dao.DevicePoints.Columns().Id, ids).Scan(&points); err != nil {
		return nil, gerror.Wrap(err, "读取采集点位名称失败")
	}
	names := map[string]string{}
	for _, point := range points {
		names[point.Id] = point.Name
	}
	return names, nil
}

package pluginmgmt

import (
	"archive/zip"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	runtimev1 "github.com/mijjjj/gcoll/api/runtime/v1"
	"github.com/mijjjj/gcoll/internal/dao"
	"github.com/mijjjj/gcoll/internal/model/do"
	"github.com/mijjjj/gcoll/internal/model/entity"
	"gopkg.in/yaml.v3"
)

// Service 提供插件管理服务。
type Service struct{}

// Manifest 描述插件清单中当前需要持久化的字段。
type Manifest struct {
	SchemaVersion int            `yaml:"schemaVersion" json:"schemaVersion"`
	Id            string         `yaml:"id"            json:"id"`
	Name          string         `yaml:"name"          json:"name"`
	Type          string         `yaml:"type"          json:"type"`
	Version       string         `yaml:"version"       json:"version"`
	Runtime       string         `yaml:"runtime"       json:"runtime"`
	Protocol      string         `yaml:"protocol"      json:"protocol"`
	Capabilities  []string       `yaml:"capabilities"  json:"capabilities"`
	Permissions   []string       `yaml:"permissions"   json:"permissions"`
	ConfigSchema  map[string]any `yaml:"configSchema"  json:"configSchema"`
}

// New 创建插件管理服务。
func New() *Service {
	return &Service{}
}

// List 返回数据库中的插件列表。
func (s *Service) List(ctx context.Context) ([]runtimev1.PluginItem, error) {
	var plugins []entity.Plugins
	if err := dao.Plugins.Ctx(ctx).OrderAsc(dao.Plugins.Columns().Id).Scan(&plugins); err != nil {
		return nil, gerror.Wrap(err, "读取插件列表失败")
	}

	items := make([]runtimev1.PluginItem, 0, len(plugins))
	for _, plugin := range plugins {
		version, err := s.activeVersion(ctx, plugin.Id, plugin.ActiveVersion)
		if err != nil {
			return nil, err
		}
		items = append(items, runtimev1.PluginItem{
			Id:          plugin.Id,
			Name:        plugin.Name,
			Type:        plugin.Type,
			Version:     plugin.ActiveVersion,
			Runtime:     plugin.Runtime,
			Protocol:    plugin.Protocol,
			Status:      plugin.Status,
			Permissions: stringSliceFromJSON(version.PermissionsJson),
			UpdatedAt:   plugin.UpdatedAt,
		})
	}
	return items, nil
}

// Import 从插件目录或 .gcpkg 包读取清单并写入插件表。
func (s *Service) Import(ctx context.Context, packagePath string) (*runtimev1.PluginItem, error) {
	if strings.TrimSpace(packagePath) == "" {
		return nil, gerror.New("插件包路径不能为空")
	}
	manifest, err := readManifest(packagePath)
	if err != nil {
		return nil, err
	}
	if err := validateManifest(manifest); err != nil {
		return nil, err
	}

	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return nil, gerror.Wrap(err, "序列化插件清单失败")
	}
	permissionsJSON, err := json.Marshal(manifest.Permissions)
	if err != nil {
		return nil, gerror.Wrap(err, "序列化插件权限失败")
	}
	capabilitiesJSON, err := json.Marshal(manifest.Capabilities)
	if err != nil {
		return nil, gerror.Wrap(err, "序列化插件能力失败")
	}
	schemaJSON, err := json.Marshal(manifest.ConfigSchema)
	if err != nil {
		return nil, gerror.Wrap(err, "序列化插件配置结构失败")
	}

	versionId := "pv-" + strings.ReplaceAll(manifest.Id, ".", "-") + "-" + strings.ReplaceAll(manifest.Version, ".", "-")
	schemaId := "pcs-" + strings.ReplaceAll(manifest.Id, ".", "-") + "-" + strings.ReplaceAll(manifest.Version, ".", "-")
	if err := dao.Plugins.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_ = ctx
		if _, err := tx.Save(dao.Plugins.Table(), do.Plugins{
			Id:            manifest.Id,
			Name:          manifest.Name,
			Type:          manifest.Type,
			Runtime:       manifest.Runtime,
			Protocol:      manifest.Protocol,
			Status:        "installed",
			ActiveVersion: manifest.Version,
			Source:        "imported",
			Enabled:       1,
		}); err != nil {
			return gerror.Wrap(err, "保存插件信息失败")
		}
		if _, err := tx.Save(dao.PluginVersions.Table(), do.PluginVersions{
			Id:               versionId,
			PluginId:         manifest.Id,
			Version:          manifest.Version,
			PackagePath:      packagePath,
			ManifestJson:     string(manifestJSON),
			PermissionsJson:  string(permissionsJSON),
			CapabilitiesJson: string(capabilitiesJSON),
			Active:           1,
		}); err != nil {
			return gerror.Wrap(err, "保存插件版本失败")
		}
		if _, err := tx.Save(dao.PluginConfigSchemas.Table(), do.PluginConfigSchemas{
			Id:              schemaId,
			PluginId:        manifest.Id,
			PluginVersionId: versionId,
			SchemaVersion:   manifest.SchemaVersion,
			SchemaJson:      string(schemaJSON),
		}); err != nil {
			return gerror.Wrap(err, "保存插件配置结构失败")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &runtimev1.PluginItem{
		Id:          manifest.Id,
		Name:        manifest.Name,
		Type:        manifest.Type,
		Version:     manifest.Version,
		Runtime:     manifest.Runtime,
		Protocol:    manifest.Protocol,
		Status:      "installed",
		Permissions: manifest.Permissions,
	}, nil
}

func (s *Service) activeVersion(ctx context.Context, pluginId string, version string) (*entity.PluginVersions, error) {
	var item entity.PluginVersions
	err := dao.PluginVersions.Ctx(ctx).
		Where(do.PluginVersions{PluginId: pluginId, Version: version}).
		Scan(&item)
	if err != nil {
		return nil, gerror.Wrapf(err, "读取插件版本失败: %s", pluginId)
	}
	if item.Id == "" {
		return nil, gerror.Newf("插件缺少活跃版本: %s", pluginId)
	}
	return &item, nil
}

func readManifest(packagePath string) (*Manifest, error) {
	if stat, err := os.Stat(packagePath); err == nil && stat.IsDir() {
		return readManifestFile(filepath.Join(packagePath, "plugin.yaml"))
	}
	if strings.EqualFold(filepath.Ext(packagePath), ".gcpkg") {
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
		if file.Name != "plugin.yaml" {
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
	if manifest.Runtime == "" {
		return gerror.New("插件清单缺少 runtime")
	}
	if manifest.Protocol == "" {
		return gerror.New("插件清单缺少 protocol")
	}
	return nil
}

func stringSliceFromJSON(raw string) []string {
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return []string{}
	}
	return values
}

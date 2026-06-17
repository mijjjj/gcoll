package pluginmgmt

import (
	"context"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
	pluginhostsvc "github.com/mijjjj/gcoll/internal/service/pluginhost"
)

// Service 提供插件管理服务，插件事实来源来自插件宿主内存注册表。
type Service struct {
	host *pluginhostsvc.Service
}

// New 创建插件管理服务。
func New() *Service {
	return &Service{host: pluginhostsvc.Instance()}
}

// List 返回内存中的插件列表。
func (s *Service) List(ctx context.Context) ([]commonv1.PluginItem, error) {
	return s.host.List(ctx)
}

// Import 将插件包解压或复制到插件目录，加载到内存并启动插件进程。
func (s *Service) Import(ctx context.Context, packagePath string) (*commonv1.PluginItem, error) {
	return s.host.ImportPackage(ctx, packagePath)
}

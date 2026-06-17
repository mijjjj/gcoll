package boot

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	i18nsvc "github.com/mijjjj/gcoll/internal/service/i18n"
	pluginhostsvc "github.com/mijjjj/gcoll/internal/service/pluginhost"
	storagesvc "github.com/mijjjj/gcoll/internal/service/storage"
)

// Init 初始化配置、日志和后续基础设施挂载点。
func Init(ctx context.Context) {
	i18nsvc.Init(ctx)
	if err := storagesvc.Init(ctx); err != nil {
		g.Log().Fatalf(ctx, "数据库初始化失败: %+v", err)
	}
	if err := pluginhostsvc.Init(ctx); err != nil {
		g.Log().Errorf(ctx, "插件宿主初始化失败: %+v", err)
	}
	g.Log().Info(ctx, "gcoll 运行时初始化完成")
}

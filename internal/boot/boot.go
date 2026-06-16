package boot

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	i18nsvc "github.com/mijjjj/gcoll/internal/service/i18n"
)

// Init 初始化配置、日志和后续基础设施挂载点。
func Init(ctx context.Context) {
	i18nsvc.Init(ctx)
	g.Log().Info(ctx, "gcoll 运行时初始化完成")
}

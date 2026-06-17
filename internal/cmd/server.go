package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"

	"github.com/mijjjj/gcoll/internal/boot"
	devicectrl "github.com/mijjjj/gcoll/internal/controller/device"
	logctrl "github.com/mijjjj/gcoll/internal/controller/log"
	pipelinectrl "github.com/mijjjj/gcoll/internal/controller/pipeline"
	pluginctrl "github.com/mijjjj/gcoll/internal/controller/plugin"
	pointcachectrl "github.com/mijjjj/gcoll/internal/controller/pointcache"
	runtimectrl "github.com/mijjjj/gcoll/internal/controller/runtime"
	targetctrl "github.com/mijjjj/gcoll/internal/controller/target"
	taskctrl "github.com/mijjjj/gcoll/internal/controller/task"
	middlewaresvc "github.com/mijjjj/gcoll/internal/service/middleware"
)

const (
	openAPITitle       = "gcoll API"
	openAPIDescription = "gcoll 管理 API，使用 GoFrame 标准路由和统一响应中间件。"
)

// RunServer 启动 gcoll 服务端 HTTP 运行时。
func RunServer(ctx context.Context) {
	boot.Init(ctx)

	s := g.Server("gcoll")
	enhanceOpenAPIDoc(s)
	registerRoutes(s)
	s.Run()
}

// registerRoutes 注册服务端管理 API 路由。
func registerRoutes(s *ghttp.Server) {
	var (
		deviceController     = devicectrl.NewV1()
		logController        = logctrl.NewV1()
		pipelineController   = pipelinectrl.NewV1()
		pluginController     = pluginctrl.NewV1()
		pointCacheController = pointcachectrl.NewV1()
		runtimeController    = runtimectrl.NewV1()
		targetController     = targetctrl.NewV1()
		taskController       = taskctrl.NewV1()
		middlewareService    = middlewaresvc.New()
	)

	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareService.Response)
		group.Bind(
			runtimeController,
			pluginController,
			deviceController,
			taskController,
			pointCacheController,
			pipelineController,
			targetController,
			logController,
		)
	})
}

// enhanceOpenAPIDoc 配置统一返回结构文档。
func enhanceOpenAPIDoc(s *ghttp.Server) {
	openAPI := s.GetOpenApi()
	openAPI.Config.CommonResponse = middlewaresvc.HandlerResponse{}
	openAPI.Config.CommonResponseDataField = "data"
	openAPI.Info = goai.Info{
		Title:       openAPITitle,
		Description: openAPIDescription,
	}
}

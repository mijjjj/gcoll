package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"

	"github.com/mijjjj/gcoll/internal/boot"
	runtimectrl "github.com/mijjjj/gcoll/internal/controller/runtime"
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
		runtimeController = runtimectrl.NewV1()
		middlewareService = middlewaresvc.New()
	)

	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareService.Response)
		group.Bind(runtimeController)
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

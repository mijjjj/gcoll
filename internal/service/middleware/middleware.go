package middleware

import (
	"context"

	i18nsvc "github.com/mijjjj/gcoll/internal/service/i18n"
)

// Service 提供 HTTP 中间件能力。
type Service struct{}

// New 创建中间件服务。
func New() *Service {
	i18nsvc.Init(context.Background())
	return &Service{}
}

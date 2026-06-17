package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// LogsReq 描述运行日志列表请求。
type LogsReq struct {
	g.Meta `path:"/logs" method:"get" tags:"Logs" summary:"获取运行日志列表"`
}

// LogsRes 描述运行日志列表响应。
type LogsRes struct {
	Items []commonv1.RuntimeEvent `json:"items"`
}

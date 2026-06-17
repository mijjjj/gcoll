package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// TargetsReq 描述北向转发目标列表请求。
type TargetsReq struct {
	g.Meta `path:"/targets" method:"get" tags:"Targets" summary:"获取北向转发目标列表"`
}

// TargetsRes 描述北向转发目标列表响应。
type TargetsRes struct {
	Items []commonv1.ForwardTargetItem `json:"items"`
}

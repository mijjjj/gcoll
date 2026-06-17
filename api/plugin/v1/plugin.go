package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// PluginsReq 描述插件列表请求。
type PluginsReq struct {
	g.Meta `path:"/plugins" method:"get" tags:"Plugins" summary:"获取插件列表"`
}

// PluginsRes 描述插件列表响应。
type PluginsRes struct {
	Items []commonv1.PluginItem `json:"items"`
}

// ImportPluginReq 描述插件导入请求。
type ImportPluginReq struct {
	g.Meta      `path:"/plugins/import" method:"post" tags:"Plugins" summary:"导入插件清单"`
	PackagePath string `json:"packagePath" v:"required#插件包路径不能为空"`
}

// ImportPluginRes 描述插件导入响应。
type ImportPluginRes struct {
	Plugin commonv1.PluginItem `json:"plugin"`
}

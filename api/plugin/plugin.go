// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package plugin

import (
	"context"

	"github.com/mijjjj/gcoll/api/plugin/v1"
)

type IPluginV1 interface {
	Plugins(ctx context.Context, req *v1.PluginsReq) (res *v1.PluginsRes, err error)
	ImportPlugin(ctx context.Context, req *v1.ImportPluginReq) (res *v1.ImportPluginRes, err error)
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package log

import (
	"context"

	"github.com/mijjjj/gcoll/api/log/v1"
)

type ILogV1 interface {
	Logs(ctx context.Context, req *v1.LogsReq) (res *v1.LogsRes, err error)
}

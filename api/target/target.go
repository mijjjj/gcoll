// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package target

import (
	"context"

	"github.com/mijjjj/gcoll/api/target/v1"
)

type ITargetV1 interface {
	Targets(ctx context.Context, req *v1.TargetsReq) (res *v1.TargetsRes, err error)
}

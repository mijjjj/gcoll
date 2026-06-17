// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

type IDeviceV1 interface {
	DevicePluginConfigPage(ctx context.Context, req *v1.DevicePluginConfigPageReq) (res *v1.DevicePluginConfigPageRes, err error)
	UpdateDevicePluginConfig(ctx context.Context, req *v1.UpdateDevicePluginConfigReq) (res *v1.UpdateDevicePluginConfigRes, err error)
	TestDevicePluginConnection(ctx context.Context, req *v1.TestDevicePluginConnectionReq) (res *v1.TestDevicePluginConnectionRes, err error)
	Devices(ctx context.Context, req *v1.DevicesReq) (res *v1.DevicesRes, err error)
	CreateDeviceGroup(ctx context.Context, req *v1.CreateDeviceGroupReq) (res *v1.CreateDeviceGroupRes, err error)
	DeleteDeviceGroup(ctx context.Context, req *v1.DeleteDeviceGroupReq) (res *v1.DeleteDeviceGroupRes, err error)
	CreateDevice(ctx context.Context, req *v1.CreateDeviceReq) (res *v1.CreateDeviceRes, err error)
	MoveDeviceToGroup(ctx context.Context, req *v1.MoveDeviceToGroupReq) (res *v1.MoveDeviceToGroupRes, err error)
	DeleteDevice(ctx context.Context, req *v1.DeleteDeviceReq) (res *v1.DeleteDeviceRes, err error)
	DevicePoints(ctx context.Context, req *v1.DevicePointsReq) (res *v1.DevicePointsRes, err error)
	CreateDevicePoint(ctx context.Context, req *v1.CreateDevicePointReq) (res *v1.CreateDevicePointRes, err error)
	UpdateDevicePoints(ctx context.Context, req *v1.UpdateDevicePointsReq) (res *v1.UpdateDevicePointsRes, err error)
}

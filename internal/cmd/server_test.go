package cmd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	storagesvc "github.com/mijjjj/gcoll/internal/service/storage"
)

func TestRegisterRoutes(t *testing.T) {
	prepareTestDatabase(t)

	s := g.Server(guid.S())
	registerRoutes(s)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		content := client.GetContent(context.Background(), "/api/v1/runtime/health")
		result := gjson.New(content)

		t.Assert(result.Get("code").Int(), 0)
		t.Assert(result.Get("message").String(), "成功")
		t.Assert(result.Get("data.status").String(), "ok")
		t.Assert(result.Get("data.service").String(), "gcoll-server")
		t.Assert(result.Get("data.mode").String(), "server")
		t.AssertNE(result.Get("data.checkedAt").String(), "")

		content = client.GetContent(context.Background(), "/api/v1/runtime/health?lang=en")
		result = gjson.New(content)

		t.Assert(result.Get("code").Int(), 0)
		t.Assert(result.Get("message").String(), "OK")

		content = client.GetContent(context.Background(), "/api/v1/runtime/overview")
		result = gjson.New(content)

		t.Assert(result.Get("code").Int(), 0)
		t.Assert(result.Get("data.metrics.0.label").String(), "运行时")
		t.Assert(result.Get("data.tasks.0.name").String(), "样例 Modbus TCP 采集链路")
		t.Assert(result.Get("data.pluginSummary.total").Int(), 1)

		content = client.GetContent(context.Background(), "/api/v1/devices")
		result = gjson.New(content)

		t.Assert(result.Get("code").Int(), 0)
		t.Assert(result.Get("data.items.0.id").String(), "dev-edge-gw-a01")

		content = client.GetContent(context.Background(), "/api/v1/devices/dev-edge-gw-a01/points")
		result = gjson.New(content)

		t.Assert(result.Get("code").Int(), 0)
		t.Assert(result.Get("data.items.0.name").String(), "TEMP_01")

		content = client.GetContent(context.Background(), "/api/v1/devices/dev-edge-gw-a01/protocol-config")
		result = gjson.New(content)

		t.Assert(result.Get("code").Int(), 0)
		t.Assert(result.Get("data.plugin.id").String(), "com.gcoll.modbus-tcp")
		t.Assert(result.Get("data.device.id").String(), "dev-edge-gw-a01")
		t.Assert(result.Get("data.readPlan.0.area").String(), "coil")
	})
}

func prepareTestDatabase(t *testing.T) {
	t.Helper()

	err := gdb.SetConfigGroup("default", gdb.ConfigGroup{{
		Type:             "sqlite",
		Link:             "sqlite::@file(:memory:)",
		MaxOpenConnCount: 1,
	}})
	if err != nil {
		t.Fatalf("设置测试数据库失败: %v", err)
	}
	if err := storagesvc.Init(context.Background()); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
}

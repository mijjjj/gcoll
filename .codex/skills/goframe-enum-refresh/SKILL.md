---
name: goframe-enum-refresh
description: gcoll GoFrame 枚举刷新流程。用于新增、删除、重命名或调整 `internal/consts` 中的枚举常量，修改 API `v:"enums"` 校验，或排查 GoFrame 枚举校验拒绝新值的问题。
---

# GoFrame Enum Refresh

当后端变更触及枚举常量或 `v:"enums"` 校验时使用。

## 必需流程

1. 在 `internal/consts` 中按既有 const/type 模式修改枚举。
2. 在仓库根目录运行：

```powershell
gf gen enums
```

3. 确认 `internal/boot/boot_enums.go` 已包含新增或调整的枚举值。
4. 运行：

```powershell
go test ./...
```

## 注意事项

- GoFrame `v:"enums"` 校验依赖生成的枚举注册。
- 只修改 `internal/consts` 不足以让新枚举值通过校验。
- 如果 API 拒绝新增状态、类型或命令值，优先检查是否遗漏 `gf gen enums`。
- 适用时在最终说明中提到已运行或未运行 `gf gen enums`。

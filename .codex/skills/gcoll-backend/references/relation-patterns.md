# GoFrame 静态关联查询模式

在列表或详情 API 需要关联表数据且不希望产生 N+1 查询时使用。

## 模型

- 业务关联模型放在 `internal/model`，不放在 Controller 或 Service 局部结构中。
- 使用 `g.Meta` 标明主表。
- 显式列出主表字段，不嵌入生成的 `entity`。
- 关联字段使用 model 层类型，并添加 `orm:"with:foreign_key=id"`。

## 查询

- 使用 `g.Model(model.XWithRelation{}).Ctx(ctx).WithAll()`。
- 先应用租户、权限、过滤、时间范围，再分页和排序。
- 扫描到关联模型后再转换为 API 响应结构。

## 响应

- API 响应使用嵌套响应结构表达关联数据。
- 不把其他模块的请求类型作为响应字段。
- 只返回当前接口需要的字段，避免泄露敏感配置。

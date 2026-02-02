---
description: 在关键位置使用 Writer 记录日志，替代传统的 fmt.Println/log.Printf
---

请将选中代码或当前文件中的关键日志改为使用 `Writer` 记录。

## 使用规则

1. **使用 Writer 的场景**：
   - 系统级错误（数据库、Redis、系统配置等）
   - 关键用户操作（登录、登出、权限变更等）
   - 定时任务执行日志

2. **不使用 Writer 的场景**：
   - 业务逻辑错误（用户不存在、密码错误等）使用普通日志 `l.Infof()`, `l.Errorf()`
   - 调试信息

## log_type 类型

| log_type | 使用场景 |
|----------|----------|
| `system` | 系统级错误（非数据库、非Redis） |
| `database` | 数据库相关操作和错误 |
| `redis` | Redis相关操作和错误 |
| `user` | 用户相关操作（需要审计追踪） |
| `permission` | 权限相关操作 |

## 代码示例

```go
// 系统级错误
l.svcCtx.Writer.Error("系统配置错误",
    writer.Field("log_type", "system"),
    writer.Field("trace", "FunctionName"),
    writer.Field("username", "system"),  // 或实际用户名
    writer.Field("error", err.Error()),
)

// 数据库错误
l.svcCtx.Writer.Error("数据库查询失败",
    writer.Field("log_type", "database"),
    writer.Field("trace", "FunctionName"),
    writer.Field("user_id", userID),
    writer.Field("error", err.Error()),
)

// 关键用户操作
l.svcCtx.Writer.Info("用户登录成功",
    writer.Field("log_type", "user"),
    writer.Field("trace", "Login"),
    writer.Field("user_id", user.ID),
    writer.Field("username", user.Name),
)

// 定时任务日志
s.writer.Info("定时任务执行完成",
    writer.Field("log_type", "system"),
    writer.Field("trace", "Scheduler.TaskHandler"),
    writer.Field("username", "system"),
    writer.Field("task_name", taskName),
)
```

## 特殊字段（会存储到独立数据库列）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `trace` | string | 链路追踪/方法名 |
| `log_type` | string | 日志类型 |
| `user_id` | int64 | 用户 ID |
| `user_name` | string | 用户名 |
| `duration` | string | 耗时 |

## 注意事项

- 避免在日志中打印敏感信息（密码、密钥等）
- `user_name` 对于系统任务设为 `"system"`
- `trace` 建议使用 `"模块.方法名"` 格式

# doc - 为 API 接口添加说明文档

为 go-zero API 文件中的接口添加标准化的说明文档。

## 功能说明

按照 go-zero 官方规范，为 API 接口添加 `@doc` 注解和字段注释，提高代码可读性和文档完整性。

## 文档规范

### 1. @doc 注解格式

每个 handler 前添加 `@doc` 注解：

```go
@doc (
    summary: "接口简短描述"
    description: "接口详细说明"
)
@handler HandlerName
post /path (Request) returns (Response)
```

### 2. 字段注释格式

为每个类型字段添加 inline 注释：

```go
type Task {
    ID     string `json:"id"` // 任务唯一标识
    Status string `json:"status"` // 任务状态: pending-待处理, processing-处理中, completed-完成
}
```

### 3. 文件级别注释

在文件开头添加 block comment：

```go
/**
 * 服务名称 API 定义
 * 服务功能描述
 */

syntax = "v1"
```

## 执行流程

1. **分析 API 文件**：读取选中的 API 文件或当前打开的 API 文件
2. **识别未文档化的部分**：
   - 缺少 `@doc` 注解的 handler
   - 缺少注释的类型字段
   - 缺少文件级别说明
3. **查询 go-zero 文档**：使用 Context7 MCP 查询 go-zero API 文档规范（如需要）
4. **添加文档**：
   - 为 handler 添加 `@doc` 注解（summary 和 description）
   - 为类型字段添加清晰的中文注释
   - 为文件添加顶部说明
5. **验证格式**：确保符合 go-zero 语法规范

## 文档编写要点

### Summary 编写规则
- 简短清晰（5-15 个字）
- 使用动词开头（获取、添加、更新、删除等）
- 准确描述接口功能

### Description 编写规则
- 详细说明接口用途
- 包含重要的业务逻辑说明
- 说明特殊参数或返回值
- 注明权限要求（如有）

### 字段注释规则
- 使用中文注释
- 对于枚举类型，列出所有可能的值
- 对于有格式要求的字段，说明格式（如时间格式、ID 格式）
- 对于可选字段，可以标注"可选"

## 示例

### 完整的 API 文档示例

```go
/**
 * 队列服务 API 定义
 * 提供任务队列管理、任务处理和业务统计功能
 */

syntax = "v1"

// 任务结构
type Task {
    ID        string `json:"id"` // 任务唯一标识
    TaskName  string `json:"task_name"` // 任务名称
    Status    string `json:"status"` // 任务状态: pending-待处理, processing-处理中, completed-完成
    CreatedAt string `json:"created_at"` // 创建时间（ISO8601格式）
}

type AddTaskRequest {
    TaskName string `json:"task_name"` // 任务名称
}

type AddTaskResponse {
    Task Task `json:"task"` // 创建的任务信息
}

@server (
    prefix: /api/v1/queue
    group: queue
    middleware: AuthMiddleware
)
service queue-api {
    @doc (
        summary: "添加任务"
        description: "创建新的任务并加入队列"
    )
    @handler AddTaskHandler
    post /tasks (AddTaskRequest) returns (AddTaskResponse)
}
```

## 注意事项

1. **保持一致性**：同一文件中的文档风格应保持一致
2. **避免冗余**：不要在注释中重复字段名本身的信息
3. **关注业务**：注释应从业务角度解释，而非技术实现
4. **及时更新**：修改接口时同步更新文档
5. **符合规范**：严格遵循 go-zero 的 API 语法规范

## 参考资源

- go-zero 官方文档：https://go-zero.dev/
- API 语法规范：https://go-zero.dev/cn/api-grammar.html
- 使用 Context7 MCP 的 `user-context7` 服务查询最新文档

---
description: 同步 .env、YAML 和 config.go 配置字段结构
customPrompt: |
  请同步 `.env` 文件、`etc/libangDataService-Api.yaml` 和 `internal/config/config.go` 三个配置文件的字段结构。
  
  ## 目标
  
  确保三个配置文件的字段结构保持一致，以便：
  1. 开发环境使用 `.env` 文件（包含实际值）
  2. 部署配置使用 `libangDataService-Api.yaml`（模板，不包含敏感值）
  3. Go 代码通过 `config.go` 读取配置
  
  ## 同步规则
  
  ### 1. 从 .env 提取配置项
  
  分析 `.env` 文件中的所有配置变量（格式：`KEY=value`），按功能分类：
  
  - **PostgreSQL** (`POSTGRES_*`)
  - **Redis** (`REDIS_*`)
  - **认证** (`AUTH_*`)
  - **服务配置** (`*_SERVICE_*`)
  - **API Key** (`*_API_KEY`)
  - **OSS** (`OSS_*`)
  
  ### 2. 映射到 YAML 和 Config 结构
  
  将 `.env` 的扁平结构映射到 YAML 的嵌套结构和 Go struct：
  
  ```env
  # .env 格式
  POSTGRES_HOST=192.168.31.128
  POSTGRES_PORT=5432
  POSTGRES_USER=libang_user
  ```
  
  ```yaml
  # YAML 格式
  Postgres:
    Host: localhost  # 保持默认值
    Port: 5432
    User: postgres   # 保持默认值
  ```
  
  ```go
  // config.go 格式
  type PostgresConfig struct {
    Host     string `json:",env=POSTGRES_HOST"`
    Port     int    `json:",env=POSTGRES_PORT"`
    User     string `json:",env=POSTGRES_USER"`
    Password string `json:",env=POSTGRES_PASSWORD"`
  }
  ```
  
  ### 3. 处理敏感信息
  
  **关键原则**：`libangDataService-Api.yaml` 中所有敏感值保持为空字符串或默认占位符。
  
  敏感字段包括：
  - 密码 (`Password`, `*_PASSWORD`)
  - 密钥 (`Secret`, `*_SECRET`, `ApiKey`, `*_API_KEY`)
  - 用户名（如果包含敏感信息）
  
  ```yaml
  # 正确：敏感值为空
  Postgres:
    Password: ""
  Auth:
    AccessSecret: ""
  Services:
    UserService:
      ApiKey: ""
  
  # 错误：不要写入实际值
  Postgres:
    Password: "libang_password"
  Auth:
    AccessSecret: "your-secret-key-change-in-production"
  ```
  
  ### 4. 处理 Host 和网络配置
  
  对于 Host/Addr 等网络配置：
  - 保持 `libangDataService-Api.yaml` 中的 `localhost` 或默认值
  - 不要将 `.env` 中的 IP 地址同步过去
  
  ```yaml
  # 正确
  Postgres:
    Host: localhost  # 不改为 192.168.31.128
  Redis:
    Addr: localhost:6379  # 不改为 192.168.31.128:6379
  Services:
    UserService:
      Host: localhost  # 不改为 IP 地址
  ```
  
  ### 5. 字段命名规范
  
  **必须统一使用 `ApiKey`**，不要使用 `ServiceApiKey`：
  
  ```yaml
  # 正确
  Services:
    UserService:
      ApiKey: ""
    QueueService:
      ApiKey: ""
  
  # 错误
  Services:
    UserService:
      ServiceApiKey: ""  # 已废弃
  ```
  
  ```go
  // 正确
  type BaseServiceConfig struct {
    ApiKey string `json:",env=USER_SERVICE_API_KEY"`
  }
  
  // 错误
  type BaseServiceConfig struct {
    ServiceApiKey string `json:",env=USER_SERVICE_API_KEY"`  // 已废弃
  }
  ```
  
  ## 同步流程
  
  1. **读取 .env 文件**
     - 提取所有配置变量
     - 按功能分类（PostgreSQL, Redis, Auth, Services, OSS）
  
  2. **读取 libangDataService-Api.yaml**
     - 分析现有结构
     - 识别缺失的配置节
  
  3. **读取 internal/config/config.go**
     - 检查 struct 定义
     - 验证 env 标签映射
  
  4. **对比差异**
     - 找出 `.env` 中有但 YAML/config.go 中缺失的配置项
     - 找出字段名称不匹配的情况（如 `ServiceApiKey` vs `ApiKey`）
     - 找出端口号、路径等配置不一致的情况
  
  5. **执行同步**
     
     **YAML 文件：**
     - 添加缺失的配置节和字段
     - 保持所有敏感值为空字符串 `""`
     - 保持网络配置为 `localhost`
     - 保留现有注释和结构
     
     **config.go 文件：**
     - 添加缺失的 struct 字段
     - 确保 `json` 标签正确（包含 `env=` 映射）
     - 可选字段添加 `optional` 标签
     - 有默认值的字段添加 `default=` 标签
     - 统一使用 `ApiKey` 而非 `ServiceApiKey`
     
     **代码引用更新：**
     - 搜索所有使用旧字段名的代码（如 `ServiceApiKey`）
     - 更新为新字段名（如 `ApiKey`）
     - 检查 `internal/request/`, `internal/svc/` 等目录
  
  6. **输出同步报告**
     - 列出新增的配置节
     - 列出新增的字段
     - 列出字段重命名操作
     - 列出更新的代码文件
     - 提示需要手动配置的环境变量
  
  ## 示例场景
  
  ### 场景 1: 新增服务配置
  
  `.env` 中新增：
  ```env
  STORAGE_SERVICE_HOST=192.168.31.128
  STORAGE_SERVICE_PORT=8004
  STORAGE_SERVICE_PATH=/api/v1/storage
  STORAGE_SERVICE_API_KEY=secret123
  ```
  
  同步到 `libangDataService-Api.yaml`：
  ```yaml
  Services:
    StorageService:
      Host: localhost
      Port: 8004
      Path: /api/v1/storage
      ApiKey: ""
  ```
  
  同步到 `config.go`：
  ```go
  type StorageServiceConfig struct {
    Host   string `json:",env=STORAGE_SERVICE_HOST"`
    Port   int    `json:",env=STORAGE_SERVICE_PORT"`
    Path   string `json:",env=STORAGE_SERVICE_PATH"`
    ApiKey string `json:",env=STORAGE_SERVICE_API_KEY"`
  }
  
  type ServicesConfig struct {
    UserService    BaseServiceConfig
    QueueService   QueueServiceConfig
    StorageService StorageServiceConfig  // 新增
  }
  ```
  
  ### 场景 2: 字段重命名（ServiceApiKey → ApiKey）
  
  发现不一致：
  ```go
  // config.go 中使用旧名称
  type QueueServiceConfig struct {
    ServiceApiKey string `json:",env=QUEUE_SERVICE_API_KEY"`
  }
  ```
  
  执行同步：
  1. 更新 `config.go` 字段名为 `ApiKey`
  2. 更新 YAML 字段名为 `ApiKey`
  3. 搜索并更新所有代码引用：
     - `internal/request/task.go`: `r.services.QueueService.ServiceApiKey` → `r.services.QueueService.ApiKey`
     - `internal/request/user.go`: `r.services.UserService.ServiceApiKey` → `r.services.UserService.ApiKey`
  
  ### 场景 3: 端口号不一致
  
  发现差异：
  - `.env`: `QUEUE_SERVICE_PORT=8003`
  - YAML: `Port: 8001`
  
  同步操作：
  - 更新 YAML 为 `Port: 8003`（以 `.env` 为准）
  - 添加 TODO 注释说明后期变更计划
  
  ## 检查清单
  
  同步完成后，确保：
  
  - [ ] 所有 `.env` 中的配置项在 YAML 和 config.go 中都有对应字段
  - [ ] YAML 中所有敏感字段值为空字符串 `""`
  - [ ] Host/Addr 配置保持为 `localhost`
  - [ ] 配置节命名符合 Go-zero 规范（首字母大写）
  - [ ] config.go 中所有字段都有正确的 `json` 标签和 `env=` 映射
  - [ ] 可选字段标记了 `optional` 标签
  - [ ] 有默认值的字段标记了 `default=` 标签
  - [ ] 统一使用 `ApiKey` 而非 `ServiceApiKey`
  - [ ] 所有代码引用已更新（搜索 `ServiceApiKey` 应无结果）
  - [ ] 保留了 YAML 中的注释和结构
  - [ ] 无 Linter 错误
  - [ ] 输出了同步报告
  
  ## 注意事项
  
  1. **不要删除** YAML 或 config.go 中已有但 `.env` 中没有的配置项（可能是默认配置）
  2. **不要修改** YAML 中的顶层配置（Name, Host, Port, Timeout, MaxBytes 等）
  3. **保持向后兼容**：确保现有配置不被破坏
  4. **统一字段命名**：
     - 服务 API Key 统一使用 `ApiKey`
     - 数据库、Redis 配置保持现有命名
  5. **更新代码引用**：修改字段名后，必须更新所有使用该字段的代码
  6. **验证修改**：使用 `ReadLints` 检查修改后的文件是否有错误
  7. 同步完成后，提示用户在实际部署时需要配置环境变量
  
  ## 相关文件
  
  同步操作涉及以下文件：
  
  - `.env` - 开发环境配置（包含实际值）
  - `etc/libangDataService-Api.yaml` - 部署配置模板（敏感值为空）
  - `internal/config/config.go` - Go 配置结构定义
  - `internal/svc/serviceContext.go` - 服务上下文初始化
  - `internal/request/*.go` - 可能包含配置字段引用
  
  ## 相关规则
  
  - `.cursor/rules/style.mdc` - 禁止 emoji，使用清晰的文本标记
  - `.cursor/rules/basics.mdc` - 项目基础规范
---


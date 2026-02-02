# Go-Zero 项目模板

基于 Go-Zero 框架的项目模板，包含完整的项目结构和配置。

## 快速开始

### 1. 初始化项目

使用 `init-project.sh` 脚本自动替换项目名称：

```bash
# 方式1: 使用当前文件夹名称作为项目名
./init-project.sh

# 方式2: 指定项目名称
./init-project.sh my-new-service
```

脚本会自动完成以下操作：

**内容替换：**
- **Go 模块名**: `go-zero-template` → `my-new-service`
- **API 服务名**: `go_zero_template-api` → `my_new_service-api`

**文件重命名：**
- `go_zero_template.api` → `my_new_service.api`
- `goZeroTemplate-Api.go` → `myNewService-Api.go`
- `etc/goZeroTemplate-Api.yaml` → `etc/myNewService-Api.yaml`

**引用更新：**
- `makefile` 中的文件名引用
- 主程序中的配置文件路径

**命名规则：**
- 模块名使用连字符分隔（如 `my-new-service`）
- 服务名自动转换为下划线分隔并添加 `-api` 后缀（如 `my_new_service-api`）
- 文件名使用驼峰命名（如 `myNewService-Api.go`）

### 2. 整理依赖

```bash
make mt
```

### 3. 生成代码

```bash
make gen
```

### 4. 启动服务

```bash
make run
```

## 项目结构

```
.
├── api/                  # API 定义文件
├── internal/
│   ├── db/              # 数据库 Repository 层
│   ├── handler/         # HTTP 处理器
│   ├── logic/           # 业务逻辑层
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── request/         # 外部请求客户端
│   └── svc/             # 服务上下文
├── makefile             # Make 命令
└── init-project.sh      # 项目初始化脚本
```

## 开发规范

详细的开发规范请参考 `.cursor/rules/` 目录下的规则文件：

- `basics.mdc` - 项目基础配置
- `api.mdc` - API 开发规范
- `db.mdc` - 数据库使用规范
- `logging.mdc` - 日志规范
- `error-handling.mdc` - 错误处理规范
- `code-reuse.mdc` - 代码复用规范
- `style.mdc` - 代码风格规范

## 常用命令

| 命令 | 说明 |
|---|---|
| `make gen` | 基于 API 文件生成 Go 代码 |
| `make format` | 格式化 API 文件 |
| `make run` | 启动服务 |
| `make mt` | 整理 Go 模块依赖 |

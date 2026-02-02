# Cursor Commands

项目自定义的 Cursor 命令，用于提高开发效率和代码质量。

## 使用方法

1. 在 Cursor 编辑器中按 `Ctrl+K`（Windows/Linux）或 `Cmd+K`（Mac）
2. 输入命令名称并回车

## 可用命令

### 🔍 check-cache - 检查缓存架构

**用途**：检查项目缓存架构是否符合规范。

**使用场景**：
- 定期检查代码质量
- Code Review 前的自检
- 合并代码前的验证

**检查内容**：
- Logic 层是否直接使用缓存（应该通过 Repository）
- Repository 层是否正确实现缓存逻辑

**如何使用**：
```
1. 按 Ctrl+K 或 Cmd+K
2. 输入 "check-cache"
3. 查看检查报告
```

### 🔧 fix-cache - 修复缓存架构问题

**用途**：自动修复缓存架构违规代码。

**使用场景**：
- 发现缓存架构问题后
- 重构旧代码时
- 修复 check-cache 发现的问题

**修复内容**：
- 移除 Logic 层的缓存直接调用
- 为 Repository 添加缓存实现
- 清理重复的缓存操作

**如何使用**：
```
1. 选中有问题的代码（或打开文件）
2. 按 Ctrl+K 或 Cmd+K
3. 输入 "fix-cache"
4. AI 会自动修复问题
```

### 🛡️ guard - 添加权限检查

**用途**：为接口函数添加权限验证代码。

**使用场景**：
- 新增需要权限控制的接口
- 补充缺失的权限检查

**如何使用**：
```
1. 选中 Logic 函数
2. 按 Ctrl+K 或 Cmd+K
3. 输入 "guard"
4. AI 会在函数开头添加权限检查代码
```

### ⚡ opt - 代码逻辑优化

**用途**：分析并优化代码逻辑。

**使用场景**：
- 优化性能瓶颈
- 简化复杂逻辑
- 提高代码可读性

**检查项目**：
- 性能优化（不必要的循环、可并行处理）
- 代码简化（冗余代码、复杂条件）
- 错误处理（缺失的检查、边界情况）
- 可读性（命名、注释、结构）

**如何使用**：
```
1. 选中需要优化的代码
2. 按 Ctrl+K 或 Cmd+K
3. 输入 "opt"
4. AI 会分析并执行优化
```

## 架构规则

所有 commands 都遵循项目的架构规则（`.cursor/rules/`）：

- `cache-architecture.mdc` - 缓存架构规范
- `code-reuse.mdc` - 代码复用规则
- `db.mdc` - 数据库操作规范
- `logging.mdc` - 日志规则
- `error-handling.mdc` - 错误处理规范
- `api.mdc` - API 定义规则
- `style.mdc` - 代码风格

## 开发流程建议

1. **编写代码前**：了解相关规则（`.cursor/rules/`）
2. **编写代码中**：使用 commands 辅助开发（guard、opt）
3. **提交代码前**：运行检查 commands（check-cache）
4. **发现问题后**：使用修复 commands（fix-cache）

## 自定义 Commands

在 `.cursor/commands/` 目录下创建 `.md` 文件：

```markdown
---
description: 命令简短描述
customPrompt: |
  详细的 AI 提示词
  说明要执行的任务和规则
---
```

## 常见问题

**Q: Commands 和 Rules 有什么区别？**

A: 
- **Rules** (`.cursor/rules/`): 持久化的规则，AI 会自动遵循
- **Commands** (`.cursor/commands/`): 主动触发的任务，需要手动调用

**Q: 如何知道哪些 commands 可用？**

A: 按 `Ctrl+K` 或 `Cmd+K` 后，输入框会显示可用的命令列表。

**Q: Commands 会修改我的代码吗？**

A: 是的，`fix-cache` 等修复类 commands 会自动修改代码。建议：
- 使用前提交代码到 Git
- 检查修改内容后再提交
- 必要时可以撤销（`Ctrl+Z`）

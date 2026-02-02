---
description: 检查项目缓存架构是否符合规范
customPrompt: |
  请检查项目的缓存架构是否符合规范要求。
  
  ## 核心架构原则
  
  **Repository 层统一处理缓存，Logic 层禁止直接操作缓存。**
  
  ```
  Logic 层（业务逻辑）
    ↓ 调用
  Repository 层（数据访问 + 缓存管理）
    ↓ 操作
  数据库 / 缓存
  ```
  
  ## 检查项目
  
  ### 1. Logic 层缓存使用检查
  
  在 `internal/logic/` 目录下搜索以下模式，如果找到则违反规则：
  
  - ❌ `svcCtx.UserCache.`
  - ❌ `svcCtx.OrgCache.`
  - ❌ `svcCtx.RegionCache.`
  - ❌ 任何 `svcCtx.*Cache.` 模式
  
  **错误示例**：
  ```go
  // ❌ Logic 层手动调用缓存（重复）
  user, err := l.svcCtx.UserCache.GetUserByID(l.ctx, req.ID, func() (*model.User, error) {
      return userRepo.GetByID(l.ctx, req.ID)  // Repository 内部已处理缓存
  })
  
  // ❌ Logic 层手动清除缓存（重复）
  err := userRepo.Update(l.ctx, user)
  l.svcCtx.UserCache.DelUserByID(l.ctx, user.ID)  // Repository 已清除
  ```
  
  **正确示例**：
  ```go
  // ✅ 直接调用 Repository，内部自动处理缓存
  user, err := userRepo.GetByID(l.ctx, req.ID)
  
  // ✅ Repository.Update 内部已清除缓存
  err := userRepo.Update(l.ctx, user)
  ```
  
  ### 2. Repository 层缓存实现检查
  
  在 `internal/db/` 目录下检查以下方法是否实现了缓存清除：
  
  #### UserRepository (`internal/db/user.go`)
  
  - ✅ `GetByID` - 应使用 `userCache.GetUserByID`
  - ✅ `Update` - 更新后应调用 `userCache.DelUserByID`
  - ✅ `UpdateMap` - 更新后应调用 `userCache.DelUserByID`
  - ✅ `Delete` - 删除后应调用 `userCache.DelUserByID`
  - ✅ `UpdatePassword` - 更新后应调用 `userCache.DelUserByID`
  
  #### OrganizationRepository (`internal/db/organization.go`)
  
  - ✅ `GetByID` - 应使用 `orgCache.GetOrgByID`
  - ✅ `GetByName` - 应使用 `orgCache.GetOrgByName`
  - ✅ `Update` - 更新后应调用 `orgCache.DelOrgByID`
  - ✅ `Delete` - 删除后应调用 `orgCache.DelOrgByID`
  
  #### RegionRepository (`internal/db/region.go`)
  
  - Region 目前没有实现缓存（未来如需要再添加）
  
  **正确的 Repository 实现示例**：
  ```go
  // 查询方法 - 使用缓存
  func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
      if r.userCache != nil {
          return r.userCache.GetUserByID(ctx, id, func() (*model.User, error) {
              return FirstOrNil[model.User](r.db.WithContext(ctx).Where("id = ?", id))
          })
      }
      return FirstOrNil[model.User](r.db.WithContext(ctx).Where("id = ?", id))
  }
  
  // 更新方法 - 清除缓存
  func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
      err := r.db.WithContext(ctx).Save(user).Error
      if err != nil {
          return err
      }
      // 更新成功后自动清除缓存
      if r.userCache != nil {
          r.userCache.DelUserByID(ctx, user.ID)
      }
      return nil
  }
  ```
  
  ## 执行流程
  
  1. **搜索 Logic 层违规**：
     - 使用 Grep 工具搜索 `internal/logic/` 目录
     - 模式：`svcCtx\..*Cache\.`
     - 如果找到，列出所有违规文件和行号
  
  2. **检查 Repository 实现**：
     - 读取 `internal/db/user.go`
     - 读取 `internal/db/organization.go`
     - 读取 `internal/db/region.go`
     - 验证查询方法是否使用缓存
     - 验证更新/删除方法是否清除缓存
  
  3. **生成报告**：
     - 列出所有违规项
     - 列出缺失的缓存实现
     - 提供修复建议
  
  4. **如果发现问题，询问是否修复**：
     - 如果用户同意，自动修复所有问题
     - 修复后重新检查
  
  ## 输出格式
  
  ```
  🔍 缓存架构检查结果
  
  ## Logic 层检查
  ✅ 未发现直接使用缓存的代码
  或
  ❌ 发现 3 处违规：
     - internal/logic/user/getUserByIDLogic.go:35
     - internal/logic/user/updateUserLogic.go:234
     - internal/logic/organization/updateOrgLogic.go:105
  
  ## Repository 层检查
  ✅ UserRepository 正确实现缓存
  ✅ OrganizationRepository 正确实现缓存
  ⚠️  RegionRepository 未实现缓存（如需要请添加）
  
  ## 总结
  ✅ 缓存架构检查通过
  或
  ❌ 发现问题需要修复
  ```
  
  ## 注意事项
  
  1. 检查时不要修改任何代码，只生成报告
  2. 只有在用户明确同意后才执行修复
  3. 修复时遵循项目现有的代码风格
  4. 参考规则文件：`.cursor/rules/cache-architecture.mdc`
  5. 使用工具而不是命令行：优先使用 Grep 工具而不是 Shell grep 命令
---

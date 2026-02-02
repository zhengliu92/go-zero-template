---
description: 修复缓存架构违规问题
customPrompt: |
  请修复当前选中代码或文件中的缓存架构违规问题。
  
  ## 修复原则
  
  **Repository 层统一处理缓存，Logic 层禁止直接操作缓存。**
  
  ## 常见问题及修复方法
  
  ### 问题 1：Logic 层手动调用缓存查询（重复查询）
  
  **违规代码**：
  ```go
  // internal/logic/user/getUserByIDLogic.go
  func (l *GetUserByIDLogic) GetUserByID(req *types.GetUserByIDRequest) (*types.GetUserByIDResponse, error) {
      userRepo := l.svcCtx.Repository.User
      
      // ❌ 第一次缓存查询
      user, err := l.svcCtx.UserCache.GetUserByID(l.ctx, req.ID, func() (*model.User, error) {
          // ❌ Repository 内部会进行第二次缓存查询（重复）
          return userRepo.GetByID(l.ctx, req.ID)
      })
      // ...
  }
  ```
  
  **修复后**：
  ```go
  func (l *GetUserByIDLogic) GetUserByID(req *types.GetUserByIDRequest) (*types.GetUserByIDResponse, error) {
      userRepo := l.svcCtx.Repository.User
      
      // ✅ 直接调用 Repository，内部自动处理缓存
      user, err := userRepo.GetByID(l.ctx, req.ID)
      // ...
  }
  ```
  
  ### 问题 2：Logic 层手动清除缓存（重复清除）
  
  **违规代码**：
  ```go
  // internal/logic/user/updateUserLogic.go
  func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) error {
      userRepo := l.svcCtx.Repository.User
      
      // Repository.UpdateMap 内部已经清除缓存
      err := userRepo.UpdateMap(l.ctx, req.ID, updateData)
      if err != nil {
          return err
      }
      
      // ❌ 重复清除缓存
      if err := l.svcCtx.UserCache.DelUserByID(l.ctx, req.ID); err != nil {
          l.svcCtx.Writer.Error("删除用户缓存失败", ...)
      }
      
      return nil
  }
  ```
  
  **修复后**：
  ```go
  func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) error {
      userRepo := l.svcCtx.Repository.User
      
      // ✅ Repository.UpdateMap 内部已自动清除缓存，无需手动清除
      err := userRepo.UpdateMap(l.ctx, req.ID, updateData)
      if err != nil {
          return err
      }
      
      return nil
  }
  ```
  
  ### 问题 3：Repository 未实现缓存清除
  
  **违规代码**：
  ```go
  // internal/db/user.go
  func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
      // ❌ 更新后未清除缓存
      return r.db.WithContext(ctx).Save(user).Error
  }
  ```
  
  **修复后**：
  ```go
  func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
      err := r.db.WithContext(ctx).Save(user).Error
      if err != nil {
          return err
      }
      // ✅ 更新成功后自动清除缓存
      if r.userCache != nil {
          if err := r.userCache.DelUserByID(ctx, user.ID); err != nil {
              // 缓存清除失败不影响业务逻辑
          }
      }
      return nil
  }
  ```
  
  ### 问题 4：Repository 未实现缓存查询
  
  **违规代码**：
  ```go
  // internal/db/user.go
  func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
      // ❌ 直接查询数据库，未使用缓存
      return FirstOrNil[model.User](r.db.WithContext(ctx).Where("id = ?", id))
  }
  ```
  
  **修复后**：
  ```go
  func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
      // ✅ 如果缓存可用，使用缓存（Cache-Aside 模式）
      if r.userCache != nil {
          return r.userCache.GetUserByID(ctx, id, func() (*model.User, error) {
              return FirstOrNil[model.User](r.db.WithContext(ctx).Where("id = ?", id))
          })
      }
      // 缓存不可用，直接查询数据库
      return FirstOrNil[model.User](r.db.WithContext(ctx).Where("id = ?", id))
  }
  ```
  
  ## 修复步骤
  
  1. **识别问题类型**：
     - Logic 层直接使用缓存？
     - Repository 层缺少缓存实现？
  
  2. **应用对应的修复模式**：
     - 使用上述示例作为参考
     - 保持代码风格一致
  
  3. **清理代码**：
     - 移除不必要的变量（如保存旧名称用于缓存清除）
     - 移除重复的缓存操作日志
     - 简化注释
  
  4. **验证修复**：
     - 确保逻辑功能不变
     - 确保错误处理完整
  
  ## 适用实体
  
  - **UserRepository** + **UserCache**
  - **OrganizationRepository** + **OrgCache**
  - **RegionRepository**（目前无缓存，如需要参考 User 实现）
  
  ## 注意事项
  
  1. **保持功能不变**：只修复架构问题，不改变业务逻辑
  2. **错误处理**：缓存操作失败不应影响主业务流程
  3. **代码风格**：遵循项目现有风格
  4. **注释更新**：修改后更新相关注释
  5. **测试**：修复后建议运行相关测试
  
  ## 执行流程
  
  1. 分析选中代码，识别违规类型
  2. 应用对应的修复模式
  3. 执行修复（使用 StrReplace 工具）
  4. 简要说明修复内容
  
  修复完成后不要创建总结文档。
---

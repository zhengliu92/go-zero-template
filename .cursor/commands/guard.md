---
description: 为接口添加权限检查（Guard）功能
customPrompt: |
  请为当前选中的 Logic 函数添加权限检查（Guard）功能。
  
  ## 要求
  
  1. **获取当前用户**：
     ```go
     currentUser := middleware.GetUserFromContext(l.ctx)
     ```
  
  2. **执行权限检查**：
     使用 `utils.RouteGuard` 进行权限验证（默认使用 `utils.ManagerGroup`）
  
  3. **记录权限拒绝日志**：
     当权限检查失败时，使用 `l.svcCtx.Writer.Error` 记录日志，包含：
     - `log_type`: "user"
     - `trace`: 操作类型（如 "DeleteUser", "UpdateUser" 等）
     - `user_id`: 当前用户 ID
     - `user_name`: 当前用户姓名
     - `role_code`: 当前用户角色
     - 操作目标信息（如 `target_user_id`, `target_org_id` 等）
  
  4. **返回错误**：
     ```go
     return nil, response.NewError(http.StatusForbidden, "无权限执行此操作")
     ```
  
  ## 可用的角色组
  
  - `utils.ManagerGroup`: `[]utils.RoleType{RoleSuperAdmin, RoleAdmin, RoleRegionManager}` **(默认推荐)**
  - `utils.AllRoles`: 所有角色
  - 自定义角色组：`[]utils.RoleType{utils.RoleSuperAdmin, utils.RoleAdmin}`
  
  ## 可用的角色类型
  
  - `utils.RoleSuperAdmin`: 超级管理员
  - `utils.RoleAdmin`: 管理员
  - `utils.RoleSuperReviewer`: 超级审核员
  - `utils.RoleOperator`: 操作员
  - `utils.RoleRegionManager`: 区域管理员
  
  ## 实现示例
  
  ```go
  // 1. 获取当前用户
  currentUser := middleware.GetUserFromContext(l.ctx)
  
  // 2. 权限检查
  ok := utils.RouteGuard(currentUser.RoleCode, utils.ManagerGroup)
  if !ok {
      // 3. 记录权限拒绝日志
      l.svcCtx.Writer.Error("无权限删除用户",
          writer.Field("log_type", "user"),
          writer.Field("trace", "DeleteUser"),
          writer.Field("user_id", currentUser.ID),
          writer.Field("username", currentUser.Name),
          writer.Field("role_code", currentUser.RoleCode),
          writer.Field("target_user_id", req.ID),
      )
      // 4. 返回错误
      return nil, response.NewError(http.StatusForbidden, "无权限删除用户")
  }
  ```
  
  ## 注意事项
  
  1. 将 guard 代码放在函数开头，业务逻辑之前
  2. **默认使用 `utils.ManagerGroup` 作为权限组**，除非有特殊业务需求
  3. 根据实际操作类型调整 `trace` 字段值
  4. 根据操作目标调整目标信息字段（如 `target_user_id`, `target_org_id` 等）
  5. 错误消息应清晰描述用户缺少的权限
  
  
  收集到这些信息后，在当前光标位置或函数开头插入完整的 guard 代码。
---

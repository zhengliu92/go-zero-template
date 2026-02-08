# Go-Zero Development Examples

This document provides concrete examples for common patterns in the Go-Zero template.

## Repository Pattern Examples

### Basic Repository Implementation

```go
// internal/db/user.go
type UserRepository struct {
    db        *gorm.DB
    userCache cache.UserCache
}

func NewUserRepository(db *gorm.DB, userCache cache.UserCache) *UserRepository {
    return &UserRepository{
        db:        db,
        userCache: userCache,
    }
}

// Get by ID with caching
func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
    if r.userCache != nil {
        return r.userCache.GetUserByID(ctx, id, func() (*model.User, error) {
            return FirstOrNil[model.User](r.db.WithContext(ctx).Where("id = ?", id))
        })
    }
    return FirstOrNil[model.User](r.db.WithContext(ctx).Where("id = ?", id))
}

// Get by login name
func (r *UserRepository) GetByLoginName(ctx context.Context, loginName string) (*model.User, error) {
    return FirstOrNil[model.User](r.db.WithContext(ctx).Where("login_name = ?", loginName))
}

// List with filters
func (r *UserRepository) List(ctx context.Context, filters UserFilters) ([]*model.User, error) {
    query := r.db.WithContext(ctx).Model(&model.User{})
    
    if filters.OrgID != nil {
        query = query.Where("org_id = ?", *filters.OrgID)
    }
    
    if filters.IsInternal != nil {
        query = query.Where("is_internal = ?", *filters.IsInternal)
    }
    
    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }
    
    var users []*model.User
    err := query.Find(&users).Error
    return users, err
}

// Create
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

// Update
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
    err := r.db.WithContext(ctx).Save(user).Error
    if err != nil {
        return err
    }
    // Clear cache after update
    if r.userCache != nil {
        r.userCache.DelUserByID(ctx, user.ID)
    }
    return nil
}

// Partial update with map
func (r *UserRepository) UpdateFields(ctx context.Context, id int, updates map[string]interface{}) error {
    err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
    if err != nil {
        return err
    }
    // Clear cache after update
    if r.userCache != nil {
        r.userCache.DelUserByID(ctx, id)
    }
    return nil
}

// Delete
func (r *UserRepository) Delete(ctx context.Context, id int) error {
    err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error
    if err != nil {
        return err
    }
    // Clear cache after delete
    if r.userCache != nil {
        r.userCache.DelUserByID(ctx, id)
    }
    return nil
}
```

### Repository Aggregation

```go
// internal/db/repo.go
type Repository struct {
    User   *UserRepository
    Org    *OrgRepository
    Region *RegionRepository
}

func NewRepository(db *gorm.DB, caches *cache.Caches) *Repository {
    return &Repository{
        User:   NewUserRepository(db, caches.UserCache),
        Org:    NewOrgRepository(db, caches.OrgCache),
        Region: NewRegionRepository(db, caches.RegionCache),
    }
}
```

## Logic Layer Examples

### Create Operation

```go
// internal/logic/user/createUserLogic.go
func (l *CreateUserLogic) CreateUser(req *types.CreateUserRequest) (*types.CreateUserResponse, error) {
    // 1. Validate input
    if req.LoginName == "" {
        return nil, errors.New("login_name is required")
    }
    
    // 2. Check if user exists
    existing, err := l.svcCtx.Repository.User.GetByLoginName(l.ctx, req.LoginName)
    if err != nil {
        l.svcCtx.Writer.Error("Database query failed",
            writer.Field("log_type", "database"),
            writer.Field("trace", "User.CreateUser"),
            writer.Field("error", err.Error()),
        )
        return nil, err
    }
    
    if existing != nil {
        return nil, errors.New("user already exists")
    }
    
    // 3. Build user model with defaults
    user := &model.User{
        LoginName:  req.LoginName,
        Name:       req.Name,
        OrgID:      req.OrgID,
        IsInternal: boolWithDefault(req.IsInternal, false),
        Status:     "active",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    
    // 4. Create user
    err = l.svcCtx.Repository.User.Create(l.ctx, user)
    if err != nil {
        l.svcCtx.Writer.Error("Failed to create user",
            writer.Field("log_type", "database"),
            writer.Field("trace", "User.CreateUser"),
            writer.Field("error", err.Error()),
        )
        return nil, err
    }
    
    // 5. Log success
    l.svcCtx.Writer.Info("User created",
        writer.Field("log_type", "user"),
        writer.Field("trace", "User.CreateUser"),
        writer.Field("user_id", user.ID),
        writer.Field("login_name", user.LoginName),
    )
    
    // 6. Return response
    return &types.CreateUserResponse{
        User: toUserResponse(user),
    }, nil
}
```

### Update Operation

```go
// internal/logic/user/updateUserLogic.go
func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (*types.UpdateUserResponse, error) {
    // 1. Get existing user
    user, err := l.svcCtx.Repository.User.GetByID(l.ctx, req.ID)
    if err != nil {
        l.svcCtx.Writer.Error("Database query failed",
            writer.Field("log_type", "database"),
            writer.Field("trace", "User.UpdateUser"),
            writer.Field("error", err.Error()),
        )
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    // 2. Build update map (only non-nil fields)
    updateData := make(map[string]interface{})
    
    if req.Name != nil {
        updateData["name"] = *req.Name
    }
    
    if req.IsInternal != nil {
        updateData["is_internal"] = *req.IsInternal
    }
    
    if req.Status != nil {
        updateData["status"] = *req.Status
    }
    
    updateData["updated_at"] = time.Now()
    
    // 3. Update user
    err = l.svcCtx.Repository.User.UpdateFields(l.ctx, req.ID, updateData)
    if err != nil {
        l.svcCtx.Writer.Error("Failed to update user",
            writer.Field("log_type", "database"),
            writer.Field("trace", "User.UpdateUser"),
            writer.Field("error", err.Error()),
        )
        return nil, err
    }
    
    // 4. Get updated user
    user, err = l.svcCtx.Repository.User.GetByID(l.ctx, req.ID)
    if err != nil {
        return nil, err
    }
    
    // 5. Log success
    l.svcCtx.Writer.Info("User updated",
        writer.Field("log_type", "user"),
        writer.Field("trace", "User.UpdateUser"),
        writer.Field("user_id", user.ID),
    )
    
    return &types.UpdateUserResponse{
        User: toUserResponse(user),
    }, nil
}
```

### List Operation with Filters

```go
// internal/logic/user/listUsersLogic.go
func (l *ListUsersLogic) ListUsers(req *types.ListUsersRequest) (*types.ListUsersResponse, error) {
    // 1. Build filters
    filters := UserFilters{
        OrgID:      req.OrgID,
        IsInternal: req.IsInternal,
        Status:     req.Status,
    }
    
    // 2. Query users
    users, err := l.svcCtx.Repository.User.List(l.ctx, filters)
    if err != nil {
        l.svcCtx.Writer.Error("Database query failed",
            writer.Field("log_type", "database"),
            writer.Field("trace", "User.ListUsers"),
            writer.Field("error", err.Error()),
        )
        return nil, err
    }
    
    // 3. Convert to response
    var userList []*types.User
    for _, user := range users {
        userList = append(userList, toUserResponse(user))
    }
    
    return &types.ListUsersResponse{
        Users: userList,
        Total: len(userList),
    }, nil
}
```

## Helper Function Examples

### User Conversion

```go
// internal/logic/user/convert.go
func toUserResponse(user *model.User) *types.User {
    return &types.User{
        ID:         user.ID,
        LoginName:  user.LoginName,
        Name:       user.Name,
        OrgID:      user.OrgID,
        IsInternal: user.IsInternal,
        Status:     user.Status,
        CreatedAt:  user.CreatedAt.Format(time.RFC3339),
        UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
    }
}

func toUserModel(req *types.User) *model.User {
    return &model.User{
        ID:         req.ID,
        LoginName:  req.LoginName,
        Name:       req.Name,
        OrgID:      req.OrgID,
        IsInternal: req.IsInternal,
        Status:     req.Status,
    }
}
```

### Bool Helpers

```go
// internal/logic/user/helper.go
func boolPtr(b bool) *bool {
    return &b
}

func boolWithDefault(b *bool, defaultValue bool) *bool {
    if b == nil {
        return &defaultValue
    }
    return b
}

func boolValue(b *bool, defaultValue bool) bool {
    if b == nil {
        return defaultValue
    }
    return *b
}
```

### Shared Business Logic

```go
// internal/logic/user/helper.go

// syncManagerInfo syncs manager information (extracted from duplicate code)
func syncManagerInfo(ctx context.Context, svcCtx *svc.ServiceContext, user *model.User, trace string) error {
    if user.ManagerID == nil {
        return nil
    }
    
    manager, err := svcCtx.Repository.User.GetByID(ctx, *user.ManagerID)
    if err != nil {
        svcCtx.Writer.Error("Failed to get manager",
            writer.Field("log_type", "database"),
            writer.Field("trace", trace),
            writer.Field("error", err.Error()),
        )
        return err
    }
    
    if manager == nil {
        return errors.New("manager not found")
    }
    
    // Update user with manager info
    updateData := map[string]interface{}{
        "manager_name": manager.Name,
        "manager_org_id": manager.OrgID,
        "updated_at": time.Now(),
    }
    
    return svcCtx.Repository.User.UpdateFields(ctx, user.ID, updateData)
}
```

## Error Handling Examples

### HTTP Client Error Handling

```go
// internal/handler/sync/userSyncHandler.go
func (h *UserSyncHandler) syncUser(userID int) error {
    // 1. Call external API
    externalUser, err := h.client.GetUser(h.token, userID)
    
    // 2. Check if 404
    var respErr *types.ResponseError
    is404 := errors.As(err, &respErr) && respErr.Code == 404
    
    // 3. Return non-404 errors
    if err != nil && !is404 {
        h.writer.Error("API call failed",
            writer.Field("log_type", "system"),
            writer.Field("trace", "UserSync.syncUser"),
            writer.Field("error", err.Error()),
        )
        return err
    }
    
    // 4. Handle not found
    if externalUser == nil || is404 {
        return h.createLocalUser(userID)
    }
    
    // 5. User exists, update local copy
    return h.updateLocalUser(externalUser)
}

func (h *UserSyncHandler) createLocalUser(userID int) error {
    // Create user logic
}

func (h *UserSyncHandler) updateLocalUser(externalUser *ExternalUser) error {
    // Update user logic
}
```

### Database Error Handling

```go
func (l *CreateUserLogic) CreateUser(req *types.CreateUserRequest) (*types.CreateUserResponse, error) {
    user := &model.User{
        LoginName: req.LoginName,
        Name:      req.Name,
    }
    
    err := l.svcCtx.Repository.User.Create(l.ctx, user)
    if err != nil {
        // Check for specific errors
        if strings.Contains(err.Error(), "duplicate key") {
            return nil, errors.New("user already exists")
        }
        
        // Log system error
        l.svcCtx.Writer.Error("Database error",
            writer.Field("log_type", "database"),
            writer.Field("trace", "User.CreateUser"),
            writer.Field("error", err.Error()),
        )
        return nil, err
    }
    
    return &types.CreateUserResponse{User: toUserResponse(user)}, nil
}
```

## API Definition Examples

### Complete API File

```go
// api/user.api
syntax = "v1"

info (
    title: "User Management API"
    desc: "User CRUD operations"
    version: "1.0"
)

type User {
    ID int64 `json:"id"` // User ID
    LoginName string `json:"login_name"` // Login name
    Name string `json:"name"` // Display name
    OrgID int64 `json:"org_id"` // Organization ID
    IsInternal *bool `json:"is_internal"` // Is internal employee
    Status string `json:"status"` // Status: active-active, inactive-inactive
    CreatedAt string `json:"created_at"` // Creation time (ISO8601)
    UpdatedAt string `json:"updated_at"` // Update time (ISO8601)
}

type CreateUserRequest {
    LoginName string `json:"login_name"` // Login name (required)
    Name string `json:"name"` // Display name (required)
    OrgID int64 `json:"org_id"` // Organization ID (required)
    IsInternal *bool `json:"is_internal,optional"` // Is internal employee
}

type CreateUserResponse {
    User User `json:"user"` // Created user
}

type UpdateUserRequest {
    ID int64 `path:"id"` // User ID
    Name *string `json:"name,optional"` // Display name
    IsInternal *bool `json:"is_internal,optional"` // Is internal employee
    Status *string `json:"status,optional"` // Status: active-active, inactive-inactive
}

type UpdateUserResponse {
    User User `json:"user"` // Updated user
}

type GetUserRequest {
    ID int64 `path:"id"` // User ID
}

type GetUserResponse {
    User User `json:"user"` // User details
}

type ListUsersRequest {
    OrgID *int64 `form:"org_id,optional"` // Filter by organization ID
    IsInternal *bool `form:"is_internal,optional"` // Filter by internal status
    Status string `form:"status,optional"` // Filter by status
}

type ListUsersResponse {
    Users []User `json:"users"` // User list
    Total int `json:"total"` // Total count
}

type DeleteUserRequest {
    ID int64 `path:"id"` // User ID
}

type DeleteUserResponse {
    Success bool `json:"success"` // Deletion success
}

@server (
    prefix: /api/v1
    group: user
)
service user-api {
    @doc (
        summary: "Create user"
        description: "Create a new user, requires admin permission"
    )
    @handler CreateUser
    post /users (CreateUserRequest) returns (CreateUserResponse)
    
    @doc (
        summary: "Update user"
        description: "Update user information, requires admin permission"
    )
    @handler UpdateUser
    put /users/:id (UpdateUserRequest) returns (UpdateUserResponse)
    
    @doc (
        summary: "Get user details"
        description: "Get user information by ID"
    )
    @handler GetUser
    get /users/:id (GetUserRequest) returns (GetUserResponse)
    
    @doc (
        summary: "List users"
        description: "List users with optional filters"
    )
    @handler ListUsers
    get /users (ListUsersRequest) returns (ListUsersResponse)
    
    @doc (
        summary: "Delete user"
        description: "Delete user by ID, requires admin permission"
    )
    @handler DeleteUser
    delete /users/:id (DeleteUserRequest) returns (DeleteUserResponse)
}
```

## Testing Examples

### Repository Test

```go
// internal/db/user_test.go
func TestUserRepository_GetByID(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db, nil)
    
    // Create test user
    user := &model.User{
        LoginName:  "test",
        Name:       "Test User",
        IsInternal: boolPtr(true),
        Status:     "active",
    }
    err := repo.Create(context.Background(), user)
    require.NoError(t, err)
    
    // Test GetByID
    found, err := repo.GetByID(context.Background(), user.ID)
    require.NoError(t, err)
    require.NotNil(t, found)
    assert.Equal(t, "test", found.LoginName)
    assert.Equal(t, true, *found.IsInternal)
}
```

### Logic Test

```go
// internal/logic/user/createUserLogic_test.go
func TestCreateUserLogic_CreateUser(t *testing.T) {
    // Setup
    ctx := context.Background()
    svcCtx := createTestServiceContext(t)
    logic := NewCreateUserLogic(ctx, svcCtx)
    
    // Test
    req := &types.CreateUserRequest{
        LoginName:  "test",
        Name:       "Test User",
        OrgID:      1,
        IsInternal: boolPtr(true),
    }
    
    resp, err := logic.CreateUser(req)
    
    // Verify
    require.NoError(t, err)
    require.NotNil(t, resp)
    assert.Equal(t, "test", resp.User.LoginName)
    assert.NotZero(t, resp.User.ID)
}
```

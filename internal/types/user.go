package types

// UserBase 用户请求体共用字段（创建/更新），不含 ID、时间戳、Password
type UserBase struct {
	Name             string `json:"name,omitempty"`
	NameOptimized    string `json:"name_optimized,omitempty"`
	LoginName        string `json:"login_name,omitempty"`
	SAPEmployeeID    int    `json:"sap_employee_id,omitempty"`
	Status           int16  `json:"status,omitempty"`
	RoleCode         string `json:"role_code,omitempty"`
	OrgID            *int   `json:"org_id,omitempty"`
	OrgName          string `json:"org_name,omitempty"`
	OrgLevel1ID      *int   `json:"org_level1_id,omitempty"`
	OrgLevel1Name    string `json:"org_level1_name,omitempty"`
	OrgLevel2ID      *int   `json:"org_level2_id,omitempty"`
	OrgLevel2Name    string `json:"org_level2_name,omitempty"`
	OrgLevel3ID      *int   `json:"org_level3_id,omitempty"`
	OrgLevel3Name    string `json:"org_level3_name,omitempty"`
	OrgLevel4ID      *int   `json:"org_level4_id,omitempty"`
	OrgLevel4Name    string `json:"org_level4_name,omitempty"`
	OrgLevel5ID      *int   `json:"org_level5_id,omitempty"`
	OrgLevel5Name    string `json:"org_level5_name,omitempty"`
	OrgLevel6ID      *int   `json:"org_level6_id,omitempty"`
	OrgLevel6Name    string `json:"org_level6_name,omitempty"`
	OrgLevel7ID      *int   `json:"org_level7_id,omitempty"`
	OrgLevel7Name    string `json:"org_level7_name,omitempty"`
	OrgLevel8ID      *int   `json:"org_level8_id,omitempty"`
	OrgLevel8Name    string `json:"org_level8_name,omitempty"`
	OrgLevel9ID      *int   `json:"org_level9_id,omitempty"`
	OrgLevel9Name    string `json:"org_level9_name,omitempty"`
	EmployeePost     int    `json:"employee_post,omitempty"`
	EmployeePostName string `json:"employee_post_name,omitempty"`
	BossEmployeeID   *int   `json:"boss_employee_id,omitempty"`
	BossName         string `json:"boss_name,omitempty"`
	BossEmployeePost *int   `json:"boss_employee_post,omitempty"`
	RegionID         *int   `json:"region_id,omitempty"`
	RegionName       string `json:"region_name,omitempty"`
	LineID           *int   `json:"line_id,omitempty"`
	LineName         string `json:"line_name,omitempty"`
	CustomerID       *int   `json:"customer_id,omitempty"`
	CustomerName     string `json:"customer_name,omitempty"`
	EmployeeQY       string `json:"employee_qy,omitempty"`
	IsInternal       bool   `json:"is_internal,omitempty"`
	IsAutoAdd        bool   `json:"is_auto_add,omitempty"`
	IsBackendSynced  bool   `json:"is_backend_synced,omitempty"`
	IsManualUpdate   bool   `json:"is_manual_update,omitempty"`
	Comment          string `json:"comment,omitempty"`
}

// User 用户信息结构（响应），复用 UserBase 并增加 ID 与时间戳
type User struct {
	UserBase
	ID           int     `json:"id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	LastActiveAt *string `json:"last_active_at"`
}

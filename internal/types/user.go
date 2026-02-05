package types

// UserBase 用户请求体共用字段（创建/更新），不含 ID、时间戳、Password
type UserBase struct {
	Name                  string  `json:"name"`                     // 用户姓名
	NameOptimized         string  `json:"name_optimized"`           // 优化后的用户姓名（用于搜索）
	LoginName             string  `json:"login_name"`               // 登录名
	SAPEmployeeID         *int    `json:"sap_employee_id"`          // SAP员工ID
	Status                int16   `json:"status"`                   // 用户状态：1-在职，2-离职
	RoleCode              string  `json:"role_code"`                // 角色编码
	OrgID                 *int    `json:"org_id"`                   // 当前组织ID
	OrgName               string  `json:"org_name"`                 // 当前组织名称
	OrgLevel1ID           *int    `json:"org_level1_id"`            // 一级组织ID
	OrgLevel1Name         string  `json:"org_level1_name"`          // 一级组织名称
	OrgLevel2ID           *int    `json:"org_level2_id"`            // 二级组织ID
	OrgLevel2Name         string  `json:"org_level2_name"`          // 二级组织名称
	OrgLevel3ID           *int    `json:"org_level3_id"`            // 三级组织ID
	OrgLevel3Name         string  `json:"org_level3_name"`          // 三级组织名称
	OrgLevel4ID           *int    `json:"org_level4_id"`            // 四级组织ID
	OrgLevel4Name         string  `json:"org_level4_name"`          // 四级组织名称
	OrgLevel5ID           *int    `json:"org_level5_id"`            // 五级组织ID
	OrgLevel5Name         string  `json:"org_level5_name"`          // 五级组织名称
	OrgLevel6ID           *int    `json:"org_level6_id"`            // 六级组织ID
	OrgLevel6Name         string  `json:"org_level6_name"`          // 六级组织名称
	OrgLevel7ID           *int    `json:"org_level7_id"`            // 七级组织ID
	OrgLevel7Name         string  `json:"org_level7_name"`          // 七级组织名称
	OrgLevel8ID           *int    `json:"org_level8_id"`            // 八级组织ID
	OrgLevel8Name         string  `json:"org_level8_name"`          // 八级组织名称
	OrgLevel9ID           *int    `json:"org_level9_id"`            // 九级组织ID
	OrgLevel9Name         string  `json:"org_level9_name"`          // 九级组织名称
	EmployeePost          *int    `json:"employee_post"`            // 员工职位代码
	EmployeePostName      *string `json:"employee_post_name"`       // 员工职位名称
	OutSourcePositionCode *string `json:"out_source_position_code"` // 外包人员职位代码
	BossEmployeeID        *int    `json:"boss_employee_id"`         // 上级员工ID
	BossName              string  `json:"boss_name"`                // 上级姓名
	BossEmployeePost      *int    `json:"boss_employee_post"`       // 上级职位代码
	LineID                *int    `json:"line_id"`                  // 产线ID
	LineName              string  `json:"line_name"`                // 产线名称
	CustomerID            *int    `json:"customer_id"`              // 客户ID
	CustomerName          string  `json:"customer_name"`            // 客户名称
	EmployeeQY            string  `json:"employee_qy"`              // 员工区域标识
	IsInternal            *bool   `json:"is_internal"`              // 是否内部员工
	IsAutoAdd             *bool   `json:"is_auto_add"`              // 是否自动添加
	IsBackendSynced       *bool   `json:"is_backend_synced"`        // 是否后端已同步
	IsManualUpdate        *bool   `json:"is_manual_update"`         // 是否手动更新
	Comment               string  `json:"comment"`                  // 备注信息
}

// User 用户信息结构（响应），复用 UserBase 并增加 ID 与时间戳
type User struct {
	UserBase
	ID           int     `json:"id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	LastActiveAt *string `json:"last_active_at"`
}

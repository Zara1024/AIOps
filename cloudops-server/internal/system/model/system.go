package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	PasswordHash   string         `gorm:"type:varchar(255);not null" json:"-"`
	Nickname       string         `gorm:"type:varchar(128)" json:"nickname"`
	Email          string         `gorm:"type:varchar(128)" json:"email"`
	Phone          string         `gorm:"type:varchar(20)" json:"phone"`
	Avatar         string         `gorm:"type:varchar(255)" json:"avatar"`
	DepartmentID   *int64         `gorm:"index" json:"department_id"`
	Status         int            `gorm:"type:smallint;default:1" json:"status"` // 1正常 0禁用
	LastLoginAt    *time.Time     `json:"last_login_at"`
	LastLoginIP    string         `gorm:"type:varchar(45)" json:"last_login_ip"`
	MFAEnabled     bool           `gorm:"default:false" json:"mfa_enabled"`
	MFASecret      string         `gorm:"type:varchar(128)" json:"-"`
	LoginFailCount int            `gorm:"default:0" json:"-"` // 登录失败计数
	LockUntil      *time.Time     `json:"-"`                  // 锁定截止时间
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Roles      []Role      `gorm:"many2many:sys_user_roles;" json:"roles,omitempty"`
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

func (User) TableName() string { return "sys_users" }

// Role 角色模型
type Role struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"role_name"`
	RoleKey     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"role_key"`
	Description string         `gorm:"type:text" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      int            `gorm:"type:smallint;default:1" json:"status"` // 1正常 0禁用
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Menus []Menu `gorm:"many2many:sys_role_menus;" json:"menus,omitempty"`
}

func (Role) TableName() string { return "sys_roles" }

// Menu 菜单模型（动态路由 + 按钮权限）
type Menu struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID   int64          `gorm:"default:0;index" json:"parent_id"`
	MenuName   string         `gorm:"type:varchar(64);not null" json:"menu_name"`
	MenuType   int            `gorm:"type:smallint" json:"menu_type"` // 1目录 2菜单 3按钮
	Path       string         `gorm:"type:varchar(255)" json:"path"`
	Component  string         `gorm:"type:varchar(255)" json:"component"`
	Icon       string         `gorm:"type:varchar(64)" json:"icon"`
	Permission string         `gorm:"type:varchar(128)" json:"permission"` // 按钮权限标识
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	Visible    bool           `gorm:"default:true" json:"visible"`
	Status     int            `gorm:"type:smallint;default:1" json:"status"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 非数据库字段
	Children []*Menu `gorm:"-" json:"children,omitempty"`
}

func (Menu) TableName() string { return "sys_menus" }

// Department 部门模型
type Department struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID  int64          `gorm:"default:0;index" json:"parent_id"`
	DeptName  string         `gorm:"type:varchar(128);not null" json:"dept_name"`
	Leader    string         `gorm:"type:varchar(64)" json:"leader"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Email     string         `gorm:"type:varchar(128)" json:"email"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    int            `gorm:"type:smallint;default:1" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 非数据库字段
	Children []*Department `gorm:"-" json:"children,omitempty"`
}

func (Department) TableName() string { return "sys_departments" }

// OperationLog 操作日志
type OperationLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	Username  string    `gorm:"type:varchar(64)" json:"username"`
	Module    string    `gorm:"type:varchar(64)" json:"module"`
	Action    string    `gorm:"type:varchar(32)" json:"action"` // CREATE/UPDATE/DELETE
	Resource  string    `gorm:"type:varchar(128)" json:"resource"`
	ResourceID string   `gorm:"type:varchar(64)" json:"resource_id"`
	Detail    string    `gorm:"type:text" json:"detail"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	UserAgent string    `gorm:"type:varchar(512)" json:"user_agent"`
	Status    int       `gorm:"type:smallint;default:1" json:"status"` // 1成功 0失败
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (OperationLog) TableName() string { return "sys_operation_logs" }

// LoginLog 登录日志
type LoginLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	Username  string    `gorm:"type:varchar(64)" json:"username"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	Location  string    `gorm:"type:varchar(128)" json:"location"`
	Browser   string    `gorm:"type:varchar(128)" json:"browser"`
	OS        string    `gorm:"type:varchar(64)" json:"os"`
	Status    int       `gorm:"type:smallint" json:"status"` // 1成功 0失败
	Message   string    `gorm:"type:varchar(255)" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (LoginLog) TableName() string { return "sys_login_logs" }

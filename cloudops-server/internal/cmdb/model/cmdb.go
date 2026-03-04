package model

import (
	"time"

	"gorm.io/gorm"
)

// Host 主机模型
type Host struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Hostname      string         `gorm:"type:varchar(128);not null" json:"hostname"`
	IP            string         `gorm:"type:varchar(45);uniqueIndex;not null" json:"ip"`
	Port          int            `gorm:"type:int;default:22" json:"port"`
	OSType        string         `gorm:"type:varchar(32);default:'linux'" json:"os_type"` // linux / windows
	OSVersion     string         `gorm:"type:varchar(128)" json:"os_version"`
	CPU           string         `gorm:"type:varchar(64)" json:"cpu"`
	Memory        string         `gorm:"type:varchar(32)" json:"memory"`
	Disk          string         `gorm:"type:varchar(128)" json:"disk"`
	Kernel        string         `gorm:"type:varchar(128)" json:"kernel"`
	Uptime        string         `gorm:"type:varchar(64)" json:"uptime"`
	AuthType      string         `gorm:"type:varchar(16);default:'password'" json:"auth_type"` // password / key
	Username      string         `gorm:"type:varchar(64);default:'root'" json:"username"`
	Password      string         `gorm:"type:varchar(512)" json:"-"`     // 加密存储
	PrivateKey    string         `gorm:"type:text" json:"-"`             // 加密存储
	GroupID       *int64         `gorm:"index" json:"group_id"`
	Status        int            `gorm:"type:smallint;default:1" json:"status"` // 1在线 2离线 3未知 0禁用
	AgentStatus   int            `gorm:"type:smallint;default:0" json:"agent_status"` // 0未安装 1正常 2异常
	Description   string         `gorm:"type:text" json:"description"`
	Labels        string         `gorm:"type:text" json:"labels"` // JSON格式标签
	CloudProvider string         `gorm:"type:varchar(32)" json:"cloud_provider"` // aliyun/tencent/aws/huawei
	CloudID       string         `gorm:"type:varchar(128)" json:"cloud_id"` // 云实例ID
	Region        string         `gorm:"type:varchar(64)" json:"region"`
	LastCheckAt   *time.Time     `json:"last_check_at"`
	CreatedBy     int64          `gorm:"index" json:"created_by"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Group *HostGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

func (Host) TableName() string { return "cmdb_hosts" }

// HostGroup 主机分组（树形结构）
type HostGroup struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID  int64          `gorm:"default:0;index" json:"parent_id"`
	GroupName string         `gorm:"type:varchar(128);not null" json:"group_name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 非数据库字段
	Children  []*HostGroup `gorm:"-" json:"children,omitempty"`
	HostCount int64        `gorm:"-" json:"host_count"`
}

func (HostGroup) TableName() string { return "cmdb_host_groups" }

// SSHRecord SSH 操作记录（审计日志）
type SSHRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	Username  string    `gorm:"type:varchar(64)" json:"username"`
	HostID    int64     `gorm:"index" json:"host_id"`
	HostIP    string    `gorm:"type:varchar(45)" json:"host_ip"`
	Command   string    `gorm:"type:text" json:"command"`
	Output    string    `gorm:"type:text" json:"output"`
	Status    int       `gorm:"type:smallint;default:1" json:"status"` // 1成功 0失败
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (SSHRecord) TableName() string { return "cmdb_ssh_records" }

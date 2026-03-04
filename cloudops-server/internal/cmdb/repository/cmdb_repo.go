package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Zara1024/AIOps/cloudops-server/internal/cmdb/model"
)

// HostRepository 主机数据访问层
type HostRepository struct {
	db *gorm.DB
}

func NewHostRepository(db *gorm.DB) *HostRepository {
	return &HostRepository{db: db}
}

// FindByID 根据ID查找主机
func (r *HostRepository) FindByID(ctx context.Context, id int64) (*model.Host, error) {
	var host model.Host
	err := r.db.WithContext(ctx).Preload("Group").First(&host, id).Error
	return &host, err
}

// FindByIP 根据IP查找主机
func (r *HostRepository) FindByIP(ctx context.Context, ip string) (*model.Host, error) {
	var host model.Host
	err := r.db.WithContext(ctx).Where("ip = ?", ip).First(&host).Error
	return &host, err
}

// Create 创建主机
func (r *HostRepository) Create(ctx context.Context, host *model.Host) error {
	return r.db.WithContext(ctx).Create(host).Error
}

// Update 更新主机
func (r *HostRepository) Update(ctx context.Context, host *model.Host) error {
	return r.db.WithContext(ctx).Save(host).Error
}

// UpdateFields 更新指定字段
func (r *HostRepository) UpdateFields(ctx context.Context, id int64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Host{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除主机
func (r *HostRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Host{}, id).Error
}

// BatchDelete 批量删除
func (r *HostRepository) BatchDelete(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Host{}).Error
}

// List 分页查询主机列表
func (r *HostRepository) List(ctx context.Context, page, pageSize int, ip, hostname string, status, groupID int) ([]model.Host, int64, error) {
	var hosts []model.Host
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Host{})

	if ip != "" {
		query = query.Where("ip LIKE ?", "%"+ip+"%")
	}
	if hostname != "" {
		query = query.Where("hostname LIKE ?", "%"+hostname+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if groupID > 0 {
		query = query.Where("group_id = ?", groupID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Group").
		Order("id DESC").
		Offset(offset).Limit(pageSize).
		Find(&hosts).Error

	return hosts, total, err
}

// CheckIPExist 检查IP是否存在
func (r *HostRepository) CheckIPExist(ctx context.Context, ip string, excludeID ...int64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Host{}).Where("ip = ?", ip)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// BatchUpdateGroup 批量更新分组
func (r *HostRepository) BatchUpdateGroup(ctx context.Context, ids []int64, groupID int64) error {
	return r.db.WithContext(ctx).Model(&model.Host{}).Where("id IN ?", ids).Update("group_id", groupID).Error
}

// HostGroupRepository 主机分组数据访问层
type HostGroupRepository struct {
	db *gorm.DB
}

func NewHostGroupRepository(db *gorm.DB) *HostGroupRepository {
	return &HostGroupRepository{db: db}
}

// FindAll 获取所有分组
func (r *HostGroupRepository) FindAll(ctx context.Context) ([]model.HostGroup, error) {
	var groups []model.HostGroup
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&groups).Error
	return groups, err
}

// FindByID 根据ID查找
func (r *HostGroupRepository) FindByID(ctx context.Context, id int64) (*model.HostGroup, error) {
	var group model.HostGroup
	err := r.db.WithContext(ctx).First(&group, id).Error
	return &group, err
}

// Create 创建分组
func (r *HostGroupRepository) Create(ctx context.Context, group *model.HostGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// Update 更新分组
func (r *HostGroupRepository) Update(ctx context.Context, group *model.HostGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

// Delete 删除分组
func (r *HostGroupRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.HostGroup{}, id).Error
}

// CountHostsByGroupID 统计分组下主机数量
func (r *HostGroupRepository) CountHostsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Host{}).Where("group_id = ?", groupID).Count(&count).Error
	return count, err
}

// SSHRecordRepository SSH记录仓库
type SSHRecordRepository struct {
	db *gorm.DB
}

func NewSSHRecordRepository(db *gorm.DB) *SSHRecordRepository {
	return &SSHRecordRepository{db: db}
}

// Create 创建SSH记录
func (r *SSHRecordRepository) Create(ctx context.Context, record *model.SSHRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// List 查询SSH操作记录
func (r *SSHRecordRepository) List(ctx context.Context, page, pageSize int, hostID int64, username string) ([]model.SSHRecord, int64, error) {
	var records []model.SSHRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.SSHRecord{})
	if hostID > 0 {
		query = query.Where("host_id = ?", hostID)
	}
	if username != "" {
		query = query.Where("username = ?", username)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	return records, total, err
}

package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Zara1024/AIOps/cloudops-server/internal/system/model"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Preload("Department").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateFields 更新指定字段
func (r *UserRepository) UpdateFields(ctx context.Context, id int64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 软删除用户
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// List 分页查询用户列表
func (r *UserRepository) List(ctx context.Context, page, pageSize int, username, phone string, status, deptID int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if deptID > 0 {
		query = query.Where("department_id = ?", deptID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("Roles").Preload("Department").
		Order("id DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// SetRoles 设置用户角色
func (r *UserRepository) SetRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	user := &model.User{ID: userID}

	// 先清除原有角色
	if err := r.db.WithContext(ctx).Model(user).Association("Roles").Clear(); err != nil {
		return err
	}

	if len(roleIDs) == 0 {
		return nil
	}

	// 设置新角色
	var roles []model.Role
	for _, id := range roleIDs {
		roles = append(roles, model.Role{ID: id})
	}
	return r.db.WithContext(ctx).Model(user).Association("Roles").Replace(roles)
}

// CheckUsernameExist 检查用户名是否存在
func (r *UserRepository) CheckUsernameExist(ctx context.Context, username string, excludeID ...int64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

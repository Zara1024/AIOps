package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Zara1024/AIOps/cloudops-server/internal/system/model"
)

// RoleRepository 角色数据访问层
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库实例
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// FindByID 根据ID查找角色
func (r *RoleRepository) FindByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Preload("Menus").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByKey 根据角色Key查找
func (r *RoleRepository) FindByKey(ctx context.Context, roleKey string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("role_key = ?", roleKey).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Create 创建角色
func (r *RoleRepository) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// List 分页查询角色列表
func (r *RoleRepository) List(ctx context.Context, page, pageSize int, roleName string, status int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Role{})

	if roleName != "" {
		query = query.Where("role_name LIKE ?", "%"+roleName+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("sort_order ASC, id ASC").
		Offset(offset).Limit(pageSize).
		Find(&roles).Error

	return roles, total, err
}

// SetMenus 设置角色权限（菜单+按钮）
func (r *RoleRepository) SetMenus(ctx context.Context, roleID int64, menuIDs []int64) error {
	role := &model.Role{ID: roleID}

	if err := r.db.WithContext(ctx).Model(role).Association("Menus").Clear(); err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}
	var menus []model.Menu
	for _, id := range menuIDs {
		menus = append(menus, model.Menu{ID: id})
	}
	return r.db.WithContext(ctx).Model(role).Association("Menus").Replace(menus)
}

// GetMenusByRoleKeys 根据角色Key列表获取菜单
func (r *RoleRepository) GetMenusByRoleKeys(ctx context.Context, roleKeys []string) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.WithContext(ctx).
		Distinct().
		Joins("JOIN sys_role_menus ON sys_role_menus.menu_id = sys_menus.id").
		Joins("JOIN sys_roles ON sys_roles.id = sys_role_menus.role_id").
		Where("sys_roles.role_key IN ? AND sys_roles.status = 1 AND sys_menus.status = 1", roleKeys).
		Order("sys_menus.sort_order ASC").
		Find(&menus).Error
	return menus, err
}

// MenuRepository 菜单数据访问层
type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// FindAll 获取所有菜单
func (r *MenuRepository) FindAll(ctx context.Context) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.WithContext(ctx).
		Where("status = 1").
		Order("sort_order ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// FindByID 根据ID查找菜单
func (r *MenuRepository) FindByID(ctx context.Context, id int64) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.WithContext(ctx).First(&menu, id).Error
	return &menu, err
}

// Create 创建菜单
func (r *MenuRepository) Create(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Menu{}, id).Error
}

// DepartmentRepository 部门数据访问层
type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// FindAll 获取所有部门
func (r *DepartmentRepository) FindAll(ctx context.Context) ([]model.Department, error) {
	var depts []model.Department
	err := r.db.WithContext(ctx).
		Where("status = 1").
		Order("sort_order ASC, id ASC").
		Find(&depts).Error
	return depts, err
}

// FindByID 根据ID查找部门
func (r *DepartmentRepository) FindByID(ctx context.Context, id int64) (*model.Department, error) {
	var dept model.Department
	err := r.db.WithContext(ctx).First(&dept, id).Error
	return &dept, err
}

// Create 创建部门
func (r *DepartmentRepository) Create(ctx context.Context, dept *model.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

// Update 更新部门
func (r *DepartmentRepository) Update(ctx context.Context, dept *model.Department) error {
	return r.db.WithContext(ctx).Save(dept).Error
}

// Delete 删除部门
func (r *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Department{}, id).Error
}

// LogRepository 日志数据访问层
type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db: db}
}

// CreateOperationLog 记录操作日志
func (r *LogRepository) CreateOperationLog(ctx context.Context, log *model.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// CreateLoginLog 记录登录日志
func (r *LogRepository) CreateLoginLog(ctx context.Context, log *model.LoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// ListOperationLogs 查询操作日志
func (r *LogRepository) ListOperationLogs(ctx context.Context, page, pageSize int, username, module string) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.OperationLog{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// ListLoginLogs 查询登录日志
func (r *LogRepository) ListLoginLogs(ctx context.Context, page, pageSize int, username string, status int) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LoginLog{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

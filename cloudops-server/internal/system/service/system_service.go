package service

import (
	"context"

	"github.com/Zara1024/AIOps/cloudops-server/internal/system/model"
	"github.com/Zara1024/AIOps/cloudops-server/internal/system/repository"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/crypto"
	appErrors "github.com/Zara1024/AIOps/cloudops-server/pkg/errors"
)

// SystemService 系统管理服务
type SystemService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
	menuRepo *repository.MenuRepository
	deptRepo *repository.DepartmentRepository
	logRepo  *repository.LogRepository
}

// NewSystemService 创建系统管理服务
func NewSystemService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	menuRepo *repository.MenuRepository,
	deptRepo *repository.DepartmentRepository,
	logRepo *repository.LogRepository,
) *SystemService {
	return &SystemService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		menuRepo: menuRepo,
		deptRepo: deptRepo,
		logRepo:  logRepo,
	}
}

// ============= 用户管理 =============

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username     string  `json:"username" binding:"required,min=2,max=64"`
	Password     string  `json:"password" binding:"required,min=8,max=128"`
	Nickname     string  `json:"nickname" binding:"max=128"`
	Email        string  `json:"email" binding:"omitempty,email"`
	Phone        string  `json:"phone" binding:"max=20"`
	DepartmentID *int64  `json:"department_id"`
	RoleIDs      []int64 `json:"role_ids"`
	Status       int     `json:"status"`
}

// CreateUser 创建用户
func (s *SystemService) CreateUser(ctx context.Context, req *CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已存在
	exist, err := s.userRepo.CheckUsernameExist(ctx, req.Username)
	if err != nil {
		return nil, appErrors.ErrInternal
	}
	if exist {
		return nil, appErrors.ErrDuplicate.WithMessage("用户名已存在")
	}

	// 哈希密码
	hash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Phone:        req.Phone,
		DepartmentID: req.DepartmentID,
		Status:       1,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, appErrors.ErrInternal.WithDetail(err.Error())
	}

	// 设置角色
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.SetRoles(ctx, user.ID, req.RoleIDs); err != nil {
			return nil, appErrors.ErrInternal.WithDetail("设置角色失败: " + err.Error())
		}
	}

	return user, nil
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname     string  `json:"nickname"`
	Email        string  `json:"email" binding:"omitempty,email"`
	Phone        string  `json:"phone"`
	DepartmentID *int64  `json:"department_id"`
	RoleIDs      []int64 `json:"role_ids"`
	Status       int     `json:"status"`
}

// UpdateUser 更新用户
func (s *SystemService) UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return appErrors.ErrNotFound.WithMessage("用户不存在")
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.DepartmentID = req.DepartmentID
	if req.Status > 0 {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return appErrors.ErrInternal
	}

	if req.RoleIDs != nil {
		if err := s.userRepo.SetRoles(ctx, id, req.RoleIDs); err != nil {
			return appErrors.ErrInternal
		}
	}

	return nil
}

// DeleteUser 删除用户
func (s *SystemService) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

// ListUsers 用户列表
func (s *SystemService) ListUsers(ctx context.Context, page, pageSize int, username, phone string, status, deptID int) ([]model.User, int64, error) {
	return s.userRepo.List(ctx, page, pageSize, username, phone, status, deptID)
}

// GetUser 获取用户详情
func (s *SystemService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// ============= 角色管理 =============

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleName    string  `json:"role_name" binding:"required,max=64"`
	RoleKey     string  `json:"role_key" binding:"required,max=64"`
	Description string  `json:"description"`
	SortOrder   int     `json:"sort_order"`
	MenuIDs     []int64 `json:"menu_ids"`
}

// CreateRole 创建角色
func (s *SystemService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*model.Role, error) {
	role := &model.Role{
		RoleName:    req.RoleName,
		RoleKey:     req.RoleKey,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      1,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, appErrors.ErrInternal.WithDetail(err.Error())
	}

	if len(req.MenuIDs) > 0 {
		if err := s.roleRepo.SetMenus(ctx, role.ID, req.MenuIDs); err != nil {
			return nil, appErrors.ErrInternal
		}
	}

	return role, nil
}

// UpdateRole 更新角色
func (s *SystemService) UpdateRole(ctx context.Context, id int64, req *CreateRoleRequest) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return appErrors.ErrNotFound
	}

	role.RoleName = req.RoleName
	role.RoleKey = req.RoleKey
	role.Description = req.Description
	role.SortOrder = req.SortOrder

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return appErrors.ErrInternal
	}

	if req.MenuIDs != nil {
		return s.roleRepo.SetMenus(ctx, id, req.MenuIDs)
	}

	return nil
}

// DeleteRole 删除角色
func (s *SystemService) DeleteRole(ctx context.Context, id int64) error {
	return s.roleRepo.Delete(ctx, id)
}

// ListRoles 角色列表
func (s *SystemService) ListRoles(ctx context.Context, page, pageSize int, roleName string, status int) ([]model.Role, int64, error) {
	return s.roleRepo.List(ctx, page, pageSize, roleName, status)
}

// ============= 菜单管理 =============

// GetMenuTree 获取菜单树
func (s *SystemService) GetMenuTree(ctx context.Context) ([]*model.Menu, error) {
	menus, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

// GetMenusByRoleKeys 根据角色获取菜单（用于动态路由）
func (s *SystemService) GetMenusByRoleKeys(ctx context.Context, roleKeys []string) ([]*model.Menu, error) {
	// 超级管理员获取所有菜单
	for _, key := range roleKeys {
		if key == "super_admin" {
			return s.GetMenuTree(ctx)
		}
	}

	menus, err := s.roleRepo.GetMenusByRoleKeys(ctx, roleKeys)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

// CreateMenu 创建菜单
func (s *SystemService) CreateMenu(ctx context.Context, menu *model.Menu) error {
	return s.menuRepo.Create(ctx, menu)
}

// UpdateMenu 更新菜单
func (s *SystemService) UpdateMenu(ctx context.Context, menu *model.Menu) error {
	return s.menuRepo.Update(ctx, menu)
}

// DeleteMenu 删除菜单
func (s *SystemService) DeleteMenu(ctx context.Context, id int64) error {
	return s.menuRepo.Delete(ctx, id)
}

// ============= 部门管理 =============

// GetDeptTree 获取部门树
func (s *SystemService) GetDeptTree(ctx context.Context) ([]*model.Department, error) {
	depts, err := s.deptRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildDeptTree(depts, 0), nil
}

// CreateDept 创建部门
func (s *SystemService) CreateDept(ctx context.Context, dept *model.Department) error {
	return s.deptRepo.Create(ctx, dept)
}

// UpdateDept 更新部门
func (s *SystemService) UpdateDept(ctx context.Context, dept *model.Department) error {
	return s.deptRepo.Update(ctx, dept)
}

// DeleteDept 删除部门
func (s *SystemService) DeleteDept(ctx context.Context, id int64) error {
	return s.deptRepo.Delete(ctx, id)
}

// ============= 树形构建工具 =============

func buildMenuTree(menus []model.Menu, parentID int64) []*model.Menu {
	var tree []*model.Menu
	for i := range menus {
		if menus[i].ParentID == parentID {
			node := &menus[i]
			node.Children = buildMenuTree(menus, node.ID)
			tree = append(tree, node)
		}
	}
	return tree
}

func buildDeptTree(depts []model.Department, parentID int64) []*model.Department {
	var tree []*model.Department
	for i := range depts {
		if depts[i].ParentID == parentID {
			node := &depts[i]
			node.Children = buildDeptTree(depts, node.ID)
			tree = append(tree, node)
		}
	}
	return tree
}

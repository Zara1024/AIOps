package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/AIOps/cloudops-server/internal/system/model"
	"github.com/Zara1024/AIOps/cloudops-server/internal/system/service"
	appErrors "github.com/Zara1024/AIOps/cloudops-server/pkg/errors"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/middleware"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/response"
)

// SystemHandler 系统管理处理器
type SystemHandler struct {
	systemService *service.SystemService
}

// NewSystemHandler 创建系统管理处理器
func NewSystemHandler(systemService *service.SystemService) *SystemHandler {
	return &SystemHandler{systemService: systemService}
}

// ============= 用户管理 =============

// ListUsers 用户列表
func (h *SystemHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	username := c.Query("username")
	phone := c.Query("phone")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	deptID, _ := strconv.Atoi(c.DefaultQuery("department_id", "0"))

	if pageSize > 100 {
		pageSize = 100
	}

	users, total, err := h.systemService.ListUsers(c.Request.Context(), page, pageSize, username, phone, status, deptID)
	if err != nil {
		response.ServerError(c, "查询用户列表失败")
		return
	}

	response.OKWithPage(c, users, total, page, pageSize)
}

// GetUser 用户详情
func (h *SystemHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.systemService.GetUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, appErrors.ErrNotFound)
		return
	}
	response.OK(c, user)
}

// CreateUser 创建用户
func (h *SystemHandler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.systemService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
		} else {
			response.ServerError(c, "创建用户失败")
		}
		return
	}
	response.OKWithMessage(c, "创建成功", user)
}

// UpdateUser 更新用户
func (h *SystemHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.systemService.UpdateUser(c.Request.Context(), id, &req); err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
		} else {
			response.ServerError(c, "更新用户失败")
		}
		return
	}
	response.OKWithMessage(c, "更新成功", nil)
}

// DeleteUser 删除用户
func (h *SystemHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.systemService.DeleteUser(c.Request.Context(), id); err != nil {
		response.ServerError(c, "删除用户失败")
		return
	}
	response.OKWithMessage(c, "删除成功", nil)
}

// ============= 角色管理 =============

// ListRoles 角色列表
func (h *SystemHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	roleName := c.Query("role_name")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	roles, total, err := h.systemService.ListRoles(c.Request.Context(), page, pageSize, roleName, status)
	if err != nil {
		response.ServerError(c, "查询角色列表失败")
		return
	}
	response.OKWithPage(c, roles, total, page, pageSize)
}

// CreateRole 创建角色
func (h *SystemHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.systemService.CreateRole(c.Request.Context(), &req)
	if err != nil {
		response.ServerError(c, "创建角色失败")
		return
	}
	response.OKWithMessage(c, "创建成功", role)
}

// UpdateRole 更新角色
func (h *SystemHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.systemService.UpdateRole(c.Request.Context(), id, &req); err != nil {
		response.ServerError(c, "更新角色失败")
		return
	}
	response.OKWithMessage(c, "更新成功", nil)
}

// DeleteRole 删除角色
func (h *SystemHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}
	if err := h.systemService.DeleteRole(c.Request.Context(), id); err != nil {
		response.ServerError(c, "删除角色失败")
		return
	}
	response.OKWithMessage(c, "删除成功", nil)
}

// ============= 菜单管理 =============

// GetMenuTree 菜单树
func (h *SystemHandler) GetMenuTree(c *gin.Context) {
	menus, err := h.systemService.GetMenuTree(c.Request.Context())
	if err != nil {
		response.ServerError(c, "获取菜单树失败")
		return
	}
	response.OK(c, menus)
}

// GetUserMenus 获取当前用户菜单（动态路由）
func (h *SystemHandler) GetUserMenus(c *gin.Context) {
	roleKeys := middleware.GetRoleKeys(c)
	menus, err := h.systemService.GetMenusByRoleKeys(c.Request.Context(), roleKeys)
	if err != nil {
		response.ServerError(c, "获取菜单失败")
		return
	}
	response.OK(c, menus)
}

// CreateMenu 创建菜单
func (h *SystemHandler) CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.systemService.CreateMenu(c.Request.Context(), &menu); err != nil {
		response.ServerError(c, "创建菜单失败")
		return
	}
	response.OKWithMessage(c, "创建成功", menu)
}

// UpdateMenu 更新菜单
func (h *SystemHandler) UpdateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	menu.ID = id
	if err := h.systemService.UpdateMenu(c.Request.Context(), &menu); err != nil {
		response.ServerError(c, "更新菜单失败")
		return
	}
	response.OKWithMessage(c, "更新成功", nil)
}

// DeleteMenu 删除菜单
func (h *SystemHandler) DeleteMenu(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.systemService.DeleteMenu(c.Request.Context(), id); err != nil {
		response.ServerError(c, "删除菜单失败")
		return
	}
	response.OKWithMessage(c, "删除成功", nil)
}

// ============= 部门管理 =============

// GetDeptTree 部门树
func (h *SystemHandler) GetDeptTree(c *gin.Context) {
	depts, err := h.systemService.GetDeptTree(c.Request.Context())
	if err != nil {
		response.ServerError(c, "获取部门树失败")
		return
	}
	response.OK(c, depts)
}

// CreateDept 创建部门
func (h *SystemHandler) CreateDept(c *gin.Context) {
	var dept model.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.systemService.CreateDept(c.Request.Context(), &dept); err != nil {
		response.ServerError(c, "创建部门失败")
		return
	}
	response.OKWithMessage(c, "创建成功", dept)
}

// UpdateDept 更新部门
func (h *SystemHandler) UpdateDept(c *gin.Context) {
	var dept model.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	dept.ID = id
	if err := h.systemService.UpdateDept(c.Request.Context(), &dept); err != nil {
		response.ServerError(c, "更新部门失败")
		return
	}
	response.OKWithMessage(c, "更新成功", nil)
}

// DeleteDept 删除部门
func (h *SystemHandler) DeleteDept(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.systemService.DeleteDept(c.Request.Context(), id); err != nil {
		response.ServerError(c, "删除部门失败")
		return
	}
	response.OKWithMessage(c, "删除成功", nil)
}

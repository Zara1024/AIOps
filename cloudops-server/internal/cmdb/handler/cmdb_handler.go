package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Zara1024/AIOps/cloudops-server/internal/cmdb/model"
	"github.com/Zara1024/AIOps/cloudops-server/internal/cmdb/service"
	appErrors "github.com/Zara1024/AIOps/cloudops-server/pkg/errors"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/middleware"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/response"
)

// CMDBHandler CMDB 处理器
type CMDBHandler struct {
	cmdbService *service.CMDBService
}

func NewCMDBHandler(cmdbService *service.CMDBService) *CMDBHandler {
	return &CMDBHandler{cmdbService: cmdbService}
}

// ============= 主机管理 =============

// ListHosts 主机列表
func (h *CMDBHandler) ListHosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	ip := c.Query("ip")
	hostname := c.Query("hostname")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	groupID, _ := strconv.Atoi(c.DefaultQuery("group_id", "0"))

	if pageSize > 100 {
		pageSize = 100
	}

	hosts, total, err := h.cmdbService.ListHosts(c.Request.Context(), page, pageSize, ip, hostname, status, groupID)
	if err != nil {
		response.ServerError(c, "查询主机列表失败")
		return
	}
	response.OKWithPage(c, hosts, total, page, pageSize)
}

// GetHost 主机详情
func (h *CMDBHandler) GetHost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的主机ID")
		return
	}

	host, err := h.cmdbService.GetHost(c.Request.Context(), id)
	if err != nil {
		response.Error(c, appErrors.ErrHostNotFound)
		return
	}
	response.OK(c, host)
}

// CreateHost 创建主机
func (h *CMDBHandler) CreateHost(c *gin.Context) {
	var req service.CreateHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	host, err := h.cmdbService.CreateHost(c.Request.Context(), &req, userID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
		} else {
			response.ServerError(c, "创建主机失败")
		}
		return
	}
	response.OKWithMessage(c, "创建成功", host)
}

// UpdateHost 更新主机
func (h *CMDBHandler) UpdateHost(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req service.UpdateHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.cmdbService.UpdateHost(c.Request.Context(), id, &req); err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
		} else {
			response.ServerError(c, "更新主机失败")
		}
		return
	}
	response.OKWithMessage(c, "更新成功", nil)
}

// DeleteHost 删除主机
func (h *CMDBHandler) DeleteHost(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.cmdbService.DeleteHost(c.Request.Context(), id); err != nil {
		response.ServerError(c, "删除主机失败")
		return
	}
	response.OKWithMessage(c, "删除成功", nil)
}

// BatchDeleteHosts 批量删除主机
func (h *CMDBHandler) BatchDeleteHosts(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.cmdbService.BatchDeleteHosts(c.Request.Context(), req.IDs); err != nil {
		response.ServerError(c, "批量删除失败")
		return
	}
	response.OKWithMessage(c, "批量删除成功", nil)
}

// BatchUpdateGroup 批量分组
func (h *CMDBHandler) BatchUpdateGroup(c *gin.Context) {
	var req struct {
		IDs     []int64 `json:"ids" binding:"required"`
		GroupID int64   `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.cmdbService.BatchUpdateGroup(c.Request.Context(), req.IDs, req.GroupID); err != nil {
		response.ServerError(c, "批量分组失败")
		return
	}
	response.OKWithMessage(c, "批量分组成功", nil)
}

// ============= 主机分组 =============

// GetGroupTree 分组树
func (h *CMDBHandler) GetGroupTree(c *gin.Context) {
	groups, err := h.cmdbService.GetGroupTree(c.Request.Context())
	if err != nil {
		response.ServerError(c, "获取分组树失败")
		return
	}
	response.OK(c, groups)
}

// CreateGroup 创建分组
func (h *CMDBHandler) CreateGroup(c *gin.Context) {
	var group model.HostGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.cmdbService.CreateGroup(c.Request.Context(), &group); err != nil {
		response.ServerError(c, "创建分组失败")
		return
	}
	response.OKWithMessage(c, "创建成功", group)
}

// UpdateGroup 更新分组
func (h *CMDBHandler) UpdateGroup(c *gin.Context) {
	var group model.HostGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	group.ID = id
	if err := h.cmdbService.UpdateGroup(c.Request.Context(), &group); err != nil {
		response.ServerError(c, "更新分组失败")
		return
	}
	response.OKWithMessage(c, "更新成功", nil)
}

// DeleteGroup 删除分组
func (h *CMDBHandler) DeleteGroup(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.cmdbService.DeleteGroup(c.Request.Context(), id); err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
		} else {
			response.ServerError(c, "删除分组失败")
		}
		return
	}
	response.OKWithMessage(c, "删除成功", nil)
}

// ============= SSH 记录 =============

// ListSSHRecords SSH操作记录
func (h *CMDBHandler) ListSSHRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	hostID, _ := strconv.ParseInt(c.DefaultQuery("host_id", "0"), 10, 64)
	username := c.Query("username")

	records, total, err := h.cmdbService.ListSSHRecords(c.Request.Context(), page, pageSize, hostID, username)
	if err != nil {
		response.ServerError(c, "查询SSH记录失败")
		return
	}
	response.OKWithPage(c, records, total, page, pageSize)
}

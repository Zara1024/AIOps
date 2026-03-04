package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Zara1024/AIOps/cloudops-server/internal/auth/service"
	appErrors "github.com/Zara1024/AIOps/cloudops-server/pkg/errors"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/middleware"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param body body service.LoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")

	result, err := h.authService.Login(c.Request.Context(), &req, ip, ua)
	if err != nil {
		response.Error(c, err.(*appErrors.AppError))
		return
	}

	response.OK(c, result)
}

// Logout 用户登出
// @Summary 用户登出
// @Tags 认证管理
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: 将 Token 加入黑名单（Redis）
	response.OKWithMessage(c, "登出成功", nil)
}

// RefreshToken 刷新 Token
// @Summary 刷新 Token
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param body body refreshTokenRequest true "刷新Token请求"
// @Success 200 {object} response.Response{data=jwt.TokenPair}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err.(*appErrors.AppError))
		return
	}

	response.OK(c, result)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证管理
// @Security Bearer
// @Success 200 {object} response.Response{data=service.UserInfo}
// @Router /api/v1/auth/userinfo [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	info, err := h.authService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err.(*appErrors.AppError))
		return
	}
	response.OK(c, info)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags 认证管理
// @Security Bearer
// @Accept json
// @Param body body changePasswordRequest true "修改密码请求"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, err.(*appErrors.AppError))
		return
	}

	response.OKWithMessage(c, "密码修改成功", nil)
}

// 请求结构体
type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Zara1024/AIOps/cloudops-server/internal/system/model"
	"github.com/Zara1024/AIOps/cloudops-server/internal/system/repository"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/config"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/crypto"
	appErrors "github.com/Zara1024/AIOps/cloudops-server/pkg/errors"
	appJwt "github.com/Zara1024/AIOps/cloudops-server/pkg/jwt"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	logRepo  *repository.LogRepository
	jwtCfg   *config.JWTConfig
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, logRepo *repository.LogRepository, jwtCfg *config.JWTConfig) *AuthService {
	return &AuthService{userRepo: userRepo, logRepo: logRepo, jwtCfg: jwtCfg}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	*appJwt.TokenPair
	User UserInfo `json:"user"`
}

// UserInfo 用户信息（脱敏）
type UserInfo struct {
	ID           int64    `json:"id"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	Avatar       string   `json:"avatar"`
	DepartmentID *int64   `json:"department_id"`
	RoleKeys     []string `json:"role_keys"`
	Roles        []string `json:"roles"`
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest, ip, ua string) (*LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		// 记录登录失败日志
		s.recordLoginLog(ctx, 0, req.Username, ip, ua, 0, "用户名或密码错误")
		return nil, appErrors.ErrLoginFailed
	}

	// 检查账户状态
	if user.Status == 0 {
		s.recordLoginLog(ctx, user.ID, user.Username, ip, ua, 0, "账户已禁用")
		return nil, appErrors.ErrAccountDisabled
	}

	// 检查账户是否锁定
	if user.LockUntil != nil && time.Now().Before(*user.LockUntil) {
		s.recordLoginLog(ctx, user.ID, user.Username, ip, ua, 0, "账户已锁定")
		return nil, appErrors.ErrAccountLocked
	}

	// 验证密码
	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		// 增加失败计数
		user.LoginFailCount++
		if user.LoginFailCount >= 5 {
			lockUntil := time.Now().Add(15 * time.Minute)
			user.LockUntil = &lockUntil
			user.LoginFailCount = 0
		}
		_ = s.userRepo.Update(ctx, user)
		s.recordLoginLog(ctx, user.ID, user.Username, ip, ua, 0, "密码错误")
		return nil, appErrors.ErrLoginFailed
	}

	// 登录成功，重置失败计数
	now := time.Now()
	user.LoginFailCount = 0
	user.LockUntil = nil
	user.LastLoginAt = &now
	user.LastLoginIP = ip
	_ = s.userRepo.Update(ctx, user)

	// 获取角色Key
	var roleKeys []string
	var roleNames []string
	for _, role := range user.Roles {
		roleKeys = append(roleKeys, role.RoleKey)
		roleNames = append(roleNames, role.RoleName)
	}

	// 生成 Token
	tokenPair, err := appJwt.GenerateTokenPair(s.jwtCfg, user.ID, user.Username, roleKeys)
	if err != nil {
		return nil, appErrors.ErrInternal.WithDetail(fmt.Sprintf("生成Token失败: %v", err))
	}

	// 记录登录成功日志
	s.recordLoginLog(ctx, user.ID, user.Username, ip, ua, 1, "登录成功")

	return &LoginResponse{
		TokenPair: tokenPair,
		User: UserInfo{
			ID:           user.ID,
			Username:     user.Username,
			Nickname:     user.Nickname,
			Email:        user.Email,
			Phone:        user.Phone,
			Avatar:       user.Avatar,
			DepartmentID: user.DepartmentID,
			RoleKeys:     roleKeys,
			Roles:        roleNames,
		},
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*appJwt.TokenPair, error) {
	// 解析 RefreshToken
	claims, err := appJwt.ParseToken(s.jwtCfg, refreshTokenStr)
	if err != nil {
		return nil, appErrors.ErrTokenExpired
	}

	if claims.TokenType != appJwt.RefreshToken {
		return nil, appErrors.ErrUnauthorized.WithMessage("Token 类型错误")
	}

	// 查找用户并获取最新角色信息
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, appErrors.ErrUnauthorized
	}

	var roleKeys []string
	for _, role := range user.Roles {
		roleKeys = append(roleKeys, role.RoleKey)
	}

	// 生成新的 TokenPair
	return appJwt.GenerateTokenPair(s.jwtCfg, user.ID, user.Username, roleKeys)
}

// GetUserInfo 获取当前用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	var roleKeys []string
	var roleNames []string
	for _, role := range user.Roles {
		roleKeys = append(roleKeys, role.RoleKey)
		roleNames = append(roleNames, role.RoleName)
	}

	return &UserInfo{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		DepartmentID: user.DepartmentID,
		RoleKeys:     roleKeys,
		Roles:        roleNames,
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, oldPwd, newPwd string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return appErrors.ErrNotFound
	}

	// 验证旧密码
	if !crypto.CheckPassword(oldPwd, user.PasswordHash) {
		return appErrors.ErrOldPassword
	}

	// 哈希新密码
	hash, err := crypto.HashPassword(newPwd)
	if err != nil {
		return appErrors.ErrInternal
	}

	return s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"password_hash": hash,
	})
}

// recordLoginLog 记录登录日志
func (s *AuthService) recordLoginLog(ctx context.Context, userID int64, username, ip, ua string, status int, msg string) {
	log := &model.LoginLog{
		UserID:   userID,
		Username: username,
		IP:       ip,
		Browser:  ua,
		Status:   status,
		Message:  msg,
	}
	if err := s.logRepo.CreateLoginLog(ctx, log); err != nil {
		slog.Error("记录登录日志失败", "error", err)
	}
}

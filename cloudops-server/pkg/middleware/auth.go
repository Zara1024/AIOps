package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	appJwt "github.com/Zara1024/AIOps/cloudops-server/pkg/jwt"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/config"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/response"
)

// 上下文中存储用户信息的 key
const (
	CtxUserIDKey   = "user_id"
	CtxUsernameKey = "username"
	CtxRoleKeysKey = "role_keys"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误，应为 Bearer {token}")
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := appJwt.ParseToken(cfg, parts[1])
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 验证 Token 类型
		if claims.TokenType != appJwt.AccessToken {
			response.Unauthorized(c, "Token 类型错误")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxUsernameKey, claims.Username)
		c.Set(CtxRoleKeysKey, claims.RoleKeys)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get(CtxUserIDKey)
	if id, ok := userID.(int64); ok {
		return id
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	username, _ := c.Get(CtxUsernameKey)
	if name, ok := username.(string); ok {
		return name
	}
	return ""
}

// GetRoleKeys 从上下文获取角色Key列表
func GetRoleKeys(c *gin.Context) []string {
	roleKeys, _ := c.Get(CtxRoleKeysKey)
	if keys, ok := roleKeys.([]string); ok {
		return keys
	}
	return nil
}

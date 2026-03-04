package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	authHandler "github.com/Zara1024/AIOps/cloudops-server/internal/auth/handler"
	authService "github.com/Zara1024/AIOps/cloudops-server/internal/auth/service"
	systemHandler "github.com/Zara1024/AIOps/cloudops-server/internal/system/handler"
	"github.com/Zara1024/AIOps/cloudops-server/internal/system/repository"
	systemService "github.com/Zara1024/AIOps/cloudops-server/internal/system/service"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/config"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/middleware"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/response"
)

// InitRouter 初始化路由
func InitRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// 设置 gin 模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.CORSMiddleware(&cfg.CORS))

	// 初始化 Repository
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	logRepo := repository.NewLogRepository(db)

	// 初始化 Service
	authSvc := authService.NewAuthService(userRepo, logRepo, &cfg.JWT)
	sysSvc := systemService.NewSystemService(userRepo, roleRepo, menuRepo, deptRepo, logRepo)

	// 初始化 Handler
	authH := authHandler.NewAuthHandler(authSvc)
	sysH := systemHandler.NewSystemHandler(sysSvc)

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 健康检查（无需认证）
		api.GET("/health", func(c *gin.Context) {
			response.OK(c, gin.H{
				"status":  "ok",
				"service": "cloudops-hub",
			})
		})

		// 认证路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authH.Login)
			auth.POST("/refresh", authH.RefreshToken)
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(&cfg.JWT))
		{
			// 认证相关
			authorized.POST("/auth/logout", authH.Logout)
			authorized.GET("/auth/userinfo", authH.GetUserInfo)
			authorized.PUT("/auth/password", authH.ChangePassword)

			// 用户菜单（动态路由）
			authorized.GET("/system/menus/user", sysH.GetUserMenus)

			// 系统管理
			system := authorized.Group("/system")
			{
				// 用户管理
				system.GET("/users", sysH.ListUsers)
				system.GET("/users/:id", sysH.GetUser)
				system.POST("/users", sysH.CreateUser)
				system.PUT("/users/:id", sysH.UpdateUser)
				system.DELETE("/users/:id", sysH.DeleteUser)

				// 角色管理
				system.GET("/roles", sysH.ListRoles)
				system.POST("/roles", sysH.CreateRole)
				system.PUT("/roles/:id", sysH.UpdateRole)
				system.DELETE("/roles/:id", sysH.DeleteRole)

				// 菜单管理
				system.GET("/menus", sysH.GetMenuTree)
				system.POST("/menus", sysH.CreateMenu)
				system.PUT("/menus/:id", sysH.UpdateMenu)
				system.DELETE("/menus/:id", sysH.DeleteMenu)

				// 部门管理
				system.GET("/departments", sysH.GetDeptTree)
				system.POST("/departments", sysH.CreateDept)
				system.PUT("/departments/:id", sysH.UpdateDept)
				system.DELETE("/departments/:id", sysH.DeleteDept)
			}
		}
	}

	return r
}

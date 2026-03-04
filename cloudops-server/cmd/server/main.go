package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	"gorm.io/gorm"

	"github.com/Zara1024/AIOps/cloudops-server/internal/router"
	"github.com/Zara1024/AIOps/cloudops-server/internal/system/model"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/config"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/crypto"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/database"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/logger"
	rds "github.com/Zara1024/AIOps/cloudops-server/pkg/redis"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 初始化配置
	cfg, err := config.InitConfig(*configPath)
	if err != nil {
		slog.Error("初始化配置失败", "error", err)
		os.Exit(1)
	}
	slog.Info("配置加载成功")

	// 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		slog.Error("初始化日志失败", "error", err)
		os.Exit(1)
	}
	slog.Info("日志初始化成功")

	// 初始化数据库
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		slog.Error("初始化数据库失败", "error", err)
		os.Exit(1)
	}
	defer database.CloseDB()

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Department{},
		&model.OperationLog{},
		&model.LoginLog{},
	); err != nil {
		slog.Error("数据库迁移失败", "error", err)
		os.Exit(1)
	}
	slog.Info("数据库迁移完成")

	// 初始化种子数据
	initSeedData(db)

	// 初始化 Redis
	_, err = rds.InitRedis(&cfg.Redis)
	if err != nil {
		slog.Warn("Redis 连接失败，部分功能可能受限", "error", err)
		// Redis 不是必须的，可以降级运行
	}
	defer rds.CloseRedis()

	// 初始化路由
	r := router.InitRouter(cfg, db)

	// 启动 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 在 goroutine 中启动服务器
	go func() {
		slog.Info("智维云枢服务启动成功",
			"port", cfg.Server.Port,
			"mode", cfg.Server.Mode,
			"url", fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("服务器启动失败", "error", err)
			os.Exit(1)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("正在关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("服务器关闭失败", "error", err)
	}
	slog.Info("服务器已关闭")
}

// initSeedData 初始化种子数据
func initSeedData(db *gorm.DB) {
	// 检查是否已有管理员用户
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	slog.Info("初始化种子数据...")

	// 创建默认角色
	roles := []model.Role{
		{RoleName: "超级管理员", RoleKey: "super_admin", Description: "拥有所有权限", SortOrder: 1, Status: 1},
		{RoleName: "管理员", RoleKey: "admin", Description: "管理权限", SortOrder: 2, Status: 1},
		{RoleName: "只读用户", RoleKey: "viewer", Description: "只读权限", SortOrder: 3, Status: 1},
	}
	for i := range roles {
		db.Create(&roles[i])
	}

	// 创建默认菜单
	menus := []model.Menu{
		// 仪表盘
		{ParentID: 0, MenuName: "仪表盘", MenuType: 2, Path: "/dashboard", Component: "views/dashboard/index", Icon: "Monitor", Permission: "", SortOrder: 1, Visible: true, Status: 1},
		// 系统管理 (目录)
		{ParentID: 0, MenuName: "系统管理", MenuType: 1, Path: "/system", Component: "", Icon: "Setting", Permission: "", SortOrder: 100, Visible: true, Status: 1},
	}
	for i := range menus {
		db.Create(&menus[i])
	}

	// 系统管理子菜单
	parentID := menus[1].ID
	subMenus := []model.Menu{
		{ParentID: parentID, MenuName: "用户管理", MenuType: 2, Path: "/system/users", Component: "views/system/users/index", Icon: "User", Permission: "system:user:list", SortOrder: 1, Visible: true, Status: 1},
		{ParentID: parentID, MenuName: "角色管理", MenuType: 2, Path: "/system/roles", Component: "views/system/roles/index", Icon: "UserFilled", Permission: "system:role:list", SortOrder: 2, Visible: true, Status: 1},
		{ParentID: parentID, MenuName: "菜单管理", MenuType: 2, Path: "/system/menus", Component: "views/system/menus/index", Icon: "Menu", Permission: "system:menu:list", SortOrder: 3, Visible: true, Status: 1},
		{ParentID: parentID, MenuName: "部门管理", MenuType: 2, Path: "/system/departments", Component: "views/system/departments/index", Icon: "OfficeBuilding", Permission: "system:dept:list", SortOrder: 4, Visible: true, Status: 1},
	}
	for i := range subMenus {
		db.Create(&subMenus[i])
	}

	// 为仪表盘和系统管理菜单创建按钮权限（示例）
	for _, sm := range subMenus {
		// 添加按钮权限
		module := ""
		switch sm.MenuName {
		case "用户管理":
			module = "system:user"
		case "角色管理":
			module = "system:role"
		case "菜单管理":
			module = "system:menu"
		case "部门管理":
			module = "system:dept"
		}
		buttons := []model.Menu{
			{ParentID: sm.ID, MenuName: "新增", MenuType: 3, Permission: module + ":create", SortOrder: 1, Status: 1},
			{ParentID: sm.ID, MenuName: "编辑", MenuType: 3, Permission: module + ":update", SortOrder: 2, Status: 1},
			{ParentID: sm.ID, MenuName: "删除", MenuType: 3, Permission: module + ":delete", SortOrder: 3, Status: 1},
		}
		for j := range buttons {
			db.Create(&buttons[j])
		}
	}

	// 为超级管理员角色分配所有菜单
	var allMenus []model.Menu
	db.Find(&allMenus)
	superAdminRole := roles[0]
	db.Model(&superAdminRole).Association("Menus").Replace(allMenus)

	// 创建默认管理员用户
	hash, _ := crypto.HashPassword("Admin@2026")
	admin := &model.User{
		Username:     "admin",
		PasswordHash: hash,
		Nickname:     "超级管理员",
		Email:        "admin@cloudops.local",
		Status:       1,
	}
	db.Create(admin)

	// 为管理员用户分配超级管理员角色
	db.Model(admin).Association("Roles").Replace([]model.Role{superAdminRole})

	slog.Info("种子数据初始化完成",
		"admin_user", "admin",
		"admin_password", "Admin@2026",
	)
}

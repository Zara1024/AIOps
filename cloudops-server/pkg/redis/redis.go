package rds

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"

	"github.com/Zara1024/AIOps/cloudops-server/pkg/config"
)

// Client 全局 Redis 客户端
var Client *redis.Client

// InitRedis 初始化 Redis 客户端
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	Client = client
	slog.Info("Redis 连接成功", "addr", cfg.Addr())

	return client, nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

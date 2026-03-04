package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Zara1024/AIOps/cloudops-server/pkg/config"
)

// 上下文 Key 类型
type contextKey string

const (
	// RequestIDKey 请求ID在上下文中的Key
	RequestIDKey contextKey = "request_id"
)

// InitLogger 初始化全局日志
func InitLogger(cfg *config.LogConfig) error {
	var writer io.Writer
	switch cfg.Output {
	case "file":
		// 确保日志目录存在
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		writer = file
	default:
		writer = os.Stdout
	}

	// 解析日志级别
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.Level == "debug",
	}

	var handler slog.Handler
	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewJSONHandler(writer, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// WithRequestID 在日志中添加请求ID
func WithRequestID(ctx context.Context) *slog.Logger {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return slog.With("request_id", reqID)
	}
	return slog.Default()
}

// Info 带上下文的Info日志
func Info(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Info(msg, args...)
}

// Error 带上下文的Error日志
func Error(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Error(msg, args...)
}

// Warn 带上下文的Warn日志
func Warn(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Warn(msg, args...)
}

// Debug 带上下文的Debug日志
func Debug(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Debug(msg, args...)
}

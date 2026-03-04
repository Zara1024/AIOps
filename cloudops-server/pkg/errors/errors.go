package errors

import "fmt"

// AppError 统一应用错误结构
type AppError struct {
	Code     int    `json:"code"`     // 业务错误码
	HTTPCode int    `json:"-"`        // HTTP 状态码
	Message  string `json:"message"`  // 用户可见消息（中文）
	Detail   string `json:"-"`        // 开发者调试信息（仅日志）
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// WithDetail 返回带详细信息的新错误（不修改原始错误）
func (e *AppError) WithDetail(detail string) *AppError {
	return &AppError{
		Code:     e.Code,
		HTTPCode: e.HTTPCode,
		Message:  e.Message,
		Detail:   detail,
	}
}

// WithMessage 返回自定义消息的新错误
func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Code:     e.Code,
		HTTPCode: e.HTTPCode,
		Message:  msg,
		Detail:   e.Detail,
	}
}

// ============= 认证相关错误 (1xxxx) =============

var (
	ErrUnauthorized    = &AppError{Code: 10001, HTTPCode: 401, Message: "未登录或登录已过期"}
	ErrTokenExpired    = &AppError{Code: 10002, HTTPCode: 401, Message: "Token 已过期"}
	ErrForbidden       = &AppError{Code: 10003, HTTPCode: 403, Message: "无权限访问"}
	ErrNotFound        = &AppError{Code: 10004, HTTPCode: 404, Message: "资源不存在"}
	ErrValidation      = &AppError{Code: 10005, HTTPCode: 400, Message: "参数校验失败"}
	ErrInternal        = &AppError{Code: 10006, HTTPCode: 500, Message: "服务器内部错误"}
	ErrBadRequest      = &AppError{Code: 10007, HTTPCode: 400, Message: "请求参数错误"}
	ErrDuplicate       = &AppError{Code: 10008, HTTPCode: 409, Message: "数据已存在"}
	ErrLoginFailed     = &AppError{Code: 10009, HTTPCode: 401, Message: "用户名或密码错误"}
	ErrAccountLocked   = &AppError{Code: 10010, HTTPCode: 403, Message: "账户已锁定，请15分钟后重试"}
	ErrAccountDisabled = &AppError{Code: 10011, HTTPCode: 403, Message: "账户已禁用"}
	ErrCaptchaInvalid  = &AppError{Code: 10012, HTTPCode: 400, Message: "验证码错误或已过期"}
	ErrPasswordWeak    = &AppError{Code: 10013, HTTPCode: 400, Message: "密码强度不够，至少8位，含大小写+数字+特殊字符"}
	ErrOldPassword     = &AppError{Code: 10014, HTTPCode: 400, Message: "原密码错误"}
)

// ============= CMDB 相关错误 (2xxxx) =============

var (
	ErrHostNotFound    = &AppError{Code: 20001, HTTPCode: 404, Message: "主机不存在"}
	ErrHostUnreachable = &AppError{Code: 20002, HTTPCode: 503, Message: "主机不可达"}
	ErrSSHConnFailed   = &AppError{Code: 20003, HTTPCode: 503, Message: "SSH 连接失败"}
)

// ============= K8s 相关错误 (3xxxx) =============

var (
	ErrClusterNotFound   = &AppError{Code: 30001, HTTPCode: 404, Message: "集群不存在"}
	ErrClusterConnFailed = &AppError{Code: 30002, HTTPCode: 503, Message: "集群连接失败"}
)

// ============= AI/LLM 相关错误 (6xxxx) =============

var (
	ErrAIServiceDown  = &AppError{Code: 60001, HTTPCode: 503, Message: "AI 服务暂时不可用"}
	ErrAITokenLimit   = &AppError{Code: 60002, HTTPCode: 429, Message: "AI Token 配额已用尽"}
	ErrAIProviderFail = &AppError{Code: 60003, HTTPCode: 502, Message: "AI 模型服务异常"}
)

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "github.com/Zara1024/AIOps/cloudops-server/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// OKWithMessage 成功响应（自定义消息）
func OKWithMessage(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: msg,
		Data:    data,
	})
}

// OKWithPage 分页成功响应
func OKWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Error 错误响应（使用 AppError）
func Error(c *gin.Context, err *appErrors.AppError) {
	c.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Message: err.Message,
	})
}

// ErrorWithMsg 自定义消息错误响应
func ErrorWithMsg(c *gin.Context, httpCode int, code int, msg string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: msg,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    appErrors.ErrBadRequest.Code,
		Message: msg,
	})
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, msg string) {
	if msg == "" {
		msg = "未登录或登录已过期"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    appErrors.ErrUnauthorized.Code,
		Message: msg,
	})
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, msg string) {
	if msg == "" {
		msg = "无权限访问"
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    appErrors.ErrForbidden.Code,
		Message: msg,
	})
}

// ServerError 500 错误
func ServerError(c *gin.Context, msg string) {
	if msg == "" {
		msg = "服务器内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    appErrors.ErrInternal.Code,
		Message: msg,
	})
}

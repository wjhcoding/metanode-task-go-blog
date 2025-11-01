package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用返回结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 标准状态码（可根据项目扩展）
const (
	CodeSuccess      = 200
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

// OkMsg 返回成功消息
func OkMsg(msg string) Response {
	return Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: nil,
	}
}

// OkData 返回成功数据
func OkData(data interface{}) Response {
	return Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	}
}

// Ok 自定义成功消息+数据
func Ok(msg string, data interface{}) Response {
	return Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	}
}

// FailMsg 返回失败消息
func FailMsg(msg string) Response {
	return Response{
		Code: CodeBadRequest,
		Msg:  msg,
		Data: nil,
	}
}

// FailCodeMsg 自定义错误码与消息
func FailCodeMsg(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// WriteJSON 快捷返回
func WriteJSON(c *gin.Context, httpCode int, res Response) {
	c.JSON(httpCode, res)
}

// 便捷封装
func Success(c *gin.Context, data interface{}) {
	WriteJSON(c, http.StatusOK, OkData(data))
}

func Error(c *gin.Context, msg string) {
	WriteJSON(c, http.StatusOK, FailMsg(msg))
}

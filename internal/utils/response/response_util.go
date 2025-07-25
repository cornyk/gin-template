package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Pagination 分页响应结构体
type Pagination struct {
	List  interface{} `json:"list"`
	Count int         `json:"count"`
}

// jsonResponse Json响应结构体
type jsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SucJson 成功响应
func SucJson(c *gin.Context, data ...interface{}) {
	switch len(data) {
	case 0:
		Json(c, CodeSuccess, MsgSuccess, nil) // 无 data 参数
	case 1:
		Json(c, CodeSuccess, MsgSuccess, data[0]) // 单参数时直接使用，避免切片包裹
	default:
		Json(c, CodeSuccess, MsgSuccess, data) // 多参数时保持切片类型
	}
}

// Json Json响应
func Json(c *gin.Context, code int, message string, data ...interface{}) {
	var returnData interface{}
	switch len(data) {
	case 0:
		returnData = nil
	case 1:
		returnData = data[0]
	default:
		returnData = data
	}

	c.JSON(http.StatusOK, jsonResponse{
		Code:    code,
		Message: message,
		Data:    returnData,
	})
}

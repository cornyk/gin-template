package response_util

import (
	"cornyk/gin-template/internal/commons/response_def"
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
		Json(c, response_def.CodeSuccess, response_def.MsgSuccess, nil)
	case 1:
		Json(c, response_def.CodeSuccess, response_def.MsgSuccess, data[0])
	default:
		Json(c, response_def.CodeSuccess, response_def.MsgSuccess, data[0], data[1])
	}
}

// Json Json响应
func Json(c *gin.Context, code int, message string, data ...interface{}) {
	var returnData interface{}
	httpCode := http.StatusOK

	switch len(data) {
	case 0:
		returnData = nil
	case 1:
		returnData = data[0]
	default:
		returnData = data[0]
		httpCode = data[1].(int)
	}

	c.JSON(httpCode, jsonResponse{
		Code:    code,
		Message: message,
		Data:    returnData,
	})
}

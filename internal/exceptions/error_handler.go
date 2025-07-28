package exceptions

import (
	"cornyk/gin-template/internal/commons/response_def"
	"cornyk/gin-template/internal/utils/response_util"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 捕获 panic
		defer func() {
			if r := recover(); r != nil {
				// 记录日志
				logger.GetLogger(c, "error").Error(fmt.Sprintf("%v\n%s", r, debug.Stack()))

				// 处理返回结果
				errorMessage := response_def.MsgSystemError
				if global.GlobalConfig.App.Debug {
					errorMessage = fmt.Sprintf("%v", r)
				}
				response_util.Json(c, response_def.CodeSystemError, errorMessage, nil, http.StatusInternalServerError)
			}
		}()

		// 执行请求
		c.Next()

		// 捕获业务逻辑中的 c.Error(err)
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case *BusinessError: // 业务日志只返回，不记录日志
				httpStatus := http.StatusOK
				if e.HttpStatus != 0 {
					httpStatus = e.HttpStatus
				}
				response_util.Json(c, e.Code, e.Message, nil, httpStatus)
				return
			case *SystemError: // 系统错误不仅返回且需要记录日志
				// 记录日志
				logger.GetLogger(c, "error").Error(e.Message)

				httpStatus := http.StatusInternalServerError
				if e.HttpStatus != 0 {
					httpStatus = e.HttpStatus
				}
				response_util.Json(c, e.Code, e.Message, nil, httpStatus)
				return
			default: // 其他为明确类型的错误，格式化返回并记录日志，Debug模式下会输出真实错误信息
				// 记录日志
				logger.GetLogger(c, "error").Error(e.Error())

				// 处理返回结果
				errorMessage := response_def.MsgSystemError
				if global.GlobalConfig.App.Debug {
					errorMessage = e.Error()
				}
				response_util.Json(c, response_def.CodeSystemError, errorMessage, nil, http.StatusInternalServerError)
				return
			}
		}
	}
}

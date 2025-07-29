package middlewares

import (
	"bytes"
	"cornyk/gin-template/pkg/logger"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求信息
		requestIp := c.ClientIP()
		requestUrl := c.Request.URL.Path
		requestMethod := c.Request.Method
		requestParams := c.Request.URL.RawQuery

		// 记录请求头
		var requestHeaders string
		for name, values := range c.Request.Header {
			for _, value := range values {
				requestHeaders = requestData2String(requestHeaders, name, value)
			}
		}

		// 记录请求体
		var requestBody string
		if c.Request.Body != nil {
			if c.ContentType() == "application/json" {
				body, _ := io.ReadAll(c.Request.Body)

				// 去掉json的空格和换行
				bodyString := strings.ReplaceAll(string(body), "\n", "")
				requestBody = strings.ReplaceAll(bodyString, " ", "")

				c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 写回 Body 供后续使用
			} else if strings.HasPrefix(c.ContentType(), "application/x-www-form-urlencoded") {
				// 解析普通的表单数据
				if err := c.Request.ParseForm(); err == nil {
					for k, v := range c.Request.PostForm {
						requestBody = requestData2String(requestBody, k, v[0])
					}
				}
			} else if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
				if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
					// 打印普通字段
					for k, v := range c.Request.PostForm {
						requestBody = requestData2String(requestBody, k, v[0])
					}
					// 打印文件字段
					for k, v := range c.Request.MultipartForm.File {
						requestBody = requestData2String(requestBody, k, v[0].Filename)
					}
				}
			} else {
				body, _ := io.ReadAll(c.Request.Body)

				var requestBodyBuffer bytes.Buffer
				requestBodyBuffer.Write(body)
				requestBody = requestBodyBuffer.String()

				c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 写回 Body 供后续使用
			}
		}

		// 初始化响应记录器
		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = recorder

		c.Next()

		// 记录响应信息
		statusCode := c.Writer.Status()
		latency := time.Since(time.Now()) // 请求耗时
		responseBody := recorder.body.String()

		// 记录请求日志
		logger.GetLogger(c, "access").Info(
			fmt.Sprintf("[%s][%v]URL:'%s', STATUS_CODE:'%d', METHOD:'%s', QUERY_PARAMS:'%s', BODY:'%s', HEADERS:'%s', RESPONSE_BODY:'%s'",
				requestIp,
				latency,
				requestUrl,
				statusCode,
				requestMethod,
				requestParams,
				requestBody,
				requestHeaders,
				responseBody,
			))
	}
}

// 添加请求日志信息
func requestData2String(requestBody string, k string, v string) string {
	if requestBody == "" {
		requestBody = fmt.Sprintf("{%s:%v}", k, v)
	} else {
		requestBody = fmt.Sprintf("%s,{%s:%v}", requestBody, k, v)
	}
	return requestBody
}

// 自定义 ResponseWriter 用于捕获响应体
type responseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseRecorder) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

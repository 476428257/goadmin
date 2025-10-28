package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

// OperationLog 操作日志中间件
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestData = string(bodyBytes)
			// 恢复请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		c.Set("requestData", requestData)
		c.Next()
	}
}

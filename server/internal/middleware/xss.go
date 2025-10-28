package middleware

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// XSSFilter 是一个防XSS的中间件
func XSSFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理URL查询参数
		if len(c.Request.URL.RawQuery) > 0 {
			cleanQuery := sanitizeValues(c.Request.URL.Query())
			c.Request.URL.RawQuery = cleanQuery.Encode()
		}

		// 处理POST、PUT等请求体
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")

			// 处理表单数据
			if strings.Contains(contentType, "application/x-www-form-urlencoded") ||
				strings.Contains(contentType, "multipart/form-data") {
				// 确保表单已解析
				if err := c.Request.ParseForm(); err == nil {
					c.Request.Form = sanitizeValues(c.Request.Form)
					c.Request.PostForm = sanitizeValues(c.Request.PostForm)
				}
			}

			// 处理JSON数据
			if strings.Contains(contentType, "application/json") {
				// 读取原始请求体
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				c.Request.Body.Close()

				// 创建一个新的请求体，替换原始请求体
				sanitizedBody := sanitizeJSON(string(bodyBytes))
				c.Request.Body = io.NopCloser(bytes.NewBufferString(sanitizedBody))

				// 更新Content-Length
				c.Request.ContentLength = int64(len(sanitizedBody))
			}
		}

		c.Next()
	}
}

// sanitizeValues 清理URL值
func sanitizeValues(values url.Values) url.Values {
	result := url.Values{}
	for key, vals := range values {
		for _, val := range vals {
			// 使用HTMLEscapeString进行HTML转义
			result.Add(key, template.HTMLEscapeString(val))
		}
	}
	return result
}

// sanitizeJSON 清理JSON字符串
func sanitizeJSON(jsonStr string) string {
	// 先尝试解析JSON
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		// 如果不是有效的JSON，则直接进行字符替换
		return template.HTMLEscapeString(jsonStr)
	}

	// 递归处理JSON数据
	sanitizedData := sanitizeJSONData(data)

	// 重新序列化为JSON
	sanitizedBytes, err := json.Marshal(sanitizedData)
	if err != nil {
		// 如果序列化失败，则使用原始字符串处理
		return template.HTMLEscapeString(jsonStr)
	}

	return string(sanitizedBytes)
}

// sanitizeJSONData 递归处理JSON数据结构
func sanitizeJSONData(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case string:
		// 对字符串类型进行HTML转义
		return template.HTMLEscapeString(v)

	case map[string]interface{}:
		// 处理对象类型
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = sanitizeJSONData(val)
		}
		return result

	case []interface{}:
		// 处理数组类型
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = sanitizeJSONData(val)
		}
		return result

	default:
		// 其他类型保持不变
		return v
	}
}

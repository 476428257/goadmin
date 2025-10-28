package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 验证规则中文映射
var validationMessages = map[string]string{
	"required":     "不能为空",
	"min":          "长度不能小于%s",
	"max":          "长度不能大于%s",
	"oneof":        "必须是以下值之一：%s",
	"email":        "邮箱格式不正确",
	"numeric":      "必须是数字",
	"alpha":        "只能包含字母",
	"alphanum":     "只能包含字母和数字",
	"len":          "长度必须为%s",
	"eq":           "必须等于%s",
	"ne":           "不能等于%s",
	"gt":           "必须大于%s",
	"gte":          "必须大于等于%s",
	"lt":           "必须小于%s",
	"lte":          "必须小于等于%s",
	"unique":       "值已存在",
	"url":          "必须是有效的URL",
	"uri":          "必须是有效的URI",
	"json":         "必须是有效的JSON格式",
	"uuid":         "必须是有效的UUID",
	"datetime":     "必须是有效的日期时间格式",
	"route_format": "路由格式不正确，必须以/开头",
}

// InitValidator 初始化验证器
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		v.RegisterValidation("route_format", validateRoute)

		// 注册字段名获取函数，优先使用name标签，然后是json标签，最后是结构体字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			// 首先尝试获取name标签
			nameTag := fld.Tag.Get("name")
			if nameTag != "" && nameTag != "-" {
				return nameTag
			}

			// 如果没有name标签，尝试获取json标签
			jsonTag := fld.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				// 提取json标签中的字段名（去掉omitempty等选项）
				name := strings.SplitN(jsonTag, ",", 2)[0]
				if name != "" {
					return name
				}
			}

			// 如果都没有，返回结构体字段名
			return fld.Name
		})
	}
}

// validateRoute 验证路由格式
func validateRoute(fl validator.FieldLevel) bool {
	route := fl.Field().String()
	return strings.HasPrefix(route, "/")
}

// formatValidationError 格式化验证错误信息
func formatValidationError(err validator.FieldError) string {
	fieldName := err.Field()
	tag := err.Tag()
	param := err.Param()

	// 获取错误消息模板
	template, exists := validationMessages[tag]
	if !exists {
		// 如果没有对应的中文模板，返回默认错误信息
		return fmt.Sprintf("%s验证失败", fieldName)
	}

	// 根据不同的验证规则格式化消息
	switch tag {
	case "required":
		return fmt.Sprintf("%s%s", fieldName, template)
	case "min", "max", "len", "eq", "ne", "gt", "gte", "lt", "lte":
		return fmt.Sprintf("%s%s", fieldName, fmt.Sprintf(template, param))
	case "oneof":
		// 处理oneof参数，将空格分隔的值转为中文描述
		options := strings.ReplaceAll(param, " ", "、")
		return fmt.Sprintf("%s%s", fieldName, fmt.Sprintf(template, options))
	case "email", "numeric", "alpha", "alphanum", "unique", "url", "uri", "json", "uuid", "datetime", "route_format":
		return fmt.Sprintf("%s%s", fieldName, template)
	default:
		// 未知验证规则，返回通用错误
		return fmt.Sprintf("%s验证失败: %s", fieldName, err.Error())
	}
}

// Validate 通用验证中间件
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加通用的验证逻辑
		c.Next()
	}
}

// ValidateJSON 验证JSON请求体的中间件
func ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			// 检查是否是验证错误
			if errs, ok := err.(validator.ValidationErrors); ok {
				// 处理验证错误，转换为中文提示
				errorMessages := make([]string, 0, len(errs))
				for _, e := range errs {
					errorMessages = append(errorMessages, formatValidationError(e))
				}

				// 返回中文错误信息
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "参数验证失败: " + strings.Join(errorMessages, "; "),
					"data": nil,
				})
				c.Abort()
				return
			}

			// 其他JSON解析错误
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "请求参数格式错误",
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ValidateForm 验证表单数据的中间件
func ValidateForm(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBind(obj); err != nil {
			// 检查是否是验证错误
			if errs, ok := err.(validator.ValidationErrors); ok {
				// 处理验证错误，转换为中文提示
				errorMessages := make([]string, 0, len(errs))
				for _, e := range errs {
					errorMessages = append(errorMessages, formatValidationError(e))
				}

				// 返回中文错误信息
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "参数验证失败: " + strings.Join(errorMessages, "; "),
					"data": nil,
				})
				c.Abort()
				return
			}

			// 其他绑定错误
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "请求参数格式错误",
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ValidateRoute 路由验证中间件
func ValidateRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/") {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "无效的路由格式",
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetValidationError 获取验证错误的中文提示（供控制器使用）
func GetValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make([]string, 0, len(errs))
		for _, e := range errs {
			errorMessages = append(errorMessages, formatValidationError(e))
		}
		return strings.Join(errorMessages, "; ")
	}
	return "参数验证失败"
}

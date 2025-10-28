package middleware

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
)

func filterEmpty(s []string) []string {
	var out []string
	for _, v := range s {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

// Authorize 权限检查中间件
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 指定方法不验证权限
		path := c.Request.URL.Path
		NoAuth := []string{"/api/v1/login", "/api/v1/auth/rule/getmenu", "/api/v1/config/getconfig"}
		if slices.Contains(NoAuth, path) {
			c.Next()
			return
		}
		// 指定方法后缀不验证权限
		NoAuthVague := []string{"getkv"}
		parts := filterEmpty(strings.Split(path, "/"))
		if len(parts) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "路由错误",
			})
			return
		}
		if slices.Contains(NoAuthVague, parts[len(parts)-1]) {
			c.Next()
			return
		}
		// 不验证权限的路径
		NoAuthSTRING := []string{"/uploads/"}
		for _, v := range NoAuthSTRING {
			if strings.Contains(c.Request.URL.Path, v) {
				c.Next()
				return
			}
		}
		// 获取当前用户ID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			c.Abort()
			return
		}

		// 获取用户角色信息
		var authUser model.AuthUser
		if err := database.DB.Debug().Preload("Role").Where("admin_id = ?", userID).First(&authUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未分配角色",
			})
			c.Abort()
			return
		}
		// 如果角色规则为*，则直接放行
		if authUser.Role.Rule == "*" {
			c.Next()
			return
		}
		// 获取当前请求的路径
		requestPath := c.Request.URL.Path
		// 查找对应的权限规则
		var rule model.AuthRule
		if err := database.DB.Where("route = ? ", requestPath).First(&rule).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限规则未配置",
			})
			c.Abort()
			return
		}

		// 检查用户角色是否有权限访问
		ruleIDs := strings.Split(authUser.Role.Rule, ",")
		hasPermission := slices.Contains(ruleIDs, strconv.Itoa(int(rule.ID)))
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "没有访问权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

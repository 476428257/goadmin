package router

import (
	"server/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册权限相关路由
func RegisterAuthRoutes(r *gin.RouterGroup) {
	authController := controller.NewAuthController()
	AuthRoleController := controller.NewAuthRoleController()
	// 权限管理相关路由
	auth := r.Group("/auth")
	{
		// 角色管理
		role := auth.Group("/role")
		{
			// 创建角色
			role.POST("/add", AuthRoleController.Create)
			role.POST("/getkv", AuthRoleController.Getkv) // 获取键值对
			// 更新角色
			role.POST("/update", AuthRoleController.Update)
			// 删除角色
			role.POST("/del", AuthRoleController.Delete)
			// 修改状态
			role.POST("/updatestatus", AuthRoleController.UpdateStatus)
			// 获取单个角色
			role.POST("/info", AuthRoleController.GetByID)
			// 获取角色树结构
			role.POST("/list", AuthRoleController.List)
		}

		// 权限规则管理
		rule := auth.Group("/rule")
		{
			rule.POST("/add", authController.Create)                // 创建规则
			rule.POST("/getkv", authController.Getkv)               // 获取键值对
			rule.POST("/update", authController.Update)             // 更新规则
			rule.POST("/updatestatus", authController.UpdateStatus) // 更新规则
			rule.POST("/del", authController.Delete)                // 删除规则
			rule.POST("/info", authController.GetByID)              // 获取单个规则
			rule.POST("/list", authController.List)                 // 获取规则树结构
			rule.POST("/getmenu", authController.GetMenu)           // 获取菜单
		}

		// 权限刷新
		auth.POST("/refresh", authController.RefreshAuth) // 刷新权限策略
	}
}

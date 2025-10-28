package router

import (
	"server/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理员相关路由
func RegisterAdminRoutes(r *gin.RouterGroup) {
	adminController := controller.NewAdminController()

	// 公开路由
	public := r.Group("")
	{
		public.POST("/login", adminController.Login) // 管理员登录
	}

	// 需要认证的管理员路由
	admin := r.Group("/admin")
	{
		admin.POST("/add", adminController.Create)                // 创建管理员
		admin.POST("/update", adminController.Update)             // 更新管理员
		admin.POST("/updatestatus", adminController.UpdateStatus) // 更新管理员状态
		admin.POST("/updateavatar", adminController.UpdateAvatar)
		admin.POST("/updatepassword", adminController.UpdatePassword)
		admin.POST("/del", adminController.Delete)   // 删除管理员
		admin.POST("/list", adminController.List)    // 管理员列表
		admin.POST("/info", adminController.GetByID) // 获取管理员信息
	}
}

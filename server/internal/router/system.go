package router

import (
	"server/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 系统接口路由
func RegisterSystemRoutes(r *gin.RouterGroup) {
	systemController := controller.NewSystemController()
	// 系统接口路由
	system := r.Group("/system")
	{
		system.POST("/aclog", systemController.Aclog)          // 操作日志
		system.POST("/dashboard", systemController.Dashboard)  // 首页
		system.POST("/upload", systemController.ProcessUpload) // 上传文件
	}
}

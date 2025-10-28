package router

import (
	"server/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterConfigRoutes 参数配置接口路由
func RegisterConfigRoutes(r *gin.RouterGroup) {
	configController := controller.NewConfigController()
	// 参数配置接口路由
	config := r.Group("/config")
	{
		config.POST("/list", configController.List)           // 配置列表
		config.POST("/update", configController.Update)       // 更新配置并写入 Redis
		config.POST("/getconfig", configController.Getconfig) // 从 Redis 获取配置
	}
}

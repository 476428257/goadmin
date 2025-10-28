package router

import (
	"server/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterArticleRoutes 文章管理相关路由
func RegisterArticleRoutes(r *gin.RouterGroup) {
	articleController := controller.NewArticleController()
	// 需要认证的文章路由
	article := r.Group("/article")
	{
		article.POST("/add", articleController.Create)                // 创建文章
		article.POST("/update", articleController.Update)             // 更新文章
		article.POST("/updatestatus", articleController.UpdateStatus) // 更新文章状态
		article.POST("/del", articleController.Delete)                // 删除文章
		article.POST("/list", articleController.List)                 // 文章列表
		article.POST("/info", articleController.GetByID)              // 获取文章信息
	}
}

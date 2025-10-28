package router

import (
	"server/config"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	// 初始化验证器
	middleware.InitValidator()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// 添加跨域处理中间件（必须在静态文件服务之前）
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RedisMiddleware())
	// 静态资源访问
	uploaddir := config.GetConfig().Upload.UploadDir
	r.Static("/"+uploaddir, "./"+uploaddir)

	// 添加XSS防护中间件
	r.Use(middleware.XSSFilter())
	// 添加JWT认证中间件
	r.Use(middleware.JWTAuth())
	// 添加权限认证中间件
	r.Use(middleware.Authorize())
	// 添加操作日志中间键
	r.Use(middleware.OperationLog())
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 注册管理员路由
		RegisterAdminRoutes(v1)

		// 注册权限管理路由
		RegisterAuthRoutes(v1)
		// 注册系统接口路由
		RegisterSystemRoutes(v1)
		// 注册文章接口路由
		RegisterArticleRoutes(v1)
		// 注册配置接口路由
		RegisterConfigRoutes(v1)
	}

	return r
}

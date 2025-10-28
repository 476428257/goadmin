package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/config"
	"server/internal/router"
	"server/pkg/database"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
)

const UPLOADDIR = "uploads"

// setGinMode 根据配置设置gin运行模式
func setGinMode(mode string) {
	switch mode {
	case "release", "production":
		gin.SetMode(gin.ReleaseMode)
		log.Println("Gin is running in release mode")
	case "test":
		gin.SetMode(gin.TestMode)
		log.Println("Gin is running in test mode")
	default:
		gin.SetMode(gin.DebugMode)
		log.Println("Gin is running in debug mode")
	}
}

func main() {
	// 初始化配置
	config.Init()
	cfg := config.GetConfig()

	// 设置gin框架运行模式（根据配置文件中的mode设置）
	setGinMode(cfg.Server.Mode)

	// 初始化日志
	logger.Init()

	// 初始化数据库
	database.Init()

	// 执行数据库迁移
	database.AutoMigrate()

	// 初始化路由
	r := router.Init()

	// 启动服务器
	startServer(r, cfg)
}

// startServer 启动HTTP/HTTPS服务器
func startServer(router *gin.Engine, cfg *config.Config) {
	// 创建HTTP服务器
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// 如果启用HTTPS
	if cfg.Server.HTTPS.Enabled {
		// 创建HTTPS服务器
		httpsServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPS.Port),
			Handler: router,
		}

		// 启动HTTPS服务器
		go func() {
			log.Printf("HTTPS Server starting on https://%s:%d", cfg.Server.Domain, cfg.Server.HTTPS.Port)
			if err := httpsServer.ListenAndServeTLS(cfg.Server.HTTPS.CertFile, cfg.Server.HTTPS.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS server startup failed: %v", err)
			}
		}()

		// 可选：HTTP重定向到HTTPS
		go func() {
			redirectServer := &http.Server{
				Addr: fmt.Sprintf(":%d", cfg.Server.Port),
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					httpsURL := fmt.Sprintf("https://%s:%d%s", cfg.Server.Domain, cfg.Server.HTTPS.Port, r.RequestURI)
					http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
				}),
			}
			log.Printf("HTTP Redirect Server starting on http://%s:%d (redirecting to HTTPS)", cfg.Server.Domain, cfg.Server.Port)
			if err := redirectServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP redirect server error: %v", err)
			}
		}()

		// 优雅关闭
		gracefulShutdown(httpsServer)
	} else {
		// 启动HTTP服务器
		go func() {
			log.Printf("HTTP Server starting on http://%s:%d", cfg.Server.Domain, cfg.Server.Port)
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP server startup failed: %v", err)
			}
		}()

		// 优雅关闭
		gracefulShutdown(httpServer)
	}
}

// gracefulShutdown 优雅关闭服务器
func gracefulShutdown(server *http.Server) {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

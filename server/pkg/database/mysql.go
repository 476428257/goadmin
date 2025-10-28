package database

import (
	"fmt"
	"time"

	"server/config"
	"server/internal/model"
	"server/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
		cfg.MySQL.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "h_", // 统一前缀
			SingularTable: true, // 禁用复数表名
		},
	})
	if err != nil {
		logger.Logger.Fatal("Failed to connect to database",
			zap.Error(err),
		)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Logger.Fatal("Failed to get database instance",
			zap.Error(err),
		)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// logger.Logger.Info("Database connected successfully")
}

// AutoMigrate 执行数据库迁移
func AutoMigrate() {
	err := DB.AutoMigrate(
		&model.Config{},
	// 在这里添加需要迁移的模型指针
	)
	if err != nil {
		logger.Logger.Fatal("Failed to migrate database",
			zap.Error(err),
		)
	}
}

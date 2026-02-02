package svc

import (
	"fmt"
	"go-zero-template/internal/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MustInitDB(pgConfig config.PostgresConfig) (*gorm.DB, string) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=Asia/Shanghai",
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.User,
		pgConfig.Password,
		pgConfig.DBName,
		pgConfig.SSLMode,
	)
	// 根据配置设置日志级别
	var logLevel logger.LogLevel
	switch pgConfig.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Error
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
		CreateBatchSize:                          100,
		NowFunc: func() time.Time {
			return time.Now().In(loc)
		},
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移（先迁移被引用的表，再迁移引用表）
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	// 检查数据库连接是否正常
	if err := PingDB(db); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("database connected and migrated successfully")
	return db, dsn
}

// PingDB 检查数据库连接是否正常
func PingDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

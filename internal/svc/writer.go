package svc

import (
	"log"
	"time"

	writer "github.com/zhengliu92/pg-log-writter"
	"gorm.io/gorm"
)

func MustInitWriter(dsn string, gormDB *gorm.DB) *writer.MultiWriter {
	writerConfig := &writer.PostgresConfig{
		TableName:     "logs",
		BufferSize:    100,
		FlushInterval: 3 * time.Second,
	}
	pgxExecutor, err := NewPgxExecutor(dsn)
	if err != nil {
		log.Fatalf("failed to create pgx executor: %v", err)
	}
	pgWriter, err := writer.NewPostgresqlWriter(pgxExecutor, writerConfig)
	if err != nil {
		log.Fatalf("failed to create pg writer: %v", err)
	}
	consoleWriter := writer.NewConsoleWriter()
	return writer.NewMultiWriter(pgWriter, consoleWriter)
}

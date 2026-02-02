package svc

import (
	"go-zero-template/internal/config"
	"go-zero-template/internal/db"
	"go-zero-template/internal/middleware"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"

	writer "github.com/zhengliu92/pg-log-writter"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Client
	Repository     *db.Repository
	Writer         *writer.MultiWriter
	AuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	gormDB, dsn := MustInitDB(c.Postgres)
	redisClient := MustInitRedis(c.Redis)
	repository := db.NewRepository(gormDB)
	writer := MustInitWriter(dsn, gormDB)

	return &ServiceContext{
		Config:         c,
		Redis:          redisClient,
		Repository:     repository,
		Writer:         writer,
		AuthMiddleware: middleware.NewAuthMiddleware(&c.Services).Handle,
	}
}

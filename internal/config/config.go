package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Postgres PostgresConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Services ServicesConfig
}

type AuthConfig struct {
	AccessSecret string `json:",env=AUTH_ACCESS_SECRET"`
}

type ServicesConfig struct {
	UserService BaseServiceConfig
}

type BaseServiceConfig struct {
	Host               string `json:",env=USER_SERVICE_HOST"`
	Port               int    `json:",env=USER_SERVICE_PORT"`
	Path               string `json:",env=USER_SERVICE_PATH"`
	SuperAdminUsername string `json:",env=USER_SERVICE_SUPER_ADMIN_USERNAME"`
	SuperAdminPassword string `json:",env=USER_SERVICE_SUPER_ADMIN_PASSWORD"`
}

type PostgresConfig struct {
	Host     string `json:",env=POSTGRES_HOST"`
	Port     int    `json:",env=POSTGRES_PORT"`
	User     string `json:",env=POSTGRES_USER"`
	Password string `json:",env=POSTGRES_PASSWORD"`
	DBName   string `json:",default=user_service,env=POSTGRES_DBNAME"`
	SSLMode  string `json:",default=disable,env=POSTGRES_SSLMODE"`
	LogLevel string `json:",default=error,env=POSTGRES_LOGLEVEL"`
}

type RedisConfig struct {
	Addr     string `json:",default=localhost:6379,env=REDIS_ADDR"`
	Password string `json:",optional,env=REDIS_PASSWORD"`
	DB       int    `json:",default=0,env=REDIS_DB"`
}

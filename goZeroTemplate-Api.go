// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"log"

	"go-zero-template/internal/config"
	"go-zero-template/internal/handler"
	"go-zero-template/internal/svc"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/goZeroTemplate-Api.yaml", "the config file")

func main() {
	flag.Parse()
	// 加载 .env 文件（如果存在）
	// 注意：.env 文件中的环境变量不会覆盖已经设置的环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，将使用系统环境变量或配置文件中的默认值")
	}
	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

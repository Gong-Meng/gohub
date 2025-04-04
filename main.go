package main

import (
	"flag"
	"fmt"

	btsConfig "github.com/gongmeng/gohub/config"
	"github.com/gongmeng/gohub/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/gongmeng/gohub/bootstrap"
)

func init() {
	// 加载 config 目录下的配置信息
	btsConfig.Initialize()
}

func main() {

	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件，如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()

	config.InitConfig(env)

	// 初始化 Logger
	bootstrap.SetupLogger()

	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)

	// 初始化gin
	router := gin.New()

	// 初始化数据库
	bootstrap.SetupDB()

	// 初始化 Redis
	bootstrap.SetupRedis()

	// 绑定路由
	bootstrap.SetupRoute(router)

	// 运行
	err := router.Run(":" + config.GetString("app.port"))
	if err != nil {
		// 处理错误
		fmt.Println(err.Error())
	}
}

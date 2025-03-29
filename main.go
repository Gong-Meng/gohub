package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gongmeng/gohub/bootstrap"
)

func main() {

	// 初始化gin
	router := gin.New()

	// 绑定路由
	bootstrap.SetupRoute(router)

	// 运行
	err := router.Run(":3000")
	if err != nil {
		// 处理错误
		fmt.Println(err.Error())
	}
}

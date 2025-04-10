package bootstrap

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gongmeng/gohub/app/http/middlewares"
	"github.com/gongmeng/gohub/routes"
)

func SetupRoute(router *gin.Engine) {

	// 注册全局中间件
	registerGlobalMiddleware(router)

	// 注册路由
	routes.RegisterAPIRoutes(router)

	// 处理404
	setup404Handler(router)
}

func registerGlobalMiddleware(router *gin.Engine) {
	// 注册中间件
	router.Use(
		middlewares.Logger(),
		middlewares.Recovery(),
	)
}

func setup404Handler(router *gin.Engine) {
	// 处理404
	router.NoRoute(func(c *gin.Context) {
		acceptString := c.Request.Header.Get("Accept")

		if strings.Contains(acceptString, "text/html") {
			// html
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			// json
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})
}

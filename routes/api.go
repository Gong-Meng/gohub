// Package routes 注册路由
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gongmeng/gohub/app/http/controllers/api/v1/auth"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	// v1
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"hello": "world",
			})
		})

		// auth
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			// 判断手机号是否注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			// 判断邮箱是否注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
		}
	}

}

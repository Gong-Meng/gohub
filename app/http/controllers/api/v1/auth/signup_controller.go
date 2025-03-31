// Package auth 处理用户身份认证相关逻辑
package auth

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/gongmeng/gohub/app/http/controllers/api/v1"
	"github.com/gongmeng/gohub/app/models/user"
	"github.com/gongmeng/gohub/app/requests"
	"github.com/gongmeng/gohub/pkg/response"
)

// SignupController 注册控制器
type SignupController struct {
	v1.BaseApiController
}

// IsPhoneExist 检测手机号是否被注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {

	// 获取请求参数，并做表单验证
	request := requests.SignupPhoneExistRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateSignupPhoneExistRequest); !ok {
		return
	}

	// 检查数据库并返回响应
	response.JSON(c, gin.H{
		"exists": user.IsPhoneExist(request.Phone),
	})
}

// IsEmailExist 检测邮箱是否已注册
func (sc *SignupController) IsEmailExist(c *gin.Context) {

	// 获取请求参数，并做表单验证
	request := requests.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateSignupEmailExistRequest); !ok {
		return
	}

	// 检查数据库并返回响应
	response.JSON(c, gin.H{
		"exists": user.IsPhoneExist(request.Email),
	})
}

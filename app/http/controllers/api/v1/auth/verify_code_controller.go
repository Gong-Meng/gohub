package auth

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/gongmeng/gohub/app/http/controllers/api/v1"
	"github.com/gongmeng/gohub/app/requests"
	"github.com/gongmeng/gohub/pkg/captcha"
	"github.com/gongmeng/gohub/pkg/logger"
	"github.com/gongmeng/gohub/pkg/response"
	"github.com/gongmeng/gohub/pkg/verifycode"
)

// VerifyCodeController 用户控制器
type VerifyCodeController struct {
	v1.BaseApiController
}

// ShowCaptcha 显示图片验证码
func (vc *VerifyCodeController) ShowCaptcha(c *gin.Context) {
	// 生成验证码
	id, b64s, _, err := captcha.NewCaptcha().GenerateCaptcha()
	// 记录错误日志，因为验证码是用户的入口，出错时应该记 error 等级的日志
	logger.LogIf(err)

	// 返回给用户
	response.JSON(c, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}

// SendUsingPhone 发送手机验证码
func (vc *VerifyCodeController) SendUsingPhone(c *gin.Context) {

	// 1. 验证表单
	request := requests.VerifyCodePhoneRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodePhone); !ok {
		return
	}

	// 2. 发送 SMS
	if code, ok := verifycode.NewVerifyCode().SendSMS(request.Phone); !ok {
		response.Abort500(c, "发送短信失败~")
	} else {
		response.Data(c, gin.H{"code": code})
	}
}

// SendUsingEmail 发送 Email 验证码
func (vc *VerifyCodeController) SendUsingEmail(c *gin.Context) {

	// 1. 验证表单
	request := requests.VerifyCodeEmailRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodeEmail); !ok {
		return
	}

	// 2. 发送邮件
	code, ok := verifycode.NewVerifyCode().SendEmail(request.Email)
	if !ok {
		response.Abort500(c, "发送 Email 验证码失败~")
	} else {
		response.Data(c, gin.H{"code": code})
	}
}

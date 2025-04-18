package verifycode

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gongmeng/gohub/app"
	"github.com/gongmeng/gohub/pkg/config"
	"github.com/gongmeng/gohub/pkg/helpers"
	"github.com/gongmeng/gohub/pkg/logger"
	"github.com/gongmeng/gohub/pkg/mail"
	"github.com/gongmeng/gohub/pkg/redis"
	"github.com/gongmeng/gohub/pkg/sms"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once

var internalVerifyCode *VerifyCode

// NewVerifyCode 单例模式获取
func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Redis,
				// 增加前缀保持数据库整洁，出问题调试时也方便
				KeyPrefix: config.GetString("app.name") + ":verifycode:",
			},
		}
	})
	return internalVerifyCode
}

// SendSMS 发送短信验证码，调用示例：
// verifycode.NewVerifyCode().SendSMS(request.Phone)
func (vc *VerifyCode) SendSMS(phone string) (string, bool) {

	// 生成验证码
	code := vc.generateVerifyCode(phone)

	// 方便本地和 API 自动测试
	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix")) {
		return code, true
	}

	// 发送短信
	result := sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})

	// 在生产环境中，不返回验证码，只返回发送状态
	if app.IsProduction() {
		return "", result
	}

	return code, result
}

// SendEmail 发送邮件验证码，调用示例：
//
//	verifycode.NewVerifyCode().SendEmail(request.Email)
func (vc *VerifyCode) SendEmail(email string) (string, bool) {

	// 生成验证码
	code := vc.generateVerifyCode(email)

	// 方便本地和 API 自动测试
	if !app.IsProduction() && strings.HasSuffix(email, config.GetString("verifycode.debug_email_suffix")) {
		return code, true
	}

	content := fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)
	// 发送邮件
	result := mail.NewMailer().Send(mail.Email{
		From: mail.From{
			Address: config.GetString("mail.from.address"),
			Name:    config.GetString("mail.from.name"),
		},
		To:      []string{email},
		Subject: "Email 验证码",
		HTML:    []byte(content),
	})

	// 在生产环境中，不返回验证码，只返回发送状态
	if app.IsProduction() {
		return "", result
	}

	return code, result
}

// CheckAnswer 检查用户提交的验证码是否正确，key 可以是手机号或者 Email
func (vc *VerifyCode) CheckAnswer(key string, answer string) bool {

	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})

	// 方便开发，在非生产环境下，具备特殊前缀的手机号和 Email后缀，会直接验证成功
	if !app.IsProduction() &&
		(strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix")) ||
			strings.HasPrefix(key, config.GetString("verifycode.debug_phone_prefix"))) {
		return true
	}

	return vc.Store.Verify(key, answer, false)
}

// generateVerifyCode 生成验证码，并放置于 Redis 中
func (vc *VerifyCode) generateVerifyCode(key string) string {
	// 生成随机码
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	// 为方便开发，本地环境使用固定验证码
	if app.IsLocal() {
		code = config.GetString("verifycode.debug_code")
	}

	logger.DebugJSON("验证码", "生成验证码", map[string]string{key: code})

	// 将验证码及 KEY（邮箱或手机号）存放到 Redis 中并设置过期时间
	vc.Store.Set(key, code)
	return code
}

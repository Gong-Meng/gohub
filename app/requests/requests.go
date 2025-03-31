// Package requests 处理请求数据和表单验证

package requests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gongmeng/gohub/pkg/response"
	"github.com/thedevsaddam/govalidator"
)

// ValidatorFunc 验证函数类型
type ValidatatorFunc func(interface{}, *gin.Context) map[string][]string

// Validate 控制器里调用示例：
//
//	if ok := requests.Validate(c, &requests.UserSaveRequest{}, requests.UserSave); !ok {
//	    return
//	}
func Validate(c *gin.Context, obj interface{}, handler ValidatatorFunc) bool {
	// 1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(obj); err != nil {
		response.BadRequest(c, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		// 打印错误信息
		fmt.Println(err.Error())
		return false
	}

	// 2. 表单验证
	errs := handler(obj, c)

	// 3. 验证失败，返回 422 状态码和错误信息
	if len(errs) > 0 {
		response.ValidationError(c, errs)
		return false
	}

	// 4. 验证成功，返回 true
	return true
}

func validate(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	// 配置初始化
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid", // 自定义验证规则的标识符 模型中的 Struct 标签标识符
		Messages:      messages,
	}

	// 执行验证
	return govalidator.New(opts).ValidateStruct()
}

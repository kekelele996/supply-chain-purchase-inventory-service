package util

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// Validate 全局 validator 实例。
var Validate *validator.Validate

func init() {
	Validate = binding.Validator.Engine().(*validator.Validate)
}

// FieldError 字段级校验错误。
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors 收集 validator 错误并转换为字段级错误列表。
func ValidationErrors(err error) []FieldError {
	if err == nil {
		return nil
	}
	var verr validator.ValidationErrors
	if !asValidationErrors(err, &verr) {
		return []FieldError{{Field: "request", Message: err.Error()}}
	}
	out := make([]FieldError, 0, len(verr))
	for _, fe := range verr {
		out = append(out, FieldError{
			Field:   fe.Field(),
			Message: validationMessage(fe),
		})
	}
	return out
}

func asValidationErrors(err error, target *validator.ValidationErrors) bool {
	if v, ok := err.(validator.ValidationErrors); ok {
		*target = v
		return true
	}
	if s, ok := err.(*validator.ValidationErrors); ok {
		*target = *s
		return true
	}
	return false
}

func validationMessage(fe validator.FieldError) string {
	tag := fe.Tag()
	switch tag {
	case "required":
		return fmt.Sprintf("%s 为必填项", fe.Field())
	case "email":
		return fmt.Sprintf("%s 必须是合法邮箱", fe.Field())
	case "min":
		return fmt.Sprintf("%s 最小值不满足要求", fe.Field())
	case "max":
		return fmt.Sprintf("%s 超出最大长度", fe.Field())
	case "gte":
		return fmt.Sprintf("%s 必须大于等于 %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s 必须大于 %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s 必须小于等于 %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s 必须是以下值之一: %s", fe.Field(), fe.Param())
	case "datetime":
		return fmt.Sprintf("%s 必须是 %s 格式日期", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s 长度必须为 %s", fe.Field(), fe.Param())
	case "alphanum":
		return fmt.Sprintf("%s 只能包含字母和数字", fe.Field())
	case "numeric":
		return fmt.Sprintf("%s 必须是数字", fe.Field())
	default:
		return fmt.Sprintf("%s 校验失败(%s=%s)", fe.Field(), tag, fe.Param())
	}
}

// BuildValidationError 构造 422 校验错误异常。
func BuildValidationError(err error) error {
	fields := ValidationErrors(err)
	msgs := make([]string, 0, len(fields))
	for _, f := range fields {
		msgs = append(msgs, f.Message)
	}
	return NewAppError(constants.CodeValidationError, 422, strings.Join(msgs, "; "), err)
}

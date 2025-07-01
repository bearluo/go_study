package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator Echo 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator 创建新的自定义验证器
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// 注册自定义验证标签
	registerCustomValidations(v)

	return &CustomValidator{
		validator: v,
	}
}

// Validate 实现 Echo 的 Validator 接口
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations(v *validator.Validate) {
	// 注册自定义验证函数
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("chinese", validateChinese)
	v.RegisterValidation("chinese_name", validateChineseName)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("idcard", validateIDCard)
	v.RegisterValidation("postal_code", validatePostalCode)
	v.RegisterValidation("credit_card", validateCreditCard)
	v.RegisterValidation("strong_password", validateStrongPassword)
	v.RegisterValidation("url", validateURL)
	v.RegisterValidation("ipv4", validateIPv4)
	v.RegisterValidation("ipv6", validateIPv6)
	v.RegisterValidation("date", validateDate)
	v.RegisterValidation("time", validateTime)
	v.RegisterValidation("datetime", validateDateTime)
	v.RegisterValidation("numeric", validateNumeric)
	v.RegisterValidation("alpha", validateAlpha)
	v.RegisterValidation("alphanumeric", validateAlphaNumeric)
	v.RegisterValidation("unique", validateUnique)
	v.RegisterValidation("starts_with", validateStartsWith)
	v.RegisterValidation("ends_with", validateEndsWith)
	v.RegisterValidation("contains", validateContains)
	v.RegisterValidation("not_contains", validateNotContains)
}

// 自定义验证函数

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if !ValidateRequired(username) {
		return false
	}

	if !ValidateMinLength(username, MinUsernameLength) {
		return false
	}

	if !ValidateMaxLength(username, MaxUsernameLength) {
		return false
	}

	return ValidateRegex(username, RegexUsername)
}

// validatePassword 验证密码
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if !ValidateRequired(password) {
		return false
	}

	if !ValidateMinLength(password, MinPasswordLength) {
		return false
	}

	if !ValidateMaxLength(password, MaxPasswordLength) {
		return false
	}

	return ValidateRegex(password, RegexPassword)
}

// validateChinese 验证中文字符
func validateChinese(fl validator.FieldLevel) bool {
	return IsChinese(fl.Field().String())
}

// validateChineseName 验证中文姓名
func validateChineseName(fl validator.FieldLevel) bool {
	return IsChineseName(fl.Field().String())
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	return IsPhone(fl.Field().String())
}

// validateIDCard 验证身份证号
func validateIDCard(fl validator.FieldLevel) bool {
	return IsIDCard(fl.Field().String())
}

// validatePostalCode 验证邮政编码
func validatePostalCode(fl validator.FieldLevel) bool {
	return IsPostalCode(fl.Field().String())
}

// validateCreditCard 验证信用卡号
func validateCreditCard(fl validator.FieldLevel) bool {
	return IsCreditCard(fl.Field().String())
}

// validateStrongPassword 验证强密码
func validateStrongPassword(fl validator.FieldLevel) bool {
	return IsStrongPassword(fl.Field().String())
}

// validateURL 验证URL
func validateURL(fl validator.FieldLevel) bool {
	return IsURL(fl.Field().String())
}

// validateIPv4 验证IPv4地址
func validateIPv4(fl validator.FieldLevel) bool {
	return IsIPv4(fl.Field().String())
}

// validateIPv6 验证IPv6地址
func validateIPv6(fl validator.FieldLevel) bool {
	return IsIPv6(fl.Field().String())
}

// validateDate 验证日期
func validateDate(fl validator.FieldLevel) bool {
	return IsDate(fl.Field().String())
}

// validateTime 验证时间
func validateTime(fl validator.FieldLevel) bool {
	return IsTime(fl.Field().String())
}

// validateDateTime 验证日期时间
func validateDateTime(fl validator.FieldLevel) bool {
	return IsDateTime(fl.Field().String())
}

// validateNumeric 验证数字
func validateNumeric(fl validator.FieldLevel) bool {
	return IsNumeric(fl.Field().String())
}

// validateAlpha 验证字母
func validateAlpha(fl validator.FieldLevel) bool {
	return IsAlpha(fl.Field().String())
}

// validateAlphaNumeric 验证字母数字
func validateAlphaNumeric(fl validator.FieldLevel) bool {
	return IsAlphaNumeric(fl.Field().String())
}

// validateUnique 验证唯一性（需要数据库查询，这里只是占位符）
func validateUnique(fl validator.FieldLevel) bool {
	// 这个验证需要数据库查询，在实际使用时需要特殊处理
	// 这里返回 true 作为占位符
	return true
}

// validateStartsWith 验证是否以指定字符串开头
func validateStartsWith(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param()
	return ValidateStartsWith(value, param)
}

// validateEndsWith 验证是否以指定字符串结尾
func validateEndsWith(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param()
	return ValidateEndsWith(value, param)
}

// validateContains 验证是否包含指定字符串
func validateContains(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param()
	return ValidateContains(value, param)
}

// validateNotContains 验证是否不包含指定字符串
func validateNotContains(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param()
	return ValidateNotContains(value, param)
}

// GetValidationErrors 获取验证错误信息
func GetValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			error := ValidationError{
				Field:   e.Field(),
				Tag:     e.Tag(),
				Value:   fmt.Sprintf("%v", e.Value()),
				Message: getErrorMessage(e),
			}
			errors = append(errors, error)
		}
	}

	return errors
}

// getErrorMessage 根据验证标签获取错误信息
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s不能为空", getFieldName(e.Field()))
	case "email":
		return "邮箱格式不正确"
	case "min":
		return fmt.Sprintf("%s长度不能少于%s个字符", getFieldName(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s长度不能超过%s个字符", getFieldName(e.Field()), e.Param())
	case "username":
		return "用户名只能包含字母、数字，长度2-10个字符"
	case "password":
		return "密码只能包含字母、数字和特殊字符，长度6-18个字符"
	case "chinese":
		return "只能包含中文字符"
	case "chinese_name":
		return "请输入有效的中文姓名"
	case "phone":
		return "请输入有效的手机号码"
	case "idcard":
		return "请输入有效的身份证号码"
	case "postal_code":
		return "请输入有效的邮政编码"
	case "credit_card":
		return "请输入有效的信用卡号码"
	case "strong_password":
		return "密码必须包含大小写字母、数字和特殊字符，且长度不少于8位"
	case "url":
		return "请输入有效的URL地址"
	case "ipv4":
		return "请输入有效的IPv4地址"
	case "ipv6":
		return "请输入有效的IPv6地址"
	case "date":
		return "请输入有效的日期格式（YYYY-MM-DD）"
	case "time":
		return "请输入有效的时间格式（HH:MM:SS）"
	case "datetime":
		return "请输入有效的日期时间格式（YYYY-MM-DD HH:MM:SS）"
	case "numeric":
		return "只能包含数字"
	case "alpha":
		return "只能包含字母"
	case "alphanumeric":
		return "只能包含字母和数字"
	case "starts_with":
		return fmt.Sprintf("必须以%s开头", e.Param())
	case "ends_with":
		return fmt.Sprintf("必须以%s结尾", e.Param())
	case "contains":
		return fmt.Sprintf("必须包含%s", e.Param())
	case "not_contains":
		return fmt.Sprintf("不能包含%s", e.Param())
	default:
		return fmt.Sprintf("%s验证失败", getFieldName(e.Field()))
	}
}

// getFieldName 获取字段的中文名称
func getFieldName(field string) string {
	fieldNames := map[string]string{
		"Name":     "用户名",
		"Email":    "邮箱",
		"Password": "密码",
		"Phone":    "手机号",
		"IDCard":   "身份证号",
		"Address":  "地址",
		"Title":    "标题",
		"Content":  "内容",
		"URL":      "链接",
		"IP":       "IP地址",
		"Date":     "日期",
		"Time":     "时间",
	}

	if name, exists := fieldNames[field]; exists {
		return name
	}

	return field
}

// SetupEchoValidator 设置 Echo 验证器
func SetupEchoValidator(e *echo.Echo) {
	e.Validator = NewCustomValidator()
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	v := NewCustomValidator()
	return v.Validate(s)
}

// ValidateStructWithErrors 验证结构体并返回详细错误信息
func ValidateStructWithErrors(s interface{}) []ValidationError {
	err := ValidateStruct(s)
	if err != nil {
		return GetValidationErrors(err)
	}
	return nil
}

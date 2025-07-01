package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Validator 验证器结构体
type Validator struct {
	Errors []string
}

// NewValidator 创建新的验证器实例
func NewValidator() *Validator {
	return &Validator{
		Errors: make([]string, 0),
	}
}

// AddError 添加错误信息
func (v *Validator) AddError(field, message string) {
	v.Errors = append(v.Errors, fmt.Sprintf("%s: %s", field, message))
}

// HasErrors 检查是否有错误
func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}

// GetErrors 获取所有错误信息
func (v *Validator) GetErrors() []string {
	return v.Errors
}

// ClearErrors 清除所有错误
func (v *Validator) ClearErrors() {
	v.Errors = make([]string, 0)
}

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsPhone 验证手机号格式（中国大陆）
func IsPhone(phone string) bool {
	// 支持 11 位手机号，以 1 开头
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsIDCard 验证身份证号格式（中国大陆）
func IsIDCard(idCard string) bool {
	// 18位身份证号验证
	if len(idCard) != 18 {
		return false
	}

	// 前17位必须是数字
	pattern := `^\d{17}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	if !matched {
		return false
	}

	// 验证校验码
	return validateIDCardChecksum(idCard)
}

// validateIDCardChecksum 验证身份证校验码
func validateIDCardChecksum(idCard string) bool {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i := 0; i < 17; i++ {
		digit, _ := strconv.Atoi(string(idCard[i]))
		sum += digit * weights[i]
	}

	checkIndex := sum % 11
	expectedCheck := checkCodes[checkIndex]
	actualCheck := strings.ToUpper(string(idCard[17]))

	return expectedCheck == actualCheck
}

// IsStrongPassword 验证密码强度
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// IsURL 验证URL格式
func IsURL(url string) bool {
	pattern := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// IsIPv4 验证IPv4地址格式
func IsIPv4(ip string) bool {
	pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	matched, _ := regexp.MatchString(pattern, ip)
	return matched
}

// IsIPv6 验证IPv6地址格式
func IsIPv6(ip string) bool {
	pattern := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`
	matched, _ := regexp.MatchString(pattern, ip)
	return matched
}

// IsDate 验证日期格式 (YYYY-MM-DD)
func IsDate(date string) bool {
	pattern := `^\d{4}-\d{2}-\d{2}$`
	matched, _ := regexp.MatchString(pattern, date)
	if !matched {
		return false
	}

	// 进一步验证日期有效性
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return false
	}

	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])

	if year < 1900 || year > 2100 {
		return false
	}

	if month < 1 || month > 12 {
		return false
	}

	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// 处理闰年
	if month == 2 && (year%4 == 0 && year%100 != 0 || year%400 == 0) {
		daysInMonth[1] = 29
	}

	if day < 1 || day > daysInMonth[month-1] {
		return false
	}

	return true
}

// IsTime 验证时间格式 (HH:MM:SS)
func IsTime(time string) bool {
	pattern := `^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`
	matched, _ := regexp.MatchString(pattern, time)
	return matched
}

// IsDateTime 验证日期时间格式 (YYYY-MM-DD HH:MM:SS)
func IsDateTime(datetime string) bool {
	pattern := `^\d{4}-\d{2}-\d{2} ([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`
	matched, _ := regexp.MatchString(pattern, datetime)
	if !matched {
		return false
	}

	parts := strings.Split(datetime, " ")
	if len(parts) != 2 {
		return false
	}

	return IsDate(parts[0]) && IsTime(parts[1])
}

// IsNumeric 验证是否为数字字符串
func IsNumeric(str string) bool {
	pattern := `^\d+$`
	matched, _ := regexp.MatchString(pattern, str)
	return matched
}

// IsAlpha 验证是否为字母字符串
func IsAlpha(str string) bool {
	pattern := `^[a-zA-Z]+$`
	matched, _ := regexp.MatchString(pattern, str)
	return matched
}

// IsAlphaNumeric 验证是否为字母数字组合
func IsAlphaNumeric(str string) bool {
	pattern := `^[a-zA-Z0-9]+$`
	matched, _ := regexp.MatchString(pattern, str)
	return matched
}

// IsChinese 验证是否包含中文字符
func IsChinese(str string) bool {
	pattern := `[一-龯]`
	matched, _ := regexp.MatchString(pattern, str)
	return matched
}

// IsChineseName 验证中文姓名
func IsChineseName(name string) bool {
	pattern := `^[一-龯]{2,4}$`
	matched, _ := regexp.MatchString(pattern, name)
	return matched
}

// IsPostalCode 验证邮政编码（中国大陆）
func IsPostalCode(code string) bool {
	pattern := `^[1-9]\d{5}$`
	matched, _ := regexp.MatchString(pattern, code)
	return matched
}

// IsCreditCard 验证信用卡号（Luhn算法）
func IsCreditCard(cardNumber string) bool {
	// 移除空格和连字符
	cardNumber = regexp.MustCompile(`[\s-]`).ReplaceAllString(cardNumber, "")

	// 检查长度
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	// 检查是否全为数字
	if !IsNumeric(cardNumber) {
		return false
	}

	// Luhn算法验证
	sum := 0
	alternate := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cardNumber[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// ValidateLength 验证字符串长度
func ValidateLength(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}

// ValidateRange 验证数值范围
func ValidateRange(value, min, max float64) bool {
	return value >= min && value <= max
}

// ValidateRequired 验证必填字段
func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

// ValidateMinLength 验证最小长度
func ValidateMinLength(str string, min int) bool {
	return len(str) >= min
}

// ValidateMaxLength 验证最大长度
func ValidateMaxLength(str string, max int) bool {
	return len(str) <= max
}

// ValidateMin 验证最小值
func ValidateMin(value, min float64) bool {
	return value >= min
}

// ValidateMax 验证最大值
func ValidateMax(value, max float64) bool {
	return value <= max
}

// ValidateIn 验证值是否在指定列表中
func ValidateIn(value string, allowedValues []string) bool {
	for _, allowed := range allowedValues {
		if value == allowed {
			return true
		}
	}
	return false
}

// ValidateNotIn 验证值是否不在指定列表中
func ValidateNotIn(value string, forbiddenValues []string) bool {
	for _, forbidden := range forbiddenValues {
		if value == forbidden {
			return false
		}
	}
	return true
}

// ValidateRegex 验证正则表达式匹配
func ValidateRegex(value, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// ValidateUnique 验证字符串中字符是否唯一
func ValidateUnique(str string) bool {
	charMap := make(map[rune]bool)
	for _, char := range str {
		if charMap[char] {
			return false
		}
		charMap[char] = true
	}
	return true
}

// ValidateStartsWith 验证字符串是否以指定前缀开头
func ValidateStartsWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// ValidateEndsWith 验证字符串是否以指定后缀结尾
func ValidateEndsWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// ValidateContains 验证字符串是否包含指定子串
func ValidateContains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// ValidateNotContains 验证字符串是否不包含指定子串
func ValidateNotContains(str, substr string) bool {
	return !strings.Contains(str, substr)
}

package utils

import (
	"testing"
)

func TestIsEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@sub.example.com", true},
		{"valid email with plus", "test+tag@example.com", true},
		{"invalid email no @", "testexample.com", false},
		{"invalid email no domain", "test@", false},
		{"invalid email no local", "@example.com", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmail(tt.email); got != tt.expected {
				t.Errorf("IsEmail() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsPhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"valid phone", "13800138000", true},
		{"valid phone 13", "13000138000", true},
		{"valid phone 15", "15000138000", true},
		{"valid phone 18", "18000138000", true},
		{"valid phone 19", "19000138000", true},
		{"invalid phone too short", "1380013800", false},
		{"invalid phone too long", "138001380000", false},
		{"invalid phone wrong prefix", "12800138000", false},
		{"invalid phone with letters", "1380013800a", false},
		{"empty phone", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPhone(tt.phone); got != tt.expected {
				t.Errorf("IsPhone() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsIDCard(t *testing.T) {
	tests := []struct {
		name     string
		idCard   string
		expected bool
	}{
		{"valid id card", "110101199001011234", true},
		{"valid id card with X", "11010119900101123X", true},
		{"invalid id card too short", "11010119900101123", false},
		{"invalid id card too long", "1101011990010112345", false},
		{"invalid id card with letters", "11010119900101123A", false},
		{"empty id card", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIDCard(tt.idCard); got != tt.expected {
				t.Errorf("IsIDCard() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"strong password", "Test123!@#", true},
		{"strong password with symbols", "MyPass123$", true},
		{"weak password no special", "Test12345", false},
		{"weak password no number", "TestPass!@#", false},
		{"weak password no upper", "test123!@#", false},
		{"weak password no lower", "TEST123!@#", false},
		{"weak password too short", "Test1!", false},
		{"empty password", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStrongPassword(tt.password); got != tt.expected {
				t.Errorf("IsStrongPassword() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"valid http url", "http://example.com", true},
		{"valid https url", "https://example.com", true},
		{"valid url with path", "https://example.com/path", true},
		{"valid url with query", "https://example.com?param=value", true},
		{"invalid url no protocol", "example.com", false},
		{"invalid url wrong protocol", "ftp://example.com", false},
		{"empty url", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURL(tt.url); got != tt.expected {
				t.Errorf("IsURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"valid ipv4", "192.168.1.1", true},
		{"valid ipv4 max values", "255.255.255.255", true},
		{"valid ipv4 min values", "0.0.0.0", true},
		{"invalid ipv4 out of range", "256.256.256.256", false},
		{"invalid ipv4 format", "192.168.1", false},
		{"invalid ipv4 with letters", "192.168.1.a", false},
		{"empty ip", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv4(tt.ip); got != tt.expected {
				t.Errorf("IsIPv4() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsDate(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected bool
	}{
		{"valid date", "2023-12-25", true},
		{"valid leap year date", "2024-02-29", true},
		{"invalid leap year date", "2023-02-29", false},
		{"invalid month", "2023-13-01", false},
		{"invalid day", "2023-12-32", false},
		{"invalid format", "2023/12/25", false},
		{"invalid year too old", "1899-12-25", false},
		{"invalid year too new", "2101-12-25", false},
		{"empty date", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDate(tt.date); got != tt.expected {
				t.Errorf("IsDate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsTime(t *testing.T) {
	tests := []struct {
		name     string
		time     string
		expected bool
	}{
		{"valid time", "12:30:45", true},
		{"valid time midnight", "00:00:00", true},
		{"valid time max", "23:59:59", true},
		{"invalid hour", "24:00:00", false},
		{"invalid minute", "12:60:00", false},
		{"invalid second", "12:30:60", false},
		{"invalid format", "12:30", false},
		{"empty time", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTime(tt.time); got != tt.expected {
				t.Errorf("IsTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsDateTime(t *testing.T) {
	tests := []struct {
		name     string
		datetime string
		expected bool
	}{
		{"valid datetime", "2023-12-25 12:30:45", true},
		{"valid datetime midnight", "2023-12-25 00:00:00", true},
		{"invalid date", "2023-13-25 12:30:45", false},
		{"invalid time", "2023-12-25 24:00:00", false},
		{"invalid format", "2023-12-25T12:30:45", false},
		{"empty datetime", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDateTime(tt.datetime); got != tt.expected {
				t.Errorf("IsDateTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{"valid numeric", "12345", true},
		{"valid single digit", "5", true},
		{"invalid with letters", "12345a", false},
		{"invalid with spaces", "123 45", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumeric(tt.str); got != tt.expected {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{"valid alpha", "HelloWorld", true},
		{"valid single letter", "A", true},
		{"invalid with numbers", "Hello123", false},
		{"invalid with spaces", "Hello World", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlpha(tt.str); got != tt.expected {
				t.Errorf("IsAlpha() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{"valid alphanumeric", "Hello123", true},
		{"valid only letters", "Hello", true},
		{"valid only numbers", "12345", true},
		{"invalid with spaces", "Hello 123", false},
		{"invalid with special chars", "Hello@123", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlphaNumeric(tt.str); got != tt.expected {
				t.Errorf("IsAlphaNumeric() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsChinese(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{"valid chinese", "你好世界", true},
		{"valid mixed", "Hello你好", true},
		{"invalid no chinese", "Hello World", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChinese(tt.str); got != tt.expected {
				t.Errorf("IsChinese() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsChineseName(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{"valid chinese name 2 chars", "张三", true},
		{"valid chinese name 3 chars", "张三丰", true},
		{"valid chinese name 4 chars", "欧阳修", true},
		{"invalid too short", "张", false},
		{"invalid too long", "张三丰李四", false},
		{"invalid with numbers", "张三123", false},
		{"invalid with english", "张三John", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChineseName(tt.str); got != tt.expected {
				t.Errorf("IsChineseName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsPostalCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"valid postal code", "100000", true},
		{"valid postal code 2", "200000", true},
		{"invalid too short", "10000", false},
		{"invalid too long", "1000000", false},
		{"invalid starts with 0", "010000", false},
		{"invalid with letters", "10000a", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPostalCode(tt.code); got != tt.expected {
				t.Errorf("IsPostalCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsCreditCard(t *testing.T) {
	tests := []struct {
		name     string
		card     string
		expected bool
	}{
		{"valid visa", "4532015112830366", true},
		{"valid mastercard", "5425233430109903", true},
		{"valid with spaces", "4532 0151 1283 0366", true},
		{"valid with dashes", "4532-0151-1283-0366", true},
		{"invalid checksum", "4532015112830367", false},
		{"invalid too short", "453201511283036", false},
		{"invalid with letters", "453201511283036a", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCreditCard(tt.card); got != tt.expected {
				t.Errorf("IsCreditCard() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidator(t *testing.T) {
	validator := NewValidator()

	// Test initial state
	if validator.HasErrors() {
		t.Error("New validator should not have errors")
	}

	if len(validator.GetErrors()) != 0 {
		t.Error("New validator should have empty errors")
	}

	// Test adding errors
	validator.AddError("email", "邮箱格式不正确")
	validator.AddError("password", "密码长度不足")

	if !validator.HasErrors() {
		t.Error("Validator should have errors after adding them")
	}

	errors := validator.GetErrors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	// Test clearing errors
	validator.ClearErrors()
	if validator.HasErrors() {
		t.Error("Validator should not have errors after clearing")
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		min      int
		max      int
		expected bool
	}{
		{"valid length", "hello", 3, 10, true},
		{"exact min", "hi", 2, 10, true},
		{"exact max", "hello world", 3, 11, true},
		{"too short", "hi", 3, 10, false},
		{"too long", "hello world", 3, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLength(tt.str, tt.min, tt.max); got != tt.expected {
				t.Errorf("ValidateLength() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateRange(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected bool
	}{
		{"valid range", 5.0, 1.0, 10.0, true},
		{"exact min", 1.0, 1.0, 10.0, true},
		{"exact max", 10.0, 1.0, 10.0, true},
		{"below min", 0.5, 1.0, 10.0, false},
		{"above max", 15.0, 1.0, 10.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRange(tt.value, tt.min, tt.max); got != tt.expected {
				t.Errorf("ValidateRange() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid required", "hello", true},
		{"invalid empty", "", false},
		{"invalid only spaces", "   ", false},
		{"invalid tabs", "\t\t", false},
		{"invalid newlines", "\n\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRequired(tt.value); got != tt.expected {
				t.Errorf("ValidateRequired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateIn(t *testing.T) {
	allowed := []string{"red", "green", "blue"}

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid in list", "red", true},
		{"valid in list 2", "green", true},
		{"invalid not in list", "yellow", false},
		{"invalid empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateIn(tt.value, allowed); got != tt.expected {
				t.Errorf("ValidateIn() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateNotIn(t *testing.T) {
	forbidden := []string{"admin", "root", "system"}

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid not in list", "user", true},
		{"valid not in list 2", "guest", true},
		{"invalid in list", "admin", false},
		{"invalid in list 2", "root", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateNotIn(tt.value, forbidden); got != tt.expected {
				t.Errorf("ValidateNotIn() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateRegex(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		pattern  string
		expected bool
	}{

		{"valid regex match", "testuser哈哈", `^[a-zA-Z0-9\p{Han}_-]+$`, true},
		{"valid regex match", "hello123", `^[a-z]+\d+$`, true},
		{"invalid regex no match", "hello", `^[a-z]+\d+$`, false},
		{"valid email pattern", "test@example.com", `^[^@]+@[^@]+\.[^@]+$`, true},
		{"invalid email pattern", "test@example", `^[^@]+@[^@]+\.[^@]+$`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRegex(tt.value, tt.pattern); got != tt.expected {
				t.Errorf("ValidateRegex() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateUnique(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{"valid unique", "abcdef", true},
		{"invalid not unique", "hello", false},
		{"valid single char", "a", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateUnique(tt.str); got != tt.expected {
				t.Errorf("ValidateUnique() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateStartsWith(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		prefix   string
		expected bool
	}{
		{"valid starts with", "hello world", "hello", true},
		{"invalid doesn't start with", "world hello", "hello", false},
		{"valid exact match", "hello", "hello", true},
		{"invalid empty prefix", "hello", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateStartsWith(tt.str, tt.prefix); got != tt.expected {
				t.Errorf("ValidateStartsWith() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateEndsWith(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		suffix   string
		expected bool
	}{
		{"valid ends with", "hello world", "world", true},
		{"invalid doesn't end with", "world hello", "world", false},
		{"valid exact match", "hello", "hello", true},
		{"invalid empty suffix", "hello", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEndsWith(tt.str, tt.suffix); got != tt.expected {
				t.Errorf("ValidateEndsWith() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateContains(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected bool
	}{
		{"valid contains", "hello world", "world", true},
		{"valid contains at start", "hello world", "hello", true},
		{"invalid doesn't contain", "hello world", "test", false},
		{"valid empty substring", "hello", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateContains(tt.str, tt.substr); got != tt.expected {
				t.Errorf("ValidateContains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateNotContains(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected bool
	}{
		{"valid doesn't contain", "hello world", "test", true},
		{"invalid contains", "hello world", "world", false},
		{"invalid contains at start", "hello world", "hello", false},
		{"valid empty substring", "hello", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateNotContains(tt.str, tt.substr); got != tt.expected {
				t.Errorf("ValidateNotContains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

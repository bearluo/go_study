package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStruct 测试结构体
type TestStruct struct {
	Name     string `validate:"required,username"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,password"`
	Phone    string `validate:"phone"`
	IDCard   string `validate:"idcard"`
}

func TestCustomValidator_Validate(t *testing.T) {
	v := NewCustomValidator()

	tests := []struct {
		name    string
		data    TestStruct
		wantErr bool
	}{
		{
			name: "有效数据",
			data: TestStruct{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Phone:    "13800138000",
				IDCard:   "110101199001011234",
			},
			wantErr: false,
		},
		{
			name: "用户名为空",
			data: TestStruct{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "用户名格式错误",
			data: TestStruct{
				Name:     "test@user",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "邮箱格式错误",
			data: TestStruct{
				Name:     "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "密码太短",
			data: TestStruct{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "123",
			},
			wantErr: true,
		},
		{
			name: "手机号格式错误",
			data: TestStruct{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Phone:    "1234567890",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetValidationErrors(t *testing.T) {
	v := NewCustomValidator()

	// 测试无效数据
	data := TestStruct{
		Name:     "",
		Email:    "invalid-email",
		Password: "123",
	}

	err := v.Validate(data)
	assert.Error(t, err)

	// 获取验证错误
	errors := GetValidationErrors(err)
	assert.NotEmpty(t, errors)

	// 验证错误信息
	for _, e := range errors {
		assert.NotEmpty(t, e.Field)
		assert.NotEmpty(t, e.Tag)
		assert.NotEmpty(t, e.Message)
	}
}

func TestValidateStruct(t *testing.T) {
	data := TestStruct{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := ValidateStruct(data)
	assert.NoError(t, err)
}

func TestValidateStructWithErrors(t *testing.T) {
	// 测试有效数据
	validData := TestStruct{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	errors := ValidateStructWithErrors(validData)
	assert.Nil(t, errors)

	// 测试无效数据
	invalidData := TestStruct{
		Name:     "",
		Email:    "invalid-email",
		Password: "123",
	}

	errors = ValidateStructWithErrors(invalidData)
	assert.NotNil(t, errors)
	assert.Len(t, errors, 3) // 应该有3个验证错误
}

func TestCustomValidationTags(t *testing.T) {
	v := NewCustomValidator()

	// 测试自定义验证标签
	type CustomTestStruct struct {
		ChineseName string `validate:"chinese_name"`
		StrongPass  string `validate:"strong_password"`
		URL         string `validate:"url"`
		IPv4        string `validate:"ipv4"`
		Date        string `validate:"date"`
		Time        string `validate:"time"`
		DateTime    string `validate:"datetime"`
		Numeric     string `validate:"numeric"`
		Alpha       string `validate:"alpha"`
		AlphaNum    string `validate:"alphanumeric"`
	}

	validData := CustomTestStruct{
		ChineseName: "张三",
		StrongPass:  "Password123!",
		URL:         "https://example.com",
		IPv4:        "192.168.1.1",
		Date:        "2024-01-01",
		Time:        "12:00:00",
		DateTime:    "2024-01-01 12:00:00",
		Numeric:     "12345",
		Alpha:       "abcde",
		AlphaNum:    "abc123",
	}

	err := v.Validate(validData)
	assert.NoError(t, err)
}

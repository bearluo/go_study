package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupEchoContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestSuccess(t *testing.T) {
	c, rec := setupEchoContext()

	data := map[string]string{"key": "value"}
	err := Success(c, data, "操作成功")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeSuccess, response.Code)
	assert.Equal(t, "操作成功", response.Message)
	assert.Equal(t, map[string]interface{}{"key": "value"}, response.Data)
	assert.Empty(t, response.Error)
}

func TestError(t *testing.T) {
	c, rec := setupEchoContext()

	err := Error(c, CodeUserNotFound, "用户不存在")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeUserNotFound, response.Code)
	assert.Equal(t, "用户不存在", response.Message)
	assert.Equal(t, "用户不存在", response.Error)
	assert.Nil(t, response.Data)
}

func TestValidationError(t *testing.T) {
	c, rec := setupEchoContext()

	details := []string{"用户名不能为空", "邮箱格式不正确"}
	err := ValidationErrorResponse(c, details)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeValidationError, response.Code)
	assert.Equal(t, "参数验证失败", response.Message)
	assert.Equal(t, "参数验证失败", response.Error)
	assert.Equal(t, []interface{}{"用户名不能为空", "邮箱格式不正确"}, response.Details)
}

func TestSystemError(t *testing.T) {
	c, rec := setupEchoContext()

	testErr := echo.NewHTTPError(http.StatusInternalServerError, "数据库连接失败")
	err := SystemError(c, testErr)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeSystemError, response.Code)
	assert.Equal(t, "系统内部错误", response.Message)
	assert.Contains(t, response.Error, "数据库连接失败")
}

func TestNotFound(t *testing.T) {
	c, rec := setupEchoContext()

	err := NotFound(c, "用户不存在")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeNotFound, response.Code)
	assert.Equal(t, "用户不存在", response.Message)
	assert.Equal(t, "用户不存在", response.Error)
}

func TestUnauthorized(t *testing.T) {
	c, rec := setupEchoContext()

	err := Unauthorized(c, "Token无效")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeUnauthorized, response.Code)
	assert.Equal(t, "Token无效", response.Message)
	assert.Equal(t, "Token无效", response.Error)
}

func TestForbidden(t *testing.T) {
	c, rec := setupEchoContext()

	err := Forbidden(c, "权限不足")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeForbidden, response.Code)
	assert.Equal(t, "权限不足", response.Message)
	assert.Equal(t, "权限不足", response.Error)
}

func TestUserNotFound(t *testing.T) {
	c, rec := setupEchoContext()

	err := UserNotFound(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeUserNotFound, response.Code)
	assert.Equal(t, "用户不存在", response.Message)
}

func TestUserExists(t *testing.T) {
	c, rec := setupEchoContext()

	err := UserExists(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodeUserExists, response.Code)
	assert.Equal(t, "用户已存在", response.Message)
}

func TestPasswordError(t *testing.T) {
	c, rec := setupEchoContext()

	err := PasswordError(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, CodePasswordError, response.Code)
	assert.Equal(t, "密码错误", response.Message)
}

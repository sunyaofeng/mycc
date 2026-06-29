package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateHandler_GetCurrentDate(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 使用固定时间进行测试
	fixedTime := time.Date(2026, 6, 29, 14, 30, 0, 0, time.Local)
	mockNow := func() time.Time {
		return fixedTime
	}

	// 创建处理器
	handler := NewDateHandlerWithNow(mockNow)

	// 创建测试路由
	r := gin.New()
	r.GET("/api/v1/date", handler.GetCurrentDate)

	// 创建测试请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/date", nil)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应体
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 验证响应字段
	assert.Equal(t, "2026-06-29", response["date"])
	assert.Equal(t, "2026-06-29 14:30:00", response["time"])
	assert.Equal(t, float64(fixedTime.Unix()), response["timestamp"])
	assert.NotEmpty(t, response["timezone"])
}

func TestDateHandler_GetCurrentDate_RealTime(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 使用真实时间的处理器
	handler := NewDateHandler()

	// 创建测试路由
	r := gin.New()
	r.GET("/api/v1/date", handler.GetCurrentDate)

	// 创建测试请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/date", nil)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应体
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 验证响应字段存在且不为空
	assert.NotEmpty(t, response["date"])
	assert.NotEmpty(t, response["time"])
	assert.NotZero(t, response["timestamp"])
	assert.NotEmpty(t, response["timezone"])

	// 验证日期格式
	dateStr, ok := response["date"].(string)
	require.True(t, ok)
	_, parseErr := time.Parse("2006-01-02", dateStr)
	assert.NoError(t, parseErr, "date should be in YYYY-MM-DD format")
}

func TestNewDateHandler(t *testing.T) {
	handler := NewDateHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.now)
}

func TestNewDateHandlerWithNow(t *testing.T) {
	mockNow := func() time.Time {
		return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	handler := NewDateHandlerWithNow(mockNow)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.now)

	// 验证 mock 时间是否生效
	now := handler.now()
	assert.Equal(t, 2026, now.Year())
	assert.Equal(t, time.January, now.Month())
	assert.Equal(t, 1, now.Day())
}

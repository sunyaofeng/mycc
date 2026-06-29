package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mycc/internal/handler"
	"mycc/internal/logger"
)

func TestRouter_HealthEndpoint(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化日志（中间件需要）
	logger.Init("info", "json", "stdout", "")

	// 创建路由
	dateHandler := handler.NewDateHandler()
	healthHandler := handler.NewHealthHandler()
	r := NewRouter(dateHandler, healthHandler)

	// 创建测试请求
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// 执行请求
	r.Engine().ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应体
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 验证响应字段
	assert.Equal(t, "ok", response["status"])
	assert.NotEmpty(t, response["timestamp"])
}

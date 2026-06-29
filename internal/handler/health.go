package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct{}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// GetHealth 健康检查
// GET /health
func (h *HealthHandler) GetHealth(c *gin.Context) {
	now := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": now.Format("2006-01-02 15:04:05"),
	})
}

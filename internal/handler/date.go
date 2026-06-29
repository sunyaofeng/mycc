package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DateHandler 日期处理器
type DateHandler struct {
	now func() time.Time
}

// NewDateHandler 创建日期处理器
func NewDateHandler() *DateHandler {
	return &DateHandler{
		now: time.Now,
	}
}

// NewDateHandlerWithNow 创建可注入时间的日期处理器（用于测试）
func NewDateHandlerWithNow(now func() time.Time) *DateHandler {
	return &DateHandler{
		now: now,
	}
}

// GetCurrentDate 获取当前日期
// GET /api/v1/date
func (h *DateHandler) GetCurrentDate(c *gin.Context) {
	now := h.now()
	c.JSON(http.StatusOK, gin.H{
		"date": now.Format("2006-01-02"),
		"time": now.Format("2006-01-02 15:04:05"),
		"timestamp": now.Unix(),
		"timezone": now.Location().String(),
	})
}

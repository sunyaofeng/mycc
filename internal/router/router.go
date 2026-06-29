package router

import (
	"github.com/gin-gonic/gin"

	"mycc/internal/middleware"
)

// Setup 设置路由
func Setup() *gin.Engine {
	// 创建路由实例
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   "",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// v1 版本
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
		}
	}

	return r
}

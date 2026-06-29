package router

import (
	"github.com/gin-gonic/gin"

	"mycc/internal/handler"
	"mycc/internal/middleware"
)

// Router 路由结构体
type Router struct {
	engine       *gin.Engine
	dateHandler  *handler.DateHandler
}

// NewRouter 创建路由实例（Wire 注入）
func NewRouter(dateHandler *handler.DateHandler) *Router {
	r := &Router{
		engine:      gin.New(),
		dateHandler: dateHandler,
	}

	// 注册全局中间件
	r.engine.Use(middleware.Logger())
	r.engine.Use(middleware.Recovery())
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.RequestID())

	// 注册路由
	r.registerRoutes()

	return r
}

// registerRoutes 注册所有路由
func (r *Router) registerRoutes() {
	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   "",
		})
	})

	// API 路由组
	api := r.engine.Group("/api")
	{
		// v1 版本
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			// 日期接口
			v1.GET("/date", r.dateHandler.GetCurrentDate)
		}
	}
}

// Engine 获取 Gin 引擎实例
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// Run 启动服务器
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

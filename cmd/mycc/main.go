package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mycc/internal/config"
	"mycc/internal/env"
	"mycc/internal/logger"
	"mycc/internal/router"
)

func main() {
	// 获取环境
	environment := env.GetEnv()
	fmt.Printf("Starting server in %s environment\n", environment)

	// 加载配置
	cfg, err := config.Load(environment.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Log.FilePath,
	); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Server starting",
		zap.String("env", environment.String()),
		zap.Int("port", cfg.Server.Port),
	)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 设置路由
	r := router.Setup()

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Server listening", zap.String("addr", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

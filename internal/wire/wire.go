package wire

import (
	"github.com/google/wire"

	"mycc/internal/handler"
	"mycc/internal/router"
)

// AppSet 应用依赖集合
var AppSet = wire.NewSet(
	// 处理器
	handler.NewDateHandler,

	// 路由
	router.NewRouter,
)

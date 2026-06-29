//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"mycc/internal/handler"
	"mycc/internal/router"
)

// InitializeRouter 使用 Wire 自动注入依赖创建 Router
func InitializeRouter() *router.Router {
	wire.Build(
		handler.NewDateHandler,
		router.NewRouter,
	)
	return &router.Router{}
}

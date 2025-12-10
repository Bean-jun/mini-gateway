package main

import (
	"context"
	"log"

	"github.com/Bean-jun/mini-gateway/config"
	"github.com/Bean-jun/mini-gateway/engine"
)

// Hnadler 启动对应协议的服务引擎
func Hnadler(ctx context.Context, config *config.ServerBlock, sign chan<- struct{}) {
	defer func() {
		// 捕获 panic 错误
		err := recover()
		if err != nil {
			log.Println("handler panic error:", err)
		}

		// 通知主服务, 该服务已退出
		log.Printf("Handler mini-gateway on protocol: %s port: %d exited!", config.Protocol, config.Port)
		sign <- struct{}{}
	}()

	var e engine.Engine
	switch config.Protocol {
	case "http":
		// 启动 HTTP 服务
		e = engine.NewHttpEngine(config)
	case "https":
		// 启动 HTTPS 服务
		e = engine.NewHttpEngine(config)
	case "tcp":
		// 启动 TCP 服务
		e = engine.NewTcpEngine(config)
	default:
		// 启动 其他协议 服务
		panic("protocol not implemented")
	}

	if e == nil {
		panic("engine is nil")
	}
	// 执行服务，阻塞等待，有问题就 panic
	err := e.Run(ctx)
	if err != nil {
		panic(err)
	}
}

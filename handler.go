package main

import (
	"context"
	"log"
)

// Handler 启动对应协议的服务引擎
func Handler(ctx context.Context, config *ServerBlock, sign chan<- struct{}) {
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

	var e Engine
	switch config.Protocol {
	case "http":
		// 启动 HTTP 服务
		e = NewHttpEngine(config)
	case "https":
		// 启动 HTTPS 服务
		e = NewHttpEngine(config)
	default:
		// 启动 其他协议 服务
		log.Printf("Unsupported protocol: %s", config.Protocol)
		return
	}

	if e == nil {
		log.Println("Engine creation failed: engine is nil")
		return
	}
	// 执行服务，阻塞等待，有问题就记录日志
	err := e.Run(ctx)
	if err != nil {
		log.Printf("Server run failed: %v", err)
		return
	}
}

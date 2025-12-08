package main

import (
	"fmt"
	"log"
	"net/http"
)

func Hnadler(config *ServerBlock, sign chan<- struct{}) {
	defer func() {
		// 捕获 panic 错误
		err := recover()
		if err != nil {
			log.Println("handler panic error:", err)
		}

		// 通知主服务, 该服务已退出
		log.Printf("Handler g-server on port %d exited!", config.Port)
		sign <- struct{}{}
	}()

	// 启动服务日志
	log.Printf("Handler g-server ... port: %d", config.Port)
	handerSrv := NewEngine(config.Port)

	// 执行服务，阻塞等待，有问题就 panic
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handerSrv)
	if err != nil {
		panic(err)
	}
}

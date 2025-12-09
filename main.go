package main

import (
	"context"
	"log"

	"github.com/Bean-jun/mini-gateway/config"
)

func main() {
	// 加载配置
	config, err := config.LoadConfigFromYAML("config.yaml")
	if err != nil {
		log.Fatalln("Load config error:", err)
		panic(err)
	}

	log.Println("Starting mini-gateway has ->", len(config.ServerBlocks), "servers...")

	stopSign := make(chan struct{}, len(config.ServerBlocks))
	ctx := context.Background()

	// 启动服务器
	for _, server_block := range config.ServerBlocks {
		go Hnadler(ctx, server_block, stopSign)
	}

	// 等待所有服务启动完成
	for i := 0; i < len(config.ServerBlocks); i++ {
		<-stopSign
		log.Println("stop server on: ", config.ServerBlocks[i].Port)
	}
	log.Println("stop all mini-gateway services!")
}

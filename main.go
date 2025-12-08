package main

import (
	"log"
)

func main() {
	// 加载配置
	config, err := LoadConfigFromYAML("config.yaml")
	if err != nil {
		log.Fatalln("Load config error:", err)
		panic(err)
	}

	log.Println("Starting mini-gateway has ->", len(config.ServerBlocks), "servers...")

	stopSign := make(chan struct{}, len(config.ServerBlocks))

	// 启动服务器
	for _, server_block := range config.ServerBlocks {
		go Hnadler(server_block, stopSign)
	}

	// 等待所有服务启动完成
	for i := 0; i < len(config.ServerBlocks); i++ {
		<-stopSign
		log.Println("stop server on: ", config.ServerBlocks[i].Port)
	}
	log.Println("stop all mini-gateway services!")
}

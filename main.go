package main

import (
	"log"
)

func main() {
	// 加载配置
	ConfigBlocks, err := LoadConfigFromYAML("config.yaml")
	if err != nil {
		log.Fatalln("Load config error:", err)
		panic(err)
	}

	log.Println("Starting g-server has ->", len(ConfigBlocks.Config), "servers...")

	stopSign := make(chan struct{}, len(ConfigBlocks.Config))

	// 启动服务器
	for _, config := range ConfigBlocks.Config {
		go Hnadler(config, stopSign)
	}

	// 等待所有服务启动完成
	for i := 0; i < len(ConfigBlocks.Config); i++ {
		<-stopSign
		log.Println("stop server on: ", ConfigBlocks.Config[i].Port)
	}
	log.Println("stop all g-server services!")
}

package engine

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Bean-jun/mini-gateway/config"
)

type HttpEngine struct {
	config *config.ServerBlock
}

// NewHttpEngine 创建 HTTP 服务引擎
func NewHttpEngine(serverBlock *config.ServerBlock) *HttpEngine {
	return &HttpEngine{config: serverBlock}
}

// Run 启动 HTTP 服务引擎
func (e *HttpEngine) Run(ctx context.Context) error {
	// 启动服务日志
	log.Printf("Handler mini-gateway HTTP server is running on port: %d", e.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", e.config.Port), e)
}

// TODO: 处理请求
func (e *HttpEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is mini-gateway! on port: %d\n", e.config.Port)
	if r.URL.Query().Get("ping") == "true" {
		fmt.Fprintf(w, "pong\n")
		panic("error on ping")
	}
}

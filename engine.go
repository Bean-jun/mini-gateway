package main

import (
	"fmt"
	"net/http"
)

type Engine struct {
	Port int
}

func NewEngine(port int) *Engine {
	return &Engine{Port: port}
}

// TODO: 处理请求
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is mini-gateway! on port: %d\n", e.Port)
	if r.URL.Query().Get("ping") == "true" {
		fmt.Fprintf(w, "pong\n")
		panic("error on ping")
	}
}

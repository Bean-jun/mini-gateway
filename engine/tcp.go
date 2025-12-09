package engine

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Bean-jun/mini-gateway/config"
)

type TcpEngine struct {
	config *config.ServerBlock
}

func NewTcpEngine(config *config.ServerBlock) *TcpEngine {
	return &TcpEngine{
		config: config,
	}
}

func (e *TcpEngine) Run(ctx context.Context) error {
	// 监听 TCP 端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", e.config.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Handler mini-gateway TCP server is running on port: %d", e.config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go e.handleConnection(ctx, conn)
	}
}

func (e *TcpEngine) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	// 处理连接
	log.Printf("Handler mini-gateway TCP server received connection from: %s", conn.RemoteAddr().String())
	select {
	case <-ctx.Done():
		return
	default:
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			// 打印客户端数据
			fmt.Printf("Received: %s", buf[:n])
			// 回复客户端
			conn.Write([]byte("Hello, client!"))
		}
	}
}

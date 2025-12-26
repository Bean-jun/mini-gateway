package main

import "context"

// Engine 服务引擎接口
type Engine interface {
	Run(ctx context.Context) error
}

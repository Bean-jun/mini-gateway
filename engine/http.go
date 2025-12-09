package engine

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Bean-jun/mini-gateway/config"
)

type LoadBalancer struct {
	proxy        []*config.ServerBlockLocationProxyPass
	currentIndex int // 当前执行的反向代理地址索引
	maxWeight    int // 最大权重
}

func NewLoadBalancer(proxy []*config.ServerBlockLocationProxyPass) *LoadBalancer {
	lb := &LoadBalancer{
		proxy:        proxy,
		currentIndex: 0,
		maxWeight:    0,
	}

	if len(proxy) != 0 {
		log.Printf("LoadBalancer initialized with %d proxy passes", len(proxy))
		weight := 0
		for _, p := range proxy {
			weight += p.Weight
		}
		lb.maxWeight = weight
	}
	return lb
}

func (lb *LoadBalancer) GetNextProxy() *config.ServerBlockLocationProxyPass {
	if len(lb.proxy) == 0 {
		return nil
	}

	if lb.maxWeight < 0 {
		lb.maxWeight = 0
	}

	if lb.maxWeight == 0 {
		proxy := lb.proxy[lb.currentIndex]
		lb.currentIndex = (lb.currentIndex + 1) % len(lb.proxy)
		return proxy
	}

	// 权重轮询算法
	// 1,3,4,4
	// 转换一下->1,4,8,12
	// 添加一个<=12的计数器，然后判断大小取
	// 1 2,3,4 5,6,7,8 9,10,11,12
	//
	// 2,5,3
	// ->2,7,10
	// 1,2, 3,4,5,6,7 8,9,10
	weightCounter := 1 + (lb.currentIndex % lb.maxWeight)
	accumulatedWeight := 0
	for _, p := range lb.proxy {
		accumulatedWeight += p.Weight
		if accumulatedWeight >= weightCounter {
			proxy := p
			lb.currentIndex = (lb.currentIndex + 1) % lb.maxWeight
			return proxy
		}
	}
	return nil
}

type HttpEngine struct {
	config   *config.ServerBlock
	patterns map[*regexp.Regexp]*LoadBalancer // 存储已注册的路径
}

// NewHttpEngine 创建 HTTP 服务引擎
func NewHttpEngine(serverBlock *config.ServerBlock) *HttpEngine {
	e := &HttpEngine{config: serverBlock, patterns: make(map[*regexp.Regexp]*LoadBalancer, len(serverBlock.Locations))}
	// 注册所有 Location 配置
	for _, location := range e.config.Locations {
		e.patterns[regexp.MustCompile(location.Path)] = NewLoadBalancer(location.ProxyPass)
	}
	return e
}

// Run 启动 HTTP 服务引擎
func (e *HttpEngine) Run(ctx context.Context) error {
	// 启动服务日志
	log.Printf("Handler mini-gateway HTTP server is running on port: %d", e.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", e.config.Port), e)
}

// TODO: 处理请求
func (e *HttpEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 遍历所有 Location 配置
	log.Println("url:", r.URL.Path, "from:", r.RemoteAddr, "RawQuery:", r.URL.RawQuery)

	var found bool
	var balancer *LoadBalancer

	for pattern, loadBalancer := range e.patterns {
		// 检查请求路径是否匹配当前 Location 的正则表达式
		if pattern.MatchString(r.URL.Path) {
			found = true
			// 找到匹配的 Location，处理请求
			balancer = loadBalancer
			break
		}
	}

	// 如果没有找到匹配的 Location，返回 404 错误
	if !found {
		http.NotFound(w, r)
		return
	}

	target := balancer.GetNextProxy()
	// 找到对应的反向代理地址
	if target == nil {
		http.Error(w, "No proxy pass configured for this location", http.StatusBadGateway)
		return
	}

	// 通过反向代理地址处理请求
	log.Printf("Proxying request %s to %s:%d", r.URL.Path, target.Host, target.Port)

	// r.URL.Path = e.config.ProxyPass + r.URL.Path
	fmt.Fprintf(w, "Hello, this is mini-gateway! on port: %d\n", e.config.Port)
	if r.URL.Query().Get("ping") == "true" {
		fmt.Fprintf(w, "pong\n")
		panic("error on ping")
	}
}

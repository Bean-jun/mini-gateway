package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

const (
	defaultMaxBodySize int64 = 64 * 1024 * 1024 // 默认最大请求体大小64MB
)

type ProxyHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// NewHttpProxyHander 创建Http反向代理处理器
func NewHttpProxyHander(schema, host string, port int) ProxyHandler {
	url := &url.URL{
		Scheme: schema,
		Host:   fmt.Sprintf("%s:%d", host, port),
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Server", "mini-gateway")
		log.Println("resp.ContentLength:", resp.ContentLength)
		return nil
	}
	return proxy
}

// NewFileProxyHandler 创建静态文件处理器
func NewFileProxyHandler(root string) ProxyHandler {
	fileServer := http.FileServer(http.Dir(root))
	return fileServer
}

type LoadBalancer struct {
	proxyPass    []*ServerBlockLocationProxyPass
	filePass     string // 静态文件根目录
	currentIndex int    // 当前执行的反向代理地址索引
	maxWeight    int    // 最大权重
}

func NewLoadBalancer(location *ServerBlockLocation) *LoadBalancer {
	lb := &LoadBalancer{
		proxyPass:    location.ProxyPass,
		filePass:     location.Root,
		currentIndex: 0,
		maxWeight:    0,
	}

	if location.ProxyPass != nil {
		log.Printf("LoadBalancer initialized with %d proxy passes", len(location.ProxyPass))
		weight := 0
		for _, p := range location.ProxyPass {
			weight += p.Weight
		}
		lb.maxWeight = weight
	}
	return lb
}

func (lb *LoadBalancer) GetNextProxy() ProxyHandler {
	if len(lb.proxyPass) == 0 && lb.filePass == "" {
		return nil
	}

	if lb.filePass != "" {
		return NewFileProxyHandler(lb.filePass)
	}

	if lb.maxWeight < 0 {
		lb.maxWeight = 0
	}

	if lb.maxWeight == 0 {
		proxy := lb.proxyPass[lb.currentIndex]
		lb.currentIndex = (lb.currentIndex + 1) % len(lb.proxyPass)
		return NewHttpProxyHander(proxy.Schema, proxy.Host, proxy.Port)
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
	for _, p := range lb.proxyPass {
		accumulatedWeight += p.Weight
		if accumulatedWeight >= weightCounter {
			proxy := p
			lb.currentIndex = (lb.currentIndex + 1) % lb.maxWeight
			return NewHttpProxyHander(proxy.Schema, proxy.Host, proxy.Port)
		}
	}
	return nil
}

type Pattern struct {
	loadBalancer []*LoadBalancer
	matchs       map[*regexp.Regexp]int
}

func NewPattern(locations []*ServerBlockLocation) *Pattern {
	p := &Pattern{
		loadBalancer: make([]*LoadBalancer, len(locations)),
		matchs:       make(map[*regexp.Regexp]int, len(locations)),
	}

	for idx, location := range locations {
		p.loadBalancer[idx] = NewLoadBalancer(location)
		regex := regexp.MustCompile(location.Path)
		p.matchs[regex] = idx
	}

	return p
}

func (p *Pattern) MatchString(r *http.Request) (*LoadBalancer, bool) {
	Index := -1
	for regex, idx := range p.matchs {
		if regex.MatchString(r.URL.Path) {
			Index = idx
			break
		}
	}
	// 如果没有找到匹配的 loadBalancer
	if Index == -1 {
		return nil, false
	}
	return p.loadBalancer[Index], true
}

type HttpEngine struct {
	Schema      string     // 协议
	Port        int        // 端口
	SSL         *SSLConfig // SSL 配置
	MaxBodySize int64      // 最大请求体大小
	pattern     *Pattern   // 路径模式匹配
}

// NewHttpEngine 创建 HTTP 服务引擎
func NewHttpEngine(serverBlock *ServerBlock) *HttpEngine {
	max_body_size := serverBlock.MaxBodySize * 1024 * 1024
	// 如果没有配置最大请求体大小，使用默认值
	if max_body_size == 0 {
		max_body_size = defaultMaxBodySize
	}
	e := &HttpEngine{
		Schema:      serverBlock.Protocol,
		Port:        serverBlock.Port,
		SSL:         serverBlock.SSL,
		MaxBodySize: max_body_size,
		pattern:     NewPattern(serverBlock.Locations),
	}
	return e
}

// Run 启动 HTTP 服务引擎
func (e *HttpEngine) Run(ctx context.Context) error {
	// 启动服务日志
	log.Printf("Handler mini-gateway HTTP server is running on %s:%d", e.Schema, e.Port)

	switch e.Schema {
	case "https":
		if e.SSL == nil {
			return fmt.Errorf("https protocol requires SSL configuration")
		}
		if e.SSL.CertFile == "" || e.SSL.KeyFile == "" {
			return fmt.Errorf("https protocol requires both cert_file and key_file in SSL configuration")
		}
		return http.ListenAndServeTLS(fmt.Sprintf(":%d", e.Port), e.SSL.CertFile, e.SSL.KeyFile, e)
	default:
		return http.ListenAndServe(fmt.Sprintf(":%d", e.Port), e)
	}
}

// TODO: 处理请求
func (e *HttpEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 遍历所有 Location 配置
	log.Println("url:", r.URL.Path, "from:", r.RemoteAddr, "RawQuery:", r.URL.RawQuery)

	// 检查请求体大小是否超过最大限制
	if r.ContentLength > e.MaxBodySize {
		http.Error(w, "Request body size exceeds maximum limit", http.StatusRequestEntityTooLarge)
		return
	}

	lb, ok := e.pattern.MatchString(r)
	// 如果没有找到匹配的 loadBalancer，返回 404 错误
	if !ok {
		http.NotFound(w, r)
		return
	}

	// 获取下一个反向代理地址
	handler := lb.GetNextProxy()
	// 找到对应的反向代理地址
	if handler == nil {
		http.Error(w, "No proxy pass configured for this location", http.StatusBadGateway)
		return
	}

	// 调用反向代理处理请求
	handler.ServeHTTP(w, r)
}

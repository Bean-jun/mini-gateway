package engine

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

type Pattern struct {
	loadBalancer []*LoadBalancer
	matchs       map[*regexp.Regexp]int
}

func NewPattern(locations []*config.ServerBlockLocation) *Pattern {
	p := &Pattern{
		loadBalancer: make([]*LoadBalancer, len(locations)),
		matchs:       make(map[*regexp.Regexp]int, len(locations)),
	}

	for idx, location := range locations {
		p.loadBalancer[idx] = NewLoadBalancer(location.ProxyPass)
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
	Schema  string            // 协议
	Port    int               // 端口
	SSL     *config.SSLConfig // SSL 配置
	pattern *Pattern          // 路径模式匹配
}

// NewHttpEngine 创建 HTTP 服务引擎
func NewHttpEngine(serverBlock *config.ServerBlock) *HttpEngine {
	e := &HttpEngine{
		Schema:  serverBlock.Protocol,
		Port:    serverBlock.Port,
		SSL:     serverBlock.SSL,
		pattern: NewPattern(serverBlock.Locations),
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
			panic("https protocol requires SSL configuration")
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

	lb, ok := e.pattern.MatchString(r)
	// 如果没有找到匹配的 loadBalancer，返回 404 错误
	if !ok {
		http.NotFound(w, r)
		return
	}

	// 获取下一个反向代理地址
	target := lb.GetNextProxy()
	// 找到对应的反向代理地址
	if target == nil {
		http.Error(w, "No proxy pass configured for this location", http.StatusBadGateway)
		return
	}

	// 通过反向代理地址处理请求
	log.Printf("Proxying request %s to %s:%d", r.URL.Path, target.Host, target.Port)

	url := &url.URL{
		Scheme: target.Schema,
		Host:   fmt.Sprintf("%s:%d", target.Host, target.Port),
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Server", "mini-gateway")
		log.Println("resp.ContentLength:", resp.ContentLength)
		return nil
	}
	// 转发请求
	proxy.ServeHTTP(w, r)
}

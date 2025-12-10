package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type SSLConfig struct {
	CertFile string `yaml:"cert_file"` // SSL 证书文件路径
	KeyFile  string `yaml:"key_file"`  // SSL 密钥文件路径
}

type ServerBlockLocationProxyPass struct {
	Schema string `yaml:"schema"` // 反向代理协议
	Host   string `yaml:"host"`   // 反向代理主机
	Port   int    `yaml:"port"`   // 反向代理端口
	Weight int    `yaml:"weight"` // 反向代理权重
}

// ServerBlockLocation 配置块位置
type ServerBlockLocation struct {
	Path      string                          `yaml:"path"`       // 匹配路径
	ProxyPass []*ServerBlockLocationProxyPass `yaml:"proxy_pass"` // 反向代理地址
}

// ServerBlock 配置块
type ServerBlock struct {
	Name        string                 `yaml:"name"`          // 服务名称
	Port        int                    `yaml:"port"`          // go-nginx 运行端口
	Protocol    string                 `yaml:"protocol"`      // 服务支持协议
	SSL         *SSLConfig             `yaml:"ssl"`           // SSL 配置
	MaxBodySize int64                  `yaml:"max_body_size"` // 最大请求体大小
	Locations   []*ServerBlockLocation `yaml:"locations"`     // 反向代理配置
}

// NewServerBlock 创建配置块
func NewServerBlock(opts ...func(*ServerBlock)) *ServerBlock {
	block := &ServerBlock{
		Name:     "default",
		Port:     7256,
		Protocol: "http",
	}
	for _, opt := range opts {
		opt(block)
	}
	return block
}

// ConfigBlock 配置块集合
type Config struct {
	ServerBlocks []*ServerBlock `yaml:"server_blocks"` // 配置列表
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		ServerBlocks: []*ServerBlock{
			{
				Name:     "default",
				Port:     7256,
				Protocol: "http",
			},
		},
	}
}

// AddServerBlock 添加配置块
func (cb *Config) AddServerBlock(opts ...func(*ServerBlock)) {
	cb.ServerBlocks = append(cb.ServerBlocks, NewServerBlock(opts...))
}

func LoadConfigFromYAML(path string) (*Config, error) {
	// 读取 yaml 文件
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// 解析 yaml 文件
	config, err := ParseYAMLConfig(yamlFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func ParseYAMLConfig(data []byte) (*Config, error) {
	// 解析 yaml 文件
	var config Config
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

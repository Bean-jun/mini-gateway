package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ServerBlock 配置块
type ServerBlock struct {
	Name string `yaml:"name"` // 服务名称
	Port int    `yaml:"port"` // go-nginx 运行端口
}

// NewServerBlock 创建配置块
func NewServerBlock(opts ...func(*ServerBlock)) *ServerBlock {
	block := &ServerBlock{
		Port: 7256,
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
				Name: "default",
				Port: 7256,
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

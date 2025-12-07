package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name string `yaml:"name"` // 服务名称
	Port int    `yaml:"port"` // go-nginx 运行端口
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig(opts ...func(*Config)) *Config {
	conf := &Config{
		Port: 7256,
	}
	for _, opt := range opts {
		opt(conf)
	}
	return conf
}

// ConfigBlock 配置块
type ConfigBlock struct {
	Config []*Config `yaml:"configs"` // 配置列表
}

// NewConfigBlock 创建配置块
func NewConfigBlock(opts ...func(*Config)) *ConfigBlock {
	return &ConfigBlock{
		Config: []*Config{NewDefaultConfig(opts...)},
	}
}

func (cb *ConfigBlock) AddConfig(opts ...func(*Config)) {
	cb.Config = append(cb.Config, NewDefaultConfig(opts...))
}

func LoadConfigFromYAML(path string) (*ConfigBlock, error) {
	// 读取 yaml 文件
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// 解析 yaml 文件
	configBlock, err := ParseYAMLConfig(yamlFile)
	if err != nil {
		return nil, err
	}
	return configBlock, nil
}

func ParseYAMLConfig(data []byte) (*ConfigBlock, error) {
	// 解析 yaml 文件
	var configBlock ConfigBlock
	err := yaml.Unmarshal(data, &configBlock)
	if err != nil {
		return nil, err
	}
	return &configBlock, nil
}

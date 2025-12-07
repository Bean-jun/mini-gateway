package main

type Config struct {
	Port int // go-nginx 运行端口
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
	Config []*Config
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

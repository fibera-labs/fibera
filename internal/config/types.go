package config

import "strings"

type Config struct {
	App      AppConfig      `yaml:"app" mapstructure:"app"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Log      LogConfig      `yaml:"log" mapstructure:"log"`
	Jwt      JwtConfig      `yaml:"jwt" mapstructure:"jwt"`
	Docs     DocsConfig     `yaml:"docs" mapstructure:"docs"`
}

type AppConfig struct {
	Port    string `mapstructure:"port"`
	Version string `mapstructure:"version"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver" mapstructure:"driver"`
	Source       string `yaml:"source" mapstructure:"source"`
	Sync         bool   `yaml:"sync" mapstructure:"sync"`
	MaxIdleConns int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Enable   bool   `yaml:"enable" mapstructure:"enable"`
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

type LogConfig struct {
	Level        string `yaml:"level" mapstructure:"level"`
	Format       string `yaml:"format" mapstructure:"format"`
	OutputMethod string `yaml:"output_method" mapstructure:"output_method"`
	OutputDir    string `yaml:"output_dir" mapstructure:"output_dir"`
	MaxSize      int    `yaml:"max_size" mapstructure:"max_size"`
	MaxBackups   int    `yaml:"max_backups" mapstructure:"max_backups"`
	MaxAge       int    `yaml:"max_age" mapstructure:"max_age"`
}

type JwtConfig struct {
	Issuer         string `yaml:"issuer" mapstructure:"issuer"`
	SigningMethod  string `yaml:"signing_method" mapstructure:"signing_method"`
	ExpirationTime string `yaml:"expiration_time" mapstructure:"expiration_time"`
}

// DocsConfig API 文档配置，接口定义由 swag 注解生成，页面由 Scalar 渲染
type DocsConfig struct {
	Enable           bool   `yaml:"enable" mapstructure:"enable"`
	Path             string `yaml:"path" mapstructure:"path"`
	BasePath         string `yaml:"base_path" mapstructure:"base_path"`
	Title            string `yaml:"title" mapstructure:"title"`
	Theme            string `yaml:"theme" mapstructure:"theme"`
	CustomStyle      string `yaml:"custom_style" mapstructure:"custom_style"`
	ProxyURL         string `yaml:"proxy_url" mapstructure:"proxy_url"`
	CacheAge         int    `yaml:"cache_age" mapstructure:"cache_age"`
	FallbackCacheAge int    `yaml:"fallback_cache_age" mapstructure:"fallback_cache_age"`
	ForceOffline     bool   `yaml:"force_offline" mapstructure:"force_offline"`
}

// 文档配置缺省值，配置项缺失或非法时按此兜底
const (
	defaultDocsPath             = "/docs"
	defaultDocsBasePath         = "/"
	defaultDocsTitle            = "Fibera API documentation"
	defaultDocsCacheAge         = 60
	defaultDocsFallbackCacheAge = 86400
)

// Normalize 补齐缺省值并规范化路径，避免把非法值透传给文档渲染器
func (d DocsConfig) Normalize() DocsConfig {
	d.Path = normalizeDocsPath(d.Path)
	if d.BasePath == "" {
		d.BasePath = defaultDocsBasePath
	}
	if d.Title == "" {
		d.Title = defaultDocsTitle
	}
	if d.CacheAge <= 0 {
		d.CacheAge = defaultDocsCacheAge
	}
	if d.FallbackCacheAge <= 0 {
		d.FallbackCacheAge = defaultDocsFallbackCacheAge
	}
	return d
}

// normalizeDocsPath 将文档路径规范为以 / 开头、不带尾部 / 的形式
func normalizeDocsPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return defaultDocsPath
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if trimmed := strings.TrimRight(path, "/"); trimmed != "" {
		return trimmed
	}
	return defaultDocsPath
}

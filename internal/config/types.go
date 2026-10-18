package config

type Config struct {
	App      AppConfig      `yaml:"app" mapstructure:"app"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Log      LogConfig      `yaml:"log" mapstructure:"log"`
	Jwt      JwtConfig      `yaml:"jwt" mapstructure:"jwt"`
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

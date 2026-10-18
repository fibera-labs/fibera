package config

import (
	"bytes"
	"fibera"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	conf    *Config
	once    sync.Once
	confErr error
)

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	once.Do(func() {
		envFile, eErr := fibera.EnvFile.ReadFile(".env")
		if eErr != nil {
			confErr = fmt.Errorf("failed to load .env file: %w", eErr)
			return
		}

		// 解析环境变量
		envMap, err := godotenv.Parse(bytes.NewReader(envFile))
		if err != nil {
			confErr = fmt.Errorf("failed to parse .env content: %w", err)
			return
		}

		// 将解析后的环境变量设置到系统环境中
		for k, v := range envMap {
			if err = os.Setenv(k, v); err != nil {
				confErr = fmt.Errorf("failed to set environment variable %s: %w", k, err)
				return
			}
		}

		// 配置文件加载路径规则（优先级从高到低）：
		// 1. 外部配置文件：通过明确指定的完整路径加载（默认项目同级目录下config.yaml）
		// 2. 内部配置文件：在程序内部根据运行环境（如"config_"+env）动态拼接路径。
		exConfigFilePath := "../config.yaml"
		inConfigFilePath := "configs/config"

		viper.SetConfigType("yaml")

		if isConfigFileExist(exConfigFilePath) {
			viper.SetConfigFile(exConfigFilePath)
			err = viper.ReadInConfig()
			if err != nil {
				confErr = fmt.Errorf("failed to read the configuration file: %w", err)
				return
			}
		} else {
			env := os.Getenv("APP_ENV")
			if env == "" {
				env = "dev"
			}

			configName := fmt.Sprintf("%s_%s.yaml", inConfigFilePath, env)
			configFile, cErr := fibera.ConfigFile.ReadFile(configName)
			if cErr != nil {
				confErr = fmt.Errorf("failed to read embedded config file %s: %w", configName, cErr)
				return
			}

			err = viper.ReadConfig(bytes.NewReader(configFile))
			if err != nil {
				confErr = fmt.Errorf("failed to parse the configuration file: %w", err)
				return
			}
		}

		conf = &Config{}
		err = viper.Unmarshal(conf)
		if err != nil {
			confErr = fmt.Errorf("failed to parse the configuration file: %w", err)
			return
		}
		// 打印系统版本号
		fmt.Printf("Fibera version:  %s\n", conf.App.Version)
	})
	return conf, confErr
}

// isConfigFileExist 检查配置文件是否存在
func isConfigFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

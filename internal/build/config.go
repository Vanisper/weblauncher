package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config 保存构建配置
type Config struct {
	AppID        string
	AppName      string
	AppVersion   string
	AppPublisher string
	AppURL       string
	OutputName   string
}

// LoadConfig 加载配置（支持 .env.local 覆盖）
// 优先级：.env.local > .env > version 文件 > 默认值
// 使用 Read 避免污染全局环境变量，确保每次构建都是干净的
func LoadConfig(envFile, projectRoot string) (*Config, error) {
	// 确定主 env 文件路径
	if envFile == "" {
		envFile = ".env"
	}
	if !filepath.IsAbs(envFile) {
		envFile = filepath.Join(projectRoot, envFile)
	}

	// 使用 Read 读取配置（不设置到全局环境变量）
	vars := make(map[string]string)

	// 读取主 .env 文件
	if envMap, err := godotenv.Read(envFile); err == nil {
		for k, v := range envMap {
			vars[k] = v
		}
	}

	// 读取 local env 覆盖
	localEnvFile := envFile + ".local"
	if localMap, err := godotenv.Read(localEnvFile); err == nil {
		for k, v := range localMap {
			vars[k] = v // 覆盖
		}
	}

	// 读取版本号（从 version 文件）
	version := readVersionFile(filepath.Join(projectRoot, "version"))

	// 构建配置（从 vars 读取，不受全局环境变量影响）
	// TODO: 支持环境变量优先，便于 ci/cd 管道覆盖配置
	cfg := &Config{
		AppID:        vars["APP_ID"], // 这里不提供默认值，强制要求用户在 .env 中设置
		AppName:      getVar(vars, "APP_NAME", "WebLauncher"),
		AppVersion:   getVar(vars, "APP_VERSION", version),
		AppPublisher: getVar(vars, "APP_PUBLISHER", "Vanisper"),
		AppURL:       getVar(vars, "APP_URL", "https://github.com/vanisper/weblauncher"),
		OutputName:   getVar(vars, "OUTPUT_NAME", "weblauncher.exe"),
	}

	return cfg, nil
}

// readVersionFile 读取 version 文件内容
func readVersionFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "0.0.0" // 默认版本
	}
	return strings.TrimSpace(string(data))
}

// getVar 从映射获取值，如果不存在返回默认值
func getVar(vars map[string]string, key, defaultValue string) string {
	if value, ok := vars[key]; ok && value != "" {
		return value
	}
	return defaultValue
}

// Validate 检查配置是否有效
func (c *Config) Validate() error {
	if c.AppID == "" {
		return fmt.Errorf("APP_ID 不能为空")
	}
	if c.AppName == "" {
		return fmt.Errorf("APP_NAME 不能为空")
	}
	if c.AppVersion == "" {
		return fmt.Errorf("APP_VERSION 不能为空")
	}
	if c.OutputName == "" {
		return fmt.Errorf("OUTPUT_NAME 不能为空")
	}
	return nil
}

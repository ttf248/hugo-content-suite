package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LMStudio LMStudioConfig `json:"lm_studio"`
	Cache    CacheConfig    `json:"cache"`
	Display  DisplayConfig  `json:"display"`
	Paths    PathsConfig    `json:"paths"`
}

type LMStudioConfig struct {
	URL     string `json:"url"`
	Model   string `json:"model"`
	Timeout int    `json:"timeout_seconds"`
}

type CacheConfig struct {
	FileName      string `json:"file_name"`
	AutoSaveCount int    `json:"auto_save_count"`
	DelayMs       int    `json:"delay_ms"`
}

type DisplayConfig struct {
	DefaultLimit int  `json:"default_limit"`
	Colors       bool `json:"colors"`
}

type PathsConfig struct {
	DefaultContentDir string `json:"default_content_dir"`
	TagsDir           string `json:"tags_dir"`
}

var defaultConfig = Config{
	LMStudio: LMStudioConfig{
		URL:     "http://172.19.192.1:2234/v1/chat/completions",
		Model:   "gemma-3-12b-it",
		Timeout: 30,
	},
	Cache: CacheConfig{
		FileName:      "tag_translations_cache.json",
		AutoSaveCount: 5,
		DelayMs:       500,
	},
	Display: DisplayConfig{
		DefaultLimit: 20,
		Colors:       true,
	},
	Paths: PathsConfig{
		DefaultContentDir: "../../content/post",
		TagsDir:           "../tags",
	},
}

func LoadConfig() (*Config, error) {
	configPath := "config.json"

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig(configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

func createDefaultConfig(configPath string) (*Config, error) {
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化默认配置失败: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("创建默认配置文件失败: %v", err)
	}

	fmt.Println("✅ 已创建默认配置文件: config.json")
	return &defaultConfig, nil
}

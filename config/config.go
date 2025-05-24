package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LMStudio    LMStudioConfig    `json:"lm_studio"`
	Cache       CacheConfig       `json:"cache"`
	Display     DisplayConfig     `json:"display"`
	Paths       PathsConfig       `json:"paths"`
	Translation TranslationConfig `json:"translation"`
	Logging     LoggingConfig     `json:"logging"`
	Language    LanguageConfig    `json:"language"`
}

type LMStudioConfig struct {
	URL     string `json:"url"`
	Model   string `json:"model"`
	Timeout int    `json:"timeout_seconds"`
}

type CacheConfig struct {
	TagFileName     string `json:"tag_file_name"`
	ArticleFileName string `json:"article_file_name"`
	AutoSaveCount   int    `json:"auto_save_count"`
	DelayMs         int    `json:"delay_ms"`
	ExpireDays      int    `json:"expire_days"`
}

type DisplayConfig struct {
	DefaultLimit int  `json:"default_limit"`
	Colors       bool `json:"colors"`
}

type PathsConfig struct {
	DefaultContentDir string `json:"default_content_dir"`
	TagsDir           string `json:"tags_dir"`
}

type TranslationConfig struct {
	RetryAttempts   int      `json:"retry_attempts"`
	DelayBetweenMs  int      `json:"delay_between_ms"`
	ValidateResult  bool     `json:"validate_result"`
	CleanupPatterns []string `json:"cleanup_patterns"`
}

type LoggingConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type LanguageConfig struct {
	TargetLanguages []string          `json:"target_languages"`
	LanguageNames   map[string]string `json:"language_names"`
}

var defaultConfig = Config{
	LMStudio: LMStudioConfig{
		URL:     "http://172.19.192.1:2234/v1/chat/completions",
		Model:   "gemma-3-12b-it",
		Timeout: 30,
	},
	Cache: CacheConfig{
		TagFileName:     "tag_translations_cache.json",
		ArticleFileName: "article_translations_cache.json",
		AutoSaveCount:   5,
		DelayMs:         500,
		ExpireDays:      30,
	},
	Display: DisplayConfig{
		DefaultLimit: 20,
		Colors:       true,
	},
	Paths: PathsConfig{
		DefaultContentDir: "../../content/post",
		TagsDir:           "../tags",
	},
	Translation: TranslationConfig{
		RetryAttempts:  2,
		DelayBetweenMs: 0,
		ValidateResult: true,
		CleanupPatterns: []string{
			"Translation:",
			"Translated:",
			"English:",
			"Result:",
			"Output:",
		},
	},
	Logging: LoggingConfig{
		Level: "DEBUG",
		File:  "tag-scanner.log",
	},
	Language: LanguageConfig{
		TargetLanguages: []string{"en"},
		LanguageNames: map[string]string{
			"en": "English",
			"ja": "Japanese",
			"ko": "Korean",
		},
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

// GetGlobalConfig 获取全局配置实例
var globalConfig *Config

func GetGlobalConfig() *Config {
	if globalConfig == nil {
		config, err := LoadConfig()
		if err != nil {
			fmt.Printf("⚠️ 加载配置失败，使用默认配置: %v\n", err)
			globalConfig = &defaultConfig
		} else {
			globalConfig = config
		}
	}
	return globalConfig
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	LMStudio    LMStudioConfig    `json:"lm_studio"`
	Cache       CacheConfig       `json:"cache"`
	Display     DisplayConfig     `json:"display"`
	Paths       PathsConfig       `json:"paths"`
	Translation TranslationConfig `json:"translation"`
	Paragraph   ParagraphConfig   `json:"paragraph"`
	Logging     LoggingConfig     `json:"logging"`
	Language    LanguageConfig    `json:"language"`
}

type LMStudioConfig struct {
	URL     string `json:"url"`
	Model   string `json:"model"`
	Timeout int    `json:"timeout_seconds"`
}

type CacheConfig struct {
	TagFileName      string `json:"article_tag_file_name"`
	ArticleFileName  string `json:"article_slug_file_name"`
	CategoryFileName string `json:"article_category_file_name"` // 新增
	AutoSaveCount    int    `json:"auto_save_count"`
	DelayMs          int    `json:"delay_ms"`
	ExpireDays       int    `json:"expire_days"`
}

type DisplayConfig struct {
	DefaultLimit int  `json:"default_limit"`
	Colors       bool `json:"colors"`
}

type PathsConfig struct {
	DefaultContentDir string `json:"default_content_dir"`
	TagsDir           string `json:"tags_dir"`
	RuntimeDir        string `json:"runtime_dir"`
}

type TranslationConfig struct {
	RetryAttempts   int      `json:"retry_attempts"`
	DelayBetweenMs  int      `json:"delay_between_ms"`
	ValidateResult  bool     `json:"validate_result"`
	CleanupPatterns []string `json:"cleanup_patterns"`
}

type ParagraphConfig struct {
	MaxLength             int  `json:"max_length"`              // 段落最大长度（字符数）
	EnableSplitting       bool `json:"enable_splitting"`        // 是否启用段落拆分
	SplitAtSentences      bool `json:"split_at_sentences"`      // 是否在句子边界拆分
	MinSplitLength        int  `json:"min_split_length"`        // 拆分后段落的最小长度
	MergeAfterTranslation bool `json:"merge_after_translation"` // 翻译后是否合并拆分的段落
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
		URL:     "http://localhost:2234/v1/chat/completions",
		Model:   "google/gemma-3-4b",
		Timeout: 30,
	},
	Cache: CacheConfig{
		TagFileName:      "tag_translations_cache.json",
		ArticleFileName:  "slug_translations_cache.json",
		CategoryFileName: "category_translations_cache.json", // 新增
		AutoSaveCount:    5,
		DelayMs:          500,
		ExpireDays:       30,
	},
	Display: DisplayConfig{
		DefaultLimit: 20,
		Colors:       true,
	},
	Paths: PathsConfig{
		DefaultContentDir: "../../content/post",
		TagsDir:           "../tags",
		RuntimeDir:        ".hugo-content-suite",
	}, Translation: TranslationConfig{
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
	Paragraph: ParagraphConfig{
		MaxLength:             800,  // 段落最大2000字符
		EnableSplitting:       true, // 默认启用段落拆分
		SplitAtSentences:      true, // 在句子边界拆分
		MinSplitLength:        200,  // 拆分后段落最小200字符
		MergeAfterTranslation: true, // 默认翻译后合并拆分的段落
	},
	Logging: LoggingConfig{
		Level: "DEBUG",
		File:  "hugo-content-suite.log",
	},
	Language: LanguageConfig{
		TargetLanguages: []string{"en", "ja", "ko", "fr", "ru", "hi"},
		LanguageNames: map[string]string{
			"en": "English",
			"ja": "Japanese",
			"ko": "Korean",
			"fr": "French",
			"ru": "Russian",
			"hi": "Hindi",
		},
	},
}

// LoadConfig 从指定文件读取配置。读取配置不应产生文件写入，避免一次拼写错误
// 在任意工作目录留下看似有效、实则未被用户确认的默认配置。
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件 %s 失败: %w", configPath, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	if err := config.ResolvePaths(filepath.Dir(configPath)); err != nil {
		return nil, err
	}
	return &config, nil
}

// ResolvePaths 将所有项目路径固定到配置文件所在目录，避免依赖启动时的工作目录。
func (c *Config) ResolvePaths(baseDir string) error {
	resolve := func(path string) (string, error) {
		if path == "" {
			return "", nil
		}
		if filepath.IsAbs(path) {
			return filepath.Clean(path), nil
		}
		return filepath.Abs(filepath.Join(baseDir, path))
	}
	var err error
	if c.Paths.DefaultContentDir, err = resolve(c.Paths.DefaultContentDir); err != nil {
		return err
	}
	if c.Paths.TagsDir, err = resolve(c.Paths.TagsDir); err != nil {
		return err
	}
	if c.Paths.RuntimeDir, err = resolve(c.Paths.RuntimeDir); err != nil {
		return err
	}
	if c.Paths.RuntimeDir == "" {
		return fmt.Errorf("paths.runtime_dir 不能为空")
	}
	c.Logging.File, err = resolve(filepath.Join(c.Paths.RuntimeDir, c.Logging.File))
	if err != nil {
		return err
	}
	for _, name := range []*string{&c.Cache.TagFileName, &c.Cache.ArticleFileName, &c.Cache.CategoryFileName} {
		*name, err = resolve(filepath.Join(c.Paths.RuntimeDir, *name))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetGlobalConfig 获取全局配置实例
var globalConfig *Config

func GetGlobalConfig() *Config {
	if globalConfig == nil {
		config, err := LoadConfig("config.json")
		if err != nil {
			fmt.Printf("⚠️ 加载配置失败，使用默认配置: %v\n", err)
			globalConfig = &defaultConfig
		} else {
			globalConfig = config
		}
	}
	return globalConfig
}

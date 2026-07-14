package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigResolvesRuntimeFilesWithoutWriting(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	content := `{"paths":{"default_content_dir":"content/post","tags_dir":"content/tags","runtime_dir":"runtime"},"cache":{"article_tag_file_name":"tags.json","article_slug_file_name":"slugs.json","article_category_file_name":"categories.json"},"logging":{"file":"suite.log"}}`
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}
	wantRuntime := filepath.Join(dir, "runtime")
	if cfg.Paths.RuntimeDir != wantRuntime || cfg.Cache.TagFileName != filepath.Join(wantRuntime, "tags.json") || cfg.Logging.File != filepath.Join(wantRuntime, "suite.log") {
		t.Fatalf("运行路径未按配置目录解析: %#v", cfg)
	}
}

func TestLoadConfigDoesNotCreateMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("缺失配置应返回错误")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("LoadConfig 不应创建配置文件，stat err=%v", err)
	}
}

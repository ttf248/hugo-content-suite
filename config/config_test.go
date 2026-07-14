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

func TestSelectedModelAndSwitch(t *testing.T) {
	cfg := &Config{ActiveModel: "local", Models: []LLMConfig{{Name: "local", APIType: "openai_chat", URL: "http://localhost", Model: "local-model"}, {Name: "remote", APIType: "anthropic_messages", URL: "https://example.test/v1/messages", Model: "remote-model"}}}
	if model, err := cfg.SelectedModel(); err != nil || model.Name != "local" {
		t.Fatalf("默认模型错误: %#v, %v", model, err)
	}
	if err := cfg.SelectModel("remote"); err != nil {
		t.Fatal(err)
	}
	if model, _ := cfg.SelectedModel(); model.APIType != "anthropic_messages" {
		t.Fatalf("切换失败: %#v", model)
	}
	if err := cfg.SelectModel("missing"); err == nil || cfg.ActiveModel != "remote" {
		t.Fatal("非法切换应保留原选择")
	}
}

func TestLoadLocalConfigMergesExample(t *testing.T) {
	dir := t.TempDir()
	example := `{"paths":{"default_content_dir":"content/post","tags_dir":"content/tags","runtime_dir":"runtime"},"cache":{"article_tag_file_name":"tags.json","article_slug_file_name":"slugs.json","article_category_file_name":"categories.json"},"logging":{"file":"suite.log"},"active_model":"base","models":[{"name":"base","api_type":"openai_chat","url":"http://example","model":"base"}]}`
	local := `{"active_model":"remote","models":[{"name":"remote","api_type":"anthropic_messages","url":"https://example.test/messages","model":"remote","api_key":"local-key"}]}`
	if err := os.WriteFile(filepath.Join(dir, "config.example.json"), []byte(example), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.local.json"), []byte(local), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(filepath.Join(dir, "config.local.json"))
	if err != nil {
		t.Fatal(err)
	}
	if model, err := cfg.SelectedModel(); err != nil || model.APIKey != "local-key" {
		t.Fatalf("本地模型覆盖失败: %#v, %v", model, err)
	}
	if cfg.Paths.RuntimeDir != filepath.Join(dir, "runtime") {
		t.Fatalf("示例路径未保留: %s", cfg.Paths.RuntimeDir)
	}
}

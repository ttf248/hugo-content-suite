package translator

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"hugo-content-suite/config"
)

func testConfig(dir, url string) *config.Config {
	return &config.Config{
		LMStudio: config.LMStudioConfig{URL: url, Model: "test-model", Timeout: 1},
		Cache:    config.CacheConfig{TagFileName: filepath.Join(dir, "tags.json"), ArticleFileName: filepath.Join(dir, "slugs.json"), CategoryFileName: filepath.Join(dir, "categories.json")},
		Language: config.LanguageConfig{LanguageNames: map[string]string{"en": "English"}},
	}
}

func TestOpenAIRequestUsesSelectedModelInsteadOfLegacyConfig(t *testing.T) {
	const selectedModel = "selected-model"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"Hello"}}]}`))
	}))
	defer server.Close()
	cfg := testConfig(t.TempDir(), "") // 模拟新配置已不再填写 legacy lm_studio。
	cfg.ActiveModel = "selected"
	cfg.Models = []config.LLMConfig{{Name: "selected", APIType: "openai_chat", URL: server.URL + "/v1/chat/completions", Model: selectedModel, Timeout: 1}}
	translator := NewTranslationUtilsWithConfig(cfg, server.Client())
	got, err := translator.TranslateToLanguage("你好", "en")
	if err != nil || got != "Hello" {
		t.Fatalf("翻译结果=%q, err=%v", got, err)
	}
}

func TestFailedTranslationIsNotCached(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "offline", http.StatusBadGateway) }))
	defer server.Close()
	translator := NewTranslationUtilsWithConfig(testConfig(t.TempDir(), server.URL), server.Client())
	if _, err := translator.translateWithCache("失败项", "en", kTagCache); err == nil {
		t.Fatal("失败请求应返回错误")
	}
	if _, found := translator.cache.Get("en:失败项", kTagCache); found {
		t.Fatal("失败翻译不得写入缓存")
	}
}

func TestAnthropicRequestUsesAPIKeyEnvAndHeaders(t *testing.T) {
	t.Setenv("MINIMAX_API_KEY", "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" || r.Header.Get("anthropic-version") != "2023-06-01" {
			http.Error(w, "headers", http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte(`{"content":[{"type":"text","text":"Hello"}]}`))
	}))
	defer server.Close()
	cfg := testConfig(t.TempDir(), server.URL)
	cfg.ActiveModel = "minimax"
	cfg.Models = []config.LLMConfig{{Name: "minimax", APIType: "anthropic_messages", URL: server.URL, Model: "MiniMax-M2.5", APIKeyEnv: "MINIMAX_API_KEY", Timeout: 1}}
	translator := NewTranslationUtilsWithConfig(cfg, server.Client())
	got, err := translator.TranslateToLanguage("你好", "en")
	if err != nil || got != "Hello" {
		t.Fatalf("Anthropic 翻译=%q, err=%v", got, err)
	}
}

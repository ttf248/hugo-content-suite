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

func TestTranslateToLanguageUsesInjectedHTTPClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"Hello"}}]}`))
	}))
	defer server.Close()
	translator := NewTranslationUtilsWithConfig(testConfig(t.TempDir(), server.URL), server.Client())
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

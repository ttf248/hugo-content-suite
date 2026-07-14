package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanArticlesReturnsFrontMatterError(t *testing.T) {
	dir := t.TempDir()
	post := filepath.Join(dir, "bad")
	if err := os.MkdirAll(post, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(post, "index.md"), []byte("---\ntitle: [\n---\nbody"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ScanArticles(dir); err == nil {
		t.Fatal("无效 front matter 应返回错误")
	}
}

func TestScanArticlesKeepsFencedCodeBlock(t *testing.T) {
	dir := t.TempDir()
	post := filepath.Join(dir, "post")
	if err := os.MkdirAll(post, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\ntitle: 示例\n---\n\n正文\n\n```go\nfmt.Println(1)\n```\n"
	if err := os.WriteFile(filepath.Join(post, "index.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	articles, err := ScanArticlesForTranslation(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 1 || len(articles[0].BodyContent) != 2 || articles[0].BodyContent[1] != "```go\nfmt.Println(1)\n```" {
		t.Fatalf("代码块边界丢失: %#v", articles)
	}
}

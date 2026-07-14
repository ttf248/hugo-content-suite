package operations

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeleteArticlesByLanguageLeavesSourceFiles(t *testing.T) {
	dir := t.TempDir()
	post := filepath.Join(dir, "post")
	if err := os.MkdirAll(post, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"index.md", "index.en.md", "index.ja.md"} {
		if err := os.WriteFile(filepath.Join(post, name), []byte("---\ntitle: test\n---"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	p := NewProcessor(dir)
	count, err := p.deleteArticlesByLanguage("en")
	if err != nil || count != 1 {
		t.Fatalf("count=%d, err=%v", count, err)
	}
	if _, err := os.Stat(filepath.Join(post, "index.md")); err != nil {
		t.Fatal("源文件不应删除")
	}
	if _, err := os.Stat(filepath.Join(post, "index.en.md")); !os.IsNotExist(err) {
		t.Fatal("目标译文应删除")
	}
}

func TestReadChoiceConsumesSingleLine(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("en\n"))
	if got := strings.TrimSpace(readChoice(reader, "")); got != "en" {
		t.Fatalf("got %q", got)
	}
}

package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteSlugOnlyChangesFrontMatter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.md")
	content := "---\ntitle: 示例\nslug: old\n---\n\n```yaml\nslug: old\n```\n"
	if err := os.WriteFile(path, []byte(content), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := writeSlug(path, "new-slug"); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(got), "slug: old") != 1 || !strings.Contains(string(got), "slug: \"new-slug\"") {
		t.Fatalf("slug 更新越过 front matter: %s", got)
	}
}

func TestWriteSlugAddsMissingField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.md")
	if err := os.WriteFile(path, []byte("---\ntitle: 示例\n---\n正文\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeSlug(path, "example"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if !strings.Contains(string(got), "slug: \"example\"\n---") {
		t.Fatalf("未在 front matter 末尾新增 slug: %s", got)
	}
}

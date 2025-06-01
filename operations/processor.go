package operations

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Processor struct {
	contentDir string
}

func NewProcessor(contentDir string) *Processor {
	return &Processor{
		contentDir: contentDir,
	}
}

func (p *Processor) confirmExecution(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}

func (p *Processor) ScanLanguages() ([]string, error) {
	// 扫描 contentDir 下所有文章，收集所有语言（假设文件名格式为 xxx.{lang}.md）
	langSet := make(map[string]struct{})
	err := filepath.Walk(p.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		parts := strings.Split(base, ".")
		// 只处理以 .md 结尾的文件
		if len(parts) >= 3 && parts[len(parts)-1] == "md" {
			lang := parts[len(parts)-2]
			langSet[lang] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	langs := make([]string, 0, len(langSet))
	for lang := range langSet {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	return langs, nil
}

func (p *Processor) DeleteArticlesByLanguage(lang string) error {
	// 删除 contentDir 下所有指定语言的文章
	return filepath.Walk(p.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		parts := strings.Split(base, ".")
		if len(parts) >= 3 && parts[len(parts)-2] == lang {
			return os.Remove(path)
		}
		return nil
	})
}

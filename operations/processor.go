package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/scanner"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
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
	// 复用 scanner.ScanArticles，收集所有文章的语言
	articles, err := scanner.ScanArticles(p.contentDir)
	if err != nil {
		return nil, err
	}
	langSet := make(map[string]struct{})
	for _, article := range articles {
		lang := extractLangFromPath(article.FilePath)
		if lang != "" {
			langSet[lang] = struct{}{}
		}
	}
	langs := make([]string, 0, len(langSet))
	for lang := range langSet {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	return langs, nil
}

// 提取语言代码（假设路径如 .../xxx.zh-cn/index.md）
func extractLangFromPath(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) < 2 {
		return ""
	}
	parent := parts[len(parts)-2]
	// 约定多语言目录为 xxx.{lang}
	if idx := strings.LastIndex(parent, "."); idx != -1 && idx < len(parent)-1 {
		return parent[idx+1:]
	}
	return ""
}

func (p *Processor) DeleteArticlesByLanguage(lang string) error {
	// 删除 contentDir 下所有指定语言的 index.md 文件夹
	return filepath.Walk(p.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == "index.md" {
			parent := filepath.Base(filepath.Dir(path))
			if strings.HasSuffix(parent, "."+lang) {
				return os.Remove(path)
			}
		}
		return nil
	})
}

// 交互式删除指定语言文章
func (p *Processor) DeleteArticlesByLanguageInteractive(reader *bufio.Reader) {
	langs, err := p.ScanLanguages()
	if err != nil {
		color.Red("扫描语言失败: %v", err)
		return
	}
	if len(langs) == 0 {
		color.Red("未检测到任何语言")
		return
	}
	color.Cyan("当前检测到的语言：")
	for i, lang := range langs {
		fmt.Printf("  %d. %s\n", i+1, lang)
	}
	choice := ""
	fmt.Print("请输入要删除的语言编号: ")
	fmt.Fscanln(reader, &choice)
	idx := -1
	fmt.Sscanf(choice, "%d", &idx)
	if idx < 1 || idx > len(langs) {
		color.Red("无效选择")
		return
	}
	langToDelete := langs[idx-1]
	fmt.Printf("确定要删除所有 [%s] 语言的文章吗？(y/N): ", langToDelete)
	confirm, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(confirm)) == "y" {
		err := p.DeleteArticlesByLanguage(langToDelete)
		if err != nil {
			color.Red("删除失败: %v", err)
		} else {
			color.Green("已删除所有 [%s] 语言的文章", langToDelete)
		}
	} else {
		color.Yellow("已取消删除操作")
	}
}

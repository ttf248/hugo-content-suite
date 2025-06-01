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

func (p *Processor) ScanLanguages() ([]string, error) {
	articles, err := scanner.ScanArticlesWithLangs(p.contentDir, true)
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

// 交互式删除指定语言文章
func (p *Processor) DeleteArticles(reader *bufio.Reader) {
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
		err := p.deleteArticlesByLanguage(langToDelete)
		if err != nil {
			color.Red("删除失败: %v", err)
		} else {
			color.Green("已删除所有 [%s] 语言的文章", langToDelete)
		}
	} else {
		color.Yellow("已取消删除操作")
	}
}

// 提取语言代码（假设路径如 markdown/index.en.md 或 markdown\index.en.md）
func extractLangFromPath(path string) string {
	base := filepath.Base(path)
	parts := strings.Split(base, ".")
	if len(parts) < 3 {
		return ""
	}
	// 形如 index.en.md，语言在倒数第二个
	return parts[len(parts)-2]
}

func (p *Processor) deleteArticlesByLanguage(lang string) error {
	articles, err := scanner.ScanArticlesWithLangs(p.contentDir, true)
	if err != nil {
		return err
	}

	var toDelete []string
	for _, article := range articles {
		if extractLangFromPath(article.FilePath) == lang {
			toDelete = append(toDelete, article.FilePath)
		}
	}

	if len(toDelete) == 0 {
		return fmt.Errorf("未找到语言为 [%s] 的文章", lang)
	}

	for _, file := range toDelete {
		err := removeFile(file)
		if err != nil {
			return fmt.Errorf("删除文件 %s 失败: %v", file, err)
		}
	}

	return nil
}

// removeFile 封装 os.Remove，便于后续扩展
func removeFile(path string) error {
	return os.Remove(path)
}

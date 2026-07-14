package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/scanner"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	choice := strings.TrimSpace(readChoice(reader, "请输入要删除的语言编号: "))
	idx, err := strconv.Atoi(choice)
	if err != nil {
		idx = -1
	}
	if idx < 1 || idx > len(langs) {
		color.Red("无效选择")
		return
	}
	langToDelete := langs[idx-1]
	if strings.TrimSpace(readChoice(reader, fmt.Sprintf("输入语言代码 %s 以确认删除: ", langToDelete))) == langToDelete {
		count, err := p.deleteArticlesByLanguage(langToDelete)
		if err != nil {
			color.Red("删除失败: %v", err)
		} else {
			color.Green("已删除 %d 个 [%s] 译文文件", count, langToDelete)
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

func readChoice(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	value, _ := reader.ReadString('\n')
	return value
}

func (p *Processor) deleteArticlesByLanguage(lang string) (int, error) {
	articles, err := scanner.ScanArticlesWithLangs(p.contentDir, true)
	if err != nil {
		return 0, err
	}

	var toDelete []string
	for _, article := range articles {
		if extractLangFromPath(article.FilePath) == lang {
			toDelete = append(toDelete, article.FilePath)
		}
	}

	if len(toDelete) == 0 {
		return 0, fmt.Errorf("未找到语言为 [%s] 的文章", lang)
	}

	for _, file := range toDelete {
		err := removeFile(file)
		if err != nil {
			return 0, fmt.Errorf("删除文件 %s 失败: %v", file, err)
		}
	}

	return len(toDelete), nil
}

// removeFile 封装 os.Remove，便于后续扩展
func removeFile(path string) error {
	return os.Remove(path)
}

package menu

import (
	"bufio"
	"fmt"
	"hugo-content-suite/models"
	"hugo-content-suite/operations"
	"hugo-content-suite/utils"
	"strings"

	"github.com/fatih/color"
)

type InteractiveMenu struct {
	reader    *bufio.Reader
	processor *operations.Processor
}

func NewInteractiveMenu(reader *bufio.Reader, contentDir string) *InteractiveMenu {
	return &InteractiveMenu{
		reader:    reader,
		processor: operations.NewProcessor(contentDir),
	}
}

func (m *InteractiveMenu) Show(tagStats []models.TagStats, categoryStats []models.CategoryStats, noTagArticles []models.Article) {
	for {
		m.displayMainMenu()
		choice := utils.GetChoice(m.reader, "请选择功能 (0-8): ")

		switch choice {
		case "1":
			m.processor.GenerateTagPages(tagStats, m.reader)
		case "2":
			m.processor.GenerateArticleSlugs(m.reader)
		case "3":
			m.processor.TranslateArticles(m.reader)
		case "4":
			m.deleteArticlesByLanguage()
		case "0":
			color.Green("感谢使用！再见！")
			return
		default:
			color.Red("⚠️  无效选择，请重新输入")
		}
	}
}

func (m *InteractiveMenu) displayMainMenu() {
	color.Cyan("\n=== Hugo 博客管理工具 ===")
	fmt.Println()

	// 主要功能模块
	color.Red("🚀 核心功能")
	// 内容管理模块
	color.Green("📝 内容管理")
	fmt.Println("  1. 生成标签页面")
	fmt.Println("  2. 生成文章Slug")
	fmt.Println("  3. 翻译文章为多语言版本")
	fmt.Println("  4. 删除指定语言的文章") // 新增菜单项
	fmt.Println()

	color.Red("  0. 退出程序")
	fmt.Println()
}

func (m *InteractiveMenu) deleteArticlesByLanguage() {
	langs, err := m.processor.ScanLanguages()
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
	choice := utils.GetChoice(m.reader, "请输入要删除的语言编号: ")
	idx := -1
	fmt.Sscanf(choice, "%d", &idx)
	if idx < 1 || idx > len(langs) {
		color.Red("无效选择")
		return
	}
	langToDelete := langs[idx-1]
	confirm := utils.GetChoice(m.reader, fmt.Sprintf("确定要删除所有 [%s] 语言的文章吗？(y/N): ", langToDelete))
	if strings.ToLower(confirm) == "y" {
		err := m.processor.DeleteArticlesByLanguage(langToDelete)
		if err != nil {
			color.Red("删除失败: %v", err)
		} else {
			color.Green("已删除所有 [%s] 语言的文章", langToDelete)
		}
	} else {
		color.Yellow("已取消删除操作")
	}
}

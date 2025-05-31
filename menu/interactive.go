package menu

import (
	"bufio"
	"fmt"
	"hugo-content-suite/models"
	"hugo-content-suite/operations"
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
		choice := m.getChoice("请选择功能 (0-7): ")

		switch choice {
		case "1":
			m.processor.QuickProcessAll(tagStats, m.reader)
		case "2":
			m.processor.GenerateTagPages(tagStats, m.reader)
		case "3":
			m.processor.GenerateArticleSlugs(m.reader)
		case "4":
			m.processor.TranslateArticles(m.reader)
		case "5":
			m.processor.ShowCacheStatus()
		case "6":
			m.processor.GenerateBulkTranslationCache(tagStats, m.reader)
		case "7":
			m.processor.ClearTranslationCache(m.reader)
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
	fmt.Println("  1. 一键处理全部 (自动执行完整博客处理流程)")
	fmt.Println()

	// 内容管理模块
	color.Green("📝 内容管理")
	fmt.Println("  2. 生成标签页面")
	fmt.Println("  3. 生成文章Slug")
	fmt.Println("  4. 翻译文章为多语言版本")
	fmt.Println()

	// 缓存管理模块
	color.Magenta("💾 缓存管理")
	fmt.Println("  5. 查看缓存状态")
	fmt.Println("  6. 生成全量翻译缓存")
	fmt.Println("  7. 清空翻译缓存")
	fmt.Println()

	color.Red("  0. 退出程序")
	fmt.Println()
}

func (m *InteractiveMenu) getChoice(prompt string) string {
	fmt.Print(prompt)
	input, _ := m.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

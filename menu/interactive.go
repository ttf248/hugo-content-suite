package menu

import (
	"bufio"
	"fmt"
	"hugo-content-suite/operations"
	"hugo-content-suite/utils"

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

func (m *InteractiveMenu) Show() {
	for {
		m.displayMainMenu()
		choice := utils.GetChoice(m.reader, "请选择功能 (0-8): ")

		switch choice {
		case "1":
			m.processor.GenerateTagPages(m.reader)
		case "2":
			m.processor.GenerateArticleSlugs(m.reader)
		case "3":
			m.processor.TranslateArticles(m.reader)
		case "4":
			m.processor.DeleteArticles(m.reader)
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

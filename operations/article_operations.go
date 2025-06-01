package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/generator"

	"github.com/fatih/color"
)

func (p *Processor) TranslateArticles(reader *bufio.Reader) {
	if p.contentDir == "" {
		color.Red("❌ 内容目录未设置")
		return
	}

	// 获取翻译状态统计
	color.Cyan("正在分析文章翻译状态...")
	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	previews, createCount, updateCount, err := articleTranslator.PrepareArticleTranslations()
	if err != nil {
		color.Red("❌ 分析失败: %v", err)
		return
	}

	p.displayTranslationStats(createCount, updateCount, len(previews))

	if createCount == 0 && updateCount == 0 {
		color.Green("✅ 所有文章都已完全翻译")
		return
	}

	// 选择翻译模式
	mode := p.selectPageMode("文章翻译", createCount, updateCount, reader)
	if mode == "" {
		return
	}

	// 根据模式筛选预览
	targetPreviews := filterTranslationsByMode(previews, mode)

	// 显示警告和确认
	p.displayTranslationWarning(mode, createCount, updateCount)

	if !p.confirmExecution(reader, "\n确认开始翻译？(y/n): ") {
		color.Yellow("❌ 已取消翻译")
		return
	}

	color.Cyan("🚀 开始翻译文章...")
	if err := articleTranslator.TranslateArticlesWithMode(targetPreviews, mode); err != nil {
		color.Red("❌ 翻译失败: %v", err)
	}
}

func (p *Processor) displayTranslationStats(createCount, updateCount, totalTasks int) {
	// 计算总文章数（去重）
	totalArticles := createCount
	if updateCount > createCount {
		totalArticles = updateCount
	}

	fmt.Printf("\n📊 翻译统计信息:\n")
	fmt.Printf("   🆕 有缺失翻译的文章: %d 篇\n", createCount)
	fmt.Printf("   ✅ 已有翻译的文章: %d 篇\n", updateCount)
	fmt.Printf("   📦 文章总数: %d 篇\n", totalArticles)
	fmt.Printf("   🌐 翻译任务总数: %d 个\n", totalTasks)

	// 显示详细的语言翻译状态
	if createCount > 0 || updateCount > 0 {
		fmt.Printf("\n💡 说明:\n")
		fmt.Printf("   • 有缺失翻译: 至少有一种目标语言缺失翻译的文章\n")
		fmt.Printf("   • 已有翻译: 至少有一种目标语言已翻译的文章\n")
		fmt.Printf("   • 翻译任务: 每篇文章的每种目标语言为一个任务\n")
	}
}

func (p *Processor) displayTranslationWarning(mode string, createCount, updateCount int) {
	fmt.Println()
	color.Yellow("⚠️  重要提示:")
	fmt.Println("• 文章翻译可能需要较长时间，建议在网络稳定时执行")
	fmt.Println("• 翻译过程中请保持网络连接稳定")
	fmt.Println("• 文章翻译会使用缓存加速重复内容的翻译")

	switch mode {
	case "create":
		fmt.Printf("• 将为 %d 篇文章补充缺失的语言翻译\n", createCount)
	case "update":
		fmt.Printf("• 将重新翻译 %d 篇已有翻译的文章\n", updateCount)
	case "all":
		fmt.Printf("• 将处理 %d 篇文章的翻译（包括新增和更新）\n", createCount+updateCount)
	}
}

// filterTranslationsByMode 根据模式筛选翻译任务
func filterTranslationsByMode(previews []generator.ArticleTranslationPreview, mode string) []generator.ArticleTranslationPreview {
	var filtered []generator.ArticleTranslationPreview

	for _, preview := range previews {
		switch mode {
		case "create":
			if preview.Status == "missing" {
				filtered = append(filtered, preview)
			}
		case "update":
			if preview.Status == "update" {
				filtered = append(filtered, preview)
			}
		case "all":
			filtered = append(filtered, preview)
		}
	}

	return filtered
}

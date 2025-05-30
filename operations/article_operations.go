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
	status, err := articleTranslator.GetTranslationStatus()
	if err != nil {
		color.Red("❌ 分析失败: %v", err)
		return
	}

	p.displayTranslationStats(status.MissingArticles, status.ExistingArticles, status.TotalArticles)

	if status.MissingArticles == 0 && status.ExistingArticles == 0 {
		color.Green("✅ 没有需要翻译的文章")
		return
	}

	// 选择翻译模式
	mode := p.selectTranslationMode(status.MissingArticles, status.ExistingArticles, reader)
	if mode == "" {
		return
	}

	// 显示警告和确认
	p.displayTranslationWarning(mode, status.MissingArticles, status.ExistingArticles)

	if !p.confirmExecution(reader, "\n确认开始翻译？(y/n): ") {
		color.Yellow("❌ 已取消翻译")
		return
	}

	color.Cyan("🚀 开始翻译文章...")
	if err := articleTranslator.TranslateArticles(mode); err != nil {
		color.Red("❌ 翻译失败: %v", err)
	}
}

func (p *Processor) displayTranslationStats(missingCount, existingCount, total int) {
	fmt.Printf("\n📊 翻译统计信息:\n")
	fmt.Printf("   🆕 需要翻译的文章: %d 篇\n", missingCount)
	fmt.Printf("   ✅ 已完全翻译的文章: %d 篇\n", existingCount)
	fmt.Printf("   📦 文章总数: %d 篇\n", total)

	// 显示详细的语言翻译状态
	if missingCount > 0 || existingCount > 0 {
		fmt.Printf("\n💡 说明:\n")
		fmt.Printf("   • 需要翻译: 至少有一种目标语言缺失翻译的文章\n")
		fmt.Printf("   • 已完全翻译: 所有目标语言都已翻译的文章\n")
	}
}

func (p *Processor) selectTranslationMode(missingCount, existingCount int, reader *bufio.Reader) string {
	fmt.Println("\n🔧 请选择翻译模式:")

	options := []string{}
	if missingCount > 0 {
		options = append(options, fmt.Sprintf("1. 仅翻译缺失的文章 (%d 篇)", missingCount))
	}
	if existingCount >= 0 {
		options = append(options, fmt.Sprintf("2. 重新翻译现有文章 (%d 篇)", existingCount))
	}
	if missingCount > 0 && existingCount >= 0 {
		options = append(options, fmt.Sprintf("3. 翻译全部文章 (%d 篇)", missingCount+existingCount))
	}

	for _, option := range options {
		fmt.Printf("   %s\n", option)
	}
	fmt.Println("   0. 取消操作")

	choice := p.getChoice(reader, "请选择: ")

	switch choice {
	case "1":
		if missingCount == 0 {
			color.Yellow("⚠️  没有需要翻译的文章")
			return ""
		}
		color.Blue("🆕 将翻译 %d 篇缺失的文章", missingCount)
		return "missing"
	case "2":
		if existingCount == 0 {
			color.Yellow("⚠️  没有现有的英文文章")
			return ""
		}
		color.Blue("🔄 将重新翻译 %d 篇现有文章", existingCount)
		return "update"
	case "3":
		if missingCount == 0 && existingCount == 0 {
			color.Yellow("⚠️  没有需要翻译的文章")
			return ""
		}
		color.Blue("📦 将翻译 %d 篇文章", missingCount+existingCount)
		return "all"
	case "0":
		color.Yellow("❌ 已取消操作")
		return ""
	default:
		color.Red("⚠️  无效选择")
		return ""
	}
}

func (p *Processor) displayTranslationWarning(mode string, missingCount, existingCount int) {
	fmt.Println()
	color.Yellow("⚠️  重要提示:")
	fmt.Println("• 文章翻译可能需要较长时间，建议在网络稳定时执行")
	fmt.Println("• 翻译过程中请保持网络连接稳定")
	fmt.Println("• 文章翻译会使用缓存加速重复内容的翻译")

	switch mode {
	case "missing":
		fmt.Printf("• 将为 %d 篇文章补充缺失的语言翻译\n", missingCount)
	case "update":
		fmt.Printf("• 将重新翻译 %d 篇已有翻译的文章\n", existingCount)
	case "all":
		fmt.Printf("• 将处理 %d 篇文章的翻译（包括新增和更新）\n", missingCount+existingCount)
	}
}

// PreviewArticleTranslations 预览文章翻译状态（简化版本）
func (p *Processor) PreviewArticleTranslations() {
	color.Cyan("=== 文章翻译预览 ===")

	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	status, err := articleTranslator.GetTranslationStatus()
	if err != nil {
		color.Red("❌ 获取翻译状态失败: %v", err)
		return
	}

	p.displayTranslationStats(status.MissingArticles, status.ExistingArticles, status.TotalArticles)
}

package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/generator"
	"strings"

	"github.com/fatih/color"
)

func (p *Processor) TranslateArticles(reader *bufio.Reader) {
	if p.contentDir == "" {
		color.Red("❌ 内容目录未设置")
		return
	}

	// 先预览以获取统计信息
	color.Cyan("正在分析文章翻译状态...")
	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	previews, err := articleTranslator.PreviewArticleTranslations()
	if err != nil {
		color.Red("❌ 分析失败: %v", err)
		return
	}

	// 修正统计逻辑：按文章维度统计
	missingCount, existingCount, totalArticles := p.countTranslationOperationsByArticle(previews)
	p.displayTranslationStats(missingCount, existingCount, totalArticles)

	if missingCount == 0 && existingCount == 0 {
		color.Green("✅ 没有需要翻译的文章")
		return
	}

	// 选择翻译模式
	mode := p.selectTranslationMode(missingCount, existingCount, reader)
	if mode == "" {
		return
	}

	// 显示警告和确认
	p.displayTranslationWarning(mode, missingCount, existingCount)

	if !p.confirmExecution(reader, "\n确认开始翻译？(y/n): ") {
		color.Yellow("❌ 已取消翻译")
		return
	}

	color.Cyan("🚀 开始翻译文章...")
	if err := articleTranslator.TranslateArticles(mode); err != nil {
		color.Red("❌ 翻译失败: %v", err)
	}
}

// countTranslationOperationsByArticle 按文章维度统计翻译状态
func (p *Processor) countTranslationOperationsByArticle(previews []generator.ArticleTranslationPreview) (int, int, int) {
	// 按原文件路径分组
	articleGroups := make(map[string][]generator.ArticleTranslationPreview)
	for _, preview := range previews {
		articleGroups[preview.OriginalFile] = append(articleGroups[preview.OriginalFile], preview)
	}

	missingCount := 0  // 有缺失翻译的文章数
	existingCount := 0 // 所有翻译都存在的文章数
	totalArticles := len(articleGroups)

	for _, group := range articleGroups {
		hasMissing := false
		hasExisting := false

		for _, preview := range group {
			if preview.Status == "missing" {
				hasMissing = true
			} else if preview.Status == "exists" {
				hasExisting = true
			}
		}

		// 如果有任何语言缺失翻译，则算作需要翻译的文章
		if hasMissing {
			missingCount++
		} else if hasExisting {
			// 只有当所有语言都存在时，才算作已翻译的文章
			existingCount++
		}
	}

	return missingCount, existingCount, totalArticles
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
	if existingCount > 0 {
		options = append(options, fmt.Sprintf("2. 重新翻译现有文章 (%d 篇)", existingCount))
	}
	if missingCount > 0 && existingCount > 0 {
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

// PreviewArticleTranslations 预览文章翻译状态（添加详细信息）
func (p *Processor) PreviewArticleTranslations() {
	color.Cyan("=== 文章翻译预览 ===")

	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	previews, err := articleTranslator.PreviewArticleTranslations()
	if err != nil {
		color.Red("❌ 获取翻译预览失败: %v", err)
		return
	}

	// 按文章分组显示详细状态
	p.displayDetailedTranslationStatus(previews)
}

// displayDetailedTranslationStatus 显示详细的翻译状态
func (p *Processor) displayDetailedTranslationStatus(previews []generator.ArticleTranslationPreview) {
	// 按原文件路径分组
	articleGroups := make(map[string][]generator.ArticleTranslationPreview)
	for _, preview := range previews {
		articleGroups[preview.OriginalFile] = append(articleGroups[preview.OriginalFile], preview)
	}

	fmt.Printf("\n📋 详细翻译状态 (共 %d 篇文章):\n", len(articleGroups))
	fmt.Println("======================================")

	articleIndex := 1
	for originalFile, group := range articleGroups {
		// 提取文章标题（去掉语言后缀）
		title := group[0].Title
		if len(group) > 0 {
			// 去掉标题中的语言标识
			titleParts := strings.Split(title, " (")
			if len(titleParts) > 0 {
				title = titleParts[0]
			}
		}

		fmt.Printf("\n%d. 📄 %s\n", articleIndex, title)
		fmt.Printf("   📁 %s\n", originalFile)

		// 显示各语言状态
		missingLangs := []string{}
		existingLangs := []string{}

		for _, preview := range group {
			// 提取语言标识
			langParts := strings.Split(preview.Title, " (")
			if len(langParts) > 1 {
				lang := strings.TrimRight(langParts[1], ")")
				if preview.Status == "missing" {
					missingLangs = append(missingLangs, lang)
				} else {
					existingLangs = append(existingLangs, lang)
				}
			}
		}

		if len(existingLangs) > 0 {
			fmt.Printf("   ✅ 已翻译: %s\n", strings.Join(existingLangs, ", "))
		}
		if len(missingLangs) > 0 {
			fmt.Printf("   ❌ 缺失翻译: %s\n", strings.Join(missingLangs, ", "))
		}

		articleIndex++
	}

	// 显示汇总统计
	missingCount, existingCount, totalArticles := p.countTranslationOperationsByArticle(previews)
	fmt.Printf("\n📊 汇总统计:\n")
	fmt.Printf("   🆕 需要翻译: %d 篇文章\n", missingCount)
	fmt.Printf("   ✅ 已完全翻译: %d 篇文章\n", existingCount)
	fmt.Printf("   📦 文章总数: %d 篇\n", totalArticles)
}

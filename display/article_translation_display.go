package display

import (
	"fmt"
	"strings"
	"tag-scanner/generator"

	"github.com/fatih/color"
)

// DisplayArticleTranslationPreview 显示文章翻译预览
func DisplayArticleTranslationPreview(previews []generator.ArticleTranslationPreview, limit int) {
	if len(previews) == 0 {
		fmt.Println("没有找到任何文章")
		return
	}

	// 统计信息
	missingCount := 0
	existingCount := 0
	totalWords := 0
	totalParagraphs := 0

	for _, preview := range previews {
		if preview.Status == "missing" {
			missingCount++
		} else {
			existingCount++
		}
		totalWords += preview.WordCount
		totalParagraphs += preview.ParagraphCount
	}

	// 显示统计
	color.Cyan("📊 文章翻译统计:")
	fmt.Printf("   📝 总文章数: %d 篇\n", len(previews))
	fmt.Printf("   🆕 需要翻译: %d 篇\n", missingCount)
	fmt.Printf("   ✅ 已有英文版: %d 篇\n", existingCount)
	fmt.Printf("   📄 总字数: %d 词\n", totalWords)
	fmt.Printf("   📋 总段落数: %d 段\n", totalParagraphs)
	fmt.Println()

	// 显示详细列表
	displayCount := limit
	if len(previews) < limit {
		displayCount = len(previews)
	}

	color.Cyan("📋 文章翻译预览 (显示前%d篇):", displayCount)
	fmt.Printf("%-4s %-50s %-8s %-8s %-8s %-15s\n", "序号", "文章标题", "状态", "字数", "段落数", "预计时间")
	fmt.Println(strings.Repeat("-", 95))

	for i := 0; i < displayCount; i++ {
		preview := previews[i]

		// 截断标题显示
		title := preview.Title
		if len(title) > 45 {
			title = title[:42] + "..."
		}

		// 状态颜色
		var statusDisplay string
		if preview.Status == "missing" {
			statusDisplay = color.RedString("需要翻译")
		} else {
			statusDisplay = color.GreenString("已存在")
		}

		fmt.Printf("%-4d %-50s %-8s %-8d %-8d %-15s\n",
			i+1,
			title,
			statusDisplay,
			preview.WordCount,
			preview.ParagraphCount,
			preview.EstimatedTime,
		)
	}

	if len(previews) > limit {
		fmt.Printf("\n... 还有 %d 篇文章未显示\n", len(previews)-limit)
	}

	fmt.Println()
	color.Yellow("💡 提示:")
	fmt.Println("• 英文文件将保存为 index.en.md")
	fmt.Println("• 翻译时间基于段落数估算，实际时间可能因网络状况而异")
	fmt.Println("• 代码块、链接等特殊内容将保持原样不翻译")
}

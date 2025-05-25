package display

import (
	"fmt"
	"hugo-content-suite/generator"
	"os"
	"path/filepath"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// DisplayArticleTranslationPreview 显示文章翻译预览
func DisplayArticleTranslationPreview(previews []generator.ArticleTranslationPreview, limit int) {
	headerColor.Println("=== 文章翻译预览 ===")

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

	// 显示统计信息
	fmt.Printf("\n📊 翻译统计信息:\n")
	fmt.Printf("   📝 总文章数: %d 篇\n", len(previews))
	fmt.Printf("   🆕 需要翻译: %d 篇\n", missingCount)
	fmt.Printf("   ✅ 已有英文版: %d 篇\n", existingCount)
	fmt.Printf("   📄 总字数: %d 词\n", totalWords)
	fmt.Printf("   📋 总段落数: %d 段\n", totalParagraphs)

	// 创建表格
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"序号", "文件名", "文章标题", "状态", "字数", "段落数", "预计时间"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetColWidth(60) // 设置列宽

	displayCount := limit
	if len(previews) < limit {
		displayCount = len(previews)
	}

	for i := 0; i < displayCount; i++ {
		preview := previews[i]

		// 获取文件名
		fileName := filepath.Base(preview.OriginalFile)

		// 截断标题显示
		title := preview.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		// 状态颜色
		var statusDisplay string
		if preview.Status == "missing" {
			statusDisplay = highColor.Sprint("需要翻译")
		} else {
			statusDisplay = lowColor.Sprint("已存在")
		}

		table.Append([]string{
			strconv.Itoa(i + 1),
			fileName,
			title,
			statusDisplay,
			strconv.Itoa(preview.WordCount),
			strconv.Itoa(preview.ParagraphCount),
			preview.EstimatedTime,
		})
	}

	table.Render()

	if len(previews) > limit {
		fmt.Printf("\n... 还有 %d 篇文章未显示\n", len(previews)-limit)
	}

	// 显示提示信息
	fmt.Printf("\n💡 翻译说明:\n")
	fmt.Printf("   📁 英文文件将保存为: index.en.md\n")
	fmt.Printf("   ⏱️  翻译时间基于段落数估算，实际时间可能因网络状况而异\n")
	fmt.Printf("   🔒 代码块、链接等特殊内容将保持原样不翻译\n")
	fmt.Printf("   🌐 翻译结果不会缓存，每次都是实时翻译\n")
	fmt.Println()
}

package display

import (
	"fmt"
	"hugo-content-suite/generator"
	"hugo-content-suite/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	titleColor  = color.New(color.FgCyan, color.Bold)
	headerColor = color.New(color.FgGreen, color.Bold)
	highColor   = color.New(color.FgRed, color.Bold)
	mediumColor = color.New(color.FgYellow, color.Bold)
	lowColor    = color.New(color.FgBlue)
)

// BulkTranslationPreview 批量翻译预览信息
type BulkTranslationPreview struct {
	TotalTags           int
	TotalArticles       int
	CachedCount         int
	MissingTranslations []string
	TagsToTranslate     []TranslationItem
	ArticlesToTranslate []TranslationItem
}

// TranslationItem 翻译项目信息
type TranslationItem struct {
	Type     string // "标签" 或 "文章"
	Original string
	Count    int
}

func DisplaySummary(articlesCount int, tagStats []models.TagStats, categoryStats []models.CategoryStats) {
	titleColor.Println("=== Hugo 博客管理工具 ===")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"统计项", "数量"})
	table.SetBorder(true)
	table.SetRowLine(true)

	data := [][]string{
		{"总文章数", strconv.Itoa(articlesCount)},
		{"标签总数", strconv.Itoa(len(tagStats))},
		{"分类总数", strconv.Itoa(len(categoryStats))},
	}

	high, medium, low := groupTagsByFrequency(tagStats)
	data = append(data, []string{"高频标签 (≥5篇)", strconv.Itoa(len(high))})
	data = append(data, []string{"中频标签 (2-4篇)", strconv.Itoa(len(medium))})
	data = append(data, []string{"低频标签 (1篇)", strconv.Itoa(len(low))})

	table.AppendBulk(data)
	table.Render()

	color.Green("\n💡 使用 '一键处理全部' 功能可自动完成所有必要的博客管理任务")
	fmt.Println()
}

func DisplayTagStats(tagStats []models.TagStats, limit int) {
	headerColor.Println("=== 标签使用统计 ===")

	if len(tagStats) == 0 {
		fmt.Println("没有找到任何标签")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"排名", "标签名", "使用次数", "频率级别"})
	table.SetBorder(true)
	table.SetRowLine(true)

	displayCount := limit
	if len(tagStats) < limit {
		displayCount = len(tagStats)
	}

	for i := 0; i < displayCount; i++ {
		stat := tagStats[i]
		rank := strconv.Itoa(i + 1)
		count := strconv.Itoa(stat.Count)

		var level string
		var levelColor *color.Color
		if stat.Count >= 5 {
			level = "高频"
			levelColor = highColor
		} else if stat.Count >= 2 {
			level = "中频"
			levelColor = mediumColor
		} else {
			level = "低频"
			levelColor = lowColor
		}

		table.Append([]string{
			rank,
			stat.Name,
			count,
			levelColor.Sprint(level),
		})
	}

	table.Render()

	if len(tagStats) > limit {
		fmt.Printf("... 还有 %d 个标签未显示\n", len(tagStats)-limit)
	}
	fmt.Println()
}

func DisplayCategoryStats(categoryStats []models.CategoryStats) {
	headerColor.Println("=== 分类统计 ===")

	if len(categoryStats) == 0 {
		fmt.Println("没有找到任何分类")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"排名", "分类名", "文章数量", "占比"})
	table.SetBorder(true)
	table.SetRowLine(true)

	total := 0
	for _, stat := range categoryStats {
		total += stat.Count
	}

	for i, stat := range categoryStats {
		rank := strconv.Itoa(i + 1)
		count := strconv.Itoa(stat.Count)
		percentage := fmt.Sprintf("%.1f%%", float64(stat.Count)/float64(total)*100)

		table.Append([]string{
			rank,
			stat.Name,
			count,
			percentage,
		})
	}

	table.Render()
	fmt.Println()
}

func DisplayNoTagArticles(articles []models.Article, limit int) {
	if len(articles) == 0 {
		return
	}

	titleColor.Printf("=== 无标签文章 (%d篇) ===\n", len(articles))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"序号", "文件名", "标题"})
	table.SetBorder(true)
	table.SetRowLine(true)

	displayCount := limit
	if len(articles) < limit {
		displayCount = len(articles)
	}

	for i := 0; i < displayCount; i++ {
		article := articles[i]
		num := strconv.Itoa(i + 1)
		fileName := filepath.Base(article.FilePath)
		title := article.Title
		if title == "" {
			title = "无标题"
		}

		table.Append([]string{num, fileName, title})
	}

	table.Render()

	if len(articles) > limit {
		fmt.Printf("... 还有 %d 篇无标签文章未显示\n", len(articles)-limit)
	}
	fmt.Println()
}

func DisplayTagDetails(tagStats []models.TagStats, tagName string) {
	var targetStat *models.TagStats
	for _, stat := range tagStats {
		if strings.EqualFold(stat.Name, tagName) {
			targetStat = &stat
			break
		}
	}

	if targetStat == nil {
		fmt.Printf("未找到标签: %s\n", tagName)
		return
	}

	titleColor.Printf("=== 标签详情: %s ===\n", targetStat.Name)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"序号", "文件名", "文件路径"})
	table.SetBorder(true)
	table.SetRowLine(true)

	for i, filePath := range targetStat.Files {
		num := strconv.Itoa(i + 1)
		fileName := filepath.Base(filePath)

		table.Append([]string{num, fileName, filePath})
	}

	table.Render()
	fmt.Println()
}

// DisplayTagPagePreview 显示标签页面生成预览
func DisplayTagPagePreview(previews []generator.TagPagePreview, limit int) {
	headerColor.Println("=== 标签页面生成预览 ===")

	if len(previews) == 0 {
		fmt.Println("没有找到任何标签")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"标签名", "当前Slug", "新Slug", "文章数", "状态"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetColWidth(40) // 设置列宽

	displayCount := limit
	if len(previews) < limit {
		displayCount = len(previews)
	}

	createCount := 0
	updateCount := 0

	for i := 0; i < displayCount; i++ {
		preview := previews[i]

		currentSlug := preview.ExistingSlug
		if currentSlug == "" {
			currentSlug = "无"
		}

		var statusColor *color.Color
		var statusText string
		switch preview.Status {
		case "create":
			statusColor = highColor
			statusText = "新建"
			createCount++
		case "update":
			statusColor = mediumColor
			statusText = "更新"
			updateCount++
		default:
			statusColor = lowColor
			statusText = "跳过"
		}

		table.Append([]string{
			preview.TagName,
			currentSlug,
			preview.Slug,
			strconv.Itoa(preview.ArticleCount),
			statusColor.Sprint(statusText),
		})
	}

	table.Render()

	if len(previews) > limit {
		fmt.Printf("... 还有 %d 个标签未显示\n", len(previews)-limit)
	}

	// 显示统计信息
	fmt.Printf("\n统计信息:\n")
	fmt.Printf("- 需要新建页面: %d 个\n", createCount)
	fmt.Printf("- 需要更新页面: %d 个\n", updateCount)
	fmt.Printf("- 总计处理: %d 个\n", len(previews))

	// 显示生成路径信息
	if len(previews) > 0 {
		fmt.Printf("\n生成路径示例:\n")
		fmt.Printf("- 目录: content/tags/[标签名]/\n")
		fmt.Printf("- 文件: content/tags/[标签名]/_index.md\n")
	}
	fmt.Println()
}

// DisplayArticleSlugPreview 显示文章slug预览
func DisplayArticleSlugPreview(previews []generator.ArticleSlugPreview, limit int) {
	headerColor.Println("=== 文章Slug生成预览 ===")

	if len(previews) == 0 {
		fmt.Println("没有找到需要处理的文章")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"文件名", "标题", "当前Slug", "新Slug", "状态"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetColWidth(50) // 设置列宽以适应长标题

	displayCount := limit
	if len(previews) < limit {
		displayCount = len(previews)
	}

	missingCount := 0
	updateCount := 0

	for i := 0; i < displayCount; i++ {
		preview := previews[i]
		fileName := filepath.Base(preview.FilePath)

		// 截断过长的标题
		title := preview.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		currentSlug := preview.CurrentSlug
		if currentSlug == "" {
			currentSlug = "无"
		}

		var statusColor *color.Color
		var statusText string
		switch preview.Status {
		case "missing":
			statusColor = highColor
			statusText = "新增"
			missingCount++
		case "update":
			statusColor = mediumColor
			statusText = "更新"
			updateCount++
		default:
			statusColor = lowColor
			statusText = "跳过"
		}

		table.Append([]string{
			fileName,
			title,
			currentSlug,
			preview.NewSlug,
			statusColor.Sprint(statusText),
		})
	}

	table.Render()

	if len(previews) > limit {
		fmt.Printf("... 还有 %d 篇文章未显示\n", len(previews)-limit)
	}

	// 显示统计信息
	fmt.Printf("\n统计信息:\n")
	fmt.Printf("- 需要新增slug: %d 篇\n", missingCount)
	fmt.Printf("- 需要更新slug: %d 篇\n", updateCount)
	fmt.Printf("- 总计处理: %d 篇\n", len(previews))
	fmt.Println()
}

// DisplayBulkTranslationPreview 显示批量翻译预览
func DisplayBulkTranslationPreview(preview *BulkTranslationPreview, limit int) {
	headerColor.Println("=== 全量翻译缓存预览 ===")

	// 显示总体统计
	fmt.Printf("\n📊 总体统计:\n")
	fmt.Printf("   🏷️  标签数量: %d 个\n", preview.TotalTags)
	fmt.Printf("   📝 文章数量: %d 篇\n", preview.TotalArticles)
	fmt.Printf("   ✅ 已缓存: %d 个\n", preview.CachedCount)
	fmt.Printf("   🔄 需翻译: %d 个\n", len(preview.MissingTranslations))

	if len(preview.MissingTranslations) == 0 {
		color.Green("✅ 所有内容都已有翻译缓存")
		return
	}

	// 显示需要翻译的标签
	if len(preview.TagsToTranslate) > 0 {
		fmt.Printf("\n🏷️  需要翻译的标签 (显示前%d个):\n", min(limit/2, len(preview.TagsToTranslate)))

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"标签名", "使用次数", "状态"})
		table.SetBorder(true)
		table.SetRowLine(true)

		displayCount := min(limit/2, len(preview.TagsToTranslate))
		for i := 0; i < displayCount; i++ {
			item := preview.TagsToTranslate[i]
			table.Append([]string{
				item.Original,
				strconv.Itoa(item.Count),
				highColor.Sprint("待翻译"),
			})
		}
		table.Render()

		if len(preview.TagsToTranslate) > displayCount {
			fmt.Printf("... 还有 %d 个标签未显示\n", len(preview.TagsToTranslate)-displayCount)
		}
	}

	// 显示需要翻译的文章
	if len(preview.ArticlesToTranslate) > 0 {
		fmt.Printf("\n📝 需要翻译的文章标题 (显示前%d个):\n", min(limit/2, len(preview.ArticlesToTranslate)))

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"文章标题", "状态"})
		table.SetBorder(true)
		table.SetRowLine(true)

		displayCount := min(limit/2, len(preview.ArticlesToTranslate))
		for i := 0; i < displayCount; i++ {
			item := preview.ArticlesToTranslate[i]

			// 截断过长的标题
			title := item.Original
			if len(title) > 50 {
				title = title[:47] + "..."
			}

			table.Append([]string{
				title,
				highColor.Sprint("待翻译"),
			})
		}
		table.Render()

		if len(preview.ArticlesToTranslate) > displayCount {
			fmt.Printf("... 还有 %d 篇文章未显示\n", len(preview.ArticlesToTranslate)-displayCount)
		}
	}

	fmt.Printf("\n💡 提示: 生成全量缓存后，后续的预览和生成操作将显著加快\n")
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateSlugPreview(tag string) string {
	translations := map[string]string{
		"人工智能":       "artificial-intelligence",
		"机器学习":       "machine-learning",
		"深度学习":       "deep-learning",
		"JavaScript": "javascript",
		"Python":     "python",
		"Go":         "golang",
		"技术":         "technology",
		"教程":         "tutorial",
	}

	if english, exists := translations[tag]; exists {
		return english
	}

	return strings.ToLower(strings.ReplaceAll(tag, " ", "-"))
}

func groupTagsByFrequency(tagStats []models.TagStats) (high, medium, low []models.TagStats) {
	for _, stat := range tagStats {
		if stat.Count >= 5 {
			high = append(high, stat)
		} else if stat.Count >= 2 {
			medium = append(medium, stat)
		} else {
			low = append(low, stat)
		}
	}
	return
}

package display

import (
	"fmt"
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

	color.Green("\n💡 使用 '一键处理全部' 功能可自动执行所有必要的博客管理任务 (无需确认)")
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

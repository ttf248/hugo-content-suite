package generator

import (
	"fmt"
	"hugo-content-suite/config"
	"strings"
	"time"
)

// FieldTranslator 字段翻译器
type FieldTranslator struct {
	translationUtils *TranslationUtils
	contentParser    *ContentParser
}

// NewFieldTranslator 创建字段翻译器
func NewFieldTranslator() *FieldTranslator {
	return &FieldTranslator{
		translationUtils: NewTranslationUtils(),
		contentParser:    NewContentParser(),
	}
}

// translateFrontMatterToLanguage 翻译前置数据到指定语言
func (a *ArticleTranslator) translateFrontMatterToLanguage(frontMatter, targetLang string) (string, error) {
	if frontMatter == "" {
		return "", nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("翻译前置数据到 %s...\n", targetLangName)

	lines := strings.Split(frontMatter, "\n")
	var translatedLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "---" {
			translatedLines = append(translatedLines, line)
			continue
		}

		// 翻译各种字段
		if translatedLine := a.translateFieldLine(line, targetLang); translatedLine != "" {
			translatedLines = append(translatedLines, translatedLine)
		} else {
			translatedLines = append(translatedLines, line)
		}
	}

	return strings.Join(translatedLines, "\n"), nil
}

// translateFieldLine 翻译字段行
func (a *ArticleTranslator) translateFieldLine(line, targetLang string) string {
	trimmedLine := strings.TrimSpace(line)

	// 处理标题字段
	if strings.HasPrefix(trimmedLine, "title:") {
		return a.translateSingleField(line, "title:", targetLang)
	}

	// 处理描述字段
	if strings.HasPrefix(trimmedLine, "description:") {
		return a.translateSingleField(line, "description:", targetLang)
	}

	// 处理副标题
	if strings.HasPrefix(trimmedLine, "subtitle:") {
		return a.translateSingleField(line, "subtitle:", targetLang)
	}

	// 处理摘要
	if strings.HasPrefix(trimmedLine, "summary:") {
		return a.translateSingleField(line, "summary:", targetLang)
	}

	// 处理slug字段
	if strings.HasPrefix(trimmedLine, "slug:") {
		return a.translateSlugField(line, targetLang)
	}

	// 处理数组字段
	if strings.HasPrefix(trimmedLine, "tags:") {
		return a.translateArrayField(line, "tags:", targetLang)
	}

	if strings.HasPrefix(trimmedLine, "categories:") {
		return a.translateArrayField(line, "categories:", targetLang)
	}

	if strings.HasPrefix(trimmedLine, "authors:") {
		return a.translateArrayField(line, "authors:", targetLang)
	}

	return ""
}

// translateSingleField 翻译单个字段
func (a *ArticleTranslator) translateSingleField(line, prefix, targetLang string) string {
	value := a.contentParser.ExtractFieldValue(line, prefix)
	if value != "" && a.translationUtils.ContainsChinese(value) {
		fmt.Printf("  %s: %s -> ", strings.TrimSuffix(prefix, ":"), value)

		// 使用缓存翻译
		translated, err := a.translationUtils.TranslateToLanguage(value, targetLang)
		if err != nil {
			fmt.Printf("翻译失败\n")
			return ""
		} else {
			translated = a.translationUtils.RemoveQuotes(translated)
			fmt.Printf("%s\n", translated)
			return fmt.Sprintf("%s \"%s\"", prefix, translated)
		}
	}
	return ""
}

// translateSlugField 翻译slug字段
func (a *ArticleTranslator) translateSlugField(line, targetLang string) string {
	slug := a.contentParser.ExtractFieldValue(line, "slug:")
	if slug != "" && a.translationUtils.ContainsChinese(slug) {
		fmt.Printf("  slug: %s -> ", slug)

		// 使用缓存翻译
		translated, err := a.translationUtils.TranslateToLanguage(slug, targetLang)
		if err != nil {
			fmt.Printf("翻译失败\n")
			return ""
		} else {
			translated = a.translationUtils.RemoveQuotes(translated)
			translated = a.translationUtils.FormatSlugField(translated)
			fmt.Printf("%s\n", translated)
			return fmt.Sprintf("slug: \"%s\"", translated)
		}
	}
	return ""
}

// translateArrayField 翻译数组字段
func (a *ArticleTranslator) translateArrayField(line, prefix, targetLang string) string {
	items := a.contentParser.ExtractArrayField(line, prefix)
	if len(items) > 0 {
		translatedItems := a.translateArrayItems(items, strings.TrimSuffix(prefix, ":"), targetLang)
		return fmt.Sprintf("%s %s", prefix, a.contentParser.FormatArrayField(translatedItems))
	}
	return ""
}

// translateArrayItems 翻译数组项目
func (a *ArticleTranslator) translateArrayItems(items []string, fieldType, targetLang string) []string {
	var translated []string

	fmt.Printf("  %s: ", fieldType)

	for _, item := range items {
		if a.translationUtils.ContainsChinese(item) {
			fmt.Printf("%s -> ", item)

			// 使用缓存翻译
			translatedItem, err := a.translationUtils.TranslateToLanguage(item, targetLang)
			if err != nil {
				fmt.Printf("失败 ")
				translated = append(translated, item)
			} else {
				translatedItem = a.translationUtils.RemoveQuotes(translatedItem)
				fmt.Printf("%s ", translatedItem)
				translated = append(translated, translatedItem)
			}
		} else {
			translated = append(translated, item)
		}
	}

	fmt.Printf("\n")
	return translated
}

// translateArticleBodyToLanguage 翻译正文到指定语言
func (a *ArticleTranslator) translateArticleBodyToLanguage(body, targetLang string) (string, error) {
	if strings.TrimSpace(body) == "" {
		return body, nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	totalChars := len(body)
	fmt.Printf("\n翻译正文到 %s (总计 %d 字符)...\n", targetLangName, totalChars)

	return a.translateContentByLinesToLanguage(body, targetLang)
}

// translateContentByLinesToLanguage 按行翻译内容到指定语言
func (a *ArticleTranslator) translateContentByLinesToLanguage(content, targetLang string) (string, error) {
	cfg := config.GetGlobalConfig()
	lines := strings.Split(content, "\n")
	var result []string

	// 翻译统计信息
	totalChars := len(content)
	translatedChars := 0
	startTime := time.Now()

	// 效率统计
	translationTimes := make([]time.Duration, 0) // 记录每次翻译的时间
	translationChars := make([]int, 0)           // 记录每次翻译的字符数
	apiCallCount := 0
	successCount := 0

	inCodeBlock := false
	translationCount := 0
	needsTranslationCount := 0

	// 预扫描计算需要翻译的行数
	for _, line := range lines {
		if !inCodeBlock && strings.TrimSpace(line) != "" && a.translationUtils.ContainsChinese(line) {
			if !strings.HasPrefix(strings.TrimSpace(line), "```") {
				needsTranslationCount++
			}
		}
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
		}
	}

	fmt.Printf("检测到 %d 行需要翻译\n", needsTranslationCount)

	inCodeBlock = false // 重置状态

	for _, line := range lines {
		// 检测代码块
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			translatedChars += len(line) + 1
			continue
		}

		// 代码块内容直接保留
		if inCodeBlock {
			result = append(result, line)
			translatedChars += len(line) + 1
			continue
		}

		// 空行直接保留
		if strings.TrimSpace(line) == "" {
			result = append(result, line)
			translatedChars += len(line) + 1
			continue
		}

		// 检查是否包含中文
		if !a.translationUtils.ContainsChinese(line) {
			result = append(result, line)
			translatedChars += len(line) + 1
			continue
		}

		// 需要翻译的行
		translationCount++
		lineStartTime := time.Now()
		lineChars := len(line)

		// 生成增强进度条
		progressBar := a.generateEnhancedProgressBar(translationCount, needsTranslationCount, 25)

		// 计算当前阶段 (每10行为一个阶段)
		currentStage := (translationCount-1)/10 + 1
		totalStages := (needsTranslationCount + 9) / 10

		fmt.Printf("  [%d/%d] Stage%d/%d %s 翻译 %d 字符...\n",
			translationCount, needsTranslationCount, currentStage, totalStages, progressBar, lineChars)

		translatedLine, err := a.translateSingleLineToLanguage(line, translationCount, targetLang)
		apiCallCount++

		if err != nil {
			fmt.Printf("    ❌ 翻译失败: %v\n", err)
			result = append(result, line)
			translatedChars += len(line) + 1
		} else {
			successCount++
			lineEndTime := time.Now()
			lineDuration := lineEndTime.Sub(lineStartTime)
			translatedChars += len(line) + 1

			// 记录翻译数据
			translationTimes = append(translationTimes, lineDuration)
			translationChars = append(translationChars, lineChars)

			// 计算多种效率指标
			totalDuration := lineEndTime.Sub(startTime)

			// 1. 实时效率 (当前行)
			realtimeEfficiency := float64(lineChars) / lineDuration.Seconds()

			// 2. 总体平均效率
			avgEfficiency := float64(translatedChars) / totalDuration.Seconds()

			// 3. 滑动窗口效率 (最近5次)
			recentEfficiency := a.calculateRecentEfficiency(translationTimes, translationChars, 5)

			// 4. 阶段效率 (最近10次)
			stageEfficiency := a.calculateRecentEfficiency(translationTimes, translationChars, 10)

			// 计算剩余估算 (使用多种方法)
			remainingChars := totalChars - translatedChars
			remainingLines := needsTranslationCount - translationCount

			// 智能预估：根据效率稳定性选择最佳预估方法
			estimatedTime := a.calculateSmartEstimate(remainingChars, avgEfficiency, recentEfficiency, stageEfficiency)

			// 进度计算
			lineProgress := float64(translationCount) * 100.0 / float64(needsTranslationCount)
			charProgress := float64(translatedChars) * 100.0 / float64(totalChars)
			successRate := float64(successCount) * 100.0 / float64(apiCallCount)

			// 效率趋势分析
			trend := a.calculateEfficiencyTrend(translationTimes, translationChars)

			// 清屏并重新显示 (动态刷新效果)
			if translationCount > 1 {
				fmt.Print("\033[6A\033[K") // 上移6行并清除
			}

			// 显示详细统计信息
			fmt.Printf("    ✅ 完成 (%.1fs) | API调用 #%d\n", lineDuration.Seconds(), apiCallCount)
			fmt.Printf("    📊 进度: 行 %.1f%% (%d/%d) | 字符 %.1f%% (%d/%d)\n",
				lineProgress, translationCount, needsTranslationCount,
				charProgress, translatedChars, totalChars)
			fmt.Printf("    ⚡ 效率: 实时%.1f | 平均%.1f | 最近%.1f | 阶段%.1f 字符/秒 %s\n",
				realtimeEfficiency, avgEfficiency, recentEfficiency, stageEfficiency, trend)
			fmt.Printf("    🎯 成功率: %.1f%% (%d/%d) | 剩余: %d行\n",
				successRate, successCount, apiCallCount, remainingLines)
			fmt.Printf("    ⏱️  预估剩余: %v | 预计完成: %v\n",
				estimatedTime.Round(time.Second),
				time.Now().Add(estimatedTime).Format("15:04:05"))
			fmt.Printf("    💾 处理速度: %.1f 行/分钟 | 总用时: %v\n",
				float64(translationCount)*60.0/totalDuration.Minutes(),
				totalDuration.Round(time.Second))

			result = append(result, translatedLine)

			// 分段统计报告 (每10行输出一次汇总)
			if translationCount%10 == 0 {
				a.printStageReport(translationCount, totalDuration, successRate, recentEfficiency)
			}
		}

		// 添加延迟避免API频率限制
		if cfg.Translation.DelayBetweenMs > 0 {
			time.Sleep(time.Duration(cfg.Translation.DelayBetweenMs) * time.Millisecond)
		}
	}

	// 输出最终统计信息
	totalDuration := time.Since(startTime)
	avgCharsPerSecond := float64(totalChars) / totalDuration.Seconds()
	finalSuccessRate := float64(successCount) * 100.0 / float64(apiCallCount)

	fmt.Printf("\n🎉 翻译完成！\n")
	fmt.Printf("   ⏱️  总用时: %v\n", totalDuration.Round(time.Second))
	fmt.Printf("   📈 平均效率: %.1f 字符/秒\n", avgCharsPerSecond)
	fmt.Printf("   📊 成功率: %.1f%% (%d/%d)\n", finalSuccessRate, successCount, apiCallCount)
	fmt.Printf("   📝 处理: %d 字符, %d 行翻译\n", totalChars, needsTranslationCount)

	return strings.Join(result, "\n"), nil
}

// generateProgressBar 生成进度条
func (a *ArticleTranslator) generateProgressBar(current, total, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}

	progress := float64(current) / float64(total)
	completed := int(progress * float64(width))

	bar := strings.Repeat("█", completed)
	remaining := strings.Repeat("░", width-completed)

	return fmt.Sprintf("[%s%s]", bar, remaining)
}

// calculateRecentEfficiency 计算最近N次翻译的效率
func (a *ArticleTranslator) calculateRecentEfficiency(times []time.Duration, chars []int, windowSize int) float64 {
	if len(times) == 0 {
		return 0
	}

	start := len(times) - windowSize
	if start < 0 {
		start = 0
	}

	var totalTime time.Duration
	var totalChars int

	for i := start; i < len(times); i++ {
		totalTime += times[i]
		totalChars += chars[i]
	}

	if totalTime.Seconds() == 0 {
		return 0
	}

	return float64(totalChars) / totalTime.Seconds()
}

// translateSingleLineToLanguage 翻译单行内容到指定语言
func (a *ArticleTranslator) translateSingleLineToLanguage(line string, lineNum int, targetLang string) (string, error) {
	// 使用缓存翻译
	return a.translationUtils.TranslateToLanguage(line, targetLang)
}

// generateEnhancedProgressBar 生成增强进度条
func (a *ArticleTranslator) generateEnhancedProgressBar(current, total, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}

	progress := float64(current) / float64(total)
	completed := int(progress * float64(width))

	// 使用不同字符表示不同进度段
	var bar strings.Builder
	for i := 0; i < width; i++ {
		if i < completed {
			if i < width/4 {
				bar.WriteString("█") // 25%以下用实心
			} else if i < width/2 {
				bar.WriteString("▓") // 25%-50%用深色
			} else if i < width*3/4 {
				bar.WriteString("▒") // 50%-75%用中色
			} else {
				bar.WriteString("░") // 75%-100%用浅色
			}
		} else {
			bar.WriteString("░")
		}
	}

	return fmt.Sprintf("[%s] %.1f%%", bar.String(), progress*100)
}

// calculateSmartEstimate 智能预估剩余时间
func (a *ArticleTranslator) calculateSmartEstimate(remainingChars int, avgEff, recentEff, stageEff float64) time.Duration {
	if remainingChars <= 0 {
		return 0
	}

	// 权重分配：最近效率50%，阶段效率30%，平均效率20%
	smartEfficiency := recentEff*0.5 + stageEff*0.3 + avgEff*0.2

	if smartEfficiency <= 0 {
		smartEfficiency = avgEff
	}

	if smartEfficiency <= 0 {
		return time.Hour // 如果无法计算，返回1小时作为默认值
	}

	return time.Duration(float64(remainingChars)/smartEfficiency) * time.Second
}

// calculateEfficiencyTrend 计算效率趋势
func (a *ArticleTranslator) calculateEfficiencyTrend(times []time.Duration, chars []int) string {
	if len(times) < 3 {
		return "📈"
	}

	// 比较最近3次和之前3次的效率
	recentEfficiency := a.calculateRecentEfficiency(times, chars, 3)
	prevEfficiency := a.calculateRecentEfficiency(times[:len(times)-3], chars[:len(chars)-3], 3)

	if recentEfficiency > prevEfficiency*1.1 {
		return "📈" // 上升
	} else if recentEfficiency < prevEfficiency*0.9 {
		return "📉" // 下降
	}
	return "📊" // 稳定
}

// printStageReport 打印阶段报告
func (a *ArticleTranslator) printStageReport(currentCount int, totalDuration time.Duration, successRate, efficiency float64) {
	stage := currentCount / 10
	fmt.Printf("\n    🏁 阶段 %d 完成 | 总计 %d 行 | 用时 %v | 成功率 %.1f%% | 效率 %.1f 字符/秒\n",
		stage, currentCount, totalDuration.Round(time.Second), successRate, efficiency)
	fmt.Printf("    ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔\n")
}

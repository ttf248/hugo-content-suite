package generator

import (
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/translator"
	"hugo-content-suite/utils"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// FieldTranslator 字段翻译器
type FieldTranslator struct {
	translationUtils *translator.TranslationUtils
	contentParser    *ContentParser
}

// NewFieldTranslator 创建字段翻译器
func NewFieldTranslator() *FieldTranslator {
	return &FieldTranslator{
		translationUtils: translator.NewTranslationUtils(),
		contentParser:    NewContentParser(),
	}
}

// translateFrontMatterToLanguage 翻译前置数据到指定语言
func (a *ArticleTranslator) translateFrontMatterToLanguage(frontMatter, targetLang string) (string, error) {
	if strings.TrimSpace(frontMatter) == "" {
		return frontMatter, nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("翻译前置数据到 %s...\n", targetLangName)

	// 解析 YAML
	var frontMatterData map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontMatter), &frontMatterData); err != nil {
		return "", fmt.Errorf("解析前置数据失败: %v", err)
	}

	// 翻译各个字段
	translatedData, err := a.translateFrontMatterFields(frontMatterData, targetLang)
	if err != nil {
		return "", fmt.Errorf("翻译前置数据字段失败: %v", err)
	}

	// 将翻译后的数据转换回 YAML
	translatedYAML, err := yaml.Marshal(translatedData)
	if err != nil {
		return "", fmt.Errorf("生成翻译后的YAML失败: %v", err)
	}

	return "---\r\n" + string(translatedYAML) + "---\r\n", nil
}

// translateFrontMatterFields 翻译前置数据的所有字段
func (a *ArticleTranslator) translateFrontMatterFields(data map[string]interface{}, targetLang string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 定义需要翻译的字段
	translatableFields := map[string]bool{
		"title":       true,
		"description": true,
	}

	// 定义需要翻译的数组字段
	translatableArrayFields := map[string]bool{
		"tags":       true,
		"categories": true,
	}

	for key, value := range data {
		switch {
		case translatableFields[key]:
			// 翻译单个字符串字段
			if strValue, ok := value.(string); ok {
				translatedValue, err := a.translateStringField(key, strValue, targetLang)
				if err != nil {
					fmt.Printf("  警告: 翻译字段 %s 失败: %v\n", key, err)
					result[key] = value // 保持原值
				} else {
					result[key] = translatedValue
				}
			} else {
				result[key] = value
			}

		case translatableArrayFields[key]:
			// 翻译数组字段
			if arrayValue, ok := value.([]interface{}); ok {
				translatedArray, err := a.translateArrayField(key, arrayValue, targetLang)
				if err != nil {
					fmt.Printf("  警告: 翻译数组字段 %s 失败: %v\n", key, err)
					result[key] = value // 保持原值
				} else {
					result[key] = translatedArray
				}
			} else {
				result[key] = value
			}

		default:
			// 其他字段保持不变
			result[key] = value
		}
	}

	return result, nil
}

// translateStringField 翻译字符串字段
func (a *ArticleTranslator) translateStringField(fieldName, value, targetLang string) (string, error) {
	if value == "" || !utils.ContainsChinese(value) {
		return value, nil
	}

	fmt.Printf("  %s: %s -> ", fieldName, value)

	// 使用缓存翻译
	translated, err := a.translationUtils.TranslateToLanguage(value, targetLang)
	if err != nil {
		fmt.Printf("翻译失败\n")
		return value, err
	}

	fmt.Printf("%s\n", translated)
	return translated, nil
}

// translateArrayField 翻译数组字段
func (a *ArticleTranslator) translateArrayField(fieldName string, items []interface{}, targetLang string) ([]interface{}, error) {
	if len(items) == 0 {
		return items, nil
	}

	fmt.Printf("  %s: ", fieldName)

	var translatedItems []interface{}
	for _, item := range items {
		if strItem, ok := item.(string); ok {
			if utils.ContainsChinese(strItem) {
				fmt.Printf("%s -> ", strItem)

				// 使用缓存翻译
				translated, err := a.translationUtils.TranslateToLanguageWithCache(strItem, targetLang)
				if err != nil {
					fmt.Printf("失败 ")
					translatedItems = append(translatedItems, item)
					continue
				}

				fmt.Printf("%s ", translated)
				translatedItems = append(translatedItems, translated)
			} else {
				fmt.Printf("%s -> %s\t", strItem, strItem)
				translatedItems = append(translatedItems, item)
			}
		} else {
			translatedItems = append(translatedItems, item)
		}
	}

	fmt.Printf("\n")
	return translatedItems, nil
}

// printParagraphStageReport 打印段落翻译阶段报告
func (a *ArticleTranslator) printParagraphStageReport(current, total int, elapsed time.Duration, success, error int) {
	stage := (current + 9) / 10
	successRate := float64(success) * 100.0 / float64(current)
	avgTime := elapsed.Seconds() / float64(current)

	fmt.Printf("\n🏁 阶段 %d 完成 | 已翻译 %d/%d 段落\n", stage, current, total)
	fmt.Printf("   ⏱️  阶段用时: %v | 平均: %.1f 秒/段落\n",
		elapsed.Round(time.Second), avgTime)
	fmt.Printf("   📊 成功率: %.1f%% (%d 成功, %d 失败)\n", successRate, success, error)
	fmt.Printf("   ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔\n")
}

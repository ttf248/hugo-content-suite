package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/translator"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// TranslationUtils 翻译工具
type TranslationUtils struct {
	translator *translator.LLMTranslator
	cache      *translator.TranslationCache
}

// NewTranslationUtils 创建翻译工具实例
func NewTranslationUtils() *TranslationUtils {
	cache := translator.NewTranslationCache()
	cache.Load() // 加载缓存

	return &TranslationUtils{
		translator: translator.NewLLMTranslator(),
		cache:      cache,
	}
}

// TestConnection 测试连接
func (t *TranslationUtils) TestConnection() error {
	return t.translator.TestConnection()
}

// ContainsChinese 检查文本是否包含中文
func (t *TranslationUtils) ContainsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

// FormatSlugField 格式化slug字段
func (t *TranslationUtils) FormatSlugField(slug string) string {
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}

// TranslateToLanguage 翻译文本到指定语言（带缓存）
func (t *TranslationUtils) TranslateToLanguage(content, targetLang string) (string, error) {
	// 先检查缓存
	cacheKey := fmt.Sprintf("%s:%s", targetLang, content)
	if cached, found := t.cache.Get(cacheKey, translator.TagCache); found {
		return cached, nil
	}

	// 缓存未命中，调用翻译服务
	result, err := t.translateWithAPI(content, targetLang)
	if err != nil {
		return "", err
	}

	// 保存到缓存
	t.cache.Set(cacheKey, result, translator.TagCache)

	return result, nil
}

// BatchTranslateWithCache 批量翻译（优先使用缓存）
func (t *TranslationUtils) BatchTranslateWithCache(texts []string, targetLang string, cacheType translator.CacheType) (map[string]string, error) {
	result := make(map[string]string)
	var missingTexts []string

	// 检查缓存
	for _, text := range texts {
		cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
		if cached, found := t.cache.Get(cacheKey, cacheType); found {
			result[text] = cached
		} else {
			missingTexts = append(missingTexts, text)
		}
	}

	fmt.Printf("📄 缓存命中: %d 个, 需要翻译: %d 个\n",
		len(texts)-len(missingTexts), len(missingTexts))

	// 翻译缺失的文本
	if len(missingTexts) > 0 {
		for i, text := range missingTexts {
			fmt.Printf("  [%d/%d] 翻译: %s -> ", i+1, len(missingTexts), text)

			translated, err := t.translateWithAPI(text, targetLang)
			if err != nil {
				fmt.Printf("失败 (%v)\n", err)
				// 翻译失败，使用原文
				translated = text
			}

			result[text] = translated

			// 保存到缓存
			cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
			t.cache.Set(cacheKey, translated, cacheType)

			fmt.Printf("%s\n", translated)

			// 添加延迟避免API频率限制
			if i < len(missingTexts)-1 {
				time.Sleep(100 * time.Millisecond)
			}
		}

		// 保存缓存到文件
		if err := t.cache.Save(); err != nil {
			fmt.Printf("⚠️ 保存缓存失败: %v\n", err)
		}
	}

	return result, nil
}

// translateWithAPI 直接调用API翻译
func (t *TranslationUtils) translateWithAPI(content, targetLang string) (string, error) {
	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]

	if targetLangName == "" {
		targetLangName = targetLang
	}

	systemContent := `
	忽略以前设置的所有指令。
	你是一位专业的技术文档翻译人员。

	请执行以下任务：
	1. 将用户提供的中文内容准确翻译为指定语言
	2. 保持原文档的markdown格式结构不变

	仅输出翻译的内容`

	// 构建历史对话，提供翻译示例
	var messages []translator.Message

	// 系统消息
	messages = append(messages, translator.Message{
		Role:    "system",
		Content: systemContent,
	})

	// 添加历史翻译示例
	switch targetLang {
	case "en":
		messages = append(messages,
			translator.Message{Role: "user", Content: "请将以下内容翻译为 English: 人工智能"},
			translator.Message{Role: "assistant", Content: "Artificial Intelligence"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 English: 机器学习"},
			translator.Message{Role: "assistant", Content: "Machine Learning"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 English: - 数据挖掘\n- 深度学习\n- 神经网络"},
			translator.Message{Role: "assistant", Content: "- Data Mining\n- Deep Learning\n- Neural Network"},
		)
	case "ja":
		messages = append(messages,
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Japanese: 人工智能"},
			translator.Message{Role: "assistant", Content: "人工知能"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Japanese: 机器学习"},
			translator.Message{Role: "assistant", Content: "機械学習"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Japanese: - 数据挖掘\n- 深度学习\n- 神经网络"},
			translator.Message{Role: "assistant", Content: "- データマイニング\n- ディープラーニング\n- ニューラルネットワーク"},
		)
	case "ko":
		messages = append(messages,
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Korean: 人工智能"},
			translator.Message{Role: "assistant", Content: "인공지능"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Korean: 机器学习"},
			translator.Message{Role: "assistant", Content: "기계학습"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Korean: - 数据挖掘\n- 深度学习\n- 神经网络"},
			translator.Message{Role: "assistant", Content: "- 데이터 마이닝\n- 딥러닝\n- 신경망"},
		)
	case "fr":
		messages = append(messages,
			translator.Message{Role: "user", Content: "请将以下内容翻译为 French: 人工智能"},
			translator.Message{Role: "assistant", Content: "Intelligence Artificielle"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 French: 机器学习"},
			translator.Message{Role: "assistant", Content: "Apprentissage Automatique"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 French: - 数据挖掘\n- 深度学习\n- 神经网络"},
			translator.Message{Role: "assistant", Content: "- Exploration de Données\n- Apprentissage Profond\n- Réseau de Neurones"},
		)
	case "ru":
		messages = append(messages,
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Russian: 人工智能"},
			translator.Message{Role: "assistant", Content: "Искусственный интеллект"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Russian: 机器学习"},
			translator.Message{Role: "assistant", Content: "Машинное обучение"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Russian: - 数据挖掘\n- 深度学习\n- 神经网络"},
			translator.Message{Role: "assistant", Content: "- Интеллектуальный анализ данных\n- Глубокое обучение\n- Нейронная сеть"},
		)
	case "hi":
		messages = append(messages,
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Hindi: 人工智能"},
			translator.Message{Role: "assistant", Content: "कृत्रिम बुद्धिमत्ता"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Hindi: 机器学习"},
			translator.Message{Role: "assistant", Content: "मशीन लर्निंग"},
			translator.Message{Role: "user", Content: "请将以下内容翻译为 Hindi: - 数据挖掘\n- 深度学习\n- 神经网络"},
			translator.Message{Role: "assistant", Content: "- डेटा माइनिंग\n- डीप लर्निंग\n- न्यूरल नेटवर्क"},
		)
	}

	// 添加当前翻译请求
	messages = append(messages, translator.Message{
		Role:    "user",
		Content: fmt.Sprintf("请将以下内容翻译为 %s: %s", targetLangName, content),
	})

	request := translator.LMStudioRequest{
		Model:            cfg.LMStudio.Model,
		Messages:         messages,
		Stream:           false,
		Temperature:      0.0,  // 设置为 0.0 可使输出更确定，适合需要精确翻译的场景。
		TopP:             1.0,  // 与 Temperature 配合使用，设置为 1.0 表示不限制采样范围。
		MaxTokens:        1000, // 根据翻译内容的长度调整，确保输出完整。
		PresencePenalty:  0.0,  // 设置为 0.0 可防止模型引入新的话题或内容，保持翻译的忠实性。
		FrequencyPenalty: 0.0,  // 设置为 0.0 可避免模型对词汇的重复使用进行惩罚，适合保持原文结构的翻译。
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to serialize request: %v", err)
	}

	client := &http.Client{Timeout: time.Duration(cfg.LMStudio.Timeout) * time.Second}
	resp, err := client.Post(cfg.LMStudio.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LM Studio returned error status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var response translator.LMStudioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no translation result received")
	}
	result := strings.TrimSpace(response.Choices[0].Message.Content)

	// 兼容思考模型，移除 <think> </think> 标签之间的内容
	thinkRegex := regexp.MustCompile(`(?s)<think>.*?</think>`)
	result = thinkRegex.ReplaceAllString(result, "")
	result = strings.TrimSpace(result)

	return result, nil
}

// SaveCache 保存缓存
func (t *TranslationUtils) SaveCache() error {
	return t.cache.Save()
}

// GetCacheStats 获取缓存统计
func (t *TranslationUtils) GetCacheStats() string {
	return t.cache.GetInfo()
}

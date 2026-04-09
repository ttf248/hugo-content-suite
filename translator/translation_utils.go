package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hugo-content-suite/config"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// LM Studio API 相关类型定义
type LMStudioRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Stream           bool      `json:"stream"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LMStudioResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// TranslationUtils 翻译工具
type TranslationUtils struct {
	translator *LLMTranslator
	cache      *TranslationCache
}

// NewTranslationUtils 创建翻译工具实例
func NewTranslationUtils() *TranslationUtils {
	cache := NewTranslationCache()
	cache.Load() // 加载缓存

	return &TranslationUtils{
		translator: NewLLMTranslator(),
		cache:      cache,
	}
}

// TestConnection 测试与LM Studio的连接
func (t *TranslationUtils) TestConnection() error {
	fmt.Println("正在测试与LM Studio的连接...")

	cfg := config.GetGlobalConfig()
	request := LMStudioRequest{
		Model: cfg.LMStudio.Model,
		Messages: []Message{
			{Role: "user", Content: "这是一个测试请求，无需处理，直接应答就行"},
		},
		Stream: false,
	}

	_, err := t.sendRequest(request)
	return err
}

func (t *TranslationUtils) TranslateTags(texts []string) (map[string]string, error) {
	return t.batchTranslateWithCache(texts, "en", kTagCache)
}

func (t *TranslationUtils) TranslateArticlesSlugs(texts []string) (map[string]string, error) {
	return t.batchTranslateWithCache(texts, "en", kSlugCache)
}

func (t *TranslationUtils) TranslateToLanguage(content, targetLang string) (string, error) {
	result, err := t.translateWithAPI(content, targetLang)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (t *TranslationUtils) TranslateCategory(content, targetLang string) (string, error) {
	result, err := t.translateWithCache(content, targetLang, kCategoryCache)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (t *TranslationUtils) TranslateTag(content, targetLang string) (string, error) {
	result, err := t.translateWithCache(content, targetLang, kTagCache)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (t *TranslationUtils) translateWithCache(text, targetLang string, cacheType CacheType) (string, error) {
	fmt.Println("\n🤖 使用AI翻译...")
	cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
	if cached, found := t.cache.Get(cacheKey, cacheType); found {
		fmt.Printf("✅ [Cache Hit] [%s] %s\n", targetLang, text)
		return cached, nil
	}

	fmt.Printf("🚀 [API Translate] [%s] %s\n", targetLang, text)
	translated, err := t.translateWithAPI(text, targetLang)
	if err != nil {
		fmt.Printf("❌ [API Error] [%s] %s: %v\n", targetLang, text, err)
		translated = text // fallback to original text on error
	}
	t.cache.Set(cacheKey, translated, cacheType)
	_ = t.cache.Save()
	fmt.Printf("✅ [Cache Set] [%s] %s\n", targetLang, text)

	return translated, err
}

func (t *TranslationUtils) batchTranslateWithCache(texts []string, targetLang string, cacheType CacheType) (map[string]string, error) {
	fmt.Println("\n🤖 使用AI翻译...")
	result := make(map[string]string)
	var missingTexts []string
	hitCount := 0

	for _, text := range texts {
		cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
		if cached, found := t.cache.Get(cacheKey, cacheType); found {
			fmt.Printf("✅ [Batch Cache Hit] [%s] %s\n", targetLang, text)
			result[text] = cached
			hitCount++
		} else {
			fmt.Printf("🚀 [Batch API Translate] [%s] %s\n", targetLang, text)
			missingTexts = append(missingTexts, text)
		}
	}

	for _, text := range missingTexts {
		translated, err := t.translateWithAPI(text, targetLang)
		if err != nil {
			fmt.Printf("❌ [Batch API Error] [%s] %s: %v\n", targetLang, text, err)
			translated = text
		}
		result[text] = translated
		cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
		t.cache.Set(cacheKey, translated, cacheType)
		fmt.Printf("✅ [Batch Cache Set] [%s] %s\n", targetLang, text)
	}

	if len(missingTexts) > 0 {
		_ = t.cache.Save()
	}

	total := len(texts)
	hitRate := float64(hitCount) / float64(total) * 100

	fmt.Printf("📊 [Batch Cache Stats] 命中率: %.2f%% (%d/%d)\n", hitRate, hitCount, total)

	return result, nil
}

// sendRequest 发送HTTP请求的通用方法
func (t *TranslationUtils) sendRequest(request LMStudioRequest) (*LMStudioResponse, error) {
	cfg := config.GetGlobalConfig()

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %v", err)
	}

	client := &http.Client{Timeout: time.Duration(cfg.LMStudio.Timeout) * time.Second}
	resp, err := client.Post(cfg.LMStudio.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LM Studio returned error status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var response LMStudioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no translation result received")
	}

	return &response, nil
}

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
	3. 保持原文的空行、换行、标题层级、列表、引用块、表格、代码围栏和缩进结构
	4. 原样保留所有 Hugo shortcode / 模板语法，例如 {{< relref "/post/xxx" >}}、{{% xxx %}}、{{ ... }}
	5. 不要翻译 shortcode 内部内容，不要把 shortcode 里的 ASCII 双引号替换成 “ ” 等智能引号

	仅输出翻译的内容`

	// 构建历史对话，提供翻译示例
	var messages []Message

	// 系统消息
	messages = append(messages, Message{
		Role:    "system",
		Content: systemContent,
	})

	// 添加历史翻译示例
	switch targetLang {
	case "en":
		messages = append(messages,
			Message{Role: "user", Content: "请将以下内容翻译为 English: 人工智能"},
			Message{Role: "assistant", Content: "Artificial Intelligence"},
			Message{Role: "user", Content: "请将以下内容翻译为 English: 机器学习"},
			Message{Role: "assistant", Content: "Machine Learning"},
			Message{Role: "user", Content: "请将以下内容翻译为 English: - 数据挖掘\n- 深度学习\n- 神经网络"},
			Message{Role: "assistant", Content: "- Data Mining\n- Deep Learning\n- Neural Network"},
		)
	case "ja":
		messages = append(messages,
			Message{Role: "user", Content: "请将以下内容翻译为 Japanese: 人工智能"},
			Message{Role: "assistant", Content: "人工知能"},
			Message{Role: "user", Content: "请将以下内容翻译为 Japanese: 机器学习"},
			Message{Role: "assistant", Content: "機械学習"},
			Message{Role: "user", Content: "请将以下内容翻译为 Japanese: - 数据挖掘\n- 深度学习\n- 神经网络"},
			Message{Role: "assistant", Content: "- データマイニング\n- ディープラーニング\n- ニューラルネットワーク"},
		)
	case "ko":
		messages = append(messages,
			Message{Role: "user", Content: "请将以下内容翻译为 Korean: 人工智能"},
			Message{Role: "assistant", Content: "인공지능"},
			Message{Role: "user", Content: "请将以下内容翻译为 Korean: 机器学习"},
			Message{Role: "assistant", Content: "기계학습"},
			Message{Role: "user", Content: "请将以下内容翻译为 Korean: - 数据挖掘\n- 深度学习\n- 神经网络"},
			Message{Role: "assistant", Content: "- 데이터 마이닝\n- 딥러닝\n- 신경망"},
		)
	case "fr":
		messages = append(messages,
			Message{Role: "user", Content: "请将以下内容翻译为 French: 人工智能"},
			Message{Role: "assistant", Content: "Intelligence Artificielle"},
			Message{Role: "user", Content: "请将以下内容翻译为 French: 机器学习"},
			Message{Role: "assistant", Content: "Apprentissage Automatique"},
			Message{Role: "user", Content: "请将以下内容翻译为 French: - 数据挖掘\n- 深度学习\n- 神经网络"},
			Message{Role: "assistant", Content: "- Exploration de Données\n- Apprentissage Profond\n- Réseau de Neurones"},
		)
	case "ru":
		messages = append(messages,
			Message{Role: "user", Content: "请将以下内容翻译为 Russian: 人工智能"},
			Message{Role: "assistant", Content: "Искусственный интеллект"},
			Message{Role: "user", Content: "请将以下内容翻译为 Russian: 机器学习"},
			Message{Role: "assistant", Content: "Машинное обучение"},
			Message{Role: "user", Content: "请将以下内容翻译为 Russian: - 数据挖掘\n- 深度学习\n- 神经网络"},
			Message{Role: "assistant", Content: "- Интеллектуальный анализ данных\n- Глубокое обучение\n- Нейронная сеть"},
		)
	case "hi":
		messages = append(messages,
			Message{Role: "user", Content: "请将以下内容翻译为 Hindi: 人工智能"},
			Message{Role: "assistant", Content: "कृत्रिम बुद्धिमत्ता"},
			Message{Role: "user", Content: "请将以下内容翻译为 Hindi: 机器学习"},
			Message{Role: "assistant", Content: "मशीन लर्निंग"},
			Message{Role: "user", Content: "请将以下内容翻译为 Hindi: - 数据挖掘\n- 深度学习\n- 神经网络"},
			Message{Role: "assistant", Content: "- डेटा माइनिंग\n- डीप लर्निंग\n- न्यूरल नेटवर्क"},
		)
	}

	// 添加当前翻译请求
	messages = append(messages, Message{
		Role:    "user",
		Content: fmt.Sprintf("请将以下内容翻译为 %s: %s", targetLangName, content),
	})

	request := LMStudioRequest{
		Model:            cfg.LMStudio.Model,
		Messages:         messages,
		Stream:           false,
		Temperature:      0.0,  // 设置为 0.0 可使输出更确定，适合需要精确翻译的场景。
		TopP:             1.0,  // 与 Temperature 配合使用，设置为 1.0 表示不限制采样范围。
		MaxTokens:        1000, // 根据翻译内容的长度调整，确保输出完整。
		PresencePenalty:  0.0,  // 设置为 0.0 可防止模型引入新的话题或内容，保持翻译的忠实性。
		FrequencyPenalty: 0.0,  // 设置为 0.0 可避免模型对词汇的重复使用进行惩罚，适合保持原文结构的翻译。
	}

	response, err := t.sendRequest(request)
	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(response.Choices[0].Message.Content)

	// 兼容思考模型，移除 <think> </think> 标签之间的内容
	thinkRegex := regexp.MustCompile(`(?s)<think>.*?</think>`)
	result = thinkRegex.ReplaceAllString(result, "")
	result = strings.TrimSpace(result)
	result = normalizeHugoShortcodeQuotes(result)

	return result, nil
}

func normalizeHugoShortcodeQuotes(content string) string {
	shortcodeRegex := regexp.MustCompile(`\{\{[<%][\s\S]*?[%>]}}`)
	return shortcodeRegex.ReplaceAllStringFunc(content, func(shortcode string) string {
		shortcode = strings.ReplaceAll(shortcode, "“", `"`)
		shortcode = strings.ReplaceAll(shortcode, "”", `"`)
		return shortcode
	})
}

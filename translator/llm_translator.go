package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	LMStudioURL = "http://172.19.192.1:2234/v1/chat/completions"
	ModelName   = "gemma-3-12b-it"
)

type LMStudioRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
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

type LLMTranslator struct {
	client  *http.Client
	baseURL string
	model   string
}

func NewLLMTranslator() *LLMTranslator {
	return &LLMTranslator{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: LMStudioURL,
		model:   ModelName,
	}
}

// TranslateToSlug 将中文标签翻译为英文slug
func (t *LLMTranslator) TranslateToSlug(tag string) (string, error) {
	// 如果已经是英文，直接处理
	if isEnglishOnly(tag) {
		return normalizeSlug(tag), nil
	}

	// 构建提示词
	prompt := fmt.Sprintf(`请将以下中文标签翻译为适合作为URL的英文slug。要求：
1. 使用小写字母
2. 单词之间用连字符(-)连接
3. 不包含特殊字符
4. 简洁准确
5. 只返回翻译结果，不要任何解释

中文标签: %s

英文slug:`, tag)

	request := LMStudioRequest{
		Model: t.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	resp, err := t.client.Post(t.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LM Studio返回错误状态: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var response LMStudioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("没有获取到翻译结果")
	}

	slug := strings.TrimSpace(response.Choices[0].Message.Content)
	return normalizeSlug(slug), nil
}

// BatchTranslate 批量翻译标签
func (t *LLMTranslator) BatchTranslate(tags []string) (map[string]string, error) {
	result := make(map[string]string)

	for i, tag := range tags {
		fmt.Printf("正在翻译标签 (%d/%d): %s", i+1, len(tags), tag)

		slug, err := t.TranslateToSlug(tag)
		if err != nil {
			fmt.Printf(" - 失败: %v\n", err)
			// 使用fallback方法
			slug = fallbackSlug(tag)
		} else {
			fmt.Printf(" -> %s\n", slug)
		}

		result[tag] = slug

		// 添加延迟避免请求过于频繁
		if i < len(tags)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return result, nil
}

// TestConnection 测试与LM Studio的连接
func (t *LLMTranslator) TestConnection() error {
	_, err := t.TranslateToSlug("测试")
	return err
}

// isEnglishOnly 检查字符串是否只包含英文字符
func isEnglishOnly(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ') {
			return false
		}
	}
	return true
}

// normalizeSlug 标准化slug格式
func normalizeSlug(s string) string {
	// 转为小写
	s = strings.ToLower(s)

	// 移除引号和其他特殊字符
	s = strings.Trim(s, "\"'`")

	// 替换空格为连字符
	s = strings.ReplaceAll(s, " ", "-")

	// 移除非法字符，只保留字母、数字和连字符
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	s = reg.ReplaceAllString(s, "")

	// 移除多个连续的连字符
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// 移除开头和结尾的连字符
	s = strings.Trim(s, "-")

	return s
}

// fallbackSlug 当翻译失败时的备用方案
func fallbackSlug(tag string) string {
	// 预定义的映射表作为备用
	fallbackTranslations := map[string]string{
		"人工智能":       "artificial-intelligence",
		"机器学习":       "machine-learning",
		"深度学习":       "deep-learning",
		"前端开发":       "frontend-development",
		"后端开发":       "backend-development",
		"JavaScript": "javascript",
		"Python":     "python",
		"Go":         "golang",
		"技术":         "technology",
		"教程":         "tutorial",
		"编程":         "programming",
		"开发":         "development",
	}

	if slug, exists := fallbackTranslations[tag]; exists {
		return slug
	}

	// 最后的备用方案：简单处理
	return normalizeSlug(tag)
}

// FallbackSlug 导出的备用slug生成函数
func FallbackSlug(tag string) string {
	return fallbackSlug(tag)
}

package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultOpenAIBaseURL = "https://api.openai.com/v1"

// OpenAIProvider implements Provider for OpenAI-compatible APIs.
type OpenAIProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	baseURL     string
	client      *http.Client
}

// NewOpenAIProvider creates a new OpenAI-compatible provider.
func NewOpenAIProvider(apiKey, model string, maxTokens int, temperature float64) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		baseURL:     defaultOpenAIBaseURL,
		client:      &http.Client{Timeout: 120 * time.Second},
	}
}

// NewOpenAICompatibleProvider creates a provider pointing at a custom base URL.
func NewOpenAICompatibleProvider(apiKey, model string, maxTokens int, temperature float64, baseURL string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = defaultOpenAIBaseURL
	}
	return &OpenAIProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		baseURL:     strings.TrimRight(baseURL, "/"),
		client:      &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *OpenAIProvider) Name() string { return "openai" }

type openaiRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
	Messages    []openaiMessage `json:"messages"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (p *OpenAIProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	model := req.Model
	if model == "" {
		model = p.model
	}
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = p.maxTokens
	}
	temp := req.Temperature
	if temp == 0 {
		temp = p.temperature
	}

	var msgs []openaiMessage
	for _, m := range req.Messages {
		msgs = append(msgs, openaiMessage{Role: m.Role, Content: m.Content})
	}

	body := openaiRequest{
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: temp,
		Messages:    msgs,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := p.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result openaiResponse
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("api error [%s]: %s", result.Error.Type, result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("empty response from api")
	}

	return &Response{
		Content:      result.Choices[0].Message.Content,
		InputTokens:  result.Usage.PromptTokens,
		OutputTokens: result.Usage.CompletionTokens,
		Model:        result.Model,
	}, nil
}

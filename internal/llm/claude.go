package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const claudeAPIURL = "https://api.anthropic.com/v1/messages"

// ClaudeProvider implements Provider for Anthropic Claude.
type ClaudeProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	client      *http.Client
}

// NewClaudeProvider creates a new Claude provider.
func NewClaudeProvider(apiKey, model string, maxTokens int, temperature float64) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		client:      &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *ClaudeProvider) Name() string { return "claude" }

type claudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
	System      string          `json:"system,omitempty"`
	Messages    []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *ClaudeProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
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

	// Separate system prompt from conversation messages
	var systemPrompt string
	var msgs []claudeMessage
	for _, m := range req.Messages {
		if m.Role == "system" {
			systemPrompt = m.Content
		} else {
			msgs = append(msgs, claudeMessage{Role: m.Role, Content: m.Content})
		}
	}

	body := claudeRequest{
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: temp,
		System:      systemPrompt,
		Messages:    msgs,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, claudeAPIURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result claudeResponse
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("claude api error [%s]: %s", result.Error.Type, result.Error.Message)
	}

	if len(result.Content) == 0 {
		return nil, fmt.Errorf("empty response from claude")
	}

	return &Response{
		Content:      result.Content[0].Text,
		InputTokens:  result.Usage.InputTokens,
		OutputTokens: result.Usage.OutputTokens,
		Model:        result.Model,
	}, nil
}

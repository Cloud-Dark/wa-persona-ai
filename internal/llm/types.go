package llm

import "context"

// Message represents a single message in a conversation.
type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// Request is the input to the LLM provider.
type Request struct {
	Messages    []Message
	MaxTokens   int
	Temperature float64
	Model       string
}

// Response is the output from the LLM provider.
type Response struct {
	Content      string
	InputTokens  int
	OutputTokens int
	Model        string
}

// Provider is the interface all LLM providers must implement.
type Provider interface {
	Generate(ctx context.Context, req *Request) (*Response, error)
	Name() string
}

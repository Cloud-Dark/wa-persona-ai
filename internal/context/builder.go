package context

import (
	"fmt"
	"strings"

	"github.com/Cloud-Dark/wa-persona-ai/internal/llm"
	"github.com/Cloud-Dark/wa-persona-ai/internal/memory"
	"github.com/Cloud-Dark/wa-persona-ai/internal/persona"
)

// Builder assembles the LLM request from persona, history, and memories.
type Builder struct {
	maxTokens       int
	includeMemories bool
	memoryMaxTokens int
}

// NewBuilder creates a new context builder.
func NewBuilder(maxTokens int, includeMemories bool, memoryMaxTokens int) *Builder {
	return &Builder{
		maxTokens:       maxTokens,
		includeMemories: includeMemories,
		memoryMaxTokens: memoryMaxTokens,
	}
}

// Build creates an LLM request from the current state.
func (b *Builder) Build(
	p *persona.Persona,
	history []*memory.ConversationMessage,
	memories []*memory.SearchResult,
	currentMsg string,
) *llm.Request {
	var messages []llm.Message

	// 1. System prompt (always included)
	systemContent := buildSystemPrompt(p, memories, b.includeMemories)
	messages = append(messages, llm.Message{
		Role:    "system",
		Content: systemContent,
	})

	// 2. Conversation history
	for _, h := range history {
		messages = append(messages, llm.Message{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	// 3. Current user message
	messages = append(messages, llm.Message{
		Role:    "user",
		Content: currentMsg,
	})

	maxTokens := b.maxTokens
	if p.LLMOverrides.MaxTokens > 0 {
		maxTokens = p.LLMOverrides.MaxTokens
	}
	temp := 0.7
	if p.LLMOverrides.Temperature > 0 {
		temp = p.LLMOverrides.Temperature
	}

	return &llm.Request{
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temp,
		Model:       p.LLMOverrides.Model,
	}
}

func buildSystemPrompt(p *persona.Persona, memories []*memory.SearchResult, includeMemories bool) string {
	var sb strings.Builder
	sb.WriteString(p.SystemPrompt)

	if includeMemories && len(memories) > 0 {
		sb.WriteString("\n\n---\nRELEVANT MEMORIES ABOUT THIS USER:\n")
		for _, m := range memories {
			sb.WriteString(fmt.Sprintf("- %s\n", m.Content))
		}
		sb.WriteString("---\nUse these memories naturally in your response when relevant.")
	}

	return sb.String()
}

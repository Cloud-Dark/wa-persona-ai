package whatsapp

import (
	"context"
	"strings"
	"time"

	"github.com/Cloud-Dark/wa-persona-ai/internal/admin"
	ctxbuilder "github.com/Cloud-Dark/wa-persona-ai/internal/context"
	"github.com/Cloud-Dark/wa-persona-ai/internal/llm"
	"github.com/Cloud-Dark/wa-persona-ai/internal/memory"
	"github.com/Cloud-Dark/wa-persona-ai/internal/persona"
	"github.com/Cloud-Dark/wa-persona-ai/internal/ratelimit"
	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// Handler orchestrates the full message processing pipeline.
type Handler struct {
	wa           *whatsmeow.Client
	personaMgr   *persona.Manager
	memoryStore  *memory.Store
	llmProvider  llm.Provider
	ctxBuilder   *ctxbuilder.Builder
	adminHandler *admin.Handler
	rateLimiter  *ratelimit.Limiter
	typingDelayMs int
	historyLimit  int
}

// HandlerConfig holds all dependencies for the Handler.
type HandlerConfig struct {
	WA            *whatsmeow.Client
	PersonaMgr    *persona.Manager
	MemoryStore   *memory.Store
	LLMProvider   llm.Provider
	CtxBuilder    *ctxbuilder.Builder
	AdminHandler  *admin.Handler
	RateLimiter   *ratelimit.Limiter
	TypingDelayMs int
	HistoryLimit  int
}

// NewHandler creates a new message handler.
func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		wa:            cfg.WA,
		personaMgr:    cfg.PersonaMgr,
		memoryStore:   cfg.MemoryStore,
		llmProvider:   cfg.LLMProvider,
		ctxBuilder:    cfg.CtxBuilder,
		adminHandler:  cfg.AdminHandler,
		rateLimiter:   cfg.RateLimiter,
		typingDelayMs: cfg.TypingDelayMs,
		historyLimit:  cfg.HistoryLimit,
	}
}

// Handle processes a WhatsApp event.
func (h *Handler) Handle(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.Message:
		h.handleMessage(evt)
	case *events.Connected:
		log.Info().Msg("WhatsApp connected")
	case *events.Disconnected:
		log.Warn().Msg("WhatsApp disconnected")
	case *events.LoggedOut:
		log.Error().Msg("WhatsApp logged out — delete session and restart")
	}
}

func (h *Handler) handleMessage(evt *events.Message) {
	// Skip messages from self
	if evt.Info.IsFromMe {
		return
	}

	// Only handle text messages
	msg := evt.Message
	var text string
	if msg.GetConversation() != "" {
		text = msg.GetConversation()
	} else if msg.GetExtendedTextMessage() != nil {
		text = msg.GetExtendedTextMessage().GetText()
	}

	if strings.TrimSpace(text) == "" {
		return
	}

	senderJID := evt.Info.Sender.String()
	chatJID := evt.Info.Chat

	log.Debug().
		Str("sender", senderJID).
		Str("chat", chatJID.String()).
		Msg("message received")

	// Rate limiting
	if h.rateLimiter != nil && !h.rateLimiter.Allow(senderJID) {
		log.Debug().Str("sender", senderJID).Msg("rate limited")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Admin command handling
	if h.adminHandler != nil && h.adminHandler.IsCommand(text) {
		if h.adminHandler.IsAdmin(senderJID) {
			response := h.adminHandler.Execute(ctx, senderJID, text)
			if response != "" {
				h.sendReply(chatJID, response)
			}
		} else {
			h.sendReply(chatJID, "Maaf, command ini hanya untuk admin.")
		}
		return
	}

	// Run full AI pipeline
	go h.processPipeline(ctx, senderJID, chatJID, text, evt)
}

func (h *Handler) processPipeline(ctx context.Context, senderJID string, chatJID types.JID, text string, evt *events.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("panic", r).Msg("recovered from panic in message pipeline")
			h.sendReply(chatJID, "Maaf, terjadi kesalahan. Silakan coba lagi.")
		}
	}()

	// 1. Get active persona
	p := h.personaMgr.GetActive(senderJID)

	// 2. Get or create conversation
	convID, err := h.memoryStore.GetOrCreateConversation(ctx, senderJID, p.Name)
	if err != nil {
		log.Error().Err(err).Msg("failed to get conversation")
		convID = ""
	}

	// 3. Get conversation history
	var history []*memory.ConversationMessage
	if convID != "" {
		history, _ = h.memoryStore.GetHistory(ctx, convID, h.historyLimit)
	}

	// 4. Retrieve relevant memories
	var memories []*memory.SearchResult
	if p.Memory.Enabled {
		memories, _ = h.memoryStore.SearchMemories(ctx, text, senderJID)
	}

	// 5. Build context
	req := h.ctxBuilder.Build(p, history, memories, text)

	// 6. Send typing indicator
	SendTyping(h.wa, chatJID, h.typingDelayMs)

	// 7. Generate AI response
	resp, err := h.llmProvider.Generate(ctx, req)
	StopTyping(h.wa, chatJID)

	if err != nil {
		log.Error().Err(err).Str("sender", senderJID).Msg("LLM generation failed")
		h.sendReply(chatJID, "Maaf, saya tidak bisa memproses pesan sekarang. Coba lagi ya~")
		return
	}

	responseText := resp.Content

	// 8. Send reply
	h.sendReply(chatJID, responseText)

	// 9. Save messages to history
	if convID != "" {
		_ = h.memoryStore.SaveMessage(ctx, convID, "user", text)
		_ = h.memoryStore.SaveMessage(ctx, convID, "assistant", responseText)
	}

	// 10. Extract and store memories asynchronously
	if p.Memory.Enabled {
		go func() {
			bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			h.memoryStore.ExtractAndSaveMemories(bgCtx, senderJID, text, responseText, p.Memory.ImportanceThreshold)
		}()
	}
}

func (h *Handler) sendReply(jid types.JID, text string) {
	if err := SendText(h.wa, jid, text); err != nil {
		log.Error().Err(err).Str("jid", jid.String()).Msg("failed to send message")
	}
}

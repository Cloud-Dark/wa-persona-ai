package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cloud-Dark/wa-persona-ai/internal/admin"
	ctxbuilder "github.com/Cloud-Dark/wa-persona-ai/internal/context"
	"github.com/Cloud-Dark/wa-persona-ai/internal/config"
	"github.com/Cloud-Dark/wa-persona-ai/internal/llm"
	"github.com/Cloud-Dark/wa-persona-ai/internal/memory"
	"github.com/Cloud-Dark/wa-persona-ai/internal/persona"
	"github.com/Cloud-Dark/wa-persona-ai/internal/ratelimit"
	"github.com/Cloud-Dark/wa-persona-ai/internal/whatsapp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	configPath := os.Getenv("WPA_CONFIG_PATH")
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup logging
	setupLogging(cfg.Logging)

	log.Info().Msg("Starting WA Persona AI...")

	// Initialize persona manager
	personaMgr, err := persona.NewManager(cfg.Persona.Dir, cfg.Persona.Default)
	if err != nil {
		log.Warn().Err(err).Msg("Persona manager initialization warning")
	}
	log.Info().Strs("personas", personaMgr.List()).Msg("Personas loaded")

	// Initialize memory store
	var memStore *memory.Store
	if cfg.Memory.Enabled {
		openAIKey := cfg.LLM.OpenAI.APIKey
		if cfg.Memory.Embedding.Provider != "openai" {
			openAIKey = ""
		}
		memStore, err = memory.NewStore(
			cfg.Memory.VectorDir,
			cfg.Memory.MetadataDB,
			cfg.Memory.Retrieval.TopK,
			cfg.Memory.Retrieval.MinSimilarity,
			openAIKey,
		)
		if err != nil {
			log.Warn().Err(err).Msg("Memory store initialization failed — running without memory")
		} else {
			log.Info().Msg("Memory store initialized")
			defer memStore.Close()
		}
	}

	if memStore == nil {
		log.Warn().Msg("Memory disabled — conversations won't persist across restarts")
		// Create a no-op memory store for graceful degradation
		memStore, _ = memory.NewStore("./data/memory/vectors", "./data/memory/metadata.db", 5, 0.7, "")
		if memStore != nil {
			defer memStore.Close()
		}
	}

	// Initialize LLM provider
	llmProvider := buildLLMProvider(cfg)

	// Initialize context builder
	ctxBldr := ctxbuilder.NewBuilder(
		cfg.Context.MaxTokens,
		cfg.Context.IncludeMemories,
		cfg.Context.MemoryMaxTokens,
	)

	// Initialize WhatsApp client
	waClient, err := whatsapp.NewClient(cfg.WhatsApp.SessionDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create WhatsApp client")
	}

	// Setup shutdown signal
	shutdownCh := make(chan struct{}, 1)
	shutdownFn := func() {
		select {
		case shutdownCh <- struct{}{}:
		default:
		}
	}

	// Initialize admin handler
	var adminHandler *admin.Handler
	if len(cfg.Admin.Numbers) > 0 {
		adminHandler = admin.NewHandler(cfg.Admin.Numbers, personaMgr, memStore, shutdownFn)
		log.Info().Strs("admins", cfg.Admin.Numbers).Msg("Admin handler initialized")
	}

	// Initialize rate limiter
	var rateLimiter *ratelimit.Limiter
	if cfg.RateLimit.Enabled {
		rateLimiter = ratelimit.NewLimiter(
			cfg.RateLimit.PerUser.MessagesPerMinute,
			cfg.RateLimit.PerUser.MessagesPerHour,
		)
	}

	// Initialize message handler
	msgHandler := whatsapp.NewHandler(whatsapp.HandlerConfig{
		WA:            waClient.WA,
		PersonaMgr:    personaMgr,
		MemoryStore:   memStore,
		LLMProvider:   llmProvider,
		CtxBuilder:    ctxBldr,
		AdminHandler:  adminHandler,
		RateLimiter:   rateLimiter,
		TypingDelayMs: cfg.WhatsApp.TypingDelayMs,
		HistoryLimit:  cfg.Context.MaxHistory,
	})

	// Register event handler
	waClient.AddEventHandler(msgHandler.Handle)

	// Connect to WhatsApp
	log.Info().Msg("Connecting to WhatsApp...")
	if err := waClient.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to WhatsApp")
	}
	log.Info().Msg("WhatsApp connection established")

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case <-shutdownCh:
		log.Info().Msg("Shutdown requested via admin command")
	}

	// Graceful shutdown
	log.Info().Msg("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = ctx

	waClient.Disconnect()
	log.Info().Msg("Shutdown complete")
}

func buildLLMProvider(cfg *config.Config) llm.Provider {
	var primary, fallback llm.Provider

	newOpenAI := func(c config.OpenAIConfig) llm.Provider {
		return llm.NewOpenAICompatibleProvider(
			c.APIKey, c.Model, c.MaxTokens, c.Temperature, c.BaseURL,
		)
	}

	switch cfg.LLM.Provider {
	case "openai":
		primary = newOpenAI(cfg.LLM.OpenAI)
		if cfg.LLM.FallbackProvider == "claude" && cfg.LLM.Claude.APIKey != "" {
			fallback = llm.NewClaudeProvider(
				cfg.LLM.Claude.APIKey,
				cfg.LLM.Claude.Model,
				cfg.LLM.Claude.MaxTokens,
				cfg.LLM.Claude.Temperature,
			)
		}
	default: // "claude"
		primary = llm.NewClaudeProvider(
			cfg.LLM.Claude.APIKey,
			cfg.LLM.Claude.Model,
			cfg.LLM.Claude.MaxTokens,
			cfg.LLM.Claude.Temperature,
		)
		if cfg.LLM.FallbackProvider == "openai" && cfg.LLM.OpenAI.APIKey != "" {
			fallback = newOpenAI(cfg.LLM.OpenAI)
		}
	}

	return llm.NewRetryProvider(
		primary,
		fallback,
		cfg.LLM.Retry.MaxAttempts,
		cfg.LLM.Retry.InitialDelayMs,
		cfg.LLM.Retry.MaxDelayMs,
		cfg.LLM.Retry.Multiplier,
	)
}

func setupLogging(cfg config.LoggingConfig) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.Format == "text" || cfg.Output == "stdout" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"})
	}
}

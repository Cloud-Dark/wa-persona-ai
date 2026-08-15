package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WhatsApp  WhatsAppConfig  `yaml:"whatsapp"`
	LLM       LLMConfig       `yaml:"llm"`
	Persona   PersonaConfig   `yaml:"persona"`
	Memory    MemoryConfig    `yaml:"memory"`
	Context   ContextConfig   `yaml:"context"`
	Admin     AdminConfig     `yaml:"admin"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Logging   LoggingConfig   `yaml:"logging"`
}

type WhatsAppConfig struct {
	SessionDir       string `yaml:"session_dir"`
	LogLevel         string `yaml:"log_level"`
	AutoReconnect    bool   `yaml:"auto_reconnect"`
	ReconnectInterval int   `yaml:"reconnect_interval"`
	TypingDelayMs    int    `yaml:"typing_delay_ms"`
}

type LLMConfig struct {
	Provider         string        `yaml:"provider"`
	FallbackProvider string        `yaml:"fallback_provider"`
	Claude           ClaudeConfig  `yaml:"claude"`
	OpenAI           OpenAIConfig  `yaml:"openai"`
	Retry            RetryConfig   `yaml:"retry"`
}

type ClaudeConfig struct {
	APIKey    string `yaml:"api_key"`
	Model     string `yaml:"model"`
	MaxTokens int    `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

type OpenAIConfig struct {
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
	BaseURL     string  `yaml:"base_url"`
}

type RetryConfig struct {
	MaxAttempts   int     `yaml:"max_attempts"`
	InitialDelayMs int   `yaml:"initial_delay_ms"`
	MaxDelayMs    int    `yaml:"max_delay_ms"`
	Multiplier    float64 `yaml:"multiplier"`
}

type PersonaConfig struct {
	Dir               string `yaml:"dir"`
	Default           string `yaml:"default"`
	HotReload         bool   `yaml:"hot_reload"`
	HotReloadInterval int    `yaml:"hot_reload_interval"`
}

type MemoryConfig struct {
	Enabled    bool            `yaml:"enabled"`
	VectorDir  string          `yaml:"vector_dir"`
	MetadataDB string          `yaml:"metadata_db"`
	Embedding  EmbeddingConfig `yaml:"embedding"`
	Retrieval  RetrievalConfig `yaml:"retrieval"`
	Cleanup    CleanupConfig   `yaml:"cleanup"`
}

type EmbeddingConfig struct {
	Provider   string `yaml:"provider"`
	Model      string `yaml:"model"`
	Dimensions int    `yaml:"dimensions"`
}

type RetrievalConfig struct {
	TopK           int     `yaml:"top_k"`
	MinSimilarity  float64 `yaml:"min_similarity"`
	RecencyWeight  float64 `yaml:"recency_weight"`
}

type CleanupConfig struct {
	Enabled             bool `yaml:"enabled"`
	RetentionDays       int  `yaml:"retention_days"`
	MaxEntriesPerUser   int  `yaml:"max_entries_per_user"`
	RunIntervalHours    int  `yaml:"run_interval_hours"`
}

type ContextConfig struct {
	MaxHistory        int  `yaml:"max_history"`
	MaxTokens         int  `yaml:"max_tokens"`
	IncludeMemories   bool `yaml:"include_memories"`
	MemoryMaxTokens   int  `yaml:"memory_section_max_tokens"`
}

type AdminConfig struct {
	Numbers       []string `yaml:"numbers"`
	CommandPrefix string   `yaml:"command_prefix"`
}

type RateLimitConfig struct {
	Enabled  bool           `yaml:"enabled"`
	PerUser  UserRateLimit  `yaml:"per_user"`
	Global   UserRateLimit  `yaml:"global"`
}

type UserRateLimit struct {
	MessagesPerMinute int `yaml:"messages_per_minute"`
	MessagesPerHour   int `yaml:"messages_per_hour"`
}

type LoggingConfig struct {
	Level  string    `yaml:"level"`
	Format string    `yaml:"format"`
	Output string    `yaml:"output"`
	File   FileLogConfig `yaml:"file"`
}

type FileLogConfig struct {
	Path        string `yaml:"path"`
	MaxSizeMB   int    `yaml:"max_size_mb"`
	MaxBackups  int    `yaml:"max_backups"`
	MaxAgeDays  int    `yaml:"max_age_days"`
	Compress    bool   `yaml:"compress"`
}

func Load(path string) (*Config, error) {
	cfg := defaultConfig()

	// Load .env file if present (simple key=value parser)
	loadDotEnv(".env")

	if path == "" {
		path = "config/config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config file: %w", err)
		}
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

// loadDotEnv reads a .env file and sets missing environment variables.
func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		// Only set if not already set
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("CLAUDE_API_KEY"); v != "" {
		cfg.LLM.Claude.APIKey = v
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.LLM.OpenAI.APIKey = v
	}
	if v := os.Getenv("WPA_LOG_LEVEL"); v != "" {
		cfg.Logging.Level = v
	}
	if v := os.Getenv("WPA_ADMIN_NUMBERS"); v != "" {
		cfg.Admin.Numbers = strings.Split(v, ",")
	}
}

func defaultConfig() *Config {
	return &Config{
		WhatsApp: WhatsAppConfig{
			SessionDir:        "./data/session",
			LogLevel:          "INFO",
			AutoReconnect:     true,
			ReconnectInterval: 30,
			TypingDelayMs:     1500,
		},
		LLM: LLMConfig{
			Provider: "claude",
			Claude: ClaudeConfig{
				Model:       "claude-haiku-4-5-20251001",
				MaxTokens:   1024,
				Temperature: 0.7,
			},
			OpenAI: OpenAIConfig{
				Model:       "gpt-4o-mini",
				MaxTokens:   1024,
				Temperature: 0.7,
			},
			Retry: RetryConfig{
				MaxAttempts:    3,
				InitialDelayMs: 1000,
				MaxDelayMs:     10000,
				Multiplier:     2.0,
			},
		},
		Persona: PersonaConfig{
			Dir:               "./persona",
			Default:           "default",
			HotReload:         true,
			HotReloadInterval: 30,
		},
		Memory: MemoryConfig{
			Enabled:    true,
			VectorDir:  "./data/memory/vectors",
			MetadataDB: "./data/memory/metadata.db",
			Embedding: EmbeddingConfig{
				Provider:   "openai",
				Model:      "text-embedding-ada-002",
				Dimensions: 1536,
			},
			Retrieval: RetrievalConfig{
				TopK:          5,
				MinSimilarity: 0.7,
				RecencyWeight: 0.2,
			},
			Cleanup: CleanupConfig{
				Enabled:           true,
				RetentionDays:     90,
				MaxEntriesPerUser: 5000,
				RunIntervalHours:  24,
			},
		},
		Context: ContextConfig{
			MaxHistory:      20,
			MaxTokens:       4096,
			IncludeMemories: true,
			MemoryMaxTokens: 1024,
		},
		Admin: AdminConfig{
			CommandPrefix: "!",
		},
		RateLimit: RateLimitConfig{
			Enabled: true,
			PerUser: UserRateLimit{
				MessagesPerMinute: 10,
				MessagesPerHour:   100,
			},
			Global: UserRateLimit{
				MessagesPerMinute: 50,
				MessagesPerHour:   500,
			},
		},
		Logging: LoggingConfig{
			Level:  "INFO",
			Format: "text",
			Output: "stdout",
			File: FileLogConfig{
				Path:       "./logs/app.log",
				MaxSizeMB:  100,
				MaxBackups: 5,
				MaxAgeDays: 30,
				Compress:   true,
			},
		},
	}
}

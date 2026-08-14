package persona

// Persona represents a complete AI personality configuration.
type Persona struct {
	Name        string      `yaml:"name"`
	DisplayName string      `yaml:"display_name"`
	Version     string      `yaml:"version"`
	Description string      `yaml:"description"`
	Author      string      `yaml:"author"`

	SystemPrompt string   `yaml:"system_prompt"`
	Traits       []string `yaml:"traits"`

	Language    LanguageConfig    `yaml:"language"`
	Constraints ConstraintsConfig `yaml:"constraints"`
	Greeting    GreetingConfig    `yaml:"greeting"`
	Memory      MemoryConfig      `yaml:"memory"`
	LLMOverrides LLMOverrides     `yaml:"llm_overrides"`
}

type LanguageConfig struct {
	Primary             string `yaml:"primary"`
	Style               string `yaml:"style"`
	Tone                string `yaml:"tone"`
	UseEmoji            bool   `yaml:"use_emoji"`
	MaxEmojiPerMessage  int    `yaml:"max_emoji_per_message"`
}

type ConstraintsConfig struct {
	AllowedTopics     []string `yaml:"allowed_topics"`
	BlockedTopics     []string `yaml:"blocked_topics"`
	MaxResponseLength int      `yaml:"max_response_length"`
	ResponseStyle     string   `yaml:"response_style"`
	SafetyLevel       string   `yaml:"safety_level"`
}

type GreetingConfig struct {
	NewUser       string `yaml:"new_user"`
	ReturningUser string `yaml:"returning_user"`
}

type MemoryConfig struct {
	Enabled             bool    `yaml:"enabled"`
	RecallStyle         string  `yaml:"recall_style"`
	ImportanceThreshold float64 `yaml:"importance_threshold"`
}

type LLMOverrides struct {
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	Model       string  `yaml:"model"`
}

// DefaultPersona returns a minimal safe persona when loading fails.
func DefaultPersona() *Persona {
	return &Persona{
		Name:        "default",
		DisplayName: "AI Assistant",
		Version:     "1.0.0",
		SystemPrompt: "You are a helpful AI assistant on WhatsApp. Be concise, friendly, and respond in the user's language.",
		Traits:      []string{"helpful", "friendly"},
		Language: LanguageConfig{
			Primary:  "auto",
			Style:    "casual",
			Tone:     "warm",
			UseEmoji: true,
		},
		Constraints: ConstraintsConfig{
			MaxResponseLength: 500,
			ResponseStyle:     "balanced",
			SafetyLevel:       "moderate",
		},
		Greeting: GreetingConfig{
			NewUser:       "Halo! Ada yang bisa saya bantu?",
			ReturningUser: "Halo lagi! Ada yang bisa saya bantu?",
		},
		Memory: MemoryConfig{
			Enabled:             true,
			RecallStyle:         "natural",
			ImportanceThreshold: 0.5,
		},
		LLMOverrides: LLMOverrides{
			Temperature: 0.7,
			MaxTokens:   1024,
		},
	}
}

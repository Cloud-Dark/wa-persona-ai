# Configuration Reference — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Configuration File

Konfigurasi utama menggunakan file YAML di `config/config.yaml`. Semua value bisa di-override menggunakan environment variable.

### Contoh Konfigurasi Lengkap

```yaml
# =============================================================================
# WA Persona AI Configuration
# =============================================================================

# WhatsApp Settings
whatsapp:
  session_dir: "./data/session"       # Directory untuk menyimpan session WhatsApp
  log_level: "INFO"                   # Log level: DEBUG, INFO, WARN, ERROR
  auto_reconnect: true                # Auto-reconnect saat disconnect
  reconnect_interval: 30              # Interval reconnect dalam detik
  typing_delay_ms: 1500               # Delay typing indicator (ms)
  read_receipt: true                  # Kirim read receipt

# LLM Provider Settings
llm:
  provider: "claude"                  # Provider: "claude" | "openai"
  fallback_provider: "openai"         # Fallback jika primary down

  claude:
    api_key: "${CLAUDE_API_KEY}"      # Via env variable
    model: "claude-sonnet-4-20250514"
    max_tokens: 1024
    temperature: 0.7

  openai:
    api_key: "${OPENAI_API_KEY}"      # Via env variable
    model: "gpt-4o"
    max_tokens: 1024
    temperature: 0.7

  retry:
    max_attempts: 3
    initial_delay_ms: 1000
    max_delay_ms: 10000
    multiplier: 2.0

# Persona Settings
persona:
  dir: "./persona"                    # Directory persona YAML files
  default: "default"                  # Nama persona default
  hot_reload: true                    # Auto-reload saat file berubah
  hot_reload_interval: 30             # Interval check perubahan (detik)

# Memory Settings
memory:
  enabled: true
  vector_dir: "./data/memory/vectors" # Directory vector DB storage
  metadata_db: "./data/memory/metadata.db"  # SQLite metadata

  embedding:
    provider: "openai"                # "openai" | "local"
    model: "text-embedding-ada-002"
    dimensions: 1536

  retrieval:
    top_k: 5                          # Jumlah memory yang di-retrieve
    min_similarity: 0.7               # Minimum cosine similarity threshold
    recency_weight: 0.2               # Bobot recency (0-1)

  cleanup:
    enabled: true
    retention_days: 90                # Hapus memory lebih lama dari N hari
    max_entries_per_user: 5000        # Maksimum memory per user
    run_interval_hours: 24            # Interval cleanup job

# Context Settings
context:
  max_history: 20                     # Jumlah pesan terakhir dalam context
  max_tokens: 4096                    # Maksimum token untuk context
  include_memories: true              # Sertakan memories dalam context
  memory_section_max_tokens: 1024     # Alokasi token untuk memories

# Admin Settings
admin:
  numbers:                            # Daftar nomor admin (format internasional)
    - "6281234567890"
  command_prefix: "!"                 # Prefix command admin

# Rate Limiting
rate_limit:
  enabled: true
  per_user:
    messages_per_minute: 10
    messages_per_hour: 100
  global:
    messages_per_minute: 50
    messages_per_hour: 500

# Logging
logging:
  level: "INFO"                       # DEBUG, INFO, WARN, ERROR
  format: "json"                      # "json" | "text"
  output: "both"                      # "stdout" | "file" | "both"
  file:
    path: "./logs/app.log"
    max_size_mb: 100
    max_backups: 5
    max_age_days: 30
    compress: true
```

## 2. Environment Variables

| Variable | Deskripsi | Required | Default |
|---|---|---|---|
| `CLAUDE_API_KEY` | API key Anthropic Claude | Yes (jika provider=claude) | - |
| `OPENAI_API_KEY` | API key OpenAI | Yes (jika provider=openai) | - |
| `WPA_CONFIG_PATH` | Path ke config file | No | `./config/config.yaml` |
| `WPA_LOG_LEVEL` | Override log level | No | `INFO` |
| `WPA_ADMIN_NUMBERS` | Override admin numbers (comma-separated) | No | Dari config |

## 3. Referensi

- Lihat [17_DEVELOPER_SETUP.md](17_DEVELOPER_SETUP.md) untuk setup guide
- Lihat [18_INSTALLATION.md](18_INSTALLATION.md) untuk installation guide

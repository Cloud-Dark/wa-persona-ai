# Developer Setup Guide — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Prerequisites

| Requirement | Version | Install |
|---|---|---|
| Go | 1.22+ | [golang.org/dl](https://golang.org/dl/) |
| Git | 2.30+ | [git-scm.com](https://git-scm.com/) |
| Make | (any) | OS package manager |
| SQLite3 | 3.35+ | Embedded via go-sqlite3 (auto) |

### Optional

| Tool | Purpose |
|---|---|
| Docker | Container deployment |
| GoLand / VS Code | IDE with Go support |
| golangci-lint | Linting |
| delve | Debugging |

## 2. Quick Start

```bash
# 1. Clone repository
git clone https://github.com/Cloud-Dark/wa-persona-ai.git
cd wa-persona-ai

# 2. Copy config
cp config/config.example.yaml config/config.yaml
cp .env.example .env

# 3. Set API keys
# Edit .env dan isi CLAUDE_API_KEY atau OPENAI_API_KEY

# 4. Install dependencies
go mod download

# 5. Build
make build
# Atau: go build -o bin/wa-persona-ai ./cmd/wa-persona-ai/

# 6. Run
make run
# Atau: ./bin/wa-persona-ai

# 7. Scan QR code yang muncul di terminal
```

## 3. Project Structure

```
wa-persona-ai/
├── cmd/wa-persona-ai/    # Entry point
├── internal/              # Private packages
│   ├── whatsapp/          # WhatsApp client
│   ├── persona/           # Persona management
│   ├── llm/               # LLM providers
│   ├── memory/            # Vector memory
│   ├── context/           # Context building
│   ├── admin/             # Admin commands
│   └── config/            # Configuration
├── persona/               # Persona YAML files
├── memory/                # Memory data (gitignored)
├── docs/                  # Documentation
├── config/                # Config files
├── scripts/               # Utility scripts
└── tests/                 # Integration tests
```

## 4. Development Workflow

### Build & Run

```bash
make build          # Build binary
make run            # Build + run
make dev            # Run with hot-reload (air)
make test           # Run all tests
make test-coverage  # Run tests with coverage
make lint           # Run linter
make clean          # Clean build artifacts
```

### Testing

```bash
# Unit tests
go test ./internal/...

# Specific package
go test ./internal/persona/...

# With coverage
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out

# Race condition detection
go test -race ./internal/...

# Integration tests (requires running services)
go test ./tests/integration/... -tags=integration
```

### Adding a New LLM Provider

1. Implement `LLMProvider` interface di `internal/llm/`
2. Register provider di `internal/llm/registry.go`
3. Add config section di `config.yaml`
4. Write tests di `internal/llm/*_test.go`

### Adding a New Persona

1. Copy `persona/default.yaml` sebagai template
2. Edit sesuai kebutuhan
3. Simpan di `persona/<nama>.yaml`
4. Hot-reload otomatis (atau kirim `!reload` via WhatsApp)

## 5. Environment Variables

```bash
# Required
export CLAUDE_API_KEY="sk-ant-..."
# Atau
export OPENAI_API_KEY="sk-..."

# Optional
export WPA_LOG_LEVEL="DEBUG"
export WPA_CONFIG_PATH="./config/config.yaml"
```

## 6. Docker Development

```bash
# Build image
docker build -t wa-persona-ai .

# Run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## 7. Debugging

### VS Code Launch Config (`.vscode/launch.json`)

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch WA Persona AI",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/wa-persona-ai",
      "envFile": "${workspaceFolder}/.env"
    }
  ]
}
```

### Common Issues

| Issue | Solusi |
|---|---|
| QR code tidak muncul | Pastikan terminal support UTF-8, coba terminal lain |
| Session expired | Hapus folder `data/session/` dan scan ulang |
| SQLite build error | Install GCC (CGO requirement untuk go-sqlite3) |
| Memory not working | Pastikan `memory.enabled: true` dan API key embedding tersedia |

## 8. Referensi

- Lihat [18_INSTALLATION.md](18_INSTALLATION.md) untuk production installation
- Lihat [16_CONFIG_REFERENCE.md](16_CONFIG_REFERENCE.md) untuk konfigurasi lengkap

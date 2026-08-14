# Code Standard — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Go Style Guide

Mengikuti [Effective Go](https://go.dev/doc/effective_go) dan [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

## 2. Formatting

- **Formatter:** `gofmt` (wajib, tanpa exception)
- **Imports:** Dikelompokkan: stdlib → third-party → internal
- **Line length:** Soft limit 100 karakter, hard limit 120
- **Indentation:** Tab (Go default)

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // Third-party
    "github.com/rs/zerolog"
    "go.mau.fi/whatsmeow"

    // Internal
    "github.com/Cloud-Dark/wa-persona-ai/internal/config"
    "github.com/Cloud-Dark/wa-persona-ai/internal/persona"
)
```

## 3. Naming Conventions

| Elemen | Convention | Contoh |
|---|---|---|
| Package | lowercase, single word | `persona`, `memory`, `llm` |
| Interface | Noun atau -er suffix | `Provider`, `Store`, `Builder` |
| Struct | PascalCase, descriptive | `PersonaManager`, `MemoryEntry` |
| Function | PascalCase (exported), camelCase (unexported) | `GetActive()`, `buildContext()` |
| Variable | camelCase | `userJID`, `topK`, `maxTokens` |
| Constant | PascalCase or ALL_CAPS | `MaxRetries`, `DefaultTopK` |
| Error | `Err` prefix | `ErrNotFound`, `ErrRateLimit` |
| File | snake_case | `persona_manager.go`, `vector_store.go` |
| Test file | `*_test.go` | `manager_test.go` |

## 4. Error Handling

```go
// DO: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to load persona %s: %w", name, err)
}

// DO: Use custom error types for known errors
var ErrPersonaNotFound = errors.New("persona not found")

// DON'T: Ignore errors
_ = file.Close() // BAD
if err := file.Close(); err != nil {
    log.Warn().Err(err).Msg("failed to close file")
}
```

## 5. Logging

```go
// Use zerolog with structured fields
log.Info().
    Str("user_jid", msg.SenderJID).
    Str("persona", persona.Name).
    Int("memory_count", len(memories)).
    Msg("processing message")

// Log levels:
// Debug — development only, verbose
// Info  — normal operations
// Warn  — recoverable issues
// Error — failures requiring attention
```

## 6. Context Usage

```go
// Always pass context as first parameter
func (s *Store) Search(ctx context.Context, query string) ([]*Result, error) {
    // Use context for cancellation and timeout
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    // ...
}
```

## 7. Testing

```go
// Table-driven tests
func TestPersonaLoad(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *Persona
        wantErr  bool
    }{
        {name: "valid persona", input: "default.yaml", expected: &Persona{Name: "default"}, wantErr: false},
        {name: "missing file", input: "nonexistent.yaml", expected: nil, wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := LoadPersona(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("LoadPersona() error = %v, wantErr %v", err, tt.wantErr)
            }
            // ...
        })
    }
}
```

## 8. Linting

```bash
# Wajib lulus sebelum commit:
golangci-lint run ./...
go vet ./...
```

## 9. Referensi

- Lihat [03_ARCHITECTURE_STANDARD.md](03_ARCHITECTURE_STANDARD.md) untuk patterns
- Lihat [05_TESTING_STANDARD.md](05_TESTING_STANDARD.md) untuk testing detail

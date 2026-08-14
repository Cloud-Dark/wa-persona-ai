# Architecture Standard — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Prinsip Arsitektur

1. **Modular** — Setiap module punya satu responsibility jelas
2. **Interface-driven** — Komunikasi antar module melalui interface
3. **Dependency Injection** — Dependencies di-inject, bukan di-instantiate internal
4. **Fail-safe** — Setiap error di-handle gracefully, sistem tetap berjalan
5. **Configurable** — Behavior bisa diubah tanpa code change

## 2. Design Patterns

### Repository Pattern (Storage)
```go
// Interface di package domain, implementasi di package storage
type MemoryRepository interface {
    Save(ctx context.Context, entry *Memory) error
    FindByUser(ctx context.Context, userJID string) ([]*Memory, error)
}
```

### Strategy Pattern (LLM Providers)
```go
// Swap provider tanpa ubah caller
type LLMProvider interface {
    Generate(ctx context.Context, req *Request) (*Response, error)
}
```

### Observer Pattern (Event Handling)
```go
// WhatsApp events dispatched ke registered handlers
type EventHandler interface {
    OnMessage(evt *MessageEvent)
    OnConnected(evt *ConnectedEvent)
}
```

## 3. Package Rules

- `cmd/` — Entry points only, no business logic
- `internal/` — Private packages, not importable externally
- No circular dependencies
- Shared types di package `types/` atau dalam interface package
- Thin packages — prefer many small packages over few large ones

## 4. Referensi

- Lihat [05_ARCHITECTURE.md](../05_ARCHITECTURE.md) untuk implementasi
- Lihat [02_CODE_STANDARD.md](02_CODE_STANDARD.md) untuk coding style

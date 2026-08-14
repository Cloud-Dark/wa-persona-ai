# Technical Requirements Document (TRD) — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Tech Stack

| Layer | Teknologi | Alasan |
|---|---|---|
| Language | Go 1.22+ | Performance, concurrency, single binary deploy |
| WhatsApp Library | whatsmeow | Go-native, multi-device support, actively maintained |
| LLM Client | HTTP client (net/http) | Langsung ke API, tanpa wrapper dependency |
| Vector Database | ChromaDB / Milvus Lite / SQLite+VSS | Embeddable, ringan untuk self-hosted |
| Embedding | OpenAI Ada-002 / local model | Generate vector embedding untuk memory |
| Configuration | YAML + envvar | Human-readable config, secret via env |
| Logging | zerolog / slog | Structured logging, minimal allocation |
| Database (metadata) | SQLite | Zero-config, embedded, cocok untuk single instance |
| Session Storage | File-based (whatsmeow default) | Persistensi session WhatsApp |

## 2. Kebutuhan Teknis

### TR-001: Performance
| Requirement | Target |
|---|---|
| Message processing latency | < 100ms (excluding LLM call) |
| LLM response time | < 5s (depending on provider) |
| Memory retrieval time | < 200ms |
| Concurrent conversations | > 50 simultaneous |
| Memory footprint | < 512MB |

### TR-002: Reliability
| Requirement | Target |
|---|---|
| Auto-reconnect WhatsApp | Within 30 seconds |
| LLM provider failover | < 2 seconds switchover |
| Data durability | No memory loss on restart |
| Graceful shutdown | Complete in-flight requests before exit |

### TR-003: Scalability
| Requirement | Target |
|---|---|
| Users per instance | Up to 1000 |
| Memory entries per user | Up to 10,000 |
| Total memory entries | Up to 1,000,000 |
| Persona count | Unlimited |

### TR-004: Security
| Requirement | Deskripsi |
|---|---|
| API key storage | Via environment variable, never in code/config files |
| Conversation data | Stored locally, encrypted at rest (optional) |
| Admin auth | Whitelist-based (phone number) |
| Input sanitization | Prevent prompt injection via user messages |
| Rate limiting | Configurable per-user rate limit |

### TR-005: Deployment
| Requirement | Deskripsi |
|---|---|
| Binary deployment | Single binary, no runtime dependency |
| Docker support | Dockerfile + docker-compose.yml |
| Environment | Linux/macOS/Windows (cross-compile) |
| Minimum specs | 1 vCPU, 2GB RAM, 10GB disk |
| Configuration | YAML file + environment variables |

## 3. Arsitektur Integrasi

### 3.1 LLM Provider Interface

```go
type LLMProvider interface {
    GenerateResponse(ctx context.Context, req LLMRequest) (LLMResponse, error)
    CountTokens(messages []Message) (int, error)
    MaxTokens() int
    Name() string
}
```

### 3.2 Vector Store Interface

```go
type VectorStore interface {
    Store(ctx context.Context, entry MemoryEntry) error
    Search(ctx context.Context, query string, opts SearchOpts) ([]MemoryResult, error)
    Delete(ctx context.Context, filter DeleteFilter) error
    Stats(ctx context.Context) (StoreStats, error)
}
```

### 3.3 Persona Interface

```go
type Persona interface {
    Name() string
    SystemPrompt() string
    Traits() []string
    Constraints() PersonaConstraints
    GreetingMessage(userName string) string
    FormatResponse(raw string) string
}
```

## 4. Dependency Matrix

| Dependency | Version | Purpose | Fallback |
|---|---|---|---|
| whatsmeow | latest | WhatsApp connection | _Tidak ada alternatif_ |
| go-sqlite3 | v1.14+ | Metadata storage | Pure Go: modernc.org/sqlite |
| zerolog | v1.31+ | Structured logging | stdlib slog |
| viper | v1.18+ | Configuration management | Manual YAML parsing |
| chromem-go | latest | Embedded vector DB | SQLite with vector extension |
| uuid | v1.6+ | ID generation | stdlib crypto/rand |

## 5. Batasan Teknis

1. **whatsmeow** adalah unofficial API — update WhatsApp bisa memerlukan update library
2. **Vector DB lokal** memiliki limit skala dibanding hosted solution
3. **Single instance** — tidak dirancang untuk horizontal scaling (fase awal)
4. **Go generics** memerlukan Go 1.18+ (target Go 1.22+)

## 6. Referensi

- Lihat [05_ARCHITECTURE.md](05_ARCHITECTURE.md) untuk diagram arsitektur
- Lihat [03_FRD.md](03_FRD.md) untuk kebutuhan fungsional
- Lihat [17_DEVELOPER_SETUP.md](17_DEVELOPER_SETUP.md) untuk setup environment

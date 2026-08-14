# Architecture Document — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. High-Level Architecture

```mermaid
graph TB
    subgraph "WhatsApp Layer"
        WA[WhatsApp Server]
        WM[whatsmeow Client]
    end

    subgraph "Core Engine"
        MH[Message Handler]
        PM[Persona Manager]
        CM[Context Manager]
        RG[Response Generator]
    end

    subgraph "AI Layer"
        LLM[LLM Provider<br/>Claude / OpenAI]
        EMB[Embedding Service]
    end

    subgraph "Storage Layer"
        VDB[(Vector DB<br/>Memory Store)]
        SQL[(SQLite<br/>Metadata)]
        FS[File System<br/>Persona Files]
    end

    subgraph "Admin"
        AC[Admin Command<br/>Handler]
    end

    WA <-->|WebSocket| WM
    WM -->|Incoming Message| MH
    MH -->|Load Persona| PM
    MH -->|Build Context| CM
    CM -->|Retrieve Memories| VDB
    CM -->|Get History| SQL
    MH -->|Generate| RG
    RG -->|Prompt| LLM
    RG -->|Reply| WM
    MH -->|Extract & Store| EMB
    EMB -->|Embed + Save| VDB
    PM -->|Read Config| FS
    AC -->|Commands| MH
    MH -->|Log| SQL
```

## 2. Component Architecture

```mermaid
graph LR
    subgraph "cmd/"
        MAIN[main.go<br/>Entry point]
    end

    subgraph "internal/"
        subgraph "whatsapp/"
            CLIENT[client.go]
            HANDLER[handler.go]
            SENDER[sender.go]
        end

        subgraph "persona/"
            LOADER[loader.go]
            MANAGER[manager.go]
            TYPES[types.go]
        end

        subgraph "llm/"
            PROVIDER[provider.go<br/>Interface]
            CLAUDE[claude.go]
            OPENAI[openai.go]
        end

        subgraph "memory/"
            STORE[store.go<br/>Interface]
            VECTOR[vector.go]
            EMBEDDER[embedder.go]
        end

        subgraph "context/"
            BUILDER[builder.go]
            PRUNER[pruner.go]
        end

        subgraph "admin/"
            COMMANDS[commands.go]
            AUTH[auth.go]
        end

        subgraph "config/"
            CONFIG[config.go]
        end
    end

    MAIN --> CLIENT
    MAIN --> CONFIG
    CLIENT --> HANDLER
    HANDLER --> MANAGER
    HANDLER --> BUILDER
    HANDLER --> PROVIDER
    HANDLER --> COMMANDS
    BUILDER --> STORE
    PROVIDER --> CLAUDE
    PROVIDER --> OPENAI
```

## 3. Project Structure

```
wa-persona-ai/
├── cmd/
│   └── wa-persona-ai/
│       └── main.go                 # Entry point
├── internal/
│   ├── whatsapp/
│   │   ├── client.go               # WhatsApp client (whatsmeow wrapper)
│   │   ├── handler.go              # Message event handler
│   │   └── sender.go               # Message sending utilities
│   ├── persona/
│   │   ├── loader.go               # Load persona from YAML files
│   │   ├── manager.go              # Persona lifecycle management
│   │   └── types.go                # Persona structs and interfaces
│   ├── llm/
│   │   ├── provider.go             # LLM provider interface
│   │   ├── claude.go               # Anthropic Claude implementation
│   │   ├── openai.go               # OpenAI GPT implementation
│   │   └── types.go                # Request/response types
│   ├── memory/
│   │   ├── store.go                # Vector store interface
│   │   ├── vector.go               # Vector DB implementation
│   │   ├── embedder.go             # Text embedding service
│   │   └── types.go                # Memory entry types
│   ├── context/
│   │   ├── builder.go              # Context assembly for LLM
│   │   └── pruner.go               # Token-aware context pruning
│   ├── admin/
│   │   ├── commands.go             # Admin command handlers
│   │   └── auth.go                 # Admin authentication
│   └── config/
│       └── config.go               # Configuration management
├── persona/                         # Persona definition files
│   ├── _schema.yaml                # Persona YAML schema
│   ├── default.yaml                # Default persona
│   └── examples/
│       ├── assistant.yaml
│       ├── customer-service.yaml
│       └── companion.yaml
├── memory/                          # Memory data directory
│   ├── vectors/                    # Vector DB storage
│   └── metadata.db                 # SQLite metadata
├── docs/                            # Documentation
├── config/
│   └── config.example.yaml         # Example configuration
├── scripts/
│   └── setup.sh                    # Setup script
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
├── .gitignore
└── README.md
```

## 4. Data Flow

### 4.1 Message Processing Flow

```mermaid
sequenceDiagram
    participant U as User (WhatsApp)
    participant W as whatsmeow Client
    participant H as Message Handler
    participant P as Persona Manager
    participant C as Context Builder
    participant M as Memory Store
    participant L as LLM Provider
    participant S as SQLite

    U->>W: Send message
    W->>H: OnMessage event
    H->>P: Get active persona
    P-->>H: Persona config
    H->>S: Get conversation history
    S-->>H: Last N messages
    H->>M: Search relevant memories
    M-->>H: Top-K memories
    H->>C: Build context
    C-->>H: Assembled prompt
    H->>L: Generate response
    L-->>H: AI response
    H->>W: Send reply
    W->>U: Deliver message
    H->>M: Extract & store new memories
    H->>S: Save to conversation history
```

### 4.2 Memory Storage Flow

```mermaid
sequenceDiagram
    participant H as Message Handler
    participant E as Embedder
    participant V as Vector Store
    participant DB as SQLite

    H->>H: Analyze conversation turn
    H->>H: Extract important facts
    H->>E: Generate embedding
    E-->>H: Vector [1536 dims]
    H->>V: Store vector + metadata
    V-->>H: memory_id
    H->>DB: Save metadata index
```

## 5. Keputusan Arsitektur

| # | Keputusan | Konteks | Alternatif yang Ditolak |
|---|---|---|---|
| AD-001 | Go sebagai bahasa utama | Performance, single binary, concurrency | Node.js (heavier runtime), Python (slower) |
| AD-002 | whatsmeow untuk WhatsApp | Native Go, multi-device, maintained | Baileys (Node.js), selenium-based |
| AD-003 | Embedded vector DB | Simplicity, no infra dependency | Hosted Pinecone (cost, dependency) |
| AD-004 | Interface-based provider | Extensibility, testability | Hardcoded single provider |
| AD-005 | Persona as YAML files | Human-readable, version-controllable | Database-stored (harder to manage) |
| AD-006 | SQLite for metadata | Zero-config, embedded, reliable | PostgreSQL (overkill for single instance) |

## 6. Referensi

- Lihat [04_TRD.md](04_TRD.md) untuk tech stack detail
- Lihat [03_FRD.md](03_FRD.md) untuk kebutuhan fungsional
- Lihat [06_IMPLEMENTATION_PLAN.md](06_IMPLEMENTATION_PLAN.md) untuk fase implementasi

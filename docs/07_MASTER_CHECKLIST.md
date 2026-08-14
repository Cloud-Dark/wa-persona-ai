# Master Checklist — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## Fase 1: Foundation

- [ ] Inisialisasi Go module (`go mod init`)
- [ ] Setup project structure (cmd/, internal/, persona/, memory/, config/)
- [ ] Implementasi WhatsApp client wrapper (whatsmeow)
- [ ] QR code login flow
- [ ] Session persistence
- [ ] Auto-reconnect handler
- [ ] Implementasi message handler (receive messages)
- [ ] Implementasi message sender (send replies)
- [ ] Typing indicator
- [ ] Quoted reply support
- [ ] Configuration management (YAML + env)
- [ ] Structured logging setup
- [ ] Echo bot (basic reply test)
- [ ] Dockerfile
- [ ] docker-compose.yml
- [ ] Makefile
- [ ] .gitignore
- [ ] .env.example

## Fase 2: AI Integration

- [ ] LLM provider interface definition
- [ ] Claude API client implementation
- [ ] OpenAI API client implementation
- [ ] Token counting utility
- [ ] Retry mechanism with exponential backoff
- [ ] Persona type definitions & YAML schema
- [ ] Persona YAML loader with validation
- [ ] Persona manager (load, switch, reload)
- [ ] Default persona YAML file
- [ ] Context builder (system prompt + history + persona)
- [ ] Message handler ← → LLM integration
- [ ] End-to-end test: send WhatsApp message → get AI reply

## Fase 3: Memory System

- [ ] Memory entry type definitions
- [ ] Embedding service (OpenAI Ada-002 / alternative)
- [ ] Vector store interface
- [ ] Vector DB implementation (chromem-go / alternative)
- [ ] Memory extraction from conversations
- [ ] Important fact detection
- [ ] Deduplication logic
- [ ] Semantic search (similarity-based retrieval)
- [ ] Top-K retrieval with relevance scoring
- [ ] Recency bias in ranking
- [ ] Per-user memory isolation
- [ ] Context builder + memory integration
- [ ] SQLite metadata store
- [ ] Memory cleanup job (configurable retention)
- [ ] End-to-end test: bot remembers previous info

## Fase 4: Admin & Polish

- [ ] Admin authentication (phone number whitelist)
- [ ] `!status` command
- [ ] `!reset` command
- [ ] `!persona list` command
- [ ] `!persona <name>` command (switch)
- [ ] `!persona info` command
- [ ] `!memory clear` command
- [ ] `!memory stats` command
- [ ] `!reload` command
- [ ] `!stop` command
- [ ] Multi-persona support (hot-switch)
- [ ] Example persona files (assistant, CS, companion)
- [ ] Per-user rate limiting
- [ ] Global rate limiting
- [ ] Error handling improvement across modules
- [ ] Graceful shutdown with in-flight request completion
- [ ] Context pruning (token-aware)
- [ ] Prompt injection protection

## Fase 5: Testing & Documentation

- [ ] Unit tests: persona module (>80% coverage)
- [ ] Unit tests: memory module (>80% coverage)
- [ ] Unit tests: LLM module (>80% coverage)
- [ ] Unit tests: context module (>80% coverage)
- [ ] Unit tests: admin module (>80% coverage)
- [ ] Integration tests
- [ ] README.md update (full documentation)
- [ ] Configuration documentation
- [ ] API documentation
- [ ] Example configurations
- [ ] Contributing guide
- [ ] Release v0.1.0

## Post-Launch

- [ ] Group chat support
- [ ] Media message handling (images, voice notes)
- [ ] Emotion detection
- [ ] Scheduled messages
- [ ] Conversation export
- [ ] Web dashboard (v2)
- [ ] Multi-language support
- [ ] Plugin system

## Referensi

- Lihat [06_IMPLEMENTATION_PLAN.md](06_IMPLEMENTATION_PLAN.md) untuk detail per-task
- Lihat [08_ROADMAP.md](08_ROADMAP.md) untuk timeline

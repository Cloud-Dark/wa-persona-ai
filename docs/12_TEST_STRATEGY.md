# Test Strategy — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Testing Overview

| Level | Scope | Target Coverage | Tool |
|---|---|---|---|
| Unit Test | Functions, methods | > 80% | Go testing + testify |
| Integration Test | Module interactions | > 60% | Go testing + testcontainers |
| E2E Test | Full message flow | Key scenarios | Manual + scripted |
| Load Test | Concurrent messages | Performance targets | vegeta / k6 |

## 2. Unit Test Strategy

### 2.1 Persona Module
**Test Cases:**
- [ ] Load valid persona YAML file
- [ ] Reject invalid persona schema
- [ ] Fallback to default on missing file
- [ ] Switch persona at runtime
- [ ] Hot-reload detection
- [ ] Multiple personas loaded simultaneously

### 2.2 LLM Module
**Test Cases:**
- [ ] Successful API call and response parsing
- [ ] Handle API timeout
- [ ] Handle rate limit (429)
- [ ] Handle server error (500)
- [ ] Retry with exponential backoff
- [ ] Token counting accuracy
- [ ] Provider failover

### 2.3 Memory Module
**Test Cases:**
- [ ] Store memory entry with embedding
- [ ] Retrieve top-K relevant memories
- [ ] Per-user memory isolation
- [ ] Deduplication logic
- [ ] Relevance scoring accuracy
- [ ] Memory cleanup (retention policy)
- [ ] Empty result handling

### 2.4 Context Module
**Test Cases:**
- [ ] Build context with persona + history + memories
- [ ] Prune context when exceeding token limit
- [ ] Priority: system prompt > recent history > memories
- [ ] Empty history handling
- [ ] Empty memory handling

### 2.5 Admin Module
**Test Cases:**
- [ ] Authenticate valid admin number
- [ ] Reject non-admin user
- [ ] Parse command with arguments
- [ ] Handle unknown command
- [ ] Each command (!status, !reset, !persona, !memory) works correctly

## 3. Integration Test Strategy

### 3.1 LLM + Context + Persona
- Assemble context from persona config and conversation history
- Send to mock LLM provider
- Verify response formatting matches persona

### 3.2 Memory + Context
- Store memories via memory module
- Build context that includes retrieved memories
- Verify relevant memories appear in LLM prompt

### 3.3 WhatsApp Handler + All Modules
- Simulate incoming WhatsApp message
- Verify full pipeline: receive → process → reply
- Mock external services (WhatsApp server, LLM API)

## 4. E2E Test Scenarios

| # | Scenario | Steps | Expected |
|---|---|---|---|
| E2E-001 | Basic conversation | Send "Hello" → Get AI reply | Reply matches persona tone |
| E2E-002 | Memory recall | Tell name → Ask name later | Bot remembers name |
| E2E-003 | Persona switch | Admin sends `!persona cs` | Tone changes to CS style |
| E2E-004 | Admin command | Non-admin sends `!status` | Polite rejection |
| E2E-005 | Rate limiting | Send 100 messages in 10s | Bot rate-limits responses |
| E2E-006 | Reconnection | Disconnect WiFi → Reconnect | Bot auto-reconnects |

## 5. Performance Targets

| Metric | Target | Test Method |
|---|---|---|
| Message processing (excl. LLM) | < 100ms | Benchmark test |
| Memory retrieval (top-5) | < 200ms | Benchmark test |
| Concurrent conversations | > 50 | Load test |
| Memory footprint | < 512MB | Runtime monitoring |
| Context building | < 50ms | Benchmark test |

## 6. Referensi

- Lihat [03_FRD.md](03_FRD.md) untuk acceptance criteria
- Lihat [04_TRD.md](04_TRD.md) untuk performance requirements

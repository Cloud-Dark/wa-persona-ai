# Requirements Traceability Matrix — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Business → Functional Traceability

| Business Req | Functional Req | Test Case | Status |
|---|---|---|---|
| BR-001 (Aksesibilitas) | FR-001 (QR Auth), FR-002 (Receive), FR-003 (Send) | E2E-001, E2E-006 | 📋 Planned |
| BR-002 (Personalisasi) | FR-004 (Load), FR-005 (Switch), FR-006 (Behavior) | E2E-003 | 📋 Planned |
| BR-003 (Continuity) | FR-010 (Store), FR-011 (Retrieve), FR-012 (Manage) | E2E-002 | 📋 Planned |
| BR-004 (Open Source) | FR-015 (Logging), Documentation | - | 📋 Planned |
| BR-005 (Cost Efficiency) | FR-008 (Context Mgmt), Rate Limiting | E2E-005 | 📋 Planned |

## 2. Functional → Technical Traceability

| Functional Req | Technical Component | File(s) | Test File(s) |
|---|---|---|---|
| FR-001 QR Auth | whatsmeow client | `internal/whatsapp/client.go` | `client_test.go` |
| FR-002 Message Receive | Event handler | `internal/whatsapp/handler.go` | `handler_test.go` |
| FR-003 Message Send | Sender utility | `internal/whatsapp/sender.go` | `sender_test.go` |
| FR-004 Persona Load | YAML loader | `internal/persona/loader.go` | `loader_test.go` |
| FR-005 Persona Switch | Persona manager | `internal/persona/manager.go` | `manager_test.go` |
| FR-006 Persona Behavior | Persona types | `internal/persona/types.go` | `types_test.go` |
| FR-007 LLM Provider | Provider interface | `internal/llm/provider.go` | `provider_test.go` |
| FR-008 Context Mgmt | Context builder | `internal/context/builder.go` | `builder_test.go` |
| FR-009 Response Gen | LLM + handler | `internal/llm/*.go` | `*_test.go` |
| FR-010 Memory Store | Vector store | `internal/memory/vector.go` | `vector_test.go` |
| FR-011 Memory Retrieve | Retriever | `internal/memory/retriever.go` | `retriever_test.go` |
| FR-012 Memory Manage | Admin commands | `internal/admin/commands.go` | `commands_test.go` |
| FR-013 Admin Auth | Auth whitelist | `internal/admin/auth.go` | `auth_test.go` |
| FR-014 Bot Control | Command handler | `internal/admin/commands.go` | `commands_test.go` |
| FR-015 Logging | Logger config | `internal/config/config.go` | `config_test.go` |

## 3. Feature → User Story → Acceptance Criteria

| Feature ID | User Story | Acceptance Criteria | Phase |
|---|---|---|---|
| F-001 | US-001 | Bot replies within 5s | Fase 1 |
| F-002 | US-001 | Receives text messages | Fase 1 |
| F-003 | US-001 | AI-powered response | Fase 2 |
| F-004 | US-002 | Configurable persona | Fase 2 |
| F-007 | US-002 | Switch persona at runtime | Fase 4 |
| F-008 | US-003 | Vector memory storage | Fase 3 |
| F-009 | US-003 | Recall previous info | Fase 3 |
| F-015 | US-004 | Admin commands work | Fase 4 |

## 4. Referensi

- Lihat [03_FRD.md](03_FRD.md) untuk detail functional requirements
- Lihat [02_BRD.md](02_BRD.md) untuk business requirements
- Lihat [12_TEST_STRATEGY.md](12_TEST_STRATEGY.md) untuk test plan

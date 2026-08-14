# Risk Register — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## Risk Matrix

| ID | Risiko | Kategori | Likelihood | Impact | Level | Mitigasi | Owner | Status |
|---|---|---|---|---|---|---|---|---|
| R-001 | WhatsApp ban nomor bot | Operational | Medium | Critical | 🔴 High | Rate limiting ketat, behavior natural, delay random antar pesan | Cloud-Dark | Open |
| R-002 | whatsmeow deprecated/breaking change | Technical | Low | High | 🟡 Medium | Pin version, monitor repo, abstraksi client layer | Cloud-Dark | Open |
| R-003 | LLM API cost membengkak | Financial | Medium | Medium | 🟡 Medium | Token budgeting, response caching, model selection configurable | Cloud-Dark | Open |
| R-004 | LLM API downtime | Technical | Low | Medium | 🟢 Low | Multi-provider fallback, graceful degradation | Cloud-Dark | Open |
| R-005 | Data privacy breach | Security | Low | Critical | 🟡 Medium | Local storage only, no cloud sync default, encryption option | Cloud-Dark | Open |
| R-006 | Prompt injection via user message | Security | Medium | High | 🔴 High | Input sanitization, system prompt hardening, output filtering | Cloud-Dark | Open |
| R-007 | Vector DB data corruption | Technical | Low | High | 🟡 Medium | Regular backup, integrity checks, recovery procedure | Cloud-Dark | Open |
| R-008 | Memory bloat (too many entries) | Operational | Medium | Medium | 🟡 Medium | Auto-cleanup policy, memory size limits, importance scoring | Cloud-Dark | Open |
| R-009 | Concurrent message race condition | Technical | Medium | Medium | 🟡 Medium | Channel-based processing, mutex on shared state | Cloud-Dark | Open |
| R-010 | AI hallucination (incorrect memory recall) | Quality | Medium | Medium | 🟡 Medium | Relevance threshold, confidence scoring, source attribution | Cloud-Dark | Open |

## Likelihood Scale

| Level | Deskripsi | Probability |
|---|---|---|
| Low | Jarang terjadi | < 20% |
| Medium | Mungkin terjadi | 20-60% |
| High | Kemungkinan besar terjadi | > 60% |

## Impact Scale

| Level | Deskripsi |
|---|---|
| Low | Minor inconvenience, workaround tersedia |
| Medium | Fungsi terganggu, perlu perbaikan |
| High | Fungsi utama terganggu, perlu perbaikan segera |
| Critical | Sistem tidak bisa digunakan |

## Referensi

- Lihat [14_SECURITY_THREAT_MODEL.md](14_SECURITY_THREAT_MODEL.md) untuk analisis keamanan
- Lihat [02_BRD.md](02_BRD.md) untuk constraint bisnis

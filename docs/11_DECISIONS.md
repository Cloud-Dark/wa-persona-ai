# Architecture Decision Records — WA Persona AI

> Status: Living Document
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## ADR Index

| # | Tanggal | Keputusan | Status |
|---|---|---|---|
| ADR-001 | 2026-08-14 | Golang sebagai bahasa utama | ✅ Accepted |
| ADR-002 | 2026-08-14 | whatsmeow sebagai WhatsApp library | ✅ Accepted |
| ADR-003 | 2026-08-14 | Embedded vector DB untuk memory | ✅ Accepted |
| ADR-004 | 2026-08-14 | YAML untuk persona configuration | ✅ Accepted |
| ADR-005 | 2026-08-14 | Interface-based LLM provider | ✅ Accepted |
| ADR-006 | 2026-08-14 | SQLite untuk metadata storage | ✅ Accepted |

---

## ADR-001: Golang sebagai Bahasa Utama

**Tanggal:** 2026-08-14
**Status:** Accepted

**Konteks:**
Dibutuhkan bahasa pemrograman yang performant, memiliki concurrency model yang baik, dan bisa di-deploy sebagai single binary.

**Keputusan:**
Menggunakan Go (1.22+) sebagai bahasa pemrograman utama.

**Alasan:**
- Single binary deployment — tidak perlu runtime
- Goroutine untuk handle concurrent WhatsApp messages
- whatsmeow ditulis dalam Go — integrasi native
- Cross-compilation untuk berbagai platform
- Memory footprint rendah, cocok untuk VPS kecil

**Alternatif yang Ditolak:**
- **Node.js** — Runtime dependency, memory usage lebih tinggi, Baileys (WhatsApp lib) less stable
- **Python** — Slower performance, GIL limitation untuk concurrency
- **Rust** — Learning curve terlalu tinggi, development speed lebih lambat

---

## ADR-002: whatsmeow sebagai WhatsApp Library

**Tanggal:** 2026-08-14
**Status:** Accepted

**Konteks:**
Dibutuhkan library untuk terhubung ke WhatsApp tanpa menggunakan WhatsApp Business API (yang berbayar dan memerlukan approval).

**Keputusan:**
Menggunakan whatsmeow (github.com/tulir/whatsmeow) — Go library untuk WhatsApp Web multi-device API.

**Alasan:**
- Native Go, konsisten dengan tech stack
- Multi-device support
- Actively maintained oleh tulir (pembuat mautrix-whatsapp bridge)
- Well-documented dengan contoh yang lengkap

**Risiko:**
- Unofficial API — bisa break jika WhatsApp update protocol
- Penggunaan harus hati-hati untuk menghindari ban

---

## ADR-003: Embedded Vector DB untuk Memory

**Tanggal:** 2026-08-14
**Status:** Accepted

**Konteks:**
Memory system membutuhkan vector database untuk semantic search. Pilihan antara hosted (Pinecone, Weaviate Cloud) vs embedded (chromem-go, SQLite+VSS).

**Keputusan:**
Menggunakan embedded vector DB (chromem-go atau SQLite with vector extension) yang berjalan in-process.

**Alasan:**
- Zero infrastructure dependency — tidak perlu service terpisah
- Data tetap lokal (privacy)
- Cukup untuk skala single-instance (ratusan ribu entries)
- Konsisten dengan prinsip "self-hosted, minimal dependency"

**Alternatif yang Ditolak:**
- **Pinecone** — Cloud dependency, biaya recurring, data di luar kontrol
- **Milvus** — Terlalu heavy untuk single instance
- **Qdrant** — Butuh service terpisah

---

## ADR-004: YAML untuk Persona Configuration

**Tanggal:** 2026-08-14
**Status:** Accepted

**Konteks:**
Persona perlu disimpan dalam format yang human-readable, mudah di-edit, dan bisa di-version control.

**Keputusan:**
Menggunakan YAML files di folder `persona/` untuk definisi persona.

**Alasan:**
- Human-readable — non-developer bisa edit
- Support multi-line string (cocok untuk system prompt)
- Bisa di-version control bersama kode
- Hot-reload tanpa restart dengan file watcher

**Alternatif yang Ditolak:**
- **Database** — Harder to manage, butuh migration, kurang readable
- **JSON** — Kurang readable untuk text panjang, tidak support comment
- **TOML** — Kurang populer, multi-line string awkward

---

## ADR-005: Interface-based LLM Provider

**Tanggal:** 2026-08-14
**Status:** Accepted

**Konteks:**
Aplikasi harus support multiple LLM provider (Claude, OpenAI, dll) dan mudah menambah provider baru.

**Keputusan:**
Mendefinisikan `LLMProvider` interface di Go, dengan implementasi terpisah per provider.

**Alasan:**
- Extensible — tambah provider baru tanpa ubah core logic
- Testable — mudah di-mock untuk unit testing
- Dependency inversion — core tidak depend pada provider spesifik
- Failover — bisa switch provider at runtime

---

## ADR-006: SQLite untuk Metadata Storage

**Tanggal:** 2026-08-14
**Status:** Accepted

**Konteks:**
Dibutuhkan database untuk menyimpan conversation history, user metadata, dan memory index.

**Keputusan:**
Menggunakan SQLite sebagai database metadata.

**Alasan:**
- Zero-config — tidak perlu install database server
- Embedded — berjalan in-process
- Reliable — battle-tested, ACID compliant
- Portable — single file, mudah backup
- Cocok untuk single-instance deployment

**Alternatif yang Ditolak:**
- **PostgreSQL** — Overkill untuk single instance, butuh service terpisah
- **MongoDB** — Butuh service terpisah, kurang cocok untuk relational data
- **BoltDB** — Terlalu low-level, kurang fitur

## Referensi

- Lihat [05_ARCHITECTURE.md](05_ARCHITECTURE.md) untuk diagram arsitektur
- Lihat [04_TRD.md](04_TRD.md) untuk tech stack lengkap

# Project Charter — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Visi

Membangun platform WhatsApp AI yang mampu merespons pesan secara otomatis dengan persona yang dapat dikustomisasi, memberikan pengalaman percakapan yang natural dan personal layaknya berbicara dengan manusia sungguhan.

## 2. Tujuan

- Membuat WhatsApp bot berbasis Go (whatsmeow) yang terhubung dengan LLM untuk menghasilkan respons cerdas.
- Menyediakan sistem persona yang fleksibel sehingga AI dapat berperan sebagai karakter berbeda (customer service, teman curhat, asisten bisnis, dll).
- Mengimplementasikan memory system berbasis vector database agar AI dapat mengingat konteks percakapan sebelumnya.
- Menyediakan arsitektur modular yang mudah di-extend dan di-maintain.

## 3. Scope

### In-Scope

- Integrasi WhatsApp menggunakan library whatsmeow (Go)
- Sistem persona yang dapat dikonfigurasi via file/database
- Memory system menggunakan vector database untuk long-term memory
- Conversation history management
- Multi-persona support (satu instance bisa handle banyak persona)
- LLM integration (Claude/OpenAI/local model)
- Rate limiting dan spam protection
- Logging dan monitoring dasar

### Out-of-Scope

- WhatsApp Business API (menggunakan whatsmeow/unofficial)
- Mobile app atau web dashboard (fase awal)
- Payment processing
- Multi-device management dari satu nomor
- End-to-end encryption key management custom

## 4. Stakeholder

| Peran | Deskripsi |
|---|---|
| Project Owner | Cloud-Dark |
| Developer | Cloud-Dark (+ kontributor open source) |
| End User | Pengguna WhatsApp yang berinteraksi dengan bot |
| Admin | Pemilik bot yang mengkonfigurasi persona dan settings |

## 5. Kriteria Sukses

| # | Kriteria | Target |
|---|---|---|
| 1 | Bot dapat menerima dan membalas pesan WhatsApp | 100% uptime saat aktif |
| 2 | Respons AI relevan dengan persona yang dipilih | > 90% response accuracy |
| 3 | Memory system menyimpan dan merecall konteks | Recall accuracy > 85% |
| 4 | Latency respons | < 5 detik untuk pesan teks biasa |
| 5 | Dapat berjalan stabil | Uptime > 99% dalam 24 jam |

## 6. Constraint

- Menggunakan Go sebagai bahasa pemrograman utama
- Library WhatsApp: whatsmeow (unofficial API)
- Harus bisa berjalan di VPS minimal (2GB RAM)
- Mematuhi rate limit WhatsApp untuk menghindari ban

## 7. Referensi

- [whatsmeow](https://github.com/tulir/whatsmeow) — Go library untuk WhatsApp Web multidevice API
- [Anthropic Claude API](https://docs.anthropic.com) — LLM provider
- Lihat [01_PRD.md](01_PRD.md) untuk detail fitur

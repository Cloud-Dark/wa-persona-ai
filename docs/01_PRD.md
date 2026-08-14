# Product Requirements Document (PRD) — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Problem Statement

Pengguna WhatsApp yang ingin memiliki asisten AI personal sering kali terbatas pada chatbot yang kaku, tidak memiliki kepribadian, dan tidak mengingat percakapan sebelumnya. Solusi yang ada di pasaran umumnya:

- Tidak mendukung persona yang bisa dikustomisasi
- Tidak memiliki long-term memory
- Menggunakan WhatsApp Business API yang mahal
- Tidak open-source dan sulit dikustomisasi

## 2. Target User

| Persona | Deskripsi | Kebutuhan Utama |
|---|---|---|
| Individual User | Pengguna yang ingin AI companion di WhatsApp | Persona yang natural, memory yang baik |
| Small Business | UMKM yang butuh CS otomatis | Persona profesional, respons cepat |
| Developer | Developer yang ingin extend/customize | API yang clean, arsitektur modular |
| Community Admin | Admin grup yang butuh bot helper | Multi-persona, keyword trigger |

## 3. Fitur Utama

### Must Have (P0)

| ID | Fitur | Deskripsi |
|---|---|---|
| F-001 | WhatsApp Connection | Koneksi ke WhatsApp via whatsmeow dengan QR code login |
| F-002 | Message Handler | Menerima dan memproses pesan masuk (teks) |
| F-003 | LLM Integration | Integrasi dengan Claude/OpenAI API untuk generate respons |
| F-004 | Basic Persona | Satu persona default yang bisa dikonfigurasi |
| F-005 | Conversation Context | Menyimpan konteks percakapan dalam sesi aktif |
| F-006 | Reply System | Mengirim balasan otomatis ke pengirim pesan |

### Should Have (P1)

| ID | Fitur | Deskripsi |
|---|---|---|
| F-007 | Multi-Persona | Support multiple persona yang bisa di-switch |
| F-008 | Vector Memory | Long-term memory menggunakan vector database |
| F-009 | Memory Recall | Mengingat informasi dari percakapan sebelumnya |
| F-010 | Persona Template | Template persona yang bisa di-load dari file |
| F-011 | Rate Limiting | Pembatasan request untuk menghindari ban |
| F-012 | Media Handler | Handle gambar/voice note (minimal dengan deskripsi) |

### Could Have (P2)

| ID | Fitur | Deskripsi |
|---|---|---|
| F-013 | Group Chat Support | Bot bisa merespons di grup (mention/reply) |
| F-014 | Scheduled Messages | Pesan terjadwal berdasarkan trigger |
| F-015 | Admin Commands | Command khusus admin (!switch, !reset, !persona) |
| F-016 | Conversation Export | Export riwayat percakapan |
| F-017 | Emotion Detection | Deteksi emosi user dan sesuaikan tone respons |
| F-018 | Multi-language | Support multi-bahasa otomatis |

### Won't Have (saat ini)

| ID | Fitur | Alasan |
|---|---|---|
| F-019 | Web Dashboard | Fokus ke backend dulu |
| F-020 | Voice Call | whatsmeow belum support |
| F-021 | Payment Integration | Di luar scope |

## 4. User Stories

### US-001: Basic Conversation
**Sebagai** pengguna WhatsApp,
**saya ingin** mengirim pesan ke nomor bot dan mendapatkan balasan AI yang natural,
**sehingga** saya merasa seperti berbicara dengan orang sungguhan.

**Acceptance Criteria:**
- [ ] User mengirim pesan teks, bot membalas dalam < 5 detik
- [ ] Respons sesuai dengan persona yang dikonfigurasi
- [ ] Respons dalam bahasa yang sama dengan pesan user

### US-002: Persona Switching
**Sebagai** admin bot,
**saya ingin** mengganti persona AI kapan saja,
**sehingga** bot bisa disesuaikan dengan kebutuhan (CS, teman, asisten).

**Acceptance Criteria:**
- [ ] Admin bisa switch persona via command
- [ ] Persona baru langsung aktif tanpa restart
- [ ] Setiap persona punya system prompt, tone, dan behavior sendiri

### US-003: Memory Recall
**Sebagai** pengguna WhatsApp,
**saya ingin** bot mengingat apa yang pernah saya ceritakan sebelumnya,
**sehingga** percakapan terasa continuous dan personal.

**Acceptance Criteria:**
- [ ] Bot mengingat nama user setelah diperkenalkan
- [ ] Bot bisa recall topik yang pernah dibahas
- [ ] Memory persist meskipun bot di-restart

### US-004: Admin Control
**Sebagai** admin bot,
**saya ingin** mengontrol behavior bot melalui command,
**sehingga** saya bisa manage bot tanpa akses server.

**Acceptance Criteria:**
- [ ] Command `!reset` untuk reset konteks percakapan
- [ ] Command `!persona <name>` untuk switch persona
- [ ] Command `!status` untuk cek status bot
- [ ] Command `!memory clear` untuk hapus memory user tertentu

## 5. Metrik Keberhasilan

| Metrik | Target | Cara Ukur |
|---|---|---|
| Response Time | < 5 detik | Log timestamp |
| Persona Accuracy | > 90% | Manual review sampling |
| Memory Recall Rate | > 85% | Test recall scenarios |
| Uptime | > 99% | Health check monitoring |
| User Satisfaction | > 4/5 | Feedback command |

## 6. Timeline Estimasi

| Fase | Durasi | Deliverable |
|---|---|---|
| Fase 1: Foundation | 2 minggu | WhatsApp connection, basic reply, single persona |
| Fase 2: Intelligence | 2 minggu | LLM integration, conversation context, multi-persona |
| Fase 3: Memory | 2 minggu | Vector DB, long-term memory, recall system |
| Fase 4: Polish | 1 minggu | Rate limiting, error handling, admin commands |
| Fase 5: Testing | 1 minggu | Integration testing, load testing, documentation |

## 7. Referensi

- Lihat [00_PROJECT_CHARTER.md](00_PROJECT_CHARTER.md) untuk visi dan scope
- Lihat [02_BRD.md](02_BRD.md) untuk kebutuhan bisnis
- Lihat [03_FRD.md](03_FRD.md) untuk kebutuhan fungsional detail

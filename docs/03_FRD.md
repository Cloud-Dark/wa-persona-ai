# Functional Requirements Document (FRD) — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Pendahuluan

Dokumen ini mendefinisikan kebutuhan fungsional detail untuk WA Persona AI. Setiap requirement memiliki ID unik dan acceptance criteria yang terukur.

## 2. Kebutuhan Fungsional

### 2.1 WhatsApp Connection Module

#### FR-001: QR Code Authentication
**Deskripsi:** Sistem harus mampu melakukan autentikasi ke WhatsApp melalui QR code scan menggunakan whatsmeow.
**Acceptance Criteria:**
- [ ] QR code ditampilkan di terminal saat pertama kali login
- [ ] Session tersimpan sehingga tidak perlu scan ulang setelah restart
- [ ] Sistem mendeteksi dan menangani disconnect/reconnect otomatis
- [ ] Log status koneksi (connected, disconnected, reconnecting)

#### FR-002: Message Receiving
**Deskripsi:** Sistem harus mampu menerima pesan masuk dari WhatsApp.
**Acceptance Criteria:**
- [ ] Menerima pesan teks dari private chat
- [ ] Menerima pesan dari group chat (dengan filter mention/reply)
- [ ] Mengekstrak metadata: sender JID, timestamp, message type, chat type
- [ ] Handle pesan media (gambar, voice note) minimal dengan fallback text

#### FR-003: Message Sending
**Deskripsi:** Sistem harus mampu mengirim pesan balasan ke WhatsApp.
**Acceptance Criteria:**
- [ ] Mengirim pesan teks ke private chat
- [ ] Reply ke pesan spesifik (quoted reply)
- [ ] Mengirim typing indicator sebelum membalas
- [ ] Mengirim pesan dengan formatting (bold, italic, monospace)

### 2.2 Persona Module

#### FR-004: Persona Loading
**Deskripsi:** Sistem harus mampu memuat konfigurasi persona dari file.
**Acceptance Criteria:**
- [ ] Load persona dari file YAML/JSON di folder `persona/`
- [ ] Validasi schema persona saat loading
- [ ] Fallback ke persona default jika file rusak/tidak ada
- [ ] Hot-reload persona tanpa restart aplikasi

#### FR-005: Persona Switching
**Deskripsi:** Admin harus bisa mengganti persona aktif melalui command.
**Acceptance Criteria:**
- [ ] Command `!persona <name>` mengganti persona aktif
- [ ] Command `!persona list` menampilkan daftar persona tersedia
- [ ] Command `!persona info` menampilkan persona aktif saat ini
- [ ] Persona baru langsung diterapkan pada pesan berikutnya

#### FR-006: Persona Behavior
**Deskripsi:** Setiap persona harus memiliki behavior yang konsisten.
**Acceptance Criteria:**
- [ ] Persona menentukan system prompt untuk LLM
- [ ] Persona menentukan tone/style bahasa
- [ ] Persona menentukan batasan topik (apa yang boleh/tidak dibahas)
- [ ] Persona menentukan greeting message untuk user baru
- [ ] Persona memiliki character traits yang konsisten

### 2.3 LLM Integration Module

#### FR-007: LLM Provider Integration
**Deskripsi:** Sistem harus mampu berkomunikasi dengan LLM API.
**Acceptance Criteria:**
- [ ] Support Claude API (Anthropic)
- [ ] Support OpenAI API (GPT)
- [ ] Abstraksi provider sehingga mudah menambah provider baru
- [ ] Konfigurasi provider via environment variable / config file
- [ ] Fallback ke provider alternatif jika primary down

#### FR-008: Context Management
**Deskripsi:** Sistem harus mengelola konteks percakapan yang dikirim ke LLM.
**Acceptance Criteria:**
- [ ] Menyertakan N pesan terakhir sebagai konteks (configurable)
- [ ] Menyertakan system prompt dari persona aktif
- [ ] Menyertakan relevant memories dari vector DB
- [ ] Token counting untuk menghindari limit
- [ ] Context pruning saat mendekati token limit

#### FR-009: Response Generation
**Deskripsi:** Sistem harus menghasilkan respons yang sesuai persona.
**Acceptance Criteria:**
- [ ] Respons sesuai dengan tone dan style persona
- [ ] Respons relevan dengan konteks percakapan
- [ ] Handle error dari LLM gracefully (timeout, rate limit, dll)
- [ ] Retry mechanism dengan exponential backoff
- [ ] Response length configurable per persona

### 2.4 Memory Module

#### FR-010: Memory Storage
**Deskripsi:** Sistem harus menyimpan informasi penting dari percakapan ke vector database.
**Acceptance Criteria:**
- [ ] Otomatis mengekstrak dan menyimpan fakta penting dari percakapan
- [ ] Menyimpan embedding vector dari setiap memory entry
- [ ] Metadata: user JID, timestamp, topic, importance score
- [ ] Deduplication untuk menghindari memory duplikat

#### FR-011: Memory Retrieval
**Deskripsi:** Sistem harus mampu mengambil memory yang relevan saat memproses pesan baru.
**Acceptance Criteria:**
- [ ] Semantic search berdasarkan similarity dengan pesan masuk
- [ ] Top-K retrieval (configurable, default K=5)
- [ ] Filter berdasarkan user JID (memory per user)
- [ ] Relevance scoring dan minimum threshold
- [ ] Recency bias (memory baru lebih diprioritaskan)

#### FR-012: Memory Management
**Deskripsi:** Admin harus bisa mengelola memory.
**Acceptance Criteria:**
- [ ] Command `!memory clear <user>` untuk hapus memory user
- [ ] Command `!memory stats` untuk melihat statistik memory
- [ ] Auto-cleanup memory lama (configurable retention period)
- [ ] Export memory ke file (backup)

### 2.5 Admin Module

#### FR-013: Admin Authentication
**Deskripsi:** Sistem harus membedakan admin dan user biasa.
**Acceptance Criteria:**
- [ ] Admin didefinisikan via config (daftar nomor WhatsApp)
- [ ] Hanya admin yang bisa menggunakan command `!`
- [ ] User biasa yang mencoba command admin mendapat respons error sopan

#### FR-014: Bot Control Commands
**Deskripsi:** Admin harus bisa mengontrol bot melalui WhatsApp.
**Acceptance Criteria:**
- [ ] `!status` — menampilkan status bot (uptime, memory usage, active persona)
- [ ] `!reset <user>` — reset konteks percakapan user tertentu
- [ ] `!reset all` — reset semua konteks percakapan
- [ ] `!reload` — reload konfigurasi dan persona
- [ ] `!stop` — graceful shutdown bot

### 2.6 Logging & Monitoring

#### FR-015: Application Logging
**Deskripsi:** Sistem harus mencatat log aktivitas.
**Acceptance Criteria:**
- [ ] Log level configurable (DEBUG, INFO, WARN, ERROR)
- [ ] Log ke file dengan rotasi
- [ ] Log ke stdout untuk container deployment
- [ ] Structured logging (JSON format)
- [ ] Tidak log isi pesan user (privacy) kecuali mode debug

## 3. Non-Functional Requirements

Lihat [04_TRD.md](04_TRD.md) untuk kebutuhan teknis non-fungsional.

## 4. Referensi

- Lihat [01_PRD.md](01_PRD.md) untuk user stories
- Lihat [04_TRD.md](04_TRD.md) untuk kebutuhan teknis
- Lihat [19_REQUIREMENTS_TRACEABILITY.md](19_REQUIREMENTS_TRACEABILITY.md) untuk traceability matrix

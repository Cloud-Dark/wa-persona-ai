# Glossary — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## Istilah

| Istilah | Definisi |
|---|---|
| **Admin** | Pengguna yang memiliki akses command `!` untuk mengontrol bot. Didefinisikan via whitelist nomor WhatsApp di config. |
| **Context** | Kumpulan informasi (system prompt, conversation history, memories) yang dikirim ke LLM sebagai input untuk generate respons. |
| **Context Pruning** | Proses memangkas context yang terlalu panjang agar tidak melebihi token limit LLM, dengan prioritas: system prompt > recent history > memories. |
| **Cosine Similarity** | Metrik pengukuran kesamaan antara dua vector, digunakan untuk semantic search pada memory retrieval. Nilai 0-1, semakin tinggi semakin mirip. |
| **Embedding** | Representasi teks dalam bentuk vector numerik berdimensi tinggi (biasanya 1536 dimensi). Digunakan untuk semantic search pada memory system. |
| **Goroutine** | Lightweight thread di Go, digunakan untuk handle concurrent message processing. |
| **Greeting Message** | Pesan sapaan yang dikirim bot saat pertama kali berkomunikasi dengan user baru. Ditentukan oleh persona aktif. |
| **Hot-Reload** | Kemampuan memuat ulang konfigurasi (persona, settings) tanpa perlu restart aplikasi. |
| **JID** | Jabber ID — identifier unik user di WhatsApp, format: `nomor@s.whatsapp.net`. |
| **LLM** | Large Language Model — model AI yang digunakan untuk generate respons (contoh: Claude, GPT-4). |
| **Memory Entry** | Satu unit informasi yang disimpan di vector database, berisi: content, embedding vector, metadata (user, timestamp, topic). |
| **Memory Recall** | Proses mengambil informasi relevan dari memory berdasarkan semantic similarity dengan pesan masuk. |
| **Persona** | Konfigurasi kepribadian AI yang menentukan system prompt, tone, traits, dan constraints. Disimpan sebagai file YAML. |
| **Prompt Injection** | Serangan di mana user memasukkan instruksi berbahaya dalam pesan untuk memanipulasi behavior LLM. |
| **QR Code Login** | Metode autentikasi WhatsApp di mana bot menampilkan QR code yang harus di-scan oleh ponsel. |
| **Rate Limiting** | Pembatasan jumlah request per periode waktu untuk menghindari ban dari WhatsApp dan mengontrol biaya LLM API. |
| **Semantic Search** | Pencarian berdasarkan makna (bukan keyword exact match), menggunakan vector similarity pada embeddings. |
| **Session** | Data koneksi WhatsApp yang tersimpan agar tidak perlu scan QR code ulang setiap restart. |
| **System Prompt** | Instruksi awal yang diberikan ke LLM untuk menentukan behavior dan kepribadian AI. Bagian inti dari persona. |
| **Token** | Unit terkecil text yang diproses oleh LLM. Rata-rata 1 token ≈ 4 karakter atau 0.75 kata bahasa Inggris. |
| **Top-K Retrieval** | Mengambil K memory entries yang paling relevan berdasarkan similarity score. Default K=5. |
| **Traits** | Sifat-sifat karakter yang melekat pada persona (contoh: "ramah", "profesional", "humoris"). |
| **Vector Database** | Database khusus yang menyimpan dan mengindex vector embeddings untuk fast similarity search. |
| **whatsmeow** | Go library open-source untuk mengakses WhatsApp Web multi-device API. Dikembangkan oleh tulir. |

## Referensi

- Lihat [01_PRD.md](01_PRD.md) untuk konteks fitur
- Lihat [04_TRD.md](04_TRD.md) untuk detail teknis

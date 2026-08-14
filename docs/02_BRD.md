# Business Requirements Document (BRD) — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Executive Summary

WA Persona AI adalah solusi open-source WhatsApp chatbot berbasis Go yang menggunakan Large Language Model (LLM) dengan sistem persona yang dapat dikustomisasi dan long-term memory. Proyek ini menjawab kebutuhan pasar akan AI assistant yang personal, cerdas, dan terjangkau di platform WhatsApp.

## 2. Kebutuhan Bisnis

### BR-001: Aksesibilitas AI Assistant
**Kebutuhan:** Menyediakan AI assistant yang bisa diakses melalui WhatsApp tanpa perlu install aplikasi tambahan.
**Justifikasi:** WhatsApp memiliki 2+ milyar pengguna aktif. Membawa AI ke platform yang sudah digunakan sehari-hari menghilangkan barrier to entry.
**Prioritas:** Critical

### BR-002: Personalisasi AI
**Kebutuhan:** AI yang bisa disesuaikan kepribadiannya sesuai use case.
**Justifikasi:** One-size-fits-all chatbot tidak efektif. CS butuh tone profesional, companion butuh tone hangat. Persona system memungkinkan satu platform untuk berbagai kebutuhan.
**Prioritas:** Critical

### BR-003: Continuity of Experience
**Kebutuhan:** AI yang mengingat percakapan sebelumnya untuk memberikan pengalaman yang konsisten.
**Justifikasi:** Chatbot tanpa memory terasa frustrating bagi user. Long-term memory meningkatkan engagement dan kepuasan user.
**Prioritas:** High

### BR-004: Open Source & Self-Hosted
**Kebutuhan:** Solusi yang bisa di-self-host dan dikustomisasi tanpa vendor lock-in.
**Justifikasi:** Banyak solusi WhatsApp bot yang proprietary dan mahal. Open source memungkinkan adopsi luas dan kontribusi komunitas.
**Prioritas:** High

### BR-005: Cost Efficiency
**Kebutuhan:** Bisa berjalan di infrastruktur minimal dengan biaya rendah.
**Justifikasi:** Target user termasuk individu dan UMKM yang budget-sensitive. Arsitektur Go yang ringan memungkinkan deployment di VPS murah.
**Prioritas:** Medium

## 3. Dampak Bisnis

| Area | Dampak | Estimasi |
|---|---|---|
| User Acquisition | AI di WhatsApp menurunkan barrier to entry | Potensi ribuan pengguna via GitHub |
| Cost Reduction | Self-hosted mengurangi biaya SaaS | Hemat 80% vs solusi berbayar |
| Customer Satisfaction | Persona + memory = pengalaman lebih baik | NPS > 50 |
| Community Growth | Open source mendorong kontribusi | Target 100+ GitHub stars |
| Time to Market | Whatsmeow + Go = development cepat | MVP dalam 4 minggu |

## 4. Constraint Bisnis

| # | Constraint | Dampak | Mitigasi |
|---|---|---|---|
| 1 | WhatsApp ToS (unofficial API) | Risiko ban nomor | Rate limiting, behavior natural |
| 2 | LLM API cost | Biaya per request | Caching, context pruning, model selection |
| 3 | Privacy concern | Data percakapan sensitif | Local storage, encryption, clear privacy policy |
| 4 | Single developer (awal) | Bandwidth terbatas | Prioritas fitur ketat, community contribution |

## 5. Success Metrics

| KPI | Target | Measurement Period |
|---|---|---|
| GitHub Stars | > 100 | 3 bulan setelah rilis |
| Active Deployments | > 20 | 3 bulan setelah rilis |
| Community Contributors | > 5 | 6 bulan setelah rilis |
| Issue Resolution Time | < 7 hari | Ongoing |
| Documentation Coverage | > 90% | Setiap rilis |

## 6. Stakeholder Analysis

| Stakeholder | Interest | Influence | Strategy |
|---|---|---|---|
| End Users | Tinggi | Rendah | Feedback loop, survey |
| Developers | Tinggi | Tinggi | Good docs, clean API, contribution guide |
| Business Users | Medium | Medium | Use case templates, easy setup |
| WhatsApp (Meta) | Rendah | Tinggi | Comply with ToS, responsible use |

## 7. Referensi

- Lihat [00_PROJECT_CHARTER.md](00_PROJECT_CHARTER.md) untuk visi proyek
- Lihat [01_PRD.md](01_PRD.md) untuk detail produk
- Lihat [10_RISK_REGISTER.md](10_RISK_REGISTER.md) untuk analisis risiko lengkap

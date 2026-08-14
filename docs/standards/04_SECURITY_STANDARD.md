# Security Standard — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Secrets Management

- **WAJIB:** API keys hanya via environment variable
- **DILARANG:** Hardcode secrets di source code
- **DILARANG:** Commit file `.env` ke Git
- `.env.example` hanya berisi placeholder, BUKAN nilai asli

## 2. Input Validation

- Semua input user di-sanitize sebelum diproses
- Limit panjang pesan (configurable, default 4096 chars)
- Filter karakter control dan injection patterns
- Validate format JID sebelum operasi database

## 3. Data Privacy

- Isi pesan user TIDAK boleh di-log kecuali mode DEBUG
- Memory storage lokal by default, TIDAK dikirim ke third-party
- Per-user memory isolation — DILARANG cross-reference
- Provide command untuk user hapus data mereka sendiri

## 4. Dependency Security

- Pin semua dependency versions di `go.sum`
- Review changelog sebelum upgrade dependency
- Run `go mod verify` di CI pipeline
- Update dependensi secara berkala (minimal bulanan)

## 5. LLM Security

- System prompt hardening terhadap prompt injection
- Output filtering untuk konten berbahaya
- Rate limiting per user untuk mencegah abuse
- Token budget per request dan per user

## 6. Referensi

- Lihat [14_SECURITY_THREAT_MODEL.md](../14_SECURITY_THREAT_MODEL.md)
- Lihat [10_RISK_REGISTER.md](../10_RISK_REGISTER.md)

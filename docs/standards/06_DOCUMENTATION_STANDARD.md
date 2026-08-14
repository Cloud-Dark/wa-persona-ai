# Documentation Standard — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. File Naming

- Dokumen inti: `NN_NAMA_UPPERCASE.md` (nomor dua digit, zero-padded)
- Standards: `NN_TOPIK_STANDARD.md` di `docs/standards/`
- Plans: `YYYY-MM-DD-slug-fitur.md` di `docs/plans/`
- Specs: `YYYY-MM-DD-slug-design.md` di `docs/specs/`

## 2. Header Template

Setiap dokumen diawali:

```markdown
# Judul Dokumen — WA Persona AI

> Status: Draft | Review | Final
> Terakhir diperbarui: YYYY-MM-DD
> Pemilik: Nama
```

## 3. Bahasa

- Dokumentasi: Bahasa Indonesia (default proyek ini)
- Code comments: Bahasa Inggris
- Commit messages: Bahasa Inggris
- Variable/function names: Bahasa Inggris

## 4. Cross-Reference

- Gunakan relative path: `Lihat [04_TRD.md](04_TRD.md)`
- Dari subfolder: `Lihat [05_ARCHITECTURE.md](../05_ARCHITECTURE.md)`

## 5. Diagram

- Gunakan Mermaid untuk diagram (rendered di GitHub)
- Sertakan dalam code block ` ```mermaid `

## 6. Update Rules

- Selalu update `31_CHANGELOG.md` saat ada perubahan
- Selalu update `07_MASTER_CHECKLIST.md` saat task selesai
- Tanggal selalu absolut (YYYY-MM-DD), BUKAN relatif

## 7. Referensi

- Lihat [00_STANDARD_INDEX.md](00_STANDARD_INDEX.md)

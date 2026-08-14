# Testing Standard — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Coverage Targets

| Module | Minimum Coverage |
|---|---|
| `internal/persona/` | 80% |
| `internal/memory/` | 80% |
| `internal/llm/` | 80% |
| `internal/context/` | 80% |
| `internal/admin/` | 80% |
| `internal/whatsapp/` | 60% (external dependency) |
| Overall | 75% |

## 2. Test Types

### Unit Tests
- Wajib untuk setiap exported function
- Table-driven tests preferred
- Mock external dependencies
- File: `*_test.go` di package yang sama

### Integration Tests
- Tag: `//go:build integration`
- Test interaksi antar module
- Boleh menggunakan real database (SQLite in-memory)

### Benchmark Tests
- Untuk critical path (message processing, memory search)
- File: `*_bench_test.go`
- Run: `go test -bench=. ./internal/...`

## 3. Naming Convention

```go
func TestFunctionName_Scenario_Expected(t *testing.T) {}

// Contoh:
func TestPersonaLoad_ValidFile_ReturnsPersona(t *testing.T) {}
func TestPersonaLoad_MissingFile_ReturnsError(t *testing.T) {}
func TestMemorySearch_NoResults_ReturnsEmptySlice(t *testing.T) {}
```

## 4. Mocking

- Gunakan interface untuk dependency injection
- Generate mock dengan `mockgen` atau hand-written
- Jangan mock yang tidak perlu (prefer real implementation untuk unit-unit kecil)

## 5. CI Integration

```yaml
# Setiap PR harus:
- go test ./... -race                    # Unit tests + race detection
- go test ./... -coverprofile=cover.out  # Coverage report
- golangci-lint run                      # Linting
```

## 6. Referensi

- Lihat [12_TEST_STRATEGY.md](../12_TEST_STRATEGY.md) untuk strategy lengkap
- Lihat [02_CODE_STANDARD.md](02_CODE_STANDARD.md) untuk code style

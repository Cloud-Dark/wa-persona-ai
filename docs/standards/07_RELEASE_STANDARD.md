# Release Standard — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Versioning

Mengikuti [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH

MAJOR — breaking changes (config format, API interface)
MINOR — new features, backward compatible
PATCH — bug fixes, backward compatible
```

## 2. Pre-Release

```
v0.x.x   — Development phase, API belum stabil
v1.0.0   — First stable release
```

## 3. Release Checklist

- [ ] All tests pass (`go test ./...`)
- [ ] Linting clean (`golangci-lint run`)
- [ ] `31_CHANGELOG.md` updated
- [ ] `08_ROADMAP.md` updated
- [ ] README.md up-to-date
- [ ] Git tag created (`git tag vX.Y.Z`)
- [ ] GitHub Release created with changelog
- [ ] Binary artifacts uploaded (linux/macOS/windows)
- [ ] Docker image pushed

## 4. Release Process

```bash
# 1. Update changelog
# 2. Commit
git commit -m "chore: prepare release vX.Y.Z"

# 3. Tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# 4. Push
git push origin main --tags

# 5. GitHub Release (auto via CI or manual)
gh release create vX.Y.Z --title "vX.Y.Z" --notes-file RELEASE_NOTES.md
```

## 5. Referensi

- Lihat [08_ROADMAP.md](../08_ROADMAP.md) untuk milestone planning
- Lihat [31_CHANGELOG.md](../31_CHANGELOG.md) untuk changelog

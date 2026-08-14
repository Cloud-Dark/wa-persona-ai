# Contribution Standard — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Getting Started

1. Fork repository di GitHub
2. Clone fork ke lokal
3. Buat branch baru dari `main`
4. Lakukan perubahan
5. Submit Pull Request

## 2. Branch Naming

```
<type>/<short-description>

Contoh:
  feat/multi-persona-support
  fix/memory-leak-vector-store
  docs/update-installation-guide
  refactor/llm-provider-interface
```

| Type | Deskripsi |
|---|---|
| `feat` | Fitur baru |
| `fix` | Bug fix |
| `docs` | Dokumentasi |
| `refactor` | Refactoring tanpa perubahan behavior |
| `test` | Menambah atau memperbaiki test |
| `chore` | Maintenance (CI, deps, config) |

## 3. Commit Message

Format: [Conventional Commits](https://www.conventionalcommits.org/)

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Contoh:
```
feat(persona): add hot-reload support for persona files

Persona files are now watched for changes using fsnotify.
When a YAML file in persona/ is modified, the persona
manager automatically reloads the affected persona.

Closes #12
```

## 4. Pull Request

### PR Template

```markdown
## Summary
Brief description of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Refactoring
- [ ] Documentation
- [ ] Other

## Checklist
- [ ] Code follows project standards
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

### Review Process

1. Minimal 1 approval required
2. All CI checks must pass
3. No merge conflicts
4. Squash merge preferred

## 5. Code of Conduct

- Be respectful and constructive
- Focus on the code, not the person
- Provide helpful feedback
- Welcome newcomers

## 6. Referensi

- Lihat [02_CODE_STANDARD.md](02_CODE_STANDARD.md) untuk coding conventions
- Lihat [05_TESTING_STANDARD.md](05_TESTING_STANDARD.md) untuk testing requirements

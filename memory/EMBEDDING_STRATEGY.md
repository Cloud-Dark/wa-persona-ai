# Memory Embedding Strategy — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. Embedding Model Comparison

| Model | Dimensions | Speed | Quality | Cost | Recommended |
|---|---|---|---|---|---|
| OpenAI text-embedding-ada-002 | 1536 | Fast | Good | $0.0001/1K tokens | ✅ Default |
| OpenAI text-embedding-3-small | 1536 | Fast | Better | $0.00002/1K tokens | ✅ Budget |
| OpenAI text-embedding-3-large | 3072 | Medium | Best | $0.00013/1K tokens | For high-accuracy |
| Local (all-MiniLM-L6-v2) | 384 | Fast (local) | Decent | Free | Offline/privacy |

## 2. Embedding Strategy

### What Gets Embedded

```
Input:  "Nama saya Budi, saya kerja di bank BCA sebagai teller"

Extracted Facts:
  1. "User's name is Budi" → embed + store
  2. "User works at bank BCA as a teller" → embed + store

NOT embedded:
  - "Halo" (no extractable fact)
  - "Ok" (no extractable fact)
  - "Makasih ya" (no extractable fact)
```

### Embedding Format

Sebelum embedding, text diformat untuk meningkatkan retrieval quality:

```
Template: "[{topic}] {content} (context: {user_name}, {date})"
Example:  "[personal_info] User's name is Budi (context: new_user, 2026-08-14)"
```

## 3. Similarity Search

### Cosine Similarity

```
similarity(A, B) = (A · B) / (||A|| × ||B||)

Range: -1 to 1 (normalized vectors: 0 to 1)
Threshold: 0.7 (default minimum)
```

### Search Flow

```
Query: "Siapa nama saya?"
  ↓
Encode: query → vector [1536 dims]
  ↓
Search: cosine_similarity(query_vector, all_memory_vectors)
  ↓
Filter: user_jid == sender_jid AND similarity >= 0.7
  ↓
Rank: (similarity * 0.6) + (recency * 0.2) + (importance * 0.1) + (frequency * 0.1)
  ↓
Return: Top-5 results
  ↓
Result: "[personal_info] User's name is Budi" (score: 0.92)
```

## 4. Optimization Strategies

### Batch Embedding
- Kumpulkan facts, embed dalam batch (max 10 per batch)
- Mengurangi API calls dan biaya

### Caching
- Cache embedding results untuk query yang sama
- LRU cache dengan TTL 1 jam

### Dimensionality Reduction (Future)
- PCA untuk mengurangi dimensi dari 1536 → 512
- Trade-off: sedikit loss quality, significant speed up

## 5. Referensi

- Lihat [memory/README.md](README.md) untuk overview memory system
- Lihat [../docs/20_DATABASE.md](../docs/20_DATABASE.md) untuk storage schema

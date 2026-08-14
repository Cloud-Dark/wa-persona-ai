package memory

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/philippgille/chromem-go"
	"github.com/rs/zerolog/log"
)

// Store manages both vector memory and conversation history.
type Store struct {
	db         *sql.DB
	collection *chromem.Collection
	chromaDB   *chromem.DB
	embedFn    chromem.EmbeddingFunc
	topK       int
	minSim     float64
}

// NewStore creates and initializes the memory store.
func NewStore(vectorDir, metadataDB string, topK int, minSim float64, openAIKey string) (*Store, error) {
	if err := os.MkdirAll(vectorDir, 0755); err != nil {
		return nil, fmt.Errorf("create vector dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(metadataDB), 0755); err != nil {
		return nil, fmt.Errorf("create metadata dir: %w", err)
	}

	db, err := sql.Open("sqlite3", metadataDB+"?_journal=WAL&_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	// Setup embedding function
	var embedFn chromem.EmbeddingFunc
	if openAIKey != "" {
		embedFn = chromem.NewEmbeddingFuncOpenAI(openAIKey, chromem.EmbeddingModelOpenAI3Small)
	} else {
		embedFn = chromem.NewEmbeddingFuncOllama("nomic-embed-text", "")
	}

	chromaDB := chromem.NewDB()
	col, err := chromaDB.GetOrCreateCollection("memories", nil, embedFn)
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}

	s := &Store{
		db:         db,
		collection: col,
		chromaDB:   chromaDB,
		embedFn:    embedFn,
		topK:       topK,
		minSim:     minSim,
	}

	return s, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS memories (
		id TEXT PRIMARY KEY,
		user_jid TEXT NOT NULL,
		content TEXT NOT NULL,
		topic TEXT DEFAULT 'general',
		importance_score REAL DEFAULT 0.5,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_accessed DATETIME DEFAULT CURRENT_TIMESTAMP,
		access_count INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_memories_user ON memories(user_jid);
	CREATE INDEX IF NOT EXISTS idx_memories_importance ON memories(importance_score);

	CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY,
		user_jid TEXT NOT NULL,
		persona_name TEXT DEFAULT 'default',
		started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_message_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		status TEXT DEFAULT 'active'
	);
	CREATE INDEX IF NOT EXISTS idx_conv_user ON conversations(user_jid, status);

	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		conversation_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (conversation_id) REFERENCES conversations(id)
	);
	CREATE INDEX IF NOT EXISTS idx_msg_conv ON messages(conversation_id, created_at);
	`)
	return err
}

// SaveMemory stores a memory entry in both vector DB and SQLite.
func (s *Store) SaveMemory(ctx context.Context, entry *Entry) error {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	// Store in vector DB
	doc := chromem.Document{
		ID:      entry.ID,
		Content: entry.Content,
		Metadata: map[string]string{
			"user_jid":  entry.UserJID,
			"topic":     entry.Topic,
			"importance": fmt.Sprintf("%.2f", entry.ImportanceScore),
		},
	}
	if err := s.collection.AddDocument(ctx, doc); err != nil {
		log.Warn().Err(err).Str("id", entry.ID).Msg("failed to add to vector store")
	}

	// Store metadata in SQLite
	_, err := s.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO memories (id, user_jid, content, topic, importance_score, created_at, last_accessed, access_count)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.UserJID, entry.Content, entry.Topic,
		entry.ImportanceScore, entry.CreatedAt, entry.LastAccessed, entry.AccessCount,
	)
	return err
}

// SearchMemories retrieves relevant memories for a user based on query.
func (s *Store) SearchMemories(ctx context.Context, query, userJID string) ([]*SearchResult, error) {
	if s.collection.Count() == 0 {
		return nil, nil
	}

	results, err := s.collection.Query(ctx, query, s.topK, map[string]string{"user_jid": userJID}, nil)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}

	var out []*SearchResult
	for _, r := range results {
		sim := float64(r.Similarity)
		if sim < s.minSim {
			continue
		}
		out = append(out, &SearchResult{
			Entry: Entry{
				ID:      r.ID,
				Content: r.Content,
				UserJID: userJID,
				Topic:   r.Metadata["topic"],
			},
			Similarity: sim,
		})
	}

	return out, nil
}

// DeleteUserMemories removes all memories for a user.
func (s *Store) DeleteUserMemories(ctx context.Context, userJID string) error {
	// Get IDs from SQLite
	rows, err := s.db.QueryContext(ctx, "SELECT id FROM memories WHERE user_jid = ?", userJID)
	if err != nil {
		return err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	rows.Close()

	// Delete from vector DB
	for _, id := range ids {
		_ = s.collection.Delete(ctx, nil, nil, id)
	}

	// Delete from SQLite
	_, err = s.db.ExecContext(ctx, "DELETE FROM memories WHERE user_jid = ?", userJID)
	return err
}

// MemoryStats returns statistics about stored memories.
func (s *Store) MemoryStats(ctx context.Context) (total int, users int, err error) {
	s.db.QueryRowContext(ctx, "SELECT COUNT(*), COUNT(DISTINCT user_jid) FROM memories").Scan(&total, &users)
	return
}

// GetOrCreateConversation returns or creates an active conversation for a user.
func (s *Store) GetOrCreateConversation(ctx context.Context, userJID, personaName string) (string, error) {
	var convID string
	err := s.db.QueryRowContext(ctx,
		"SELECT id FROM conversations WHERE user_jid = ? AND status = 'active' ORDER BY last_message_at DESC LIMIT 1",
		userJID,
	).Scan(&convID)

	if err == sql.ErrNoRows {
		convID = uuid.New().String()
		_, err = s.db.ExecContext(ctx,
			"INSERT INTO conversations (id, user_jid, persona_name) VALUES (?, ?, ?)",
			convID, userJID, personaName,
		)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return convID, nil
}

// SaveMessage stores a conversation message.
func (s *Store) SaveMessage(ctx context.Context, convID, role, content string) error {
	id := uuid.New().String()
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO messages (id, conversation_id, role, content) VALUES (?, ?, ?, ?)",
		id, convID, role, content,
	)
	if err != nil {
		return err
	}

	// Update conversation last_message_at
	_, _ = s.db.ExecContext(ctx,
		"UPDATE conversations SET last_message_at = CURRENT_TIMESTAMP WHERE id = ?",
		convID,
	)
	return nil
}

// GetHistory returns recent messages for a conversation.
func (s *Store) GetHistory(ctx context.Context, convID string, limit int) ([]*ConversationMessage, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, conversation_id, role, content, created_at FROM messages
		 WHERE conversation_id = ? ORDER BY created_at DESC LIMIT ?`,
		convID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*ConversationMessage
	for rows.Next() {
		m := &ConversationMessage{}
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			continue
		}
		msgs = append(msgs, m)
	}

	// Reverse to chronological order
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

// ResetConversation archives the current conversation for a user.
func (s *Store) ResetConversation(ctx context.Context, userJID string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE conversations SET status = 'reset' WHERE user_jid = ? AND status = 'active'",
		userJID,
	)
	return err
}

// ExtractAndSaveMemories analyzes a conversation turn and saves important facts.
func (s *Store) ExtractAndSaveMemories(ctx context.Context, userJID, userMsg, botMsg string, minImportance float64) {
	facts := extractFacts(userMsg)
	for _, fact := range facts {
		if fact.score < minImportance {
			continue
		}
		entry := &Entry{
			UserJID:         userJID,
			Content:         fact.content,
			Topic:           fact.topic,
			ImportanceScore: fact.score,
		}
		if err := s.SaveMemory(ctx, entry); err != nil {
			log.Warn().Err(err).Str("user", userJID).Msg("failed to save memory")
		}
	}
}

type fact struct {
	content string
	topic   string
	score   float64
}

// extractFacts applies simple heuristics to find memorable facts in a message.
func extractFacts(text string) []fact {
	var facts []fact
	lower := strings.ToLower(text)

	// Name introduction
	namePatterns := []string{"nama saya", "namaku", "nama aku", "my name is", "i'm ", "i am "}
	for _, p := range namePatterns {
		if idx := strings.Index(lower, p); idx >= 0 {
			rest := text[idx+len(p):]
			if len(rest) > 2 {
				words := strings.Fields(rest)
				if len(words) > 0 {
					facts = append(facts, fact{
						content: fmt.Sprintf("User's name is %s", words[0]),
						topic:   "personal_info",
						score:   0.9,
					})
				}
			}
		}
	}

	// Work/job
	workPatterns := []string{"kerja di", "bekerja di", "work at", "work for", "karyawan", "pegawai"}
	for _, p := range workPatterns {
		if strings.Contains(lower, p) {
			facts = append(facts, fact{
				content: text,
				topic:   "personal_info",
				score:   0.8,
			})
			break
		}
	}

	// Location
	locPatterns := []string{"tinggal di", "live in", "stay in", "pindah ke", "moved to"}
	for _, p := range locPatterns {
		if strings.Contains(lower, p) {
			facts = append(facts, fact{
				content: text,
				topic:   "personal_info",
				score:   0.7,
			})
			break
		}
	}

	// Preferences
	prefPatterns := []string{"suka", "favorit", "like", "prefer", "love", "hate", "tidak suka"}
	for _, p := range prefPatterns {
		if strings.Contains(lower, p) {
			facts = append(facts, fact{
				content: text,
				topic:   "preference",
				score:   0.6,
			})
			break
		}
	}

	return facts
}

// Close releases resources.
func (s *Store) Close() error {
	return s.db.Close()
}

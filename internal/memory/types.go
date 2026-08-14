package memory

import "time"

// Entry represents a single memory record.
type Entry struct {
	ID              string    `json:"id"`
	UserJID         string    `json:"user_jid"`
	Content         string    `json:"content"`
	Topic           string    `json:"topic"`
	ImportanceScore float64   `json:"importance_score"`
	CreatedAt       time.Time `json:"created_at"`
	LastAccessed    time.Time `json:"last_accessed"`
	AccessCount     int       `json:"access_count"`
}

// SearchResult is a memory entry with a similarity score.
type SearchResult struct {
	Entry
	Similarity float64 `json:"similarity"`
}

// ConversationMessage is a single turn in the conversation history.
type ConversationMessage struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

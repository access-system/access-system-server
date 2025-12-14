package domain

import (
	"github.com/pgvector/pgvector-go"
)

// Embedding represents a vector embedding with an ID, Name, and Vector.
type Embedding struct {
	ID       int64           `json:"id"`
	Name     string          `json:"name"`
	Vector   pgvector.Vector `json:"vector"`
	Accuracy float32         `json:"accuracy,omitempty"`
}

package domain

import (
	"github.com/pgvector/pgvector-go"
)

// Embedding represents a vector embedding with an ID, Name, and Vector.
type Embedding struct {
	ID       int64
	Name     string
	Vector   pgvector.Vector
	Accuracy float32
}

package repository

import (
	"context"
	"database/sql"
	"errors"

	"access-system-api/internal/domain"

	"github.com/pgvector/pgvector-go"
)

// EmbeddingRepository defines the methods for managing embeddings in the database.
type EmbeddingRepository interface {
	CreateEmbedding(ctx context.Context, embedding *domain.Embedding) error
	GetEmbeddingById(ctx context.Context, id int64) (*domain.Embedding, error)
	ListEmbeddings(ctx context.Context) ([]*domain.Embedding, error)
	GetSimilarEmbeddingByVector(ctx context.Context, vector pgvector.Vector) (*domain.Embedding, error)
	UpdateEmbedding(ctx context.Context, embedding *domain.Embedding) error
	DeleteEmbeddingById(ctx context.Context, id int64) error
}

// embeddingRepository implements EmbeddingRepository.
type embeddingRepository struct {
	db *sql.DB
}

// NewEmbeddingsRepository creates a new instance of embeddingRepository.
func NewEmbeddingsRepository(db *sql.DB) EmbeddingRepository {
	return &embeddingRepository{db: db}
}

// CreateEmbedding inserts a new embedding into the database.
func (r *embeddingRepository) CreateEmbedding(ctx context.Context, embedding *domain.Embedding) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	const query = "INSERT INTO embedding (name, vector_) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, embedding.Name, embedding.Vector)
	if err != nil {
		return err
	}

	return nil
}

func (r *embeddingRepository) GetEmbeddingById(ctx context.Context, id int64) (*domain.Embedding, error) {
	if err := r.db.Ping(); err != nil {
		return nil, err
	}

	const query = "SELECT id, name, vector_ FROM embedding WHERE id = $1"
	embedding := &domain.Embedding{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&embedding.ID, &embedding.Name, &embedding.Vector)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return embedding, nil
}

func (r *embeddingRepository) ListEmbeddings(ctx context.Context) ([]*domain.Embedding, error) {
	if err := r.db.Ping(); err != nil {
		return nil, err
	}

	const query = "SELECT id, name, vector_ FROM embedding"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var embeddings []*domain.Embedding
	for rows.Next() {
		embedding := &domain.Embedding{}
		if err := rows.Scan(&embedding.ID, &embedding.Name, &embedding.Vector); err != nil {
			return nil, err
		}
		embeddings = append(embeddings, embedding)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return embeddings, nil
}

// GetSimilarEmbeddingByVector retrieves the most similar embedding from the database based on the provided vector.
func (r *embeddingRepository) GetSimilarEmbeddingByVector(ctx context.Context, vector pgvector.Vector) (*domain.Embedding, error) {
	if err := r.db.Ping(); err != nil {
		return nil, err
	}

	const query = "SELECT id, name, vector_, (1 - (vector_ <=> $1)) AS accuracy FROM embedding WHERE (1 - (vector_ <=> $1)) > 0.58 ORDER BY (vector_ <=> $1) ASC LIMIT 1;"
	embedding := &domain.Embedding{}
	err := r.db.QueryRowContext(ctx, query, vector).Scan(&embedding.ID, &embedding.Name, &embedding.Vector, &embedding.Accuracy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return embedding, nil
}

func (r *embeddingRepository) UpdateEmbedding(ctx context.Context, embedding *domain.Embedding) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	const query = "UPDATE embedding SET name = $1, vector_ = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, embedding.Name, embedding.Vector, embedding.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteEmbeddingById removes an embedding from the database by its ID.
func (r *embeddingRepository) DeleteEmbeddingById(ctx context.Context, id int64) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	const query = "DELETE FROM embedding WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"context"
	"fmt"

	"access-system-api/internal/domain"
	"access-system-api/internal/repository"

	"github.com/pgvector/pgvector-go"
)

//go:generate mockgen -destination=../mocks/service/embedding_mock.go -package=mocks . EmbeddingService

// EmbeddingService defines the interface for managing embeddings.
type EmbeddingService interface {
	AddEmbedding(ctx context.Context, name string, vector []float32) error
	ValidateEmbedding(ctx context.Context, vector []float32) (*domain.Embedding, error)
	DeleteEmbedding(ctx context.Context, id int64) error
}

// embeddingService is the concrete implementation of EmbeddingService.
type embeddingService struct {
	embeddingRepo repository.EmbeddingRepository
}

// NewEmbeddingService creates a new instance of EmbeddingService.
func NewEmbeddingService(embeddingRepo repository.EmbeddingRepository) EmbeddingService {
	return &embeddingService{embeddingRepo: embeddingRepo}
}

// AddEmbedding adds a new embedding to the repository.
func (s *embeddingService) AddEmbedding(ctx context.Context, name string, vector []float32) error {
	if len(vector) != 512 {
		return fmt.Errorf("vector size must be 512, got %d", len(vector))
	}
	embedding := &domain.Embedding{
		Name:   name,
		Vector: pgvector.NewVector(vector),
	}
	return s.embeddingRepo.CreateEmbedding(ctx, embedding)
}

// ValidateEmbedding checks if a similar embedding exists in the repository.
func (s *embeddingService) ValidateEmbedding(ctx context.Context, vector []float32) (*domain.Embedding, error) {
	if len(vector) != 512 {
		return nil, fmt.Errorf("vector size must be 512, got %d", len(vector))
	}
	embedding, err := s.embeddingRepo.GetSimilarEmbeddingByVector(ctx, pgvector.NewVector(vector))
	if err != nil {
		return nil, err
	}
	return embedding, nil
}

// DeleteEmbedding removes an embedding from the repository by its ID.
func (s *embeddingService) DeleteEmbedding(ctx context.Context, id int64) error {
	return s.embeddingRepo.DeleteEmbeddingById(ctx, id)
}

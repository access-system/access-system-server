package service

import (
	"context"
	"testing"

	"access-system-api/internal/domain"
	"access-system-api/internal/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
)

func TestEmbeddingService_AddEmbedding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	name := "test"
	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}

	embedding := &domain.Embedding{
		Name:   name,
		Vector: pgvector.NewVector(vector),
	}

	repo.EXPECT().CreateEmbedding(ctx, embedding).Return(nil)

	err := service.AddEmbedding(ctx, name, vector)
	assert.NoError(t, err)
}

func TestEmbeddingService_AddEmbedding_InvalidVectorSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	name := "test"
	vector := make([]float32, 100) // Invalid size

	err := service.AddEmbedding(ctx, name, vector)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vector size must be 512")
}

func TestEmbeddingService_ValidateEmbedding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}

	repo.EXPECT().GetSimilarEmbeddingByVector(ctx, pgvector.NewVector(vector)).Return(&domain.Embedding{}, nil)

	_, err := service.ValidateEmbedding(ctx, vector)
	assert.NoError(t, err)
}

func TestEmbeddingService_ValidateEmbedding_InvalidVectorSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	vector := make([]float32, 100) // Invalid size

	_, err := service.ValidateEmbedding(ctx, vector)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vector size must be 512")
}

func TestEmbeddingService_DeleteEmbedding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	id := int64(123)

	repo.EXPECT().DeleteEmbeddingById(ctx, id).Return(nil)

	err := service.DeleteEmbedding(ctx, id)
	assert.NoError(t, err)
}

func TestEmbeddingService_GetEmbedding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	id := int64(123)

	repo.EXPECT().GetEmbeddingById(ctx, id).Return(&domain.Embedding{}, nil)

	embedding, err := service.GetEmbedding(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, embedding)
}

func TestEmbeddingService_ListEmbeddings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()

	repo.EXPECT().ListEmbeddings(ctx).Return([]*domain.Embedding{}, nil)

	embeddings, err := service.ListEmbeddings(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, embeddings)
}

func TestEmbeddingService_UpdateEmbedding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	id := int64(123)
	name := "updated"
	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}

	embedding := &domain.Embedding{
		ID:     id,
		Name:   name,
		Vector: pgvector.NewVector(vector),
	}

	repo.EXPECT().UpdateEmbedding(ctx, embedding).Return(nil)

	err := service.UpdateEmbedding(ctx, id, name, vector)
	assert.NoError(t, err)
}

func TestEmbeddingService_UpdateEmbedding_InvalidVectorSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	id := int64(123)
	name := "updated"
	vector := make([]float32, 100) // Invalid size

	err := service.UpdateEmbedding(ctx, id, name, vector)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vector size must be 512")
}

func TestEmbeddingService_ValidateEmbedding_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockEmbeddingRepository(ctrl)
	service := NewEmbeddingService(repo)

	ctx := context.Background()
	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}

	repo.EXPECT().GetSimilarEmbeddingByVector(ctx, pgvector.NewVector(vector)).Return(nil, assert.AnError)

	emb, err := service.ValidateEmbedding(ctx, vector)
	assert.Error(t, err)
	assert.Nil(t, emb)
}

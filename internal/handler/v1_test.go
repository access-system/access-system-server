package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"access-system-api/internal/domain"
	mocks "access-system-api/internal/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pgvector/pgvector-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler V1Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/add", handler.AddEmbeddingHandler)
	r.POST("/validate", handler.ValidateEmbeddingHandler)
	r.POST("/delete", handler.DeleteEmbeddingHandler)
	return r
}

func TestAddEmbeddingHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}
	body, _ := json.Marshal(map[string]interface{}{
		"name":   "test",
		"vector": vector,
	})

	service.EXPECT().AddEmbedding(gomock.Any(), "test", vector).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddEmbeddingHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	// Missing 'name' field
	body, _ := json.Marshal(map[string]interface{}{
		"vector": make([]float32, 512),
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddEmbeddingHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}
	body, _ := json.Marshal(map[string]interface{}{
		"name":   "test",
		"vector": vector,
	})
	service.EXPECT().AddEmbedding(gomock.Any(), "test", vector).Return(assert.AnError)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddEmbeddingHandler_InvalidVectorSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	vector := make([]float32, 100) // Invalid size
	body, _ := json.Marshal(map[string]interface{}{
		"name":   "test",
		"vector": vector,
	})
	service.EXPECT().AddEmbedding(gomock.Any(), "test", vector).Return(assert.AnError)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestValidateEmbeddingHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}
	body, _ := json.Marshal(map[string]interface{}{
		"vector": vector,
	})

	// Return a valid embedding and no error
	service.EXPECT().ValidateEmbedding(gomock.Any(), vector).Return(&domain.Embedding{
		ID:       1,
		Name:     "test",
		Vector:   pgvector.NewVector(vector),
		Accuracy: 0.99,
	}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateEmbeddingHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	// Missing 'vector' field
	body, _ := json.Marshal(map[string]interface{}{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateEmbeddingHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}
	body, _ := json.Marshal(map[string]interface{}{
		"vector": vector,
	})
	// Return error
	service.EXPECT().ValidateEmbedding(gomock.Any(), vector).Return(nil, assert.AnError)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestValidateEmbeddingHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	vector := make([]float32, 512)
	for i := range vector {
		vector[i] = float32(i)
	}
	body, _ := json.Marshal(map[string]interface{}{
		"vector": vector,
	})
	// Return not found error
	service.EXPECT().ValidateEmbedding(gomock.Any(), vector).Return(nil, sql.ErrNoRows)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestValidateEmbeddingHandler_InvalidVectorSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	vector := make([]float32, 100) // Invalid size
	body, _ := json.Marshal(map[string]interface{}{
		"vector": vector,
	})
	// Return size validation error
	service.EXPECT().ValidateEmbedding(gomock.Any(), vector).Return(nil, assert.AnError)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteEmbeddingHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	body, _ := json.Marshal(map[string]interface{}{
		"id": int64(123),
	})

	service.EXPECT().DeleteEmbedding(gomock.Any(), int64(123)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteEmbeddingHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	body, _ := json.Marshal(map[string]interface{}{}) // Missing 'id'
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteEmbeddingHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockEmbeddingService(ctrl)
	log := logrus.New()
	handler := NewV1Handler(service, log)
	r := setupRouter(handler)

	body, _ := json.Marshal(map[string]interface{}{
		"id": int64(123),
	})
	service.EXPECT().DeleteEmbedding(gomock.Any(), int64(123)).Return(assert.AnError)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

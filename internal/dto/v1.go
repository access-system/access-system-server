package dto

type AddEmbeddingRequest struct {
	Name   string    `json:"name" encrypt:"name"`
	Vector []float32 `json:"vector" encrypt:"vector"`
}

type ValidateEmbeddingRequest struct {
	Vector []float32 `json:"vector" encrypt:"vector"`
}

type ValidateEmbeddingResponse struct {
	ID       int64     `json:"id" encrypt:"id"`
	Name     string    `json:"name" encrypt:"name"`
	Vector   []float32 `json:"vector" encrypt:"vector"`
	Accuracy float32   `json:"accuracy" encrypt:"accuracy"`
}

type DeleteEmbeddingRequest struct {
	ID int64 `json:"id" encrypt:"id"`
}

package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"access-system-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AdminHandler interface {
	AddEmbeddingHandler(c *gin.Context)
	GetEmbeddingHandler(c *gin.Context)
	ListEmbeddingsHandler(c *gin.Context)
	UpdateEmbeddingHandler(c *gin.Context)
	DeleteEmbeddingHandler(c *gin.Context)
}

type adminHandler struct {
	embeddingService service.EmbeddingService
	log              *logrus.Logger
}

func NewAdminHandler(embeddingService service.EmbeddingService, log *logrus.Logger) AdminHandler {
	return &adminHandler{
		embeddingService: embeddingService,
		log:              log,
	}
}

func (h *adminHandler) AddEmbeddingHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var data struct {
		Name   string    `json:"name" binding:"required"`
		Vector []float32 `json:"vector" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		h.log.Errorln("Error binding JSON:", err)
		c.String(http.StatusBadRequest, "Bad Request: %v", err)
		return
	}

	err := h.embeddingService.AddEmbedding(ctx, data.Name, data.Vector)
	if err != nil {
		h.log.Errorln("Error adding embedding:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error: %v", err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *adminHandler) GetEmbeddingHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		h.log.Errorln("Invalid ID parameter:", err)
		c.String(http.StatusBadRequest, "Bad Request: invalid ID parameter")
		return
	}

	embedding, err := h.embeddingService.GetEmbedding(ctx, intId)
	if err != nil {
		h.log.Errorln("Error getting embedding:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error: %v", err)
		return
	}

	c.JSON(http.StatusOK, embedding)
}

func (h *adminHandler) ListEmbeddingsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	embeddings, err := h.embeddingService.ListEmbeddings(ctx)
	if err != nil {
		h.log.Errorln("Error listing embeddings:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error: %v", err)
		return
	}

	var response []gin.H
	for _, embedding := range embeddings {
		response = append(response, gin.H{
			"id":     embedding.ID,
			"name":   embedding.Name,
			"vector": embedding.Vector,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *adminHandler) UpdateEmbeddingHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var data struct {
		ID     int64     `json:"id" binding:"required"`
		Name   string    `json:"name" binding:"required"`
		Vector []float32 `json:"vector" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		h.log.Errorln("Error binding JSON:", err)
		c.String(http.StatusBadRequest, "Bad Request: %v", err)
		return
	}

	err := h.embeddingService.UpdateEmbedding(ctx, data.ID, data.Name, data.Vector)
	if err != nil {
		h.log.Errorln("Error updating embedding:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error: %v", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *adminHandler) DeleteEmbeddingHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var data struct {
		ID int64 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		h.log.Errorln("Error binding JSON:", err)
		c.String(http.StatusBadRequest, "Bad Request: %v", err)
		return
	}

	err := h.embeddingService.DeleteEmbedding(ctx, data.ID)
	if err != nil {
		h.log.Errorln("Error deleting embedding:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error: %v", err)
		return
	}

	c.Status(http.StatusOK)
}

package router

import (
	"net/http"

	"access-system-api/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Router struct to hold the Gin engine and handlers
type Router struct {
	engine *gin.Engine
	v1     handler.V1Handler
	admin  handler.AdminHandler
	log    *logrus.Logger
}

// NewRouter initializes a new Router instance
func NewRouter(v1 handler.V1Handler, admin handler.AdminHandler, log *logrus.Logger) *Router {
	return &Router{
		engine: gin.Default(),
		v1:     v1,
		admin:  admin,
		log:    log,
	}
}

// Run starts the Gin server and sets up the routes
func (r *Router) Run() {
	v1 := r.engine.Group("/api/v1")
	{
		v1.POST("/embedding", r.v1.AddEmbeddingHandler)
		v1.POST("/embedding/validate", r.v1.ValidateEmbeddingHandler)
		v1.DELETE("/embedding", r.v1.DeleteEmbeddingHandler)
	}

	admin := v1.Group("/admin")
	{
		admin.POST("/embedding", r.admin.AddEmbeddingHandler)
		admin.GET("/embedding/:id", r.admin.GetEmbeddingHandler)
		admin.GET("/embeddings", r.admin.ListEmbeddingsHandler)
		admin.PUT("/embedding", r.admin.UpdateEmbeddingHandler)
		admin.DELETE("/embedding", r.admin.DeleteEmbeddingHandler)
	}

	r.engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	gin.SetMode(gin.ReleaseMode)
	if err := r.engine.Run(":8081"); err != nil {
		r.log.Fatalf("Failed to start server: %v", err)
	}
}

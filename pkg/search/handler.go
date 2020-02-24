package search

import (
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ivan-marquez/es-mdb/pkg/middleware/logger"
)

// Handler type with middleware props
type Handler struct {
	*logger.Middleware
	*elasticsearch.Client
}

// Search handler sends query to ElasticSearch and returns
// search term results
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// TODO: search term on ElasticSearch and write response
	w.Write([]byte("Initial response"))
}

// SetupRoutes registers handles for API endpoints
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", h.Logger(h.Search))
}

// NewHandlers "constructor" function to inject dependencies
func NewHandlers(lgr *log.Logger, es *elasticsearch.Client) *Handler {
	l := logger.NewLogger(lgr)

	return &Handler{l, es}
}

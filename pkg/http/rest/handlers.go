package rest

import (
	"net/http"

	"github.com/ivan-marquez/es-mdb/pkg/domain"
)

// Handler for HTTP
type Handler struct {
	mw *Middleware
	us domain.UserService
}

// Search handles a request to retrieve users by term
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {

}

// SetupRoutes registers handles for API endpoints
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", h.mw.LoggerMiddleware(h.Search))
}

// NewHandler returns a new HTTP handler
func NewHandler(mw *Middleware, us domain.UserService) *Handler {
	return &Handler{mw, us}
}

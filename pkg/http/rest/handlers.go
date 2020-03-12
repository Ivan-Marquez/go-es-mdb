package rest

import (
	"encoding/json"
	"fmt"
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	keys, ok := r.URL.Query()["term"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request. Missing term parameter")
		return
	}

	term := keys[0]

	switch r.Method {
	case http.MethodGet:
		users, err := h.us.GetByTerm(term)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(&users)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request")
	}
}

// SetupRoutes registers handles for API endpoints
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", h.mw.LoggerMiddleware(h.Search))
}

// NewHandler returns a new HTTP handler
func NewHandler(mw *Middleware, us domain.UserService) *Handler {
	return &Handler{mw, us}
}

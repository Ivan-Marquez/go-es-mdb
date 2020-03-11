package rest

import (
	"log"
	"net/http"
	"time"
)

// Middleware for HTTP handlers
type Middleware struct {
	logger *log.Logger
}

// LoggerMiddleware returns an HTTP handler with logging middleware
func (l *Middleware) LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer l.logger.Printf("request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

// NewMiddleware returns a logger middleware with logger implementation
func NewMiddleware(logger *log.Logger) *Middleware {
	return &Middleware{logger}
}

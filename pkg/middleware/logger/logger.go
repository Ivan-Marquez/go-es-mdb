package logger

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	logger *log.Logger
}

func (l *Middleware) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer l.logger.Printf("request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func NewLogger(logger *log.Logger) *Middleware {
	return &Middleware{logger}
}

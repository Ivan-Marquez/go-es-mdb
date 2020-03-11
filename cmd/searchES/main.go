package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ivan-marquez/es-mdb/pkg/domain"
	"github.com/ivan-marquez/es-mdb/pkg/http/rest"
	"github.com/ivan-marquez/es-mdb/pkg/storage"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	logger := log.New(os.Stdout, "searchES ", log.LstdFlags|log.Lshortfile)
	store, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}

	var us domain.UserService
	us = domain.NewUserService(store)
	mw := rest.NewMiddleware(logger)
	h := rest.NewHandler(mw, us)

	mux := http.NewServeMux()
	h.SetupRoutes(mux)

	srv := rest.NewServer(mux, ":8080")

	logger.Println("server starting")
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
}

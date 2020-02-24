package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ivan-marquez/es-mdb/pkg/search"
	"github.com/ivan-marquez/es-mdb/pkg/server"
	"github.com/ivan-marquez/es-mdb/pkg/storage"
)

func main() {
	logger := log.New(os.Stdout, "searchService ", log.LstdFlags|log.Lshortfile)
	store, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	h := search.NewHandlers(logger, store.ESClient)

	mux := http.NewServeMux()
	h.SetupRoutes(mux)

	srv := server.New(mux, ":8080")

	logger.Println("server starting")
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
}

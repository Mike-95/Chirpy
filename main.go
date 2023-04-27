package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {

	const filepathRoot = "."
	const port = "8080"
	cfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()
	r.Mount("/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))

	mainRouter := chi.NewRouter()
	mainRouter.Get("/healthz", handlerReadiness)
	r.Mount("/api", mainRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.handlerMetrics)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

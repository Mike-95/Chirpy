package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"www.github.com/Mike-95/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	cfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	r := chi.NewRouter()
	r.Mount("/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Post("/chirps", cfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", cfg.handlerChirpsRetrieve)
	r.Mount("/api", apiRouter)

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

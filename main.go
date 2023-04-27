package main

import (
	"encoding/json"
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
	mainRouter.Post("/validate_chirp", handlerValidate)
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

func handlerValidate(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
		return
	}
	respondWithJSON(writer, http.StatusOK, returnVals{
		Valid: true,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})

}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

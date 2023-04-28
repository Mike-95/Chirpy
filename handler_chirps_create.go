package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}
	chirp, err := cfg.DB.CreateChirpy(cleaned)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}
	respondWithJSON(writer, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})

}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanBody(body, badWords)
	return cleaned, nil

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

func getCleanBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned

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

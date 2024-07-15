package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseWithJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	id_string := r.PathValue("chirpID")
	id, err := strconv.Atoi(id_string)
	if err != nil {
		log.Println("Error:", err)
		responseWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
	}
	if id > len(chirps) {
		responseWithError(w, http.StatusNotFound, "Chirp does not exist")
		return
	}
	fmt.Println(id)
	responseWithJson(w, http.StatusOK, chirps[id-1])
}

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := removeProfanity(params.Body)

	chirp, err := cfg.db.CreateChirp(cleaned)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseWithJson(w, http.StatusCreated, chirp)
}

func removeProfanity(body string) string {
	words := strings.Split(body, " ")
	for i, w := range words {
		if strings.ToLower(w) == "kerfuffle" || strings.ToLower(w) == "sharbert" || strings.ToLower(w) == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

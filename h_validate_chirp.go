package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

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

	type returnVals struct {
		// the key will be the name of struct field unless you give it an explicit JSON tag
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := returnVals{
		CleanedBody: removeProfanity(params.Body),
	}
	responseWithJson(w, http.StatusOK, respBody)

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

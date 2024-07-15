package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerPostUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newUser, err := cfg.db.CreateUser(params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseWithJson(w, http.StatusCreated, newUser)
}

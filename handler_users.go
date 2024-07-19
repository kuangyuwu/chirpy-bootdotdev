package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kuangyuwu/chirpy-bootdev/internal/auth"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerPostUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		responseWithError(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	newUser, err := cfg.db.CreateUser(params.Email, hashed)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		responseWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	responseWithJson(w, http.StatusCreated, User{
		Id:    newUser.Id,
		Email: newUser.Email,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	tokenString, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	id, err := auth.ValidateToken(tokenString, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		responseWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		responseWithError(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	newUser, err := cfg.db.UpdateUser(id, params.Email, hashed)
	if err != nil {
		log.Printf("Error updating user: %s", err)
		responseWithError(w, http.StatusInternalServerError, "error updating user")
		return
	}

	responseWithJson(w, http.StatusOK, User{
		Id:    newUser.Id,
		Email: newUser.Email,
	})
}

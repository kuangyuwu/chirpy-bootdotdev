package main

import (
	"encoding/json"
	"log"
	"net/http"

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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		log.Printf("Error: %s", err)
		responseWithError(w, http.StatusUnauthorized, "User not found")
		return
	}

	err = auth.CheckPassword(user.Hashed, params.Password)
	if err != nil {
		log.Printf("Error: %s", err)
		responseWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	responseWithJson(w, http.StatusOK, User{
		Id:    user.Id,
		Email: user.Email,
	})
}

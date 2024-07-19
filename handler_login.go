package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kuangyuwu/chirpy-bootdev/internal/auth"
)

type UserWithToken struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ExpiresInSecond int    `json:"expires_in_seconds"`
		Email           string `json:"email"`
		Password        string `json:"password"`
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

	token, err := auth.GetToken(user.Id, params.ExpiresInSecond, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error generating token: %s", err)
		responseWithError(w, http.StatusUnauthorized, "Error generating token")
		return
	}

	responseWithJson(w, http.StatusOK, UserWithToken{
		Id:    user.Id,
		Email: user.Email,
		Token: token,
	})
}

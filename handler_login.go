package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kuangyuwu/chirpy-bootdev/internal/auth"
)

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
		responseWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	err = auth.CheckPassword(user.Hashed, params.Password)
	if err != nil {
		log.Printf("Error: %s", err)
		responseWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := auth.GetToken(user.Id, params.ExpiresInSecond, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error generating token: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(user.Id)
	if err != nil {
		log.Printf("Error generating refresh token: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	type data struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	responseWithJson(w, http.StatusOK, data{
		Id:           user.Id,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	refreshTokenString, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	userId, err := cfg.db.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		log.Printf("Error validating refreshing token: %s", err)
		responseWithError(w, http.StatusUnauthorized, "Invalid refreshing token")
		return
	}

	token, err := auth.GetToken(userId, 3600, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error generating token: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	type data struct {
		Token string `json:"token"`
	}
	responseWithJson(w, http.StatusOK, data{
		Token: token,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	refreshTokenString, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	err := cfg.db.DeleteRefreshToken(refreshTokenString)
	if err != nil {
		log.Printf("Error revoking refreshing token: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Error revoking refreshing token")
		return
	}

	responseWithNoContent(w)
}

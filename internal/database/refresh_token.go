package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type RefreshToken struct {
	UserId           int   `json:"user_id"`
	ExpiresAtInMicro int64 `json:"expires_at"`
}

func (db *DB) CreateRefreshToken(userID int) (string, error) {

	length := 32
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(token)

	dbStructure, err := db.loadDB()
	if err != nil {
		return "", err
	}

	expires_in_days := 60
	expires := time.Duration(expires_in_days*24) * time.Hour
	expiresAtInMicro := time.Now().Add(expires).UnixMicro()
	dbStructure.RefreshTokens[tokenString] = RefreshToken{
		UserId:           userID,
		ExpiresAtInMicro: expiresAtInMicro,
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (db *DB) ValidateRefreshToken(refreshTokenString string) (userId int, err error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return 0, err
	}

	rt, ok := dbStructure.RefreshTokens[refreshTokenString]
	if !ok {
		return 0, errors.New("refresh token not found")
	}

	expiresAt := time.UnixMicro(rt.ExpiresAtInMicro)
	if time.Now().After(expiresAt) {
		err = db.DeleteRefreshToken(refreshTokenString)
		return 0, errors.Join(errors.New("refresh token expired"), err)
	}

	return rt.UserId, nil
}

func (db *DB) DeleteRefreshToken(refreshTokenString string) error {

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.RefreshTokens[refreshTokenString]
	if !ok {
		return errors.New("refresh token not found")
	}
	delete(dbStructure.RefreshTokens, refreshTokenString)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

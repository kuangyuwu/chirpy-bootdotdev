package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetToken(id, expires_in_seconds int, verificationKey string) (string, error) {

	if expires_in_seconds == 0 || expires_in_seconds > 86400 {
		expires_in_seconds = 86400
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expires_in_seconds))),
		Subject:   strconv.Itoa(id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(verificationKey))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateToken(tokenString string, verificationKey string) (int, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(verificationKey), nil
	})
	if err != nil {
		return 0, err
	}
	idString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}
	return id, nil
}

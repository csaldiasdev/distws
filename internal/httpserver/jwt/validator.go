package jwt

import (
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrInvalidIssuer = errors.New("issuer has invalid")
)

func ValidateToken(tokenString string, issuer string, secret string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}

	jwtClaims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, jwtClaims, keyFunc)

	if err != nil {
		return nil, err
	}

	if !jwtClaims.VerifyIssuer(issuer, true) {
		return nil, ErrInvalidIssuer
	}

	test, _ := json.Marshal(jwtClaims)

	payload := Payload{}

	json.Unmarshal(test, &payload)

	return &payload, nil
}

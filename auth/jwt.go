package auth

import (
	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

func MakeToken(email string) (string, error) {
	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{"email": email})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

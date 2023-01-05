package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/jgsheppa/mongo-go/errors"
	"github.com/lestrrat-go/jwx/jwt"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			response := errors.Unauthorized(err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			response := errors.Unauthorized(errors.ErrNoToken)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

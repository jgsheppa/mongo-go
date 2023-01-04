package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/jgsheppa/mongo-go/auth"
	"github.com/jgsheppa/mongo-go/errors"
	"github.com/jgsheppa/mongo-go/models"
)

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	us models.UserService
}

func NewUser(ms models.UserService) *User {
	return &User{
		ms,
	}
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var login LoginForm

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Printf("err from ReadAll: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &login); err != nil {
		fmt.Printf("err from unmarshall: %v", err)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := u.us.Authenticate(login.Email, login.Password)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	token, err := auth.MakeToken(user.Email)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Secure:   true,
		Name:     "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
		Value:    token,
	})

	http.Redirect(w, r, "/magazines", http.StatusFound)
}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/magazines", http.StatusFound)
}

func (u *User) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		responseError := errors.InternalError("JSON web token context failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	email := claims["email"].(string)
	user, err := u.us.ByEmail(email)
	if err != nil {
		responseError := errors.NotFound(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

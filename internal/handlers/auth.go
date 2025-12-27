package handlers

import (
	db "dynoc-registry/internal/db/gen"
	"dynoc-registry/internal/jwt"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"time"

	"dynoc-registry/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}
	if req.Username.String == "" || req.Password.String == "" {
		writeJSON(w, http.StatusBadRequest, models.MissingFieldError)
		return
	}

	pool := getDB(r)
	q := db.New(pool)
	n := rand.IntN(10_000) + 5_000 // to mitigate timing attacks

	row, err := q.GetSecurityInfo(r.Context(), req.Username)
	if err != nil {
		time.Sleep(time.Duration(n) * time.Millisecond)
		writeJSON(w, http.StatusUnauthorized, models.Response{Message: "invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(row.Password.String), []byte(req.Password.String))
	if err != nil {
		time.Sleep(time.Duration(n) * time.Millisecond)
		writeJSON(w, http.StatusUnauthorized, models.Response{Message: "invalid username or password"})
		return
	}

	token, err := jwt.CreateToken(row.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	resp := models.TokenResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: 3600,
	}

	writeJSON(w, http.StatusOK, resp)
}

func RenewToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		writeJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
	}

	userId, err := jwt.GetTokenClaims(tokenString)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	pool := getDB(r)
	q := db.New(pool)

	row, err := q.GetUserByID(r.Context(), userId)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	token, err := jwt.CreateToken(row.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	resp := models.TokenResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: 3600,
	}

	writeJSON(w, http.StatusOK, resp)
}

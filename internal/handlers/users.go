package handlers

import (
	db "dynoc-registry/internal/db/gen"
	"dynoc-registry/internal/jwt"
	"dynoc-registry/internal/models"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}
	if req.Username.String == "" || req.Email.String == "" || req.Password.String == "" {
		writeJSON(w, http.StatusBadRequest, models.MissingFieldError)
		return
	}

	pool := getDB(r)
	q := db.New(pool)

	uRow, err := q.GetUserByUsername(r.Context(), req.Username)
	if err == nil && uRow.ID.String() != "" {
		writeJSON(w, http.StatusBadRequest, models.Response{Message: "username already taken"})
		return
	}

	eRow, err := q.GetUserByEmail(r.Context(), req.Email)
	if err == nil && eRow.ID.String() != "" {
		writeJSON(w, http.StatusBadRequest, models.Response{Message: "email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password.String), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	var password pgtype.Text
	password.String = string(hashedPassword)
	password.Valid = true

	row, err := q.CreateUser(r.Context(), db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: password,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	resp := models.RegisterResponse{
		UserID:    row.ID,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	writeJSON(w, http.StatusCreated, resp)
}

func ReadAccount(w http.ResponseWriter, r *http.Request) {
	var username pgtype.Text
	err := username.Scan(r.PathValue("name"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	pool := getDB(r)
	q := db.New(pool)

	row, err := q.GetUserByUsername(r.Context(), username)
	if err != nil || row.ID.String() == "" {
		writeJSON(w, http.StatusNotFound, models.NotFoundError)
		return
	}

	resp := models.AccountResponse{
		UserID:    row.ID,
		Username:  row.Username,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, resp)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
	}

	userId, err := jwt.GetTokenClaims(token)
	if err != nil {
		return
	}

	var req models.UpdateAccountRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	pool := getDB(r)
	q := db.New(pool)

	user, err := q.GetUserByID(r.Context(), userId)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	updateUsername := r.PathValue("name")
	if user.ID != userId || user.Username.String != updateUsername {
		writeJSON(w, http.StatusForbidden, models.ForbiddenError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password.String), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	err = q.UpdateUser(r.Context(), db.UpdateUserParams{
		ID: userId,
		Username: pgtype.Text{
			String: req.Username.String,
			Valid:  req.Username.String != "",
		},
		Email: pgtype.Text{
			String: req.Email.String,
			Valid:  req.Email.String != "",
		},
		Password: pgtype.Text{
			String: string(hashedPassword),
			Valid:  req.Password.String != "",
		},
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
	}

	userId, err := jwt.GetTokenClaims(token)
	if err != nil {
		return
	}

	pool := getDB(r)
	q := db.New(pool)

	user, err := q.GetUserByID(r.Context(), userId)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	deleteUsername := r.PathValue("name")
	if user.ID != userId || user.Username.String != deleteUsername {
		writeJSON(w, http.StatusForbidden, models.ForbiddenError)
		return
	}

	err = q.DeleteUser(r.Context(), userId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, models.SuccessResponse)
}

package handlers

import (
	"dynoc-registry/internal/commons"
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
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}
	if req.Username.String == "" || req.Email.String == "" || req.Password.String == "" {
		commons.WriteJSON(w, http.StatusBadRequest, models.MissingFieldError)
		return
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	uRow, err := q.GetUserByUsername(r.Context(), req.Username)
	if err == nil && uRow.UserID.String() != "" {
		commons.WriteJSON(w, http.StatusConflict, models.Response{Message: "username already taken"})
		return
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	eRow, err := q.GetUserByEmail(r.Context(), req.Email)
	if err == nil && eRow.UserID.String() != "" {
		commons.WriteJSON(w, http.StatusConflict, models.Response{Message: "email already registered"})
		return
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password.String), bcrypt.DefaultCost)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
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
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	resp := models.RegisterResponse{
		UserID:    row.ID,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	commons.WriteJSON(w, http.StatusCreated, resp)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	var username pgtype.Text
	err := username.Scan(r.PathValue("name"))
	if err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	row, err := q.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			commons.WriteJSON(w, http.StatusNotFound, models.NotFoundError)
			return
		}
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	resp := models.AccountResponse{
		UserID:    row.UserID,
		Username:  row.UserName,
		Email:     row.UserEmail,
		CreatedAt: row.UserCreatedAt,
		UpdatedAt: row.UserUpdatedAt,
	}

	commons.WriteJSON(w, http.StatusOK, resp)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
	}

	userId, err := jwt.GetTokenClaims(token)
	if err != nil {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	var req models.UpdateAccountRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	user, err := q.GetUserByID(r.Context(), userId)
	if err != nil {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	updateUsername := r.PathValue("name")
	if user.UserID != userId || user.UserName.String != updateUsername {
		commons.WriteJSON(w, http.StatusForbidden, models.ForbiddenError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password.String), bcrypt.DefaultCost)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
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
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	commons.WriteJSON(w, http.StatusOK, models.SuccessResponse)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
	}

	userId, err := jwt.GetTokenClaims(token)
	if err != nil {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	user, err := q.GetUserByID(r.Context(), userId)
	if err != nil {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	deleteUsername := r.PathValue("name")
	if user.UserID != userId || user.UserName.String != deleteUsername {
		commons.WriteJSON(w, http.StatusForbidden, models.ForbiddenError)
		return
	}

	err = q.DeleteUser(r.Context(), userId)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	commons.WriteJSON(w, http.StatusOK, models.SuccessResponse)
}

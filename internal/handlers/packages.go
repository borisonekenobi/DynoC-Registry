package handlers

import (
	"dynoc-registry/internal/commons"
	db "dynoc-registry/internal/db/gen"
	"dynoc-registry/internal/jwt"
	"dynoc-registry/internal/models"
	"encoding/json"
	"net/http"
)

func CreatePackage(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	userId, err := jwt.GetTokenClaims(token)
	if err != nil {
		commons.WriteJSON(w, http.StatusUnauthorized, models.UnauthorizedError)
		return
	}

	var req models.CreatePackageRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	row, err := q.CreatePackage(r.Context(), db.CreatePackageParams{
		Name:        req.Name,
		Description: req.Description,
		Visibility:  req.Visibility,
		OwnerID:     userId,
	})
	if err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	resp := models.PackageResponse{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Visibility:  row.Visibility,
		Owner:       row.OwnerUsername,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}

	commons.WriteJSON(w, http.StatusCreated, resp)
}

func CreatePackageVersion(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

func GetLatest(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

func GetVersions(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

func GetPackage(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

func UpdatePackage(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

func UpdatePackageVersion(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

func DeletePackage(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

func DeletePackageVersion(w http.ResponseWriter, r *http.Request) {
	commons.WriteJSON(w, http.StatusNotImplemented, models.NotImplementedError)
}

package handlers

import (
	"dynoc-registry/internal/commons"
	db "dynoc-registry/internal/db/gen"
	"dynoc-registry/internal/jwt"
	"dynoc-registry/internal/models"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
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

	var req models.VersionMetadata
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	if !req.PackageID.Valid || !req.Version.Valid {
		commons.WriteJSON(w, http.StatusBadRequest, models.MissingFieldError)
		return
	}

	cs, err := commons.CalculateSHA256([]byte(""))
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}
	checksum := pgtype.Text{String: cs, Valid: true}

	pool := commons.GetDB(r)
	q := db.New(pool)

	pkg, err := q.GetPackageByID(r.Context(), req.PackageID)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}
	if pkg.OwnerID != userId {
		commons.WriteJSON(w, http.StatusForbidden, models.ForbiddenError)
		return
	}

	_, err = q.GetPackageByVersion(r.Context(), db.GetPackageByVersionParams{
		Name:    pkg.Name,
		Version: req.Version,
	})
	if err != nil {
		msg := err.Error()
		if msg != "no rows in result set" {
			commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
			return
		}
	}
	if err == nil {
		commons.WriteJSON(w, http.StatusConflict, models.ConflictError)
		return
	}

	v, err := q.CreatePackageVersion(r.Context(), db.CreatePackageVersionParams{
		PackageID: req.PackageID,
		Version:   req.Version,
		Checksum:  checksum,
		SizeBytes: pgtype.Int8{Int64: 0, Valid: true},
		Location:  pgtype.Text{String: "s3://bucket/path/to/package", Valid: true},
	})
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	for depName, depVersion := range req.Dependencies {
		err = q.CreateDependency(r.Context(), db.CreateDependencyParams{
			VersionID:      v.ID,
			DependencyName: pgtype.Text{String: depName, Valid: true},
			ConstraintExpr: depVersion,
		})
		if err != nil {
			commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
			return
		}
	}

	resp := models.VersionResponse{
		Name:         pkg.Name,
		Version:      v.Version,
		Checksum:     v.Checksum,
		Size:         v.SizeBytes,
		Dependencies: req.Dependencies,
		DownloadURL:  v.Location,
	}

	commons.WriteJSON(w, http.StatusCreated, resp)
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

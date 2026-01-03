package handlers

import (
	"dynoc-registry/internal/commons"
	db "dynoc-registry/internal/db/gen"
	"dynoc-registry/internal/jwt"
	"dynoc-registry/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

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

	var pkgName pgtype.Text
	err = pkgName.Scan(r.PathValue("name"))
	if err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	var req models.VersionMetadata
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	if !req.Version.Valid {
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

	pkg, err := q.GetPackageByName(r.Context(), pkgName)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}
	if pkg.PackageOwnerID != userId {
		commons.WriteJSON(w, http.StatusForbidden, models.ForbiddenError)
		return
	}

	_, err = q.GetPackageByVersion(r.Context(), db.GetPackageByVersionParams{
		Name:    pkg.PackageName,
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
		PackageID: pkg.PackageID,
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
			Name:           pgtype.Text{String: depName, Valid: true},
			ConstraintExpr: depVersion,
		})
		if err != nil {
			commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
			return
		}
	}

	resp := models.VersionResponse{
		Name:         pkg.PackageName,
		Version:      v.Version,
		Checksum:     v.Checksum,
		Size:         v.SizeBytes,
		Dependencies: req.Dependencies,
		DownloadURL:  v.Location,
	}

	commons.WriteJSON(w, http.StatusCreated, resp)
}

func GetLatest(w http.ResponseWriter, r *http.Request) {
	var pkgName pgtype.Text
	err := pkgName.Scan(r.PathValue("name"))
	if err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	v, err := q.GetLatestPackageVersion(r.Context(), pkgName)
	if err != nil {
		if err.Error() == "no rows in result set" {
			commons.WriteJSON(w, http.StatusNotFound, models.NotFoundError)
			return
		}
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	pkg, err := q.GetPackageByID(r.Context(), v.PackageID)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	depsRows, err := q.GetDependenciesByVersionID(r.Context(), v.PackageVersionID)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	dependencies := make(map[string]pgtype.Text)
	for _, dep := range depsRows {
		dependencies[dep.DependencyName.String] = dep.ConstraintExpr
	}

	resp := models.VersionResponse{
		Name:         pkg.PackageName,
		Version:      v.PackageVersion,
		Checksum:     v.PackageVersionChecksum,
		Size:         v.PackageVersionSizeBytes,
		Dependencies: dependencies,
		DownloadURL:  v.PackageVersionLocation,
	}

	commons.WriteJSON(w, http.StatusOK, resp)
}

func GetVersions(w http.ResponseWriter, r *http.Request) {
	var pkgName pgtype.Text
	err := pkgName.Scan(r.PathValue("name"))
	if err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	skip := commons.DefaultSkip
	skipParam := r.URL.Query().Get("skip")
	if skipParam != "" {
		i, err := strconv.ParseInt(skipParam, 10, 64)
		if err != nil || i < 0 {
			commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
			return
		}
		skip = pgtype.Int4{Int32: int32(i), Valid: true}
	}

	take := commons.DefaultTake
	takeParam := r.URL.Query().Get("skip")
	if takeParam != "" {
		i, err := strconv.ParseInt(takeParam, 10, 32)
		if err != nil || i < 0 {
			commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
			return
		}
		take = pgtype.Int4{Int32: int32(i), Valid: true}
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	versionsRows, err := q.GetPackageVersionsByName(r.Context(), db.GetPackageVersionsByNameParams{
		Name:   pkgName,
		Limit:  take,
		Offset: skip,
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			commons.WriteJSON(w, http.StatusNotFound, models.NotFoundError)
			return
		}
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	resp := models.VersionsListResponse{
		Name:     pkgName,
		Versions: make([]pgtype.Text, len(versionsRows)),
	}
	for i, v := range versionsRows {
		resp.Versions[i] = v.PackageVersion
	}

	commons.WriteJSON(w, http.StatusOK, resp)
}

func GetPackage(w http.ResponseWriter, r *http.Request) {
	var pkgName pgtype.Text
	err := pkgName.Scan(r.PathValue("name"))
	if err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	var pkgVersion pgtype.Text
	err = pkgVersion.Scan(r.PathValue("version"))
	if err != nil {
		commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
		return
	}

	pool := commons.GetDB(r)
	q := db.New(pool)

	pkgVer, err := q.GetPackageByVersion(r.Context(), db.GetPackageByVersionParams{
		Name:    pkgName,
		Version: pkgVersion,
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			commons.WriteJSON(w, http.StatusNotFound, models.NotFoundError)
			return
		}
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	depsRows, err := q.GetDependenciesByVersionID(r.Context(), pkgVer.PackageVersionID)
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	dependencies := make(map[string]pgtype.Text)
	for _, dep := range depsRows {
		dependencies[dep.DependencyName.String] = dep.ConstraintExpr
	}

	resp := models.VersionResponse{
		Name:         pkgVer.PackageName,
		Version:      pkgVer.PackageVersion,
		Checksum:     pkgVer.PackageVersionChecksum,
		Size:         pkgVer.PackageVersionSizeBytes,
		Dependencies: dependencies,
		DownloadURL:  pkgVer.PackageVersionLocation,
	}

	commons.WriteJSON(w, http.StatusOK, resp)
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

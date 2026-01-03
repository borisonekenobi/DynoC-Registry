package handlers

import (
	"dynoc-registry/internal/commons"
	db "dynoc-registry/internal/db/gen"
	"dynoc-registry/internal/models"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func FindPackages(w http.ResponseWriter, r *http.Request) {
	query := pgtype.Text{String: r.URL.Query().Get("q"), Valid: true}
	if query.String == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}

	skip := commons.DefaultSkip
	skipParam := r.URL.Query().Get("skip")
	if skipParam != "" {
		i, err := strconv.ParseInt(skipParam, 10, 32)
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

	rows, err := q.FindPackages(r.Context(), db.FindPackagesParams{
		Column1: query,
		Limit:   take,
		Offset:  skip,
	})
	if err != nil {
		commons.WriteJSON(w, http.StatusInternalServerError, models.InternalServerError)
		return
	}

	resp := make([]models.PackageResponse, len(rows))
	for i, row := range rows {
		resp[i] = models.PackageResponse{
			Name:        row.PackageName,
			Description: row.PackageDescription,
			Visibility:  row.PackageVisibility,
			Owner:       row.PackageOwnerUsername,
			CreatedAt:   row.PackageCreatedAt,
			UpdatedAt:   row.PackageUpdatedAt,
		}
	}

	commons.WriteJSON(w, http.StatusOK, resp)
}

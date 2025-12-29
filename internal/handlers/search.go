package handlers

import (
	"dynoc-registry/internal/commons"
	db "dynoc-registry/internal/db/gen"
	"dynoc-registry/internal/models"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

const defaultSkip int32 = 0
const defaultTake int32 = 50

func FindPackages(w http.ResponseWriter, r *http.Request) {
	query := pgtype.Text{String: r.URL.Query().Get("q"), Valid: true}
	if query.String == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}

	skip := defaultSkip
	skipParam := r.URL.Query().Get("skip")
	if skipParam != "" {
		i, err := strconv.ParseInt(skipParam, 10, 64)
		if err != nil || i < 0 {
			commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
			return
		}
		skip = int32(i)
	}

	take := defaultTake
	takeParam := r.URL.Query().Get("skip")
	if takeParam != "" {
		i, err := strconv.ParseInt(takeParam, 10, 32)
		if err != nil || i < 0 {
			commons.WriteJSON(w, http.StatusBadRequest, models.BadRequestError)
			return
		}
		take = int32(i)
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

	commons.WriteJSON(w, http.StatusOK, rows)
}

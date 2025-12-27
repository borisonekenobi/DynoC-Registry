package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
}

func getDB(r *http.Request) *pgxpool.Pool {
	return r.Context().Value("db").(*pgxpool.Pool)
}

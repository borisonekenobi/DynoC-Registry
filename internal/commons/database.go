package commons

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDB(r *http.Request) *pgxpool.Pool {
	return r.Context().Value("db").(*pgxpool.Pool)
}

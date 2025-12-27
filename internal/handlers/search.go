package handlers

import (
	"net/http"
)

func FindPackages(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"query": q,
		"results": []map[string]string{
			{
				"name":           "dynoc-json",
				"description":    "JSON support for DynoC",
				"latest_version": "1.0.0",
			},
		},
	}

	writeJSON(w, http.StatusOK, resp)
}

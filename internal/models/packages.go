package models

import (
	db "dynoc-registry/internal/db/gen"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreatePackageRequest struct {
	Name        pgtype.Text   `json:"name"`
	Description pgtype.Text   `json:"description"`
	Visibility  db.Visibility `json:"visibility"`
}

type PackageResponse struct {
	Name        pgtype.Text        `json:"name"`
	Description pgtype.Text        `json:"description"`
	Visibility  db.Visibility      `json:"visibility"`
	Owner       pgtype.Text        `json:"owner"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type VersionMetadata struct {
	Version      pgtype.Text            `json:"version"`
	Dependencies map[string]pgtype.Text `json:"dependencies"`
}

type VersionResponse struct {
	Name         pgtype.Text            `json:"name"`
	Version      pgtype.Text            `json:"version"`
	Checksum     pgtype.Text            `json:"checksum"`
	Size         pgtype.Int8            `json:"size"`
	Dependencies map[string]pgtype.Text `json:"dependencies"`
	DownloadURL  pgtype.Text            `json:"download_url"`
}

type VersionsListResponse struct {
	Name     pgtype.Text   `json:"name"`
	Versions []pgtype.Text `json:"versions"`
}

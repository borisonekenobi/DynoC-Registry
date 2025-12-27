package models

import "time"

type CreatePackageRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

type PackageResponse struct {
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}

type VersionMetadata struct {
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

type VersionResponse struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Checksum     string            `json:"checksum"`
	Size         int64             `json:"size"`
	Dependencies map[string]string `json:"dependencies"`
	DownloadURL  string            `json:"download_url"`
}

type VersionsListResponse struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type Visibility int

const (
	VisibilityPublic Visibility = iota
	VisibilityPrivate
)

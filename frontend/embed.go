package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist/frontend/browser/*
var dist embed.FS

func FS() (fs.FS, error) {
	fsys, err := fs.Sub(dist, "dist/frontend/browser")

	return fsys, err
}

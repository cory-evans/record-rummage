package api

import (
	"errors"
	"net/http"
	"os"
)

type SpaFileSystem struct {
	root http.FileSystem
}

func NewSpaFileSystem(root http.FileSystem) *SpaFileSystem {
	return &SpaFileSystem{root}
}

func (s *SpaFileSystem) Open(name string) (http.File, error) {
	f, err := s.root.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return s.root.Open("index.html")
	}

	return f, err
}

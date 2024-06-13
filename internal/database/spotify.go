package database

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type SpotifyRepository struct {
	db *sqlx.DB
}

func NewSpotifyRepository(db *sqlx.DB) *SpotifyRepository {
	return &SpotifyRepository{db: db}
}

const createLoginStateQuery = `INSERT INTO spotify_Login_StateT (state) VALUES ($1)`

func (r *SpotifyRepository) CreateLoginState(state string) error {
	_, err := r.db.Exec(createLoginStateQuery, state)

	return err
}

const getLoginStateQuery = `SELECT COUNT(*) FROM spotify_Login_StateT WHERE state = $1`

func (r *SpotifyRepository) CheckLoginState(state string, remove bool) error {
	var count int
	err := r.db.QueryRow(getLoginStateQuery, state).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("invalid state")
	}

	if remove {
		return r.RemoveLoginState(state)
	}

	return nil
}

func (r *SpotifyRepository) RemoveLoginState(state string) error {
	_, err := r.db.Exec(`DELETE FROM spotify_Login_StateT WHERE state = $1`, state)
	return err
}

package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type SpotifyRepository struct {
	db *pgx.Conn
}

func NewSpotifyRepository(db *pgx.Conn) *SpotifyRepository {
	return &SpotifyRepository{db: db}
}

const createLoginStateQuery = `INSERT INTO spotify_Login_StateT (state) VALUES ($1)`

func (r *SpotifyRepository) CreateLoginState(ctx context.Context, state string) error {
	_, err := r.db.Exec(ctx, createLoginStateQuery, state)

	return err
}

const getLoginStateQuery = `SELECT COUNT(*) FROM spotify_Login_StateT WHERE state = $1`

func (r *SpotifyRepository) CheckLoginState(ctx context.Context, state string, remove bool) error {
	var count int
	err := r.db.QueryRow(ctx, getLoginStateQuery, state).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("invalid state")
	}

	if remove {
		return r.RemoveLoginState(ctx, state)
	}

	return nil
}

func (r *SpotifyRepository) RemoveLoginState(ctx context.Context, state string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM spotify_Login_StateT WHERE state = $1`, state)
	return err
}

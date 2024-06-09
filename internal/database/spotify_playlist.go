package database

import (
	"context"
	"encoding/json"

	"github.com/cory-evans/record-rummage/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/zmb3/spotify/v2"
)

func (r *SpotifyRepository) GetPlaylistSnapshot(ctx context.Context, id spotify.ID) (string, error) {
	var snapshotID string
	err := r.db.QueryRow(ctx, `SELECT snapshot_id FROM spotify_PlaylistT WHERE id = $1`, id).Scan(&snapshotID)
	if err != nil {
		return "", err
	}

	return snapshotID, nil
}

func (r *SpotifyRepository) GetAddedBy(ctx context.Context, playlistID string, trackID string) ([]string, error) {
	var addedBy = make([]string, 0)
	rows, err := r.db.Query(ctx, `
SELECT DISTINCT added_by
FROM spotify_Playlist_TrackT
WHERE playlist_id = $1 AND track_id = $2`, playlistID, trackID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		addedBy = append(addedBy, id)
	}

	return addedBy, nil
}

func (r *SpotifyRepository) CreateOrUpdatePlaylist(ctx context.Context, playlist models.SpotifyPlaylist) error {
	sql := `
        MERGE INTO spotify_PlaylistT AS target
        USING (
            VALUES
                ($1, $2, $3, $4::json)
        ) AS source (id, snapshot_id, name, images)
        ON target.id = source.id
        WHEN MATCHED THEN
            UPDATE SET
                snapshot_id = source.snapshot_id,
                name = source.name,
                images = source.images
        WHEN NOT MATCHED THEN
            INSERT (id, snapshot_id, name, images)
            VALUES (source.id, source.snapshot_id, source.name, source.images);
    `

	imageJson, err := json.Marshal(playlist.Images)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, playlist.ID, playlist.SnapshotID, playlist.Name, imageJson)

	return err
}

func (r *SpotifyRepository) Refresh(ctx context.Context, playlistID string, tracks []spotify.PlaylistItem) error {
	// start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// delete all tracks for the playlist
	_, err = tx.Exec(context.Background(), `DELETE FROM spotify_Playlist_TrackT WHERE playlist_id = $1`, playlistID)
	if err != nil {
		return err
	}

	// insert all tracks for the playlist
	stmt, err := tx.Prepare(context.Background(), "insert_spotify_Playlist_TrackT", `INSERT INTO spotify_Playlist_TrackT (playlist_id, track_id, added_at, added_by) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}

	for _, t := range tracks {
		track := t.Track.Track

		if track == nil {
			continue
		}

		createOrUpdateTrack(ctx, tx, models.SpotifyTrack{
			ID:      track.ID.String(),
			Name:    track.Name,
			Artists: models.SpotifyArtistFromSlice(track.Artists),
			Album: models.SpotifyAlbum{
				ID:     track.Album.ID.String(),
				Name:   track.Album.Name,
				Images: models.SpotifyImageFromList(track.Album.Images),
			},
		})

		_, err = tx.Exec(ctx, stmt.SQL, playlistID, t.Track.Track.ID.String(), t.AddedAt, t.AddedBy.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func createOrUpdateTrack(ctx context.Context, tx pgx.Tx, track models.SpotifyTrack) error {
	sql := `
MERGE INTO spotify_TrackT AS target
USING (
    VALUES
        ($1, $2, $3::json, $4::json)
) AS source (id, name, album, artists)
ON target.id = source.id
WHEN MATCHED THEN
    UPDATE SET
        name = source.name,
        album = source.album,
        artists = source.artists
WHEN NOT MATCHED THEN
    INSERT (id, name, album, artists)
    VALUES (source.id, source.name, source.album, source.artists);`

	artistJson, err := json.Marshal(track.Artists)
	if err != nil {
		return err
	}

	albumJson, err := json.Marshal(track.Album)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, track.ID, track.Name, albumJson, artistJson)

	return err
}

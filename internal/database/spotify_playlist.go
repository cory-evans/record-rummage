package database

import (
	"encoding/json"

	"github.com/cory-evans/record-rummage/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
)

func (r *SpotifyRepository) GetPlaylistSnapshot(id spotify.ID) (string, error) {
	var snapshotID string
	err := r.db.QueryRow(`SELECT snapshot_id FROM spotify_PlaylistT WHERE id = $1`, id).Scan(&snapshotID)
	if err != nil {
		return "", err
	}

	return snapshotID, nil
}

func (r *SpotifyRepository) GetPlaylistSnapshotBulk(ids []spotify.ID) (map[spotify.ID]string, error) {
	var snapshotIDs = make(map[spotify.ID]string)
	rows, err := r.db.Queryx(`SELECT id, snapshot_id FROM spotify_PlaylistT WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id spotify.ID
		var snapshotID string
		err = rows.Scan(&id, &snapshotID)
		if err != nil {
			return nil, err
		}

		snapshotIDs[id] = snapshotID
	}

	return snapshotIDs, nil
}

func (r *SpotifyRepository) GetAddedBy(playlistID string, trackID string) ([]string, error) {
	var addedBy = make([]string, 0)
	err := r.db.Select(&addedBy, `
SELECT DISTINCT added_by
FROM spotify_Playlist_TrackT
WHERE playlist_id = $1 AND track_id = $2`, playlistID, trackID)

	if err != nil {
		return nil, err
	}

	return addedBy, nil
}

func (r *SpotifyRepository) CreateOrUpdatePlaylist(playlist models.SpotifyPlaylist) error {
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

	_, err = r.db.Exec(sql, playlist.ID, playlist.SnapshotID, playlist.Name, imageJson)

	return err
}

func (r *SpotifyRepository) Refresh(playlistID string, tracks []spotify.PlaylistItem, removeAll bool) error {
	// start a transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if removeAll {
		// delete all tracks for the playlist
		_, err = tx.Exec(`DELETE FROM spotify_Playlist_TrackT WHERE playlist_id = $1`, playlistID)
		if err != nil {
			return err
		}
	}

	// insert all tracks for the playlist
	insertStmt, err := tx.Prepare(`INSERT INTO spotify_Playlist_TrackT (playlist_id, track_id, added_at, added_by) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}

	mergeStmt, err := tx.Prepare(`
MERGE INTO spotify_Playlist_TrackT AS target
USING (
    VALUES
        ($1, $2, $3::timestamp, $4)
) AS source (playlist_id, track_id, added_at, added_by)
ON target.playlist_id = source.playlist_id AND target.track_id = source.track_id AND target.added_by = source.added_by
WHEN MATCHED THEN
    UPDATE SET
        added_at = source.added_at
WHEN NOT MATCHED THEN
    INSERT (playlist_id, track_id, added_at, added_by)
    VALUES (source.playlist_id, source.track_id, source.added_at, source.added_by);`)
	if err != nil {
		return err
	}

	for _, t := range tracks {
		track := t.Track.Track

		if track == nil {
			continue
		}

		createOrUpdateTrack(tx, models.SpotifyTrack{
			ID:      track.ID.String(),
			Name:    track.Name,
			Artists: models.SpotifyArtistFromSlice(track.Artists),
			Album: models.SpotifyAlbum{
				ID:     track.Album.ID.String(),
				Name:   track.Album.Name,
				Images: models.SpotifyImageFromList(track.Album.Images),
			},
		})

		if removeAll {
			_, err = insertStmt.Exec(playlistID, t.Track.Track.ID.String(), t.AddedAt, t.AddedBy.ID)
		} else {
			_, err = mergeStmt.Exec(playlistID, t.Track.Track.ID.String(), t.AddedAt, t.AddedBy.ID)
		}
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func createOrUpdateTrack(tx *sqlx.Tx, track models.SpotifyTrack) error {
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

	_, err = tx.Exec(sql, track.ID, track.Name, albumJson, artistJson)

	return err
}

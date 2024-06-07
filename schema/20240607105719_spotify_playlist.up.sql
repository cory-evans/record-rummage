CREATE TABLE spotify_UserT (
	id TEXT PRIMARY KEY,
	display_name TEXT NULL,
	images JSON NULL
);

CREATE TABLE spotify_TrackT (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	album JSON NULL,
	artists JSON NULL
);

CREATE TABLE spotify_PlaylistT (
	id TEXT PRIMARY KEY,
	snapshot_id TEXT NULL,
	name TEXT NOT NULL,
	images JSON NULL
);

CREATE TABLE spotify_Playlist_TrackT (
	playlist_id TEXT NOT NULL,
	track_id TEXT NOT NULL,
	added_at TIMESTAMP NOT NULL,
	added_by TEXT NULL,
	FOREIGN KEY (playlist_id) REFERENCES spotify_PlaylistT(id),
	FOREIGN KEY (track_id) REFERENCES spotify_TrackT(id)
);
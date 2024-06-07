package database

import (
	"context"
	"encoding/json"

	"github.com/cory-evans/record-rummage/internal/models"
)

func (r *SpotifyRepository) GetUser(ctx context.Context, id string) (*models.SpotifyUser, error) {
	sql := `SELECT id, display_name, images FROM spotify_UserT WHERE id = $1;`

	row := r.db.QueryRow(ctx, sql, id)

	u := &models.SpotifyUser{}
	imageJson := ""

	err := row.Scan(&u.ID, &u.DisplayName, &imageJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(imageJson), &u.Images)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *SpotifyRepository) CreateOrUpdateUser(ctx context.Context, user *models.SpotifyUser) error {
	sql := `
MERGE INTO spotify_UserT AS target
USING (
    VALUES
        ($1, $2, $3::json)
) AS source (id, display_name, images)
ON target.id = source.id
WHEN MATCHED THEN
    UPDATE SET
        display_name = source.display_name,
        images = source.images
WHEN NOT MATCHED THEN
    INSERT (id, display_name, images)
    VALUES (source.id, source.display_name, source.images);	
`

	imageJson, err := json.Marshal(user.Images)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, user.ID, user.DisplayName, imageJson)
	if err != nil {
		return err
	}

	return nil
}

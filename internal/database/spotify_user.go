package database

import (
	"encoding/json"

	"github.com/cory-evans/record-rummage/internal/models"
)

func (r *SpotifyRepository) GetUser(id string) (*models.SpotifyUser, error) {
	sql := `SELECT id, display_name, images FROM spotify_UserT WHERE id = $1;`

	u := &models.SpotifyUser{}
	err := r.db.QueryRowx(sql, id).StructScan(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *SpotifyRepository) CreateOrUpdateUser(user *models.SpotifyUser) error {
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

	_, err = r.db.Exec(sql, user.ID, user.DisplayName, imageJson)
	if err != nil {
		return err
	}

	return nil
}

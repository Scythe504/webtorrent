package postgresdb

import "time"

type STATUS string

const (
	PROCESSING  STATUS = "processing"
	DOWNLOADING STATUS = "downloading"
	FAILED      STATUS = "failed"
)

type Video struct {
	Id         string    `db:"id" json:"id"`
	MagnetLink string    `db:"magnet_link" json:"magnet_link"`
	Status     STATUS    `db:"status" json:"status"`
	FilePath   string    `db:"file_path" json:"file_path"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	Deleted    bool      `db:"deleted" json:"deleted"`
}

func (s *service) CreateVideo(video Video) error {
	stmt := `
		INSERT INTO videos (
			id,
			magnet_link,
			status,
			file_path,
			created_at,
			deleted
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.Exec(stmt,
		video.Id,
		video.MagnetLink,
		video.Status,
		video.FilePath,
		video.CreatedAt,
		video.Deleted,
	)

	return err
}

func (s *service) GetVideo(videoId string) (Video, error) {
	var video Video

	stmt := `
		SELECT 
			id, 
			magnet_link, 
			status, 
			file_path, 
			created_at, 
			deleted
		FROM videos
		WHERE id = $1
	`

	row := s.db.QueryRow(stmt, videoId)

	err := row.Scan(&video.Id, &video.MagnetLink, &video.Status, &video.FilePath, &video.CreatedAt, &video.Deleted)

	return video, err
}

func (s *service) UpdateStatus(status STATUS, videoId string) error {
	stmt := `
		UPDATE videos
		SET status = $1
		WHERE id = $2
	`

	_, err := s.db.Exec(stmt, status, videoId)
	return err
}

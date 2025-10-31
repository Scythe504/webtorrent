package postgresdb

import (
	"fmt"
	"time"
)

type STATUS string

const (
	PROCESSING  STATUS = "processing"
	DOWNLOADING STATUS = "downloading"
	DOWNLOADED  STATUS = "downloaded"
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

func (s *service) GetAllVideos() ([]Video, error) {
	stmt := `
		SELECT 
			id, 
			magnet_link, 
			status, 
			file_path, 
			created_at, 
			deleted
		FROM videos
		WHERE deleted = FALSE
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video

	for rows.Next() {
		var v Video
		err := rows.Scan(&v.Id, &v.MagnetLink, &v.Status, &v.FilePath, &v.CreatedAt, &v.Deleted)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	// handle possible iteration error
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}

func (s *service) UpdateStatus(status STATUS, videoId string, filePath *string) error {
	query := `
		UPDATE videos
		SET status = $1
		%s
		WHERE id = $2
	`

	// If filePath is provided, include it in the update
	var stmt string
	var args []any

	if filePath != nil {
		stmt = fmt.Sprintf(query, ", file_path = $3")
		args = []any{status, videoId, *filePath}
	} else {
		stmt = fmt.Sprintf(query, "")
		args = []any{status, videoId}
	}

	_, err := s.db.Exec(stmt, args...)
	return err
}

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/scythe504/webtorrent/internal"
	postgresdb "github.com/scythe504/webtorrent/internal/postgres-db"
	redisdb "github.com/scythe504/webtorrent/internal/redis-db"
)

func (s *Server) startVideo(w http.ResponseWriter, r *http.Request) {
	// 1. Get magnet link from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[StartVideo] Invalid Request body", err)
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var link struct {
		MagnetLink string `json:"magnet_link"`
	}

	if err = json.Unmarshal(body, &link); err != nil {
		log.Println("[StartVideo] Invalid Json Body", err)
		http.Error(w, "Failed to Parse JSON", http.StatusBadRequest)
		return
	}

	videoId := internal.RandomId()

	video := postgresdb.Video{
		Id:         videoId,
		MagnetLink: link.MagnetLink,
		Status:     postgresdb.PROCESSING,
		FilePath:   "",
		CreatedAt:  time.Now().UTC(),
		Deleted:    false,
	}

	if err = s.postgresClient.CreateVideo(video); err != nil {
		log.Println("[StartVideo] Failed to generate video", err)
		http.Error(w, "Failed to get video", http.StatusInternalServerError)
		return
	}

	job := redisdb.Job{
		Id:   videoId,
		Link: link.MagnetLink,
	}

	if err = s.redisClient.PublishJob(r.Context(), job); err != nil {
		log.Println("[StartVideo] Failed to publish job", err)
		http.Error(w, "Failed to get video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"video_id\": \"%s\" }", videoId)))
}

func (s *Server) serveVideo(w http.ResponseWriter, r *http.Request) {
	// videoId := mux.Vars(r)["videoId"]
	// log.Println("Getting VideoId Reader")
	// // reader := s.torrentClient.GetReader(videoId)
	// if reader == nil {
	// 	log.Println("[ServeVideo] Magnet Link Doesnt exist")
	// 	http.Error(w, "Please insert magnet link again, content doesn't exist", http.StatusNotFound)
	// 	return
	// }
	// defer reader.Close()

	rangeHeader := r.Header.Get("Range")
	log.Println("Request: %s", rangeHeader)

	// http.ServeContent(w, r, "video.mkv", time.Now(), reader)	
}

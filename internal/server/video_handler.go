package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/scythe504/webtorrent/internal"
	postgresdb "github.com/scythe504/webtorrent/internal/postgres-db"
	redisdb "github.com/scythe504/webtorrent/internal/redis-db"
)

func (s *Server) saveVideoForLater(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[StartVideo] Invalid Request body", err)
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var link struct {
		VideoId string `json:"video_id"`
	}

	if err = json.Unmarshal(body, &link); err != nil {
		log.Println("[StartVideo] Invalid Json Body", err)
		http.Error(w, "Failed to Parse JSON", http.StatusBadRequest)
		return
	}

	magnetLink := s.t.GetMagnetLink(link.VideoId)

	if magnetLink == nil {
		log.Println("[StartVideo] Could not get magnet link", err)
		http.Error(w, "Failed to get magnet link, please renter the magnet link to get the video", http.StatusNotFound)
		return
	}

	video := postgresdb.Video{
		Id:         link.VideoId,
		MagnetLink: *magnetLink,
		Status:     postgresdb.PROCESSING,
		FilePath:   "",
		CreatedAt:  time.Now().UTC(),
		Deleted:    false,
	}

	if err = s.db.CreateVideo(video); err != nil {
		log.Println("[StartVideo] Failed to generate video", err)
		http.Error(w, "Failed to get video", http.StatusInternalServerError)
		return
	}

	job := redisdb.Job{
		Id:   link.VideoId,
		Link: *magnetLink,
	}

	if err = s.rdb.PublishJob(r.Context(), job); err != nil {
		log.Println("[StartVideo] Failed to publish job", err)
		http.Error(w, "Failed to get video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"video_id\": \"%s\", \"message\": \"%s\" }", link.VideoId, "We will notify you when the video is ready to watch!")))
}

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
		MagnetLink   string `json:"magnet_link"`
		SaveForLater bool   `json:"save_for_later"`
	}

	if err = json.Unmarshal(body, &link); err != nil {
		log.Println("[StartVideo] Invalid Json Body", err)
		http.Error(w, "Failed to Parse JSON", http.StatusBadRequest)
		return
	}

	videoId := internal.RandomId()

	if link.SaveForLater {

		video := postgresdb.Video{
			Id:         videoId,
			MagnetLink: link.MagnetLink,
			Status:     postgresdb.PROCESSING,
			FilePath:   "",
			CreatedAt:  time.Now().UTC(),
			Deleted:    false,
		}

		if err = s.db.CreateVideo(video); err != nil {
			log.Println("[StartVideo] Failed to generate video", err)
			http.Error(w, "Failed to get video", http.StatusInternalServerError)
			return
		}

		job := redisdb.Job{
			Id:   videoId,
			Link: link.MagnetLink,
		}

		if err = s.rdb.PublishJob(r.Context(), job); err != nil {
			log.Println("[StartVideo] Failed to publish job", err)
			http.Error(w, "Failed to get video", http.StatusInternalServerError)
			return
		}

	} else {
		if err = s.t.AddMagnet(videoId, link.MagnetLink); err != nil {
			log.Println("[StartVideo] failed to get the magnet link")
			http.Error(w, "failed to get video", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"video_id\": \"%s\" }", videoId)))
}

func (s *Server) serveVideo(w http.ResponseWriter, r *http.Request) {
	videoId := mux.Vars(r)["videoId"]
	log.Println("Getting VideoId Reader")
	reader := s.t.GetReader(videoId)
	if reader == nil {
		log.Println("[ServeVideo] Magnet Link Doesnt exist")
		http.Error(w, "Please insert magnet link again, content doesn't exist", http.StatusNotFound)
		return
	}
	defer reader.Close()

	http.ServeContent(w, r, "video.mkv", time.Now(), reader)
}

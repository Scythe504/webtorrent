package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	postgresdb "github.com/scythe504/webtorrent/internal/postgres-db"
	redisdb "github.com/scythe504/webtorrent/internal/redis-db"
	"github.com/scythe504/webtorrent/internal/tor"
)

type Server struct {
	port           int
	torrentClient  *tor.Torrent
	redisClient    redisdb.Service
	postgresClient postgresdb.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	ctx := context.Background()
	NewServer := &Server{
		port:           port,
		torrentClient:  tor.New(),
		redisClient:    redisdb.New(ctx),
		postgresClient: postgresdb.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

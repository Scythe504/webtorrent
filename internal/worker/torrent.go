package worker

import (
	"context"
	"log"

	postgresdb "github.com/scythe504/webtorrent/internal/postgres-db"
	redisdb "github.com/scythe504/webtorrent/internal/redis-db"
	"github.com/scythe504/webtorrent/internal/storage"
	"github.com/scythe504/webtorrent/internal/tor"
)

type TorrentWorker struct {
	rdb        redisdb.Service
	postgresdb postgresdb.Service
	tor        tor.Torrent
	st         storage.Service
	jobsChan   chan redisdb.Job
	numWorker  int
	errChan    chan WorkerError
	ctx        context.Context
}

type WorkerError struct {
	JobId string
	Err   error
	Phase ErrPhase
}

type ErrPhase string

const (
	MAGNET           ErrPhase = "magnet_link_add"
	UPDATE_FAILED    ErrPhase = "update_failed"
	DOWNLOAD_FAILED  ErrPhase = "download_failed"
	BUCKET_WRITE_ERR ErrPhase = "bucket_write_err"
)

func NewTorrentWorker(worker int) *TorrentWorker {
	ctx := context.Background()

	jobsChan := make(chan redisdb.Job, worker)
	errChan := make(chan WorkerError, worker)
	tw := &TorrentWorker{
		rdb:        redisdb.New(ctx),
		postgresdb: postgresdb.New(),
		tor:        tor.New(),
		st:         storage.New(),
		jobsChan:   jobsChan,
		ctx:        ctx,
		numWorker:  worker,
		errChan:    errChan,
	}

	return tw
}

func (tw *TorrentWorker) Start(consumerName string) {
	for {
		job, err := tw.rdb.ConsumeJob(tw.ctx, consumerName)

		if err != nil {
			log.Printf("[%s] ConsumeJob error: %v\n", consumerName, err)
			continue
		}

		if job == nil {
			continue
		}

		if err := tw.postgresdb.UpdateStatus(postgresdb.DOWNLOADED, job.Id, nil); err != nil {
			log.Printf("[%s] UpdateStatus DOWNLOADING failed: %v\n", consumerName, err)
			continue
		}

		tw.jobsChan <- *job
	}
}

func (tw *TorrentWorker) DownloadWorker(i int) {
	for job := range tw.jobsChan {
		tw.processJob(job)
	}
}

func (tw *TorrentWorker) processJob(job redisdb.Job) {
	err := tw.tor.AddMagnet(job.Id, job.Link)

	if err != nil {
		tw.errChan <- WorkerError{
			JobId: job.Id,
			Err:   err,
			Phase: MAGNET,
		}

		return
	}

	reader := tw.tor.GetReader(job.Id)
	if reader == nil {
		tw.errChan <- WorkerError{
			JobId: job.Id,
			Err:   err,
			Phase: DOWNLOAD_FAILED,
		}

		return
	}

	filepath, err := tw.tor.GetFileName(job.Id)

	if err != nil {
		tw.errChan <- WorkerError{
			JobId: job.Id,
			Err:   err,
			Phase: BUCKET_WRITE_ERR,
		}

		return
	}

	err = tw.st.SaveForLater(job.Id, reader)

	if err != nil {
		tw.errChan <- WorkerError{
			JobId: job.Id,
			Err:   err,
			Phase: BUCKET_WRITE_ERR,
		}

		return
	}

	err = tw.postgresdb.UpdateStatus(postgresdb.DOWNLOADED, job.Id, &filepath)
	if err != nil {
		tw.errChan <- WorkerError{
			JobId: job.Id,
			Err:   err,
			Phase: UPDATE_FAILED,
		}
		return
	}

	// TODO: Cleanup torrent connection
}

func (tw *TorrentWorker) HandleErrors() {
	for we := range tw.errChan {
		log.Printf("[ERROR] JobId: %s, Phase: %s, Error: %v\n", we.JobId, we.Phase, we.Err)

		// Update DB status to FAILED
		if err := tw.postgresdb.UpdateStatus(postgresdb.FAILED, we.JobId, nil); err != nil {
			log.Printf("[ERROR] Failed to update status to FAILED for jobId %s: %v\n", we.JobId, err)
		}
	}
}

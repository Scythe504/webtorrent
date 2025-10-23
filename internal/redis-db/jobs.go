package redisdb

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Job struct {
	Id   string `json:"id" redis:"id" db:"id"`
	Link string `json:"link" redis:"link" db:"magnet_link"`
}

func (s *service) PublishJob(ctx context.Context, job Job) error {
	_, err := s.db.XAdd(ctx, &redis.XAddArgs{
		Stream: "job:magnet-link",
		Values: map[string]any{
			"id":   job.Id,
			"link": job.Link,
		},
	}).Result()

	if err != nil {
		return err
	}

	return nil
}

func (s *service) ConsumeJob(ctx context.Context, consumerName string) (*Job, error) {
	log.Println("[DEBUG]", ctx, consumerName, s.db)
	res, err := s.db.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{"job:magnet-link", ">"},
		Count:    1,
		Block:    0,
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(res) == 0 || len(res[0].Messages) == 0 {
		return nil, nil
	}

	msg := res[0].Messages[0]

	Job := Job{
		Id:   msg.Values["id"].(string),
		Link: msg.Values["link"].(string),
	}

	s.db.XAck(ctx, "job:magnet-link", consumerGroup, msg.ID)

	return &Job, nil
}

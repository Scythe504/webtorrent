package redisdb

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type MagnetLink struct {
	Id   string `json:"id" redis:"id"`
	Link string `json:"link" redis:"link"`
}

func (s *service) PublishJob(ctx context.Context, job MagnetLink) error {
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

func (s *service) ConsumeJob(ctx context.Context, consumerName string) (*MagnetLink, error) {
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

	magnetLink := MagnetLink{
		Id:   msg.Values["id"].(string),
		Link: msg.Values["link"].(string),
	}

	s.db.XAck(ctx, "job:magnet-link", consumerGroup, msg.ID)

	return &magnetLink, nil
}

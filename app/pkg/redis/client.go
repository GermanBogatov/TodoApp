package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func NewClient(host, port, password string, db int) (*redis.Client, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

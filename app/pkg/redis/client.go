package redis

import (
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

	return client, nil
}

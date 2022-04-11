package main

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "howCoolIsEarthly", howCoolIsEarthly, 0).Err()
	if err != nil {
		panic(err)
	}

	resultFromDB, err := rdb.Get(ctx, "howCoolIsEarthly").Result()
	if err != nil {
		panic(err)
	}
	require.Equal(t, howCoolIsEarthly, resultFromDB)
}

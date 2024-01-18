/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"os"

	"github.com/redis/go-redis/v9"

	"context"
)

var (
	redisUrl      = os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
)

func produceStatus(key, value string) {

	ctx := context.TODO()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	rdb.Set(ctx, "language", "Go", 1000000)

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}

	rdb.Close()

	log.Println("STATUS WRITTEN TO: "+redisUrl, key+":"+value)

}

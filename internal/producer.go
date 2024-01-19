/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
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

	log.Println("STATUS WRITTEN TO: "+redisUrl, key+": "+value)

}

func checkStageStatus(pipelineRunLabels map[string]string) {

	fmt.Println(pipelineRunLabels)
	// ctx := context.TODO()
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     redisUrl,
	// 	Password: redisPassword, // no password set
	// 	DB:       0,             // use default DB
	// })

}

// pipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)
// log.Info("ALL PIPELEINRUNS OF THIS STAGE: ", pipelineRuns)

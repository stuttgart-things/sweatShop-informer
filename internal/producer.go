/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"

	goredis "github.com/redis/go-redis/v9"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"context"
)

var (
	redisUrl      = os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisClient   = goredis.NewClient(&goredis.Options{Addr: redisUrl, Password: redisPassword, DB: 0})
)

func produceStatus(key, value string) {

	ctx := context.TODO()
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     redisUrl,
	// 	Password: redisPassword, // no password set
	// 	DB:       0,             // use default DB
	// })

	redisClient.Set(ctx, "language", "Go", 1000000)

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}

	redisClient.Close()

	log.Println("STATUS WRITTEN TO: "+redisUrl, key+": "+value)

}

func checkStageStatus(pipelineRunLabels map[string]string) {

	fmt.Println(pipelineRunLabels)
	stageKey := pipelineRunLabels["stagetime/date"] + "-" + pipelineRunLabels["stagetime/commit"] + "-" + pipelineRunLabels["stagetime/stage"]

	stagePipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageKey)

	fmt.Println("ALL PIPELEINRUNS OF THIS STAGE: ", stagePipelineRuns)

}

// pipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)
// log.Info("ALL PIPELEINRUNS OF THIS STAGE: ", pipelineRuns)

// │ map[stagetime/author:patrick-hermann-sva stagetime/commit:3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0 stagetime/repo:stuttgart-things stagetime/stage:0 tekton.dev/pipeline:st-0-simu │

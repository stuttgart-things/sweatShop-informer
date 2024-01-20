/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/stageTime-server/server"

	goredis "github.com/redis/go-redis/v9"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"context"
)

var (
	redisUrl         = os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT")
	redisPassword    = os.Getenv("REDIS_PASSWORD")
	redisClient      = goredis.NewClient(&goredis.Options{Addr: redisUrl, Password: redisPassword, DB: 0})
	redisJSONHandler = rejson.NewReJSONHandler()
)

func produceStatus(key, value string) {

	ctx := context.TODO()
	rc := goredis.NewClient(&goredis.Options{
		Addr:     redisUrl,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	rc.Set(ctx, "language", "Go", 1000000)

	err := rc.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}

	rc.Close()

	log.Println("STATUS WRITTEN TO: "+redisUrl, key+": "+value)

}

func checkStageStatus(pipelineRunLabels map[string]string) {

	redisJSONHandler.SetGoRedisClient(redisClient)

	fmt.Println(pipelineRunLabels)
	stageKey := pipelineRunLabels["stagetime/date"] + "-" + pipelineRunLabels["stagetime/commit"] + "-" + pipelineRunLabels["stagetime/stage"]

	stagePipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageKey)

	fmt.Println("ALL PIPELEINRUNS OF THIS STAGE: ", stagePipelineRuns)
	// sthingsCli.AddValueToRedisSet(redisClient, prInformation["stagetime/date"]+"-"+prInformation["stagetime/commit"]+"-"+prInformation["stagetime/stage"], prInformation["name"])

	stageStatus := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunLabels["name"]+"-status")
	fmt.Println(stageStatus)

	stageStatusFromRedis := server.StageStatus{}
	err := json.Unmarshal(stageStatus, &stageStatusFromRedis)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
	}

	stageStatusFromRedis.Status = "TESTED"

	server.PrintTable(stageStatusFromRedis)

	// IF STOP FOUND MARK REVISIONRUN AS FAILED
	// IF ALL CONTINUE MARK STAGE AS SUCCESSFULL

}

// pipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)
// log.Info("ALL PIPELEINRUNS OF THIS STAGE: ", pipelineRuns)

// │ map[stagetime/author:patrick-hermann-sva stagetime/commit:3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0 stagetime/repo:stuttgart-things stagetime/stage:0 tekton.dev/pipeline:st-0-simu │

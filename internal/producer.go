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
)

var (
	redisUrl         = os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT")
	redisPassword    = os.Getenv("REDIS_PASSWORD")
	redisClient      = goredis.NewClient(&goredis.Options{Addr: redisUrl, Password: redisPassword, DB: 0})
	redisJSONHandler = rejson.NewReJSONHandler()
)

func setPipelineRunStatus(pipelineRunLabels map[string]string) {

	jsonKey := pipelineRunLabels["name"] + "-status"
	redisJSONHandler.SetGoRedisClient(redisClient)

	// PIPELINERUN STATUS
	pipelineRunStatus := sthingsCli.GetRedisJSON(redisJSONHandler, jsonKey)

	pipelineRunStatusFromRedis := server.PipelineRunStatus{}
	err := json.Unmarshal(pipelineRunStatus, &pipelineRunStatusFromRedis)
	if err != nil {
		log.Fatalf("FAILED TO JSON UNMARSHAL")
	}

	pipelineRunStatusFromRedis.Status = pipelineRunLabels["status"]

	server.PrintTable(pipelineRunStatusFromRedis)

	sthingsCli.SetRedisJSON(redisJSONHandler, pipelineRunStatusFromRedis, jsonKey)

}

func setStageStatus(pipelineRunLabels map[string]string) {

	jsonKey := pipelineRunLabels["stagetime/commit"] + pipelineRunLabels["stagetime/stage"]
	redisJSONHandler.SetGoRedisClient(redisClient)

	// STAGE STATUS
	stageStatus := sthingsCli.GetRedisJSON(redisJSONHandler, jsonKey)

	stageStatusFromRedis := server.StageStatus{}
	err := json.Unmarshal(stageStatus, &stageStatusFromRedis)
	if err != nil {
		log.Fatalf("FAILED TO JSON UNMARSHAL")
	}

	stageStatusFromRedis.Status = pipelineRunLabels["status"]

	server.PrintTable(stageStatusFromRedis)

	sthingsCli.SetRedisJSON(redisJSONHandler, stageStatusFromRedis, jsonKey)

}

func checkStageStatus(pipelineRunLabels map[string]string) {

	fmt.Println("LABLES", pipelineRunLabels)
	fmt.Println(pipelineRunLabels["name"] + "-status")

	redisJSONHandler.SetGoRedisClient(redisClient)

	fmt.Println(pipelineRunLabels)
	stageKey := pipelineRunLabels["stagetime/date"] + "-" + pipelineRunLabels["stagetime/commit"] + "-" + pipelineRunLabels["stagetime/stage"]

	stagePipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageKey)

	fmt.Println("ALL PIPELEINRUNS OF THIS STAGE: ", stagePipelineRuns)
	// sthingsCli.AddValueToRedisSet(redisClient, prInformation["stagetime/date"]+"-"+prInformation["stagetime/commit"]+"-"+prInformation["stagetime/stage"], prInformation["name"])

	// STAGE STATUS
	stageStatus := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunLabels["stagetime/commit"]+pipelineRunLabels["stagetime/stage"])

	stageStatusFromRedis := server.StageStatus{}
	err := json.Unmarshal(stageStatus, &stageStatusFromRedis)
	if err != nil {
		log.Fatalf("FAILED TO JSON UNMARSHAL")
	}

	stageStatusFromRedis.Status = "TESTED"

	server.PrintTable(stageStatusFromRedis)

	// IF STOP FOUND MARK REVISIONRUN AS FAILED
	// IF ALL CONTINUE MARK STAGE AS SUCCESSFULL

}

// pipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)
// log.Info("ALL PIPELEINRUNS OF THIS STAGE: ", pipelineRuns)

// │ map[stagetime/author:patrick-hermann-sva stagetime/commit:3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0 stagetime/repo:stuttgart-things stagetime/stage:0 tekton.dev/pipeline:st-0-simu │

/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"

	"github.com/nitishm/go-rejson/v4"
	server "github.com/stuttgart-things/stageTime-server/server"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

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

	pipelineRunStatusFromRedis := server.GetPipelineRunStatus(jsonKey, redisJSONHandler)

	pipelineRunStatusFromRedis.Status = pipelineRunLabels["annotation"]

	server.PrintTable(pipelineRunStatusFromRedis)

	sthingsCli.SetRedisJSON(redisJSONHandler, pipelineRunStatusFromRedis, jsonKey)

}

func setStageStatus(pipelineRunLabels map[string]string) {

	jsonKey := pipelineRunLabels["stagetime/commit"] + pipelineRunLabels["stagetime/stage"]
	redisJSONHandler.SetGoRedisClient(redisClient)

	stageStatusFromRedis := server.GetStageStatus(jsonKey, redisJSONHandler)

	stageStatusFromRedis.Status = pipelineRunLabels["status"]

	server.PrintTable(stageStatusFromRedis)

	stagePipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)
	fmt.Println("ALL PRS: ", stagePipelineRuns)

	var prStatus []string

	for _, name := range stagePipelineRuns {
		fmt.Println(name)
		pipelineRunStatusFromRedis := server.GetPipelineRunStatus(name+"-status", redisJSONHandler)
		prStatus = append(prStatus, fmt.Sprintln(pipelineRunStatusFromRedis))
	}
	fmt.Println("STTTAUUS", prStatus)

	if sthingsBase.CheckForStringInSlice(prStatus, "STOP") {
		fmt.Println("STAGE IS DEAD", jsonKey)
	}

}

// func checkStageStatus(pipelineRunLabels map[string]string) {

// 	fmt.Println("LABLES", pipelineRunLabels)
// 	fmt.Println(pipelineRunLabels["name"] + "-status")

// 	redisJSONHandler.SetGoRedisClient(redisClient)

// 	fmt.Println(pipelineRunLabels)
// 	stageKey := pipelineRunLabels["stagetime/date"] + "-" + pipelineRunLabels["stagetime/commit"] + "-" + pipelineRunLabels["stagetime/stage"]

// 	stagePipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageKey)

// 	fmt.Println("ALL PIPELEINRUNS OF THIS STAGE: ", stagePipelineRuns)
// 	// sthingsCli.AddValueToRedisSet(redisClient, prInformation["stagetime/date"]+"-"+prInformation["stagetime/commit"]+"-"+prInformation["stagetime/stage"], prInformation["name"])

// 	// STAGE STATUS
// 	stageStatus := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunLabels["stagetime/commit"]+pipelineRunLabels["stagetime/stage"])

// 	stageStatusFromRedis := server.StageStatus{}
// 	err := json.Unmarshal(stageStatus, &stageStatusFromRedis)
// 	if err != nil {
// 		log.Fatalf("FAILED TO JSON UNMARSHAL")
// 	}

// 	stageStatusFromRedis.Status = "TESTED"

// 	server.PrintTable(stageStatusFromRedis)

// 	// IF STOP FOUND MARK REVISIONRUN AS FAILED
// 	// IF ALL CONTINUE MARK STAGE AS SUCCESSFULL

// }

// pipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)
// log.Info("ALL PIPELEINRUNS OF THIS STAGE: ", pipelineRuns)

// │ map[stagetime/author:patrick-hermann-sva stagetime/commit:3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0 stagetime/repo:stuttgart-things stagetime/stage:0 tekton.dev/pipeline:st-0-simu │

// func GetPipelineRunStatus(jsonKey string) server.PipelineRunStatus {

// 	pipelineRunStatusJson := sthingsCli.GetRedisJSON(redisJSONHandler, jsonKey)
// 	pipelineRunStatus := server.PipelineRunStatus{}

// 	err := json.Unmarshal(pipelineRunStatusJson, &pipelineRunStatus)
// 	if err != nil {
// 		fmt.Println(err)
// 		log.Fatalf("FAILED TO JSON UNMARSHAL")
// 	}

// 	return pipelineRunStatus
// }

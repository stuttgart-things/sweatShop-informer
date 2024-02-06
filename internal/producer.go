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

	var prStatus []string
	jsonKey := pipelineRunLabels["stagetime/commit"] + pipelineRunLabels["stagetime/stage"]
	redisJSONHandler.SetGoRedisClient(redisClient)

	// GET CURRENT STAGE STATUS
	stageStatusFromRedis := server.GetStageStatus(jsonKey, redisJSONHandler)
	stagePipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)

	// GET CURRENT PIPELINERUN STATUS
	for _, name := range stagePipelineRuns {
		// fmt.Println(name)
		pipelineRunStatusFromRedis := server.GetPipelineRunStatus(name+"-status", redisJSONHandler)
		prStatus = append(prStatus, fmt.Sprintln(pipelineRunStatusFromRedis))
	}

	// CHECK IF STAGE IS SUCCESSFULL, FAILED OR STILL RUNNING
	if sthingsBase.CheckForStringInSlice(prStatus, "STOP") {
		fmt.Println("STAGE IS DEAD", jsonKey)
	}

	// SET STAGE STATUS
	stageStatusFromRedis.Status = pipelineRunLabels["status"]

	fmt.Println("STAGE STATUS: ", pipelineRunLabels["status"])
	// PRINT UPDATED STAGE STATUS
	server.PrintTable(stageStatusFromRedis)

	fmt.Println("ALL PRS: ", stagePipelineRuns)

	// func SendStageToMessageQueue(stageID string) {

	// 	streamValues := map[string]interface{}{
	// 		"stage": stageID,
	// 	}

	// 	sthingsCli.EnqueueDataInRedisStreams(redisAddress+":"+redisPort, redisPassword, redisStream, streamValues)
	// 	fmt.Println("STREAM", redisStream)
	// 	fmt.Println("VALUES", streamValues)
	// }

}

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

	// PRINT ALL PRS FROM STAGE
	fmt.Println("ALL PRS: ", stagePipelineRuns)

	// SET STAGE STATUS
	stageStatusFromRedis.Status = pipelineRunLabels["status"]
	// PRINT UPDATED STAGE STATUS
	server.PrintTable(stageStatusFromRedis)

	if pipelineRunLabels["status"] == "SUCCEEDED" {

		revisionRunStatus := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunLabels["stagetime/commit"]+"-status")
		revisionRunFromRedis := server.RevisionRunStatus{}

		err := json.Unmarshal(revisionRunStatus, &revisionRunFromRedis)
		if err != nil {
			log.Fatalf("FAILED TO JSON UNMARSHAL REVISIONRUN STATUS")
		}
		server.PrintTable(revisionRunFromRedis)

		countCurrentStage := sthingsBase.ConvertStringToInteger(pipelineRunLabels["stagetime/stage"])

		fmt.Println("CURRENT STAGE:", countCurrentStage)
		fmt.Println("COUNT STAGES:", revisionRunFromRedis.CountStages)

		if countCurrentStage > revisionRunFromRedis.CountStages {
			fmt.Println("NEXT STAGE LETS GOOO")
		} else {
			fmt.Println("REVISION RUN FINISHED")
		}

	}

	// func SendStageToMessageQueue(stageID string) {

	// 	streamValues := map[string]interface{}{
	// 		"stage": stageID,
	// 	}

	// 	sthingsCli.EnqueueDataInRedisStreams(redisAddress+":"+redisPort, redisPassword, redisStream, streamValues)
	// 	fmt.Println("STREAM", redisStream)
	// 	fmt.Println("VALUES", streamValues)
	// }

}

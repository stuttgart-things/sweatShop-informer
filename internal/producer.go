/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"
	"strings"

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
	redisStream      = os.Getenv("REDIS_STREAM")
)

func getPipelineRunStatus(pipelineRunLabels map[string]string) {

	jsonKey := pipelineRunLabels["name"] + "-status"
	redisJSONHandler.SetGoRedisClient(redisClient)

	pipelineRunStatusFromRedis := server.GetPipelineRunStatus(jsonKey, redisJSONHandler)

	pipelineRunStatusFromRedis.Status = pipelineRunLabels["annotation"]

	server.PrintTable(pipelineRunStatusFromRedis)

	sthingsCli.SetRedisJSON(redisJSONHandler, pipelineRunStatusFromRedis, jsonKey)

}

func setPipelineRunStatus(pipelineRunLabels map[string]string) {

	jsonKey := pipelineRunLabels["name"] + "-status"
	redisJSONHandler.SetGoRedisClient(redisClient)

	pipelineRunStatusFromRedis := server.GetPipelineRunStatus(jsonKey, redisJSONHandler)

	server.PrintTable(pipelineRunStatusFromRedis)

	// ONLY SET STATUS IF NO PREVIOUS STATUS WAS SET TO THE PIPELINERUN
	if !strings.Contains(pipelineRunStatusFromRedis.Status, "STOP") || !strings.Contains(pipelineRunStatusFromRedis.Status, "CONTINUE") {

		pipelineRunStatusFromRedis.Status = pipelineRunLabels["annotation"]

		server.PrintTable(pipelineRunStatusFromRedis)

		sthingsCli.SetRedisJSON(redisJSONHandler, pipelineRunStatusFromRedis, jsonKey)
	}
}

func setStageStatus(pipelineRunLabels map[string]string) (stageFinished, continueRun bool, stageID string, currentStageNumber int) {

	// VARS DECLARATION
	var prStatus []string
	var stageStatus string
	continueRun = true

	stageID = pipelineRunLabels["stagetime/commit"] + pipelineRunLabels["stagetime/stage"]
	redisJSONHandler.SetGoRedisClient(redisClient)
	currentStageNumber = sthingsBase.ConvertStringToInteger(pipelineRunLabels["stagetime/stage"])

	// GET CURRENT STAGE STATUS
	stageStatusFromRedis := server.GetStageStatus(stageID, redisJSONHandler)
	stagePipelineRuns := sthingsCli.GetValuesFromRedisSet(redisClient, stageStatusFromRedis.StageID)

	// CHECK STATUS OF ALL PRS OF STAGE
	for _, name := range stagePipelineRuns {
		pipelineRunStatusFromRedis := server.GetPipelineRunStatus(name+"-status", redisJSONHandler)
		prStatus = append(prStatus, fmt.Sprintln(pipelineRunStatusFromRedis))
	}

	// CHECK IF STAGE IS SUCCESSFULL, FAILED OR STILL RUNNING
	if sthingsBase.CheckForStringInSlice(prStatus, "STOP") {
		stageStatus = "STAGE FAILED"
		fmt.Println("STAGE FAILED", stageID)
		stageFinished = true
		continueRun = false

	} else {
		stageStatus = "SUCCESSFUL"
		fmt.Println("STAGE SUCCEEDED", stageID)
		stageFinished = true
		continueRun = true
	}

	server.SetStageStatusInRedis(redisJSONHandler, stageID+"-status", stageStatus, stageStatusFromRedis, true)

	return
}

func setRevisionRunStatus(revisionRunID, StageID string, succeeded bool) {

	var revisionRunStatus string
	revisionRunFromRedis := server.GetRevisionRunFromRedis(redisJSONHandler, revisionRunID+"-status", true)

	if succeeded {
		revisionRunStatus = "REVISIONRUN SUCCEEDED AT STAGE " + StageID
	} else {
		revisionRunStatus = "REVISIONRUN FAILED AT STAGE " + StageID
	}

	server.SetRevisionRunStatusInRedis(redisJSONHandler, revisionRunID+"-status", revisionRunStatus, revisionRunFromRedis, true)

}

func checkForNextStage(stageID, revisionRunID string, nextStage int) bool {

	revisionRunFromRedis := server.GetRevisionRunFromRedis(redisJSONHandler, revisionRunID+"-status", true)

	fmt.Println("NEXT STAGE:", nextStage)
	fmt.Println("COUNT STAGES:", revisionRunFromRedis.CountStages)

	if nextStage > revisionRunFromRedis.CountStages {
		fmt.Println("NEXT STAGE LETS GO")
		return true

	} else {
		fmt.Println("NO NEXT STAGE")
		return false

	}

}

// 	// SET STAGE STATUS
// 	stageStatusFromRedis.Status = pipelineRunLabels["status"]
// 	server.PrintTable(stageStatusFromRedis)

// 	// GET REVISIONRUN STATUS
// 	revisionRunStatus := sthingsCli.GetRedisJSON(redisJSONHandler, pipelineRunLabels["stagetime/commit"]+"-status")
// 	revisionRunFromRedis := server.RevisionRunStatus{}

// 	err := json.Unmarshal(revisionRunStatus, &revisionRunFromRedis)
// 	if err != nil {
// 		log.Fatalf("FAILED TO JSON UNMARSHAL REVISIONRUN STATUS")
// 	}

// 	server.PrintTable(revisionRunFromRedis)

// 	// CALL SERVER GET REVISIONRUN STATUS
// 	// IF STATUS NOT ALREADY SET BY INFORMER
// 	// SET STATUS BY INFORMER

// 	countCurrentStage := sthingsBase.ConvertStringToInteger(pipelineRunLabels["stagetime/stage"])

// 	fmt.Println("CURRENT STAGE:", countCurrentStage)
// 	fmt.Println("COUNT STAGES:", revisionRunFromRedis.CountStages)

// if revisionRunFromRedis.CountStages > countCurrentStage {

// 	if pipelineRunLabels["status"] != "SUCCEEDED" {

// 	fmt.Println("NEXT STAGE LETS GOOO")

// 	currentStageID := stageStatusFromRedis.StageID
// 	nextStageIDBuilder := strings.LastIndex(currentStageID, "-")

// 	nextStageID := replaceLastOccurrenceInSubstring(stageStatusFromRedis.StageID[:nextStageIDBuilder]+"+"+sthingsBase.ConvertIntegerToString(countCurrentStage+1), "-", "+")

// 	fmt.Println("NEXT STAGE!?", nextStageID)
// 	SendStageToMessageQueue(nextStageID)

// 	} else {
// 		fmt.Println("REVISION RUN FINISHED", pipelineRunLabels["stagetime/stage"])
// 		server.SetRevisionRunStatusInRedis(redisJSONHandler, pipelineRunLabels["stagetime/commit"]+"-status", "REVISIONRUN SUCCESSFUL", revisionRunFromRedis, true)
// 	}

// 	// }

func SendStageToMessageQueue(stageID string) {

	streamValues := map[string]interface{}{
		"stage": stageID,
	}

	sthingsCli.EnqueueDataInRedisStreams(redisUrl, redisPassword, redisStream, streamValues)
	fmt.Println("STREAM", redisStream)
	fmt.Println("VALUES", streamValues)
}

func replaceLastOccurrenceInSubstring(subString, searchFor, replaceWith string) (x2 string) {
	i := strings.LastIndex(subString, searchFor)
	if i == -1 {
		return subString
	}
	return subString[:i] + replaceWith + subString[i+len(searchFor):]
}

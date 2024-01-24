/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

func verifyJobCompletionStatus(prStatus, regexPattern string) (jobStatus string) {

	jobStatusMessage, _ := sthingsBase.GetRegexSubMatch(prStatus, regexPattern)
	fmt.Println(jobStatusMessage)

	switch jobStatusMessage {

	case "True":
		jobStatus = "SUCCEEDED"
	case "False":
		jobStatus = "FAILED"
	default:
		jobStatus = "RUNNING"
	}

	return

}

func verifyInformerStatus(kind, function string, obj interface{}) {

	switch kind {

	case "pipelineruns":
		annotation := ":CONTINUE"

		pipelineRun := CreatePipelineRunFromUnstructuredObj(obj)
		pipelineRunStatusMessage := verifyJobCompletionStatus(fmt.Sprintln(pipelineRun.Status), `Succeeded\s(\w+)`)

		pipelineRunLabels := pipelineRun.Labels
		pipelineRunLabels["name"] = pipelineRun.Name
		pipelineRunLabels["status"] = pipelineRunStatusMessage

		pipelineRunAnnotations := pipelineRun.Annotations

		if pipelineRunAnnotations["canfail"] == "false" && pipelineRunStatusMessage == "FAILED" {
			annotation = "STOP"
			fmt.Println(pipelineRunLabels["name"] + " can fail: " + pipelineRunAnnotations["canfail"] + "and annotation is: " + annotation)
		}
		pipelineRunLabels["annotation"] = annotation

		setPipelineRunStatus(pipelineRunLabels)
		setStageStatus(pipelineRunLabels)
	}
}

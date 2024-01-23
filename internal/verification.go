/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import "fmt"

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
			annotation = ":STOP"
		}
		pipelineRunLabels["annotation"] = annotation

		setPipelineRunStatus(pipelineRunLabels)
		setStageStatus(pipelineRunLabels)
	}
}

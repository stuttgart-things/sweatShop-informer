/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

func CreatePipelineRunFromUnstructuredObj(obj interface{}) (pr *tekton.PipelineRun) {

	pr = new(tekton.PipelineRun)

	createdUnstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		log.Fatal(err)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(createdUnstructuredObj, &pr)
	if err != nil {
		log.Fatal(err)
	}

	return

}

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
		pipelineRun := CreatePipelineRunFromUnstructuredObj(obj)
		fmt.Println(pipelineRun)
		fmt.Println("STATUS!", pipelineRun.Status)

		pipelineRunStatusMessage := verifyJobCompletionStatus(fmt.Sprintln(pipelineRun.Status), `Succeeded\s(\w+)`)
		fmt.Println(pipelineRunStatusMessage)
		produceStatus("job-"+pipelineRun.Name, pipelineRunStatusMessage)

	}
}

/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
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

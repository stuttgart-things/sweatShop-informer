/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	batchv1 "k8s.io/api/batch/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

func CreateJobFromUnstructuredObj(obj interface{}) (job *batchv1.Job) {

	job = new(batchv1.Job)

	createdUnstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		log.Fatal(err)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(createdUnstructuredObj, &job)
	if err != nil {
		log.Fatal(err)
	}

	return

}

func verifyJobCompletionStatus(prStatus string) (jobStatus string) {

	jobStatusMessage, _ := sthingsBase.GetRegexSubMatch(prStatus, `Complete\s(\w+)`)

	if jobStatusMessage != "True" {
		jobStatus = "running"
	} else {
		jobStatus = "finished"
	}

	return

}

func verifyInformerStatus(kind, function string, obj interface{}) {

	switch kind {

	case "jobs":
		job := CreateJobFromUnstructuredObj(obj)
		log.Println("job " + function + ": " + job.Name)
		jobStatusMessage := verifyJobCompletionStatus(fmt.Sprintln(job.Status))
		produceStatus(job.Name, jobStatusMessage)
	case "configmaps":
		fmt.Println("FOUND CONFIGMAP!")
	}

}

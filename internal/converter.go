/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
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

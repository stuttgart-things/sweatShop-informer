/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"testing"
)

func TestVerifyJobCompletionStatus(t *testing.T) {

	prStatus := "{{0 [{Succeeded False  {2024-01-18 12:28:56 +0000 UTC} Failed Tasks Completed: 2 (Failed: 1, Cancelled 0)"
	regexPattern := `Succeeded\s(\w+)`
	jobStatus := verifyJobCompletionStatus(prStatus, regexPattern)
	fmt.Println(jobStatus)
}

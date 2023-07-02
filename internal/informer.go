/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/dynamic"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

func InformResoureceStatus() {

	clusterConfig, _ := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))
	clusterClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(clusterClient)

}

/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"k8s.io/client-go/dynamic/dynamicinformer"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

var (
	wg             sync.WaitGroup
	log            = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	logfilePath    = "/tmp/stageTime-informer.log"
	namespace      = os.Getenv("INFORMING_NAMESPACE")
	informingKinds = os.Getenv("INFORMING_KINDS") // jobs
	kinds          = strings.Split(informingKinds, ";")
	apiResource    = map[string]map[string]string{
		"jobs": {
			"group":   "batch",
			"version": "v1",
		},
		"configmaps": {
			"group":   "",
			"version": "v1",
		},
	}
)

func InformResoureceStatus() {

	// STATUS OUTPUT
	log.Println("INFORMING_KINDS", kinds)
	log.Println("INFORMING_NAMESPACE", namespace)

	// GET CLUSTER CONFIG
	clusterConfig, _ := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))
	clusterClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// START INFORMING
	for k := range kinds {

		wg.Add(1)

		kind := kinds[k]
		log.Println("start informing for", kind)
		log.Println("group:", apiResource[kind]["group"])
		log.Println("version:", apiResource[kind]["version"])

		go func() {
			defer wg.Done()

			log.Println("namespace:", namespace)

			resourceDefinition := schema.GroupVersionResource{Group: apiResource[kind]["group"], Version: apiResource[kind]["version"], Resource: kind}
			factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient, time.Minute, namespace, nil)
			informer := factory.ForResource(resourceDefinition).Informer()

			mux := &sync.RWMutex{}
			synced := false
			informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					mux.RLock()
					defer mux.RUnlock()
					if !synced {
						return
					}

					verifyInformerStatus(kind, "added", obj)

				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					mux.RLock()
					defer mux.RUnlock()

					if !synced {
						return
					}

					verifyInformerStatus(kind, "update", newObj)

				},
				DeleteFunc: func(obj interface{}) {
					mux.RLock()
					defer mux.RUnlock()

					if !synced {
						return
					}

					verifyInformerStatus(kind, "deleted", obj)
				},
			})

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			go informer.Run(ctx.Done())

			isSynced := cache.WaitForCacheSync(ctx.Done(), informer.HasSynced)
			mux.Lock()
			synced = isSynced
			mux.Unlock()

			if !isSynced {
				log.Fatal("failed to sync")
			}

			<-ctx.Done()

		}()

	}

	wg.Wait()

}

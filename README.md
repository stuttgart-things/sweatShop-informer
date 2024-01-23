# stuttgart-things/stageTime-informer

informs dynamic of stageTime resource status

## DEV-TASKS

```bash
task --list: Available tasks for this project
* build:               Build the app
* build-image:         Build image
* git-push:            Commit & push the module
* lint:                Lint code
* package:             Update Chart.yaml and package archive
* push:                Push to registry
* run:                 Run app
* run-container:       Run container
* tag:                 Commit, push & tag the module
* test:                Test code
* vcluster:            Test deploy on vcluster
```

## TEST SERVICE LOCALLY (OUTSIDE CLUSTER)

<details><summary><b>START CONSUMER</b></summary>

```
export KUBECONFIG=~/.kube/dev11
export INFORMING_KINDS="jobs;configmaps"
export INFORMING_NAMESPACE=machine-shop-packer
export REDIS_PASSWORD=<SET-ME>
export REDIS_SERVER=redis-pve.labul.sva.de
export REDIS_PORT=6379
task run
```

</details>

## LICENSE

<details><summary><b>APACHE 2.0</b></summary>

Copyright 2023 patrick hermann.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

</details>

Author Information
------------------
Patrick Hermann, stuttgart-things 07/2023

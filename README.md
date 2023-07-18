# stuttgart-things/sweatShop-informer

distributes dynamic kubernetes resource status

## DEPLOY TO CLUSTER

<details><summary><b>REDIS</b></summary>

</details>

<details><summary><b>DEPLOYMENT</b></summary>

</details>


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

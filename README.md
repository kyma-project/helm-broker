# Helm Broker

[![Go Report Card](https://goreportcard.com/badge/github.com/kyma-project/helm-broker)](https://goreportcard.com/report/github.com/kyma-project/helm-broker)
[![Sourcegraph](https://sourcegraph.com/github.com/kyma-project/helm-broker/-/badge.svg)](https://sourcegraph.com/github.com/kyma-project/helm-broker?badge)

## Overview

Helm Broker is a [Service Broker](https://kyma-project.io/docs/master/components/service-catalog/#overview-service-brokers) that exposes Helm charts as Service Classes in [Service Catalog](https://kyma-project.io/docs/master/components/service-catalog/#overview-service-catalog). To do so, Helm Broker uses the concept of addons. An addon is an abstraction layer over a Helm chart which provides all information required to convert the chart into a Service Class.

You can install Helm Broker either as a standalone project, or as part of [Kyma](https://kyma-project.io/). In Kyma, you can use addons to install the following Service Brokers:

* Azure Service Broker
* AWS Service Broker
* GCP Service Broker

>**NOTE:** Starting from Kyma 2.0, Helm Broker will no longer be supported.

To see all addons that Helm Broker provides, go to the [`addons`](https://github.com/kyma-project/addons) repository.

Helm Broker implements the [Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/profile.md#service-metadata) (OSB API). To be compliant with the Service Catalog version used in Kyma, the Helm Broker supports only the following OSB API versions:
- v2.13
- v2.12
- v2.11

> **NOTE:** The Helm Broker does not implement the OSB API update operation.


## Installation

Follow these steps to install Helm Broker locally.

### Prerequisites

* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) 1.16
* [Helm CLI](https://github.com/kubernetes/helm#install) 3.2.0
* [Docker](https://docs.docker.com/install/) 19.03
* [Kind](https://github.com/kubernetes-sigs/kind#installation-and-usage) 0.5

>**NOTE:** For non-local installation, use Kubernetes v1.15.

### Install Helm Broker with Service Catalog

ead [this](#install-helm-broker-from-chart) section to learn how to install the Helm Broker from the chart together with the Service Catalog

To run the Helm Broker, you need a Kubernetes cluster with Service Catalog. Run the `./hack/run-dev-kind.sh` script, or follow these steps to set up the Helm Broker on Kind with all necessary dependencies:

1. Create a local cluster on Kind:
```bash
kind create cluster
```

3. Install Service Catalog as a Helm chart:

```bash
helm repo add svc-cat https://kubernetes-sigs.github.io/service-catalog
helm install catalog svc-cat/catalog --namespace catalog --set asyncBindingOperationsEnabled=true
```

4. Clone the Helm Broker repository:
```bash
git clone git@github.com:kyma-project/helm-broker.git
```

5. Install the Helm Broker chart from the cloned repository:
```bash
helm install charts/helm-broker helm-broker --namespace helm-broker
```

### Install Helm Broker manually

Rand [this](#install-helm-broker-manually) section to learn how to install the Helm Broker manually as a standalone component. To run the Helm Broker without building a binary file, follow these steps:

1. Start Minikube:
```bash
minikube start
```

2. Create necessary CRDs:
```bash
kubectl apply -f config/crds/
```

3. Start etcd in the Docker container:
```bash
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  -d \
  quay.io/coreos/etcd:v3.3 \
  /usr/local/bin/etcd \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380
```

4. Start the Broker:
```bash
APP_KUBECONFIG_PATH=/Users/$User/.kube/config \
APP_CONFIG_FILE_NAME=hack/examples/local-etcd-config.yaml \
go run cmd/broker/main.go
```

Now you can test the Broker using the **/v2/catalog** endpoint.

```bash
curl -H "X-Broker-API-Version: 2.13" localhost:8080/cluster/v2/catalog
```

5. Start the Controller:
```bash
APP_KUBECONFIG_PATH=/Users/$User/.kube/config \
APP_DOCUMENTATION_ENABLED=false \
APP_TMP_DIR=/tmp APP_NAMESPACE=default \
APP_SERVICE_NAME=helm-broker \
APP_CONFIG_FILE_NAME=hack/examples/local-etcd-config.yaml \
APP_CLUSTER_SERVICE_BROKER_NAME=helm-broker \
APP_DEVELOP_MODE=true \
go run cmd/controller/main.go -metrics-addr ":8081"
```

>**NOTE:** Not all features are available when you run the Helm Broker locally. All features which perform actions with Tiller do not work. Moreover, the Controller performs operations on ClusterServiceBroker/ServiceBroker resources, which needs the Service Catalog to work properly.

You can run the Controller and the Broker configured with the in-memory storage, but then the Broker cannot read data stored by the Controller. To run the Broker and the Controller without etcd, run these commands:

```bash
APP_KUBECONFIG_PATH=/Users/$User/.kube/config \
APP_CONFIG_FILE_NAME=hack/examples/minimal-config.yaml \
APP_NAMESPACE=kyma-system go run cmd/broker/main.go
```

```bash
APP_KUBECONFIG_PATH=/Users/$User/.kube/config \
APP_DOCUMENTATION_ENABLED=false \
APP_TMP_DIR=/tmp APP_NAMESPACE=default \
APP_SERVICE_NAME=helm-broker \
APP_CONFIG_FILE_NAME=hack/examples/minimal-config.yaml \
APP_CLUSTER_SERVICE_BROKER_NAME=helm-broker \
APP_DEVELOP_MODE=true \
go run cmd/controller/main.go -metrics-addr ":8081"
```

## Next steps

# Helm Broker

[![Go Report Card](https://goreportcard.com/badge/github.com/kyma-project/helm-broker)](https://goreportcard.com/report/github.com/kyma-project/helm-broker)
[![Sourcegraph](https://sourcegraph.com/github.com/kyma-project/helm-broker/-/badge.svg)](https://sourcegraph.com/github.com/kyma-project/helm-broker?badge)

## Overview

Helm Broker is a [Service Broker](https://kyma-project.io/docs/master/components/service-catalog/#overview-service-brokers) that exposes Helm charts as Service Classes in [Service Catalog](https://kyma-project.io/docs/master/components/service-catalog/#overview-service-catalog). To do so, Helm Broker uses the concept of addons. An addon is an abstraction layer over a Helm chart which provides all information required to convert the chart into a Service Class.

Helm Broker fetches default cluster-wide addons defined by the [helm-repos-urls](https://github.com/kyma-project/kyma/blob/master/resources/helm-broker/templates/default-addons-cfg.yaml) custom resource (CR). This CR contains URLs that point to the release of the [`addons`](https://github.com/kyma-project/addons/releases) repository compatible with a given [Kyma release](https://github.com/kyma-project/kyma/releases). You can also configure the Helm Broker to fetch addons definitions from other addons repositories.

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

To run the Helm Broker, you need a Kubernetes cluster with Service Catalog. Run the `./hack/run-dev-kind.sh` script, or follow these steps to set up the Helm Broker on Kind with all necessary dependencies:

1. Create a local cluster on Kind:
```bash
kind create cluster
```

2. Install Service Catalog as a Helm chart:
```bash
helm repo add svc-cat https://kubernetes-sigs.github.io/service-catalog
helm install catalog svc-cat/catalog --namespace catalog --set asyncBindingOperationsEnabled=true
```

3. Clone the Helm Broker repository:
```bash
git clone git@github.com:kyma-project/helm-broker.git
```

4. Install the Helm Broker chart from the cloned repository:
```bash
helm install charts/helm-broker helm-broker --namespace helm-broker
```

### Install Helm Broker as a standalone component

Follow these steps to run the Helm Broker without building a binary file:

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

4. Start the Broker component:
```bash
APP_KUBECONFIG_PATH=/Users/$User/.kube/config \
APP_CONFIG_FILE_NAME=hack/examples/local-etcd-config.yaml \
go run cmd/broker/main.go
```

Now you can test the Broker using the **/v2/catalog** endpoint:
```bash
curl -H "X-Broker-API-Version: 2.13" localhost:8080/cluster/v2/catalog
```

5. Start the Controller component:
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

>**NOTE:** Not all features are available when you run Helm Broker locally. All features that perform actions with Tiller do not work. Moreover, the Controller performs operations on ClusterServiceBroker/ServiceBroker resources, which needs the Service Catalog to work properly.

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

## Development

Follow these steps to develop the project.

### Prerequisites

* [Go](https://golang.org/dl/) 1.12
* [Dep](https://github.com/golang/dep) 0.5
* [Docker](https://www.docker.com/)

>**NOTE:** The versions of Go and Dep are compliant with the `buildpack` used by Prow. For more details, read [this](https://github.com/kyma-project/test-infra/blob/main/prow/images/buildpack-golang/README.md) document.

### Run tests

Before each commit, use the `before-commit.sh` script. The script runs unit tests that check your changes and build binaries.
You can also run integration tests that check if all parts of the Helm Broker work together.
These are the prerequisites for integration tests:

- [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) 1.0.8
- [Etcd](https://github.com/etcd-io/etcd#etcd) 3.4
- [Minio](https://min.io/download) RELEASE.2019-10-12T01-39-57Z

Run integration tests using this command:
```bash
make integration-test
```

### Update chart's images tag

To change the chart's tags version, run this command:
```bash
make VERSION=v0.0.1 DIR=/pr tag-chart-images
```

This command overrides the images tag in the `charts/helm-broker/values.yaml` file to:
```
eu.gcr.io/kyma-project/helm-broker/pr:v0.0.1
```

### Build Docker images

If you want to build Docker images with your changes and push them to a registry, follow these steps:

1. Run tests and build binaries:
```bash
make build
```

2. Build Docker images:
```bash
make build-image
```

3. Configure environment variables pointing to your registry, for example:
```bash
export DOCKER_PUSH_REPOSITORY=eu.gcr.io/
export DOCKER_PUSH_DIRECTORY=your-project
export DOCKER_TAG=latest
```

4. Push the image to the registry:
```bash
make push-image
```

5. Install the Helm Broker with your custom image using the following command:
```bash
helm install charts/helm-broker \
 --name helm-broker \
 --namespace helm-broker \
 --set global.helm_broker.image="${DOCKER_PUSH_REPOSITORY}${DOCKER_PUSH_DIRECTORY}/helm-broker" \
 --set global.helm_broker.version=${DOCKER_TAG} \
 --set global.helm_controller.image="${DOCKER_PUSH_REPOSITORY}${DOCKER_PUSH_DIRECTORY}/helm-controller" \
 --set global.helm_controller.version=${DOCKER_TAG}
```

If you already have the Helm Broker installed, you can upgrade it to use new images:
```bash
helm upgrade helm-broker charts/helm-broker \
 --set global.helm_broker.image="${DOCKER_PUSH_REPOSITORY}${DOCKER_PUSH_DIRECTORY}/helm-broker" \
 --set global.helm_broker.version=${DOCKER_TAG} \
 --set global.helm_controller.image="${DOCKER_PUSH_REPOSITORY}${DOCKER_PUSH_DIRECTORY}/helm-controller" \
 --set global.helm_controller.version=${DOCKER_TAG}
```

## Next steps

For more details, tutorials, and troubleshooting, go to the [documentation](./docs) directory.

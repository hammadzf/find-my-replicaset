# ReplicaSet Finder
A toy Kubernetes controller that labels K8s Deployments with the names of their corresponding ReplicaSets. 

## Features
This simple controller is built using the [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime) library and operates on Kubernetes Deployments. When the controller is running, it "reconciles" Kubernetes Deployment resources (in all namespaces) using the custom reconciliation logic.

> This controller works as an add-on to the native Deployment controller (part of kubernetes-controller-manager) running in the K8s Control Plane. It does not . The only 'reconcilliation' this controller provides is to look for the ReplicaSet owned by a Deployment resource and, after it finds it, to add a label to the Deployment with the name of the found ReplicaSet.

## Source Code (Go Packages)

- (package main)`cmd/main.go`: creates a new manager, a controller that operates on Deployment resources, reconciles them with the implemented reconciler, and starts the created manager.
- (package controller)`internal/controller/controller.go`: implements reconcilliation logic for the Deployment resource, i.e., finding owned ReplicaSet and adding its name as a label to the reconciled Deployment resource. 

> The Ginkgo test suite in the controller package is designed from the perspective of local development. They can only run against a locally running K8s cluster with kubeconfig in the default location. Just run `ginkgo run -v` in the `internal/controller` directory to run the test suite (pre-requisite: ginkgo installed locally).

## Usage

### Local

**Pre-requisites**
- Go v1.24.0
- Access to K8s cluster (looks for the config file in .kube/config)

**Steps**
The controller can run locally on your machine with access to a local or remote Kubernetes cluster. Just build the controller using `go build cmd/main.go` and run the generated executable.

### K8s cluster

**Pre-requisites**
- Kubernetes cluster (tested using Kubernetes v1.32.0 on minikube)

**Steps**
1. Build the image:
```
docker build -t <NAME>:<TAG> .
```
> The image name and tag used in the provided K8s manifests is `rs-finder:latest`. Don't forget to customize the manifests accordingly if you choose a different name/tag while building your image.

2. Load the built image to your local cluster, e.g., in a minikube cluster. Refer to [Loading an Image Into Your Cluster](https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster) if you are using kind.
```
minikube load image <NAME>:<TAG>
```
3. Create required Kubernetes resources (RBAC, namespace, controller deployment) in the [manifests](./manifests/) directory. The controller should be running as a deployment in the `rsfinder-system` namespace.
4. Deployments running in the cluster (across all namespaces) should now have a new label  `owned-replica-set: <RS_NAME>` as part of their metadata.

## Contributions
Feel free to use this basic controller as your starting journey into building custom Kubernetes controllers using the controller-runtime library. If you think of enhancing this basic controller further, don't hesitate to create a PR.
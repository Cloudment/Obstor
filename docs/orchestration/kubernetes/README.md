# Deploy Obstor on Kubernetes [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord)  [![Docker Pulls](https://img.shields.io/docker/pulls/obstor/obstor.svg?maxAge=604800)](https://hub.docker.com/r/obstor/obstor/)

Obstor is a high performance distributed object storage server, designed for large-scale private cloud infrastructure. Obstor is designed in a cloud-native manner to scale sustainably in multi-tenant environments. Orchestration platforms like Kubernetes provide perfect cloud-native environment to deploy and scale Obstor.

## Obstor Deployment on Kubernetes

There are multiple options to deploy Obstor on Kubernetes:

- Obstor-Operator: Operator offers seamless way to create and update highly available distributed Obstor clusters. Refer [Obstor Operator documentation](https://github.com/minio/minio-operator/blob/master/README.md) for more details.

- Helm Chart: Obstor Helm Chart offers customizable and easy Obstor deployment with a single command. Refer [Obstor Helm Chart documentation](https://github.com/minio/charts) for more details.

## Monitoring Obstor in Kubernetes

Obstor server exposes un-authenticated liveness endpoints so Kubernetes can natively identify unhealthy Obstor containers. Obstor also exposes Prometheus compatible data on a different endpoint to enable Prometheus users to natively monitor their Obstor deployments.

## Explore Further

- [Obstor Erasure Code QuickStart Guide](https://pgg.net/docs/obstor/obstor-erasure-code-quickstart-guide)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Helm package manager for kubernetes](https://helm.sh/)

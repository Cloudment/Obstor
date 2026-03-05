# Deploy ObStor on Kubernetes [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord)  [![Docker Pulls](https://img.shields.io/docker/pulls/obstor/obstor.svg?maxAge=604800)](https://hub.docker.com/r/obstor/obstor/)

ObStor is a high performance distributed object storage server, designed for large-scale private cloud infrastructure. ObStor is designed in a cloud-native manner to scale sustainably in multi-tenant environments. Orchestration platforms like Kubernetes provide perfect cloud-native environment to deploy and scale ObStor.

## ObStor Deployment on Kubernetes

There are multiple options to deploy ObStor on Kubernetes:

- ObStor-Operator: Operator offers seamless way to create and update highly available distributed ObStor clusters. Refer [ObStor Operator documentation](https://github.com/minio/minio-operator/blob/master/README.md) for more details.

- Helm Chart: ObStor Helm Chart offers customizable and easy ObStor deployment with a single command. Refer [ObStor Helm Chart documentation](https://github.com/minio/charts) for more details.

## Monitoring ObStor in Kubernetes

ObStor server exposes un-authenticated liveness endpoints so Kubernetes can natively identify unhealthy ObStor containers. ObStor also exposes Prometheus compatible data on a different endpoint to enable Prometheus users to natively monitor their ObStor deployments.

## Explore Further

- [ObStor Erasure Code QuickStart Guide](https://pgg.net/docs/obstor/obstor-erasure-code-quickstart-guide)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Helm package manager for kubernetes](https://helm.sh/)

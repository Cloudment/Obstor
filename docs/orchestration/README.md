# ObStor Deployment Quickstart Guide [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord) [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

ObStor is a cloud-native application designed to scale in a sustainable manner in multi-tenant environments. Orchestration platforms provide perfect launchpad for ObStor to scale. Below is the list of ObStor deployment documents for various orchestration platforms:

| Orchestration platforms|
|:---|
| [`Docker Swarm`](https://pgg.net/docs/obstor/deploy-minio-on-docker-swarm) |
| [`Docker Compose`](https://pgg.net/docs/obstor/deploy-minio-on-docker-compose) |
| [`Kubernetes`](https://pgg.net/docs/obstor/deploy-minio-on-kubernetes) |

## Why is ObStor cloud-native?
The term cloud-native revolves around the idea of applications deployed as micro services, that scale well. It is not about just retrofitting monolithic applications onto modern container based compute environment. A cloud-native application is portable and resilient by design, and can scale horizontally by simply replicating. Modern orchestration platforms like Swarm, Kubernetes and DC/OS make replicating and managing containers in huge clusters easier than ever.

While containers provide isolated application execution environment, orchestration platforms allow seamless scaling by helping replicate and manage containers. ObStor extends this by adding isolated storage environment for each tenant.

ObStor is built ground up on the cloud-native premise. With features like erasure-coding, distributed and shared setup, it focuses only on storage and does it very well. While, it can be scaled by just replicating ObStor instances per tenant via an orchestration platform.

> In a cloud-native environment, scalability is not a function of the application but the orchestration platform.

In a typical modern infrastructure deployment, application, database, key-store, etc. already live in containers and are managed by orchestration platforms. ObStor brings robust, scalable, AWS S3 compatible object storage to the lot.

![Cloud-native](https://github.com/cloudment/obstor/blob/master/docs/screenshots/Minio_Cloud_Native_Arch.jpg?raw=true)

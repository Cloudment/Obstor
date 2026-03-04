# Deploy ObStor on Docker Compose [![Discord](https://discord.pgg.net/discord?type=svg)](https://discord.pgg.net)  [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

Docker Compose allows defining and running single host, multi-container Docker applications.

With Compose, you use a Compose file to configure ObStor services. Then, using a single command, you can create and launch all the Distributed ObStor instances from your configuration. Distributed ObStor instances will be deployed in multiple containers on the same host. This is a great way to set up development, testing, and staging environments, based on Distributed ObStor.

## 1. Prerequisites

* Familiarity with [Docker Compose](https://docs.docker.com/compose/overview/).
* Docker installed on your machine. Download the relevant installer from [here](https://www.docker.com/community-edition#/download).

## 2. Run Distributed ObStor on Docker Compose

To deploy Distributed ObStor on Docker Compose, please download [docker-compose.yaml](https://github.com/cloudment/obstor/blob/master/docs/orchestration/docker-compose/docker-compose.yaml?raw=true) and [nginx.conf](https://github.com/cloudment/obstor/blob/master/docs/orchestration/docker-compose/nginx.conf?raw=true) to your current working directory. Note that Docker Compose pulls the ObStor Docker image, so there is no need to explicitly download ObStor binary. Then run one of the below commands

### GNU/Linux and macOS

```sh
docker-compose pull
docker-compose up
```

### Windows

```sh
docker-compose.exe pull
docker-compose.exe up
```

Distributed instances are now accessible on the host at ports 9000, proceed to access the Web browser at http://127.0.0.1:9000/. Here 4 ObStor server instances are reverse proxied through Nginx load balancing.

### Notes

* By default the Docker Compose file uses the Docker image for latest ObStor server release. You can change the image tag to pull a specific [ObStor Docker image](https://hub.docker.com/r/minio/minio/).

* There are 4 minio distributed instances created by default. You can add more ObStor services (up to total 16) to your ObStor Compose deployment. To add a service
  * Replicate a service definition and change the name of the new service appropriately.
  * Update the command section in each service.
  * Add a new ObStor server instance to the upstream directive in the Nginx configuration file.

  Read more about distributed ObStor [here](https://pgg.net/docs/obstor/distributed-minio-quickstart-guide).

### Explore Further
- [Overview of Docker Compose](https://docs.docker.com/compose/overview/)
- [ObStor Docker Quickstart Guide](https://pgg.net/docs/obstor/minio-docker-quickstart-guide)
- [Deploy ObStor on Docker Swarm](https://pgg.net/docs/obstor/deploy-minio-on-docker-swarm)
- [ObStor Erasure Code QuickStart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide)

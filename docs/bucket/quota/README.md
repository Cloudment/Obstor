# Bucket Quota Configuration Quickstart Guide [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord) [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

![quota](https://raw.githubusercontent.com/cloudment/obstor/master/docs/bucket/quota/bucketquota.png)

Buckets can be configured to have one of two types of quota configuration - FIFO and Hard quota.

- `Hard` quota disallows writes to the bucket after configured quota limit is reached.
- `FIFO` quota automatically deletes oldest content until bucket usage falls within configured limit while permitting writes.

> NOTE: Bucket quotas are not supported under gateway or standalone single disk deployments.

## Prerequisites
- Install ObStor - [ObStor Quickstart Guide](https://pgg.net/docs/obstor/minio-quickstart-guide).
- [Use `mc` with ObStor Server](https://pgg.net/docs/obstor/minio-client-quickstart-guide)

## Set bucket quota configuration

### Set a hard quota of 1GB for a bucket `mybucket` on ObStor object storage:

```sh
$ mc admin bucket quota myminio/mybucket --hard 1gb
```

### Set FIFO quota of 5GB for a bucket "mybucket" on ObStor to allow automatic deletion of older content to ensure bucket usage remains within 5GB

```sh
$ mc admin bucket quota myminio/mybucket --fifo 5gb
```

### Verify the quota configured on `mybucket` on ObStor

```sh
$ mc admin bucket quota myminio/mybucket
```

### Clear bucket quota configuration for `mybucket` on ObStor

```sh
$ mc admin bucket quota myminio/mybucket --clear
```

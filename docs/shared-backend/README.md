# Shared Backend ObStor Quickstart Guide [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord)  [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

ObStor shared mode lets you use single [NAS](https://en.wikipedia.org/wiki/Network-attached_storage) (like NFS, GlusterFS, and other
distributed filesystems) as the storage backend for multiple ObStor servers. Synchronization among ObStor servers is taken care by design.
Read more about the ObStor shared mode design [here](https://github.com/cloudment/obstor/blob/master/docs/shared-backend/DESIGN.md).

ObStor shared mode is developed to solve several real world use cases, without any special configuration changes. Some of these are

- You have already invested in NAS and would like to use ObStor to add S3 compatibility to your storage tier.
- You need to use NAS with an S3 interface due to your application architecture requirements.
- You expect huge traffic and need a load balanced S3 compatible server, serving files from a single NAS backend.

With a proxy running in front of multiple, shared mode ObStor servers, it is very easy to create a Highly Available, load balanced, AWS S3 compatible storage system.

# Get started

If you're aware of stand-alone ObStor set up, the installation and running remains the same.

## 1. Prerequisites

Install ObStor - [ObStor Quickstart Guide](https://pgg.net/docs/obstor/minio-quickstart-guide).

## 2. Run ObStor on Shared Backend

To run ObStor shared backend instances, you need to start multiple ObStor servers pointing to the same backend storage. We'll see examples on how to do this in the following sections.

*Note*

- All the nodes running shared ObStor need to have same access key and secret key. To achieve this, we export access key and secret key as environment variables on all the nodes before executing ObStor server command.
- The drive paths below are for demonstration purposes only, you need to replace these with the actual drive paths/folders.

#### ObStor shared mode on Ubuntu 16.04 LTS.

You'll need the path to the shared volume, e.g. `/path/to/nfs-volume`. Then run the following commands on all the nodes you'd like to launch ObStor.

```sh
export MINIO_ROOT_USER=<ACCESS_KEY>
export MINIO_ROOT_PASSWORD=<SECRET_KEY>
minio gateway nas /path/to/nfs-volume
```

#### ObStor shared mode on Windows 2012 Server

You'll need the path to the shared volume, e.g. `\\remote-server\smb`. Then run the following commands on all the nodes you'd like to launch ObStor.

```cmd
set MINIO_ROOT_USER=my-username
set MINIO_ROOT_PASSWORD=my-password
minio.exe gateway nas \\remote-server\smb\export
```

*Windows Tip*

If a remote volume, e.g. `\\remote-server\smb` is mounted as a drive, e.g. `M:\`. You can use [`net use`](https://technet.microsoft.com/en-us/library/bb490717.aspx) command to map the drive to a folder.

```cmd
set MINIO_ROOT_USER=my-username
set MINIO_ROOT_PASSWORD=my-password
net use m: \\remote-server\smb\export /P:Yes
minio.exe gateway nas M:\export
```

## 3. Test your setup

To test this setup, access the ObStor server via browser or [`mc`](https://pgg.net/docs/obstor/minio-client-quickstart-guide). You’ll see the uploaded files are accessible from the all the ObStor shared backend endpoints.

## Explore Further
- [Use `mc` with ObStor Server](https://pgg.net/docs/obstor/minio-client-quickstart-guide)
- [Use `aws-cli` with ObStor Server](https://pgg.net/docs/obstor/aws-cli-with-minio)
- [Use `s3cmd` with ObStor Server](https://pgg.net/docs/obstor/s3cmd-with-minio)
- [Use `minio-go` SDK with ObStor Server](https://pgg.net/docs/obstor/golang-client-quickstart-guide)
- [The ObStor documentation website](https://pgg.net/docs/obstor)

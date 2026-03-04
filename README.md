# ObStor Quickstart Guide
[![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord) [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

[![ObStor](https://raw.githubusercontent.com/minio/minio/master/.github/logo.svg?sanitize=true)](https://pgg.net)

ObStor is a High Performance Object Storage released under Apache License v2.0. It is API compatible with Amazon S3 cloud storage service. Use ObStor to build high performance infrastructure for machine learning, analytics and application data workloads.

This README provides quickstart instructions on running ObStor on baremetal hardware, including Docker-based installations. For Kubernetes environments,
use the [ObStor Kubernetes Operator](https://github.com/minio/operator/blob/master/README.md).

# Docker Installation

Use the following commands to run a standalone ObStor server on a Docker container.

Standalone ObStor servers are best suited for early development and evaluation. Certain features such as versioning, object locking, and bucket replication
require distributed deploying ObStor with Erasure Coding. For extended development and production, deploy ObStor with Erasure Coding enabled - specifically,
with a *minimum* of 4 drives per ObStor server. See [ObStor Erasure Code Quickstart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide.html)
for more complete documentation.

## Stable

Run the following command to run the latest stable image of ObStor on a Docker container using an ephemeral data volume:

```sh
docker run -p 9000:9000 minio/minio server /data
```

The ObStor deployment starts using default root credentials `minioadmin:minioadmin`. You can test the deployment using the ObStor Browser, an embedded
web-based object browser built into ObStor Server. Point a web browser running on the host machine to http://127.0.0.1:9000 and log in with the
root credentials. You can use the Browser to create buckets, upload objects, and browse the contents of the ObStor server.

You can also connect using any S3-compatible tool, such as the ObStor Client `mc` commandline tool. See
[Test using ObStor Client `mc`](#test-using-minio-client-mc) for more information on using the `mc` commandline tool. For application developers,
see https://pgg.net/docs/obstor/ and click **MINIO SDKS** in the navigation to view ObStor SDKs for supported languages.


> NOTE: To deploy ObStor on Docker with persistent storage, you must map local persistent directories from the host OS to the container using the
  `docker -v` option. For example, `-v /mnt/data:/data` maps the host OS drive at `/mnt/data` to `/data` on the Docker container.

## Edge

Run the following command to run the bleeding-edge image of ObStor on a Docker container using an ephemeral data volume:

```
docker run -p 9000:9000 minio/minio:edge server /data
```

The ObStor deployment starts using default root credentials `minioadmin:minioadmin`. You can test the deployment using the ObStor Browser, an embedded
web-based object browser built into ObStor Server. Point a web browser running on the host machine to http://127.0.0.1:9000 and log in with the
root credentials. You can use the Browser to create buckets, upload objects, and browse the contents of the ObStor server.

You can also connect using any S3-compatible tool, such as the ObStor Client `mc` commandline tool. See
[Test using ObStor Client `mc`](#test-using-minio-client-mc) for more information on using the `mc` commandline tool. For application developers,
see https://pgg.net/docs/obstor/ and click **MINIO SDKS** in the navigation to view ObStor SDKs for supported languages.


> NOTE: To deploy ObStor on Docker with persistent storage, you must map local persistent directories from the host OS to the container using the
  `docker -v` option. For example, `-v /mnt/data:/data` maps the host OS drive at `/mnt/data` to `/data` on the Docker container.

# macOS

Use the following commands to run a standalone ObStor server on macOS.

Standalone ObStor servers are best suited for early development and evaluation. Certain features such as versioning, object locking, and bucket replication
require distributed deploying ObStor with Erasure Coding. For extended development and production, deploy ObStor with Erasure Coding enabled - specifically,
with a *minimum* of 4 drives per ObStor server. See [ObStor Erasure Code Quickstart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide.html)
for more complete documentation.

## Homebrew (recommended)

Run the following command to install the latest stable ObStor package using [Homebrew](https://brew.sh/). Replace ``/data`` with the path to the drive or directory in which you want ObStor to store data.

```sh
brew install minio/stable/minio
minio server /data
```

> NOTE: If you previously installed minio using `brew install minio` then it is recommended that you reinstall minio from `minio/stable/minio` official repo instead.

```sh
brew uninstall minio
brew install minio/stable/minio
```

The ObStor deployment starts using default root credentials `minioadmin:minioadmin`. You can test the deployment using the ObStor Browser, an embedded
web-based object browser built into ObStor Server. Point a web browser running on the host machine to http://127.0.0.1:9000 and log in with the
root credentials. You can use the Browser to create buckets, upload objects, and browse the contents of the ObStor server.

You can also connect using any S3-compatible tool, such as the ObStor Client `mc` commandline tool. See
[Test using ObStor Client `mc`](#test-using-minio-client-mc) for more information on using the `mc` commandline tool. For application developers,
see https://pgg.net/docs/obstor/ and click **MINIO SDKS** in the navigation to view ObStor SDKs for supported languages.

## Binary Download

Use the following command to download and run a standalone ObStor server on macOS. Replace ``/data`` with the path to the drive or directory in which you want ObStor to store data.

```sh
wget https://dl.pgg.net/server/minio/release/darwin-amd64/minio
chmod +x minio
./minio server /data
```

The ObStor deployment starts using default root credentials `minioadmin:minioadmin`. You can test the deployment using the ObStor Browser, an embedded
web-based object browser built into ObStor Server. Point a web browser running on the host machine to http://127.0.0.1:9000 and log in with the
root credentials. You can use the Browser to create buckets, upload objects, and browse the contents of the ObStor server.

You can also connect using any S3-compatible tool, such as the ObStor Client `mc` commandline tool. See
[Test using ObStor Client `mc`](#test-using-minio-client-mc) for more information on using the `mc` commandline tool. For application developers,
see https://pgg.net/docs/obstor/ and click **MINIO SDKS** in the navigation to view ObStor SDKs for supported languages.


# GNU/Linux

Use the following command to run a standalone ObStor server on Linux hosts running 64-bit Intel/AMD architectures. Replace ``/data`` with the path to the drive or directory in which you want ObStor to store data.

```sh
wget https://dl.pgg.net/server/minio/release/linux-amd64/minio
chmod +x minio
./minio server /data
```

Replace ``/data`` with the path to the drive or directory in which you want ObStor to store data.

The following table lists supported architectures. Replace the `wget` URL with the architecture for your Linux host.

| Architecture                   | URL                                                        |
| --------                       | ------                                                     |
| 64-bit Intel/AMD               | https://dl.pgg.net/server/minio/release/linux-amd64/minio   |
| 64-bit ARM                     | https://dl.pgg.net/server/minio/release/linux-arm64/minio   |
| 64-bit PowerPC LE (ppc64le)    | https://dl.pgg.net/server/minio/release/linux-ppc64le/minio |
| IBM Z-Series (S390X)           | https://dl.pgg.net/server/minio/release/linux-s390x/minio   |

The ObStor deployment starts using default root credentials `minioadmin:minioadmin`. You can test the deployment using the ObStor Browser, an embedded
web-based object browser built into ObStor Server. Point a web browser running on the host machine to http://127.0.0.1:9000 and log in with the
root credentials. You can use the Browser to create buckets, upload objects, and browse the contents of the ObStor server.

You can also connect using any S3-compatible tool, such as the ObStor Client `mc` commandline tool. See
[Test using ObStor Client `mc`](#test-using-minio-client-mc) for more information on using the `mc` commandline tool. For application developers,
see https://pgg.net/docs/obstor/ and click **MINIO SDKS** in the navigation to view ObStor SDKs for supported languages.


> NOTE: Standalone ObStor servers are best suited for early development and evaluation. Certain features such as versioning, object locking, and bucket replication
require distributed deploying ObStor with Erasure Coding. For extended development and production, deploy ObStor with Erasure Coding enabled - specifically,
with a *minimum* of 4 drives per ObStor server. See [ObStor Erasure Code Quickstart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide.html)
for more complete documentation.

# Microsoft Windows

To run ObStor on 64-bit Windows hosts, download the ObStor executable from the following URL:

```sh
https://dl.pgg.net/server/minio/release/windows-amd64/minio.exe
```

Use the following command to run a standalone ObStor server on the Windows host. Replace ``D:\`` with the path to the drive or directory in which you want ObStor to store data. You must change the terminal or powershell directory to the location of the ``minio.exe`` executable, *or* add the path to that directory to the system ``$PATH``:

```sh
minio.exe server D:\
```

The ObStor deployment starts using default root credentials `minioadmin:minioadmin`. You can test the deployment using the ObStor Browser, an embedded
web-based object browser built into ObStor Server. Point a web browser running on the host machine to http://127.0.0.1:9000 and log in with the
root credentials. You can use the Browser to create buckets, upload objects, and browse the contents of the ObStor server.

You can also connect using any S3-compatible tool, such as the ObStor Client `mc` commandline tool. See
[Test using ObStor Client `mc`](#test-using-minio-client-mc) for more information on using the `mc` commandline tool. For application developers,
see https://pgg.net/docs/obstor/ and click **MINIO SDKS** in the navigation to view ObStor SDKs for supported languages.

> NOTE: Standalone ObStor servers are best suited for early development and evaluation. Certain features such as versioning, object locking, and bucket replication
require distributed deploying ObStor with Erasure Coding. For extended development and production, deploy ObStor with Erasure Coding enabled - specifically,
with a *minimum* of 4 drives per ObStor server. See [ObStor Erasure Code Quickstart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide.html)
for more complete documentation.

# FreeBSD

ObStor does not provide an official FreeBSD binary. However, FreeBSD maintains an [upstream release](https://www.freshports.org/www/minio) using [pkg](https://github.com/freebsd/pkg):

```sh
pkg install minio
sysrc minio_enable=yes
sysrc minio_disks=/home/user/Photos
service minio start
```

# Install from Source

Use the following commands to compile and run a standalone ObStor server from source. Source installation is only intended for developers and advanced users. If you do not have a working Golang environment, please follow [How to install Golang](https://golang.org/doc/install). Minimum version required is [go1.16](https://golang.org/dl/#stable)

```sh
GO111MODULE=on go get github.com/cloudment/obstor
```

The ObStor deployment starts using default root credentials `minioadmin:minioadmin`. You can test the deployment using the ObStor Browser, an embedded
web-based object browser built into ObStor Server. Point a web browser running on the host machine to http://127.0.0.1:9000 and log in with the
root credentials. You can use the Browser to create buckets, upload objects, and browse the contents of the ObStor server.

You can also connect using any S3-compatible tool, such as the ObStor Client `mc` commandline tool. See
[Test using ObStor Client `mc`](#test-using-minio-client-mc) for more information on using the `mc` commandline tool. For application developers,
see https://pgg.net/docs/obstor/ and click **MINIO SDKS** in the navigation to view ObStor SDKs for supported languages.


> NOTE: Standalone ObStor servers are best suited for early development and evaluation. Certain features such as versioning, object locking, and bucket replication
require distributed deploying ObStor with Erasure Coding. For extended development and production, deploy ObStor with Erasure Coding enabled - specifically,
with a *minimum* of 4 drives per ObStor server. See [ObStor Erasure Code Quickstart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide.html)
for more complete documentation.

ObStor strongly recommends *against* using compiled-from-source ObStor servers for production environments.

# Deployment Recommendations

## Allow port access for Firewalls

By default ObStor uses the port 9000 to listen for incoming connections. If your platform blocks the port by default, you may need to enable access to the port.

### ufw

For hosts with ufw enabled (Debian based distros), you can use `ufw` command to allow traffic to specific ports. Use below command to allow access to port 9000

```sh
ufw allow 9000
```

Below command enables all incoming traffic to ports ranging from 9000 to 9010.

```sh
ufw allow 9000:9010/tcp
```

### firewall-cmd

For hosts with firewall-cmd enabled (CentOS), you can use `firewall-cmd` command to allow traffic to specific ports. Use below commands to allow access to port 9000

```sh
firewall-cmd --get-active-zones
```

This command gets the active zone(s). Now, apply port rules to the relevant zones returned above. For example if the zone is `public`, use

```sh
firewall-cmd --zone=public --add-port=9000/tcp --permanent
```

Note that `permanent` makes sure the rules are persistent across firewall start, restart or reload. Finally reload the firewall for changes to take effect.

```sh
firewall-cmd --reload
```

### iptables

For hosts with iptables enabled (RHEL, CentOS, etc), you can use `iptables` command to enable all traffic coming to specific ports. Use below command to allow
access to port 9000

```sh
iptables -A INPUT -p tcp --dport 9000 -j ACCEPT
service iptables restart
```

Below command enables all incoming traffic to ports ranging from 9000 to 9010.

```sh
iptables -A INPUT -p tcp --dport 9000:9010 -j ACCEPT
service iptables restart
```

## Pre-existing data
When deployed on a single drive, ObStor server lets clients access any pre-existing data in the data directory. For example, if ObStor is started with the command  `minio server /mnt/data`, any pre-existing data in the `/mnt/data` directory would be accessible to the clients.

The above statement is also valid for all gateway backends.

# Test ObStor Connectivity

## Test using ObStor Browser
ObStor Server comes with an embedded web based object browser. Point your web browser to http://127.0.0.1:9000 to ensure your server has started successfully.

![Screenshot](https://github.com/cloudment/obstor/blob/master/docs/screenshots/minio-browser.png?raw=true)

## Test using ObStor Client `mc`
`mc` provides a modern alternative to UNIX commands like ls, cat, cp, mirror, diff etc. It supports filesystems and Amazon S3 compatible cloud storage services. Follow the ObStor Client [Quickstart Guide](https://pgg.net/docs/obstor/minio-client-quickstart-guide) for further instructions.

# Upgrading ObStor
ObStor server supports rolling upgrades, i.e. you can update one ObStor instance at a time in a distributed cluster. This allows upgrades with no downtime. Upgrades can be done manually by replacing the binary with the latest release and restarting all servers in a rolling fashion. However, we recommend all our users to use [`mc admin update`](https://pgg.net/docs/obstor/minio-admin-complete-guide.html#update) from the client. This will update all the nodes in the cluster simultaneously and restart them, as shown in the following command from the ObStor client (mc):

```
mc admin update <minio alias, e.g., myminio>
```

> NOTE: some releases might not allow rolling upgrades, this is always called out in the release notes and it is generally advised to read release notes before upgrading. In such a situation `mc admin update` is the recommended upgrading mechanism to upgrade all servers at once.

## Important things to remember during ObStor upgrades

- `mc admin update` will only work if the user running ObStor has write access to the parent directory where the binary is located, for example if the current binary is at `/usr/local/bin/minio`, you would need write access to `/usr/local/bin`.
- `mc admin update` updates and restarts all servers simultaneously, applications would retry and continue their respective operations upon upgrade.
- `mc admin update` is disabled in kubernetes/container environments, container environments provide their own mechanisms to rollout of updates.
- In the case of federated setups `mc admin update` should be run against each cluster individually. Avoid updating `mc` to any new releases until all clusters have been successfully updated.
- If using `kes` as KMS with ObStor, just replace the binary and restart `kes` more information about `kes` can be found [here](https://github.com/minio/kes/wiki)
- If using Vault as KMS with ObStor, ensure you have followed the Vault upgrade procedure outlined here: https://www.vaultproject.io/docs/upgrading/index.html
- If using etcd with ObStor for the federation, ensure you have followed the etcd upgrade procedure outlined here: https://github.com/etcd-io/etcd/blob/master/Documentation/upgrades/upgrading-etcd.md

# Explore Further
- [ObStor Erasure Code QuickStart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide)
- [Use `mc` with ObStor Server](https://pgg.net/docs/obstor/minio-client-quickstart-guide)
- [Use `aws-cli` with ObStor Server](https://pgg.net/docs/obstor/aws-cli-with-minio)
- [Use `s3cmd` with ObStor Server](https://pgg.net/docs/obstor/s3cmd-with-minio)
- [Use `minio-go` SDK with ObStor Server](https://pgg.net/docs/obstor/golang-client-quickstart-guide)
- [The ObStor documentation website](https://pgg.net/docs/obstor)

# Contribute to ObStor Project
Please follow ObStor [Contributor's Guide](https://github.com/cloudment/obstor/blob/master/CONTRIBUTING.md)

# License
Use of ObStor is governed by the Apache 2.0 License found at [LICENSE](https://github.com/cloudment/obstor/blob/master/LICENSE).

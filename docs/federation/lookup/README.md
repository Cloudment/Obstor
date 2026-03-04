*Federation feature is deprecated and should be avoided for future deployments*

# Federation Quickstart Guide [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord)
This document explains how to configure ObStor with `Bucket lookup from DNS` style federation.

## Get started

### 1. Prerequisites
Install ObStor - [ObStor Quickstart Guide](https://pgg.net/docs/obstor/minio-quickstart-guide).

### 2. Run ObStor in federated mode
Bucket lookup from DNS federation requires two dependencies

- etcd (for bucket DNS service records)
- CoreDNS (for DNS management based on populated bucket DNS service records, optional)

## Architecture

![bucket-lookup](https://github.com/cloudment/obstor/blob/master/docs/federation/lookup/bucket-lookup.png?raw=true)

### Environment variables

#### MINIO_ETCD_ENDPOINTS

This is comma separated list of etcd servers that you want to use as the ObStor federation back-end. This should
be same across the federated deployment, i.e. all the ObStor instances within a federated deployment should use same
etcd back-end.

#### MINIO_DOMAIN

This is the top level domain name used for the federated setup. This domain name should ideally resolve to a load-balancer
running in front of all the federated ObStor instances. The domain name is used to create sub domain entries to etcd. For
example, if the domain is set to `domain.com`, the buckets `bucket1`, `bucket2` will be accessible as `bucket1.domain.com`
and `bucket2.domain.com`.

#### MINIO_PUBLIC_IPS

This is comma separated list of IP addresses to which buckets created on this ObStor instance will resolve to. For example,
a bucket `bucket1` created on current ObStor instance will be accessible as `bucket1.domain.com`, and the DNS entry for
`bucket1.domain.com` will point to IP address set in `MINIO_PUBLIC_IPS`.

*Note*

- This field is mandatory for standalone and erasure code ObStor server deployments, to enable federated mode.
- This field is optional for distributed deployments. If you don't set this field in a federated setup, we use the IP addresses of
hosts passed to the ObStor server startup and use them for DNS entries.

### Run Multiple Clusters

> cluster1

```sh
export MINIO_ETCD_ENDPOINTS="http://remote-etcd1:2379,http://remote-etcd2:4001"
export MINIO_DOMAIN=domain.com
export MINIO_PUBLIC_IPS=44.35.2.1,44.35.2.2,44.35.2.3,44.35.2.4
minio server http://rack{1...4}.host{1...4}.domain.com/mnt/export{1...32}
```

> cluster2

```sh
export MINIO_ETCD_ENDPOINTS="http://remote-etcd1:2379,http://remote-etcd2:4001"
export MINIO_DOMAIN=domain.com
export MINIO_PUBLIC_IPS=44.35.1.1,44.35.1.2,44.35.1.3,44.35.1.4
minio server http://rack{5...8}.host{5...8}.domain.com/mnt/export{1...32}
```

In this configuration you can see `MINIO_ETCD_ENDPOINTS` points to the etcd backend which manages ObStor's
`config.json` and bucket DNS SRV records. `MINIO_DOMAIN` indicates the domain suffix for the bucket which
will be used to resolve bucket through DNS. For example if you have a bucket such as `mybucket`, the
client can use now `mybucket.domain.com` to directly resolve itself to the right cluster. `MINIO_PUBLIC_IPS`
points to the public IP address where each cluster might be accessible, this is unique for each cluster.

NOTE: `mybucket` only exists on one cluster either `cluster1` or `cluster2` this is random and
is decided by how `domain.com` gets resolved, if there is a round-robin DNS on `domain.com` then
it is randomized which cluster might provision the bucket.

### 3. Upgrading to `etcdv3` API

Users running ObStor federation from release `RELEASE.2018-06-09T03-43-35Z` to `RELEASE.2018-07-10T01-42-11Z`, should migrate the existing bucket data on etcd server to `etcdv3` API, and update CoreDNS version to `1.2.0` before updating their ObStor server to the latest version.

Here is some background on why this is needed - ObStor server release `RELEASE.2018-06-09T03-43-35Z` to `RELEASE.2018-07-10T01-42-11Z` used etcdv2 API to store bucket data to etcd server. This was due to `etcdv3` support not available for CoreDNS server. So, even if ObStor used `etcdv3` API to store bucket data, CoreDNS wouldn't be able to read and serve it as DNS records.

Now that CoreDNS [supports etcdv3](https://coredns.io/2018/07/11/coredns-1.2.0-release/), ObStor server uses `etcdv3` API to store bucket data to etcd server. As `etcdv2` and `etcdv3` APIs are not compatible, data stored using `etcdv2` API is not visible to the `etcdv3` API. So, bucket data stored by previous ObStor version will not be visible to current ObStor version, until a migration is done.

CoreOS team has documented the steps required to migrate existing data from `etcdv2` to `etcdv3` in [this blog post](https://coreos.com/blog/migrating-applications-etcd-v3.html). Please refer the post and migrate etcd data to `etcdv3` API.

### 4. Test your setup

To test this setup, access the ObStor server via browser or [`mc`](https://pgg.net/docs/obstor/minio-client-quickstart-guide). You’ll see the uploaded files are accessible from the all the ObStor endpoints.

# Explore Further

- [Use `mc` with ObStor Server](https://pgg.net/docs/obstor/minio-client-quickstart-guide)
- [Use `aws-cli` with ObStor Server](https://pgg.net/docs/obstor/aws-cli-with-minio)
- [Use `s3cmd` with ObStor Server](https://pgg.net/docs/obstor/s3cmd-with-minio)
- [Use `minio-go` SDK with ObStor Server](https://pgg.net/docs/obstor/golang-client-quickstart-guide)
- [The ObStor documentation website](https://pgg.net/docs/obstor)

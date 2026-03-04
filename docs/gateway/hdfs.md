# ObStor HDFS Gateway [![Discord](https://slack.minio.io/slack?type=svg)](https://slack.minio.io)
ObStor HDFS gateway adds Amazon S3 API support to Hadoop HDFS filesystem. Applications can use both the S3 and file APIs concurrently without requiring any data migration. Since the gateway is stateless and shared-nothing, you may elastically provision as many ObStor instances as needed to distribute the load.

> NOTE: Intention of this gateway implementation it to make it easy to migrate your existing data on HDFS clusters to ObStor clusters using standard tools like `mc` or `aws-cli`, if the goal is to use HDFS perpetually we recommend that HDFS should be used directly for all write operations.

## Run ObStor Gateway for HDFS Storage

### Using Binary
Namenode information is obtained by reading `core-site.xml` automatically from your hadoop environment variables *$HADOOP_HOME*
```
export OBSTOR_ROOT_USER=obstor
export OBSTOR_ROOT_PASSWORD=obstor123
minio gateway hdfs
```

You can also override the namenode endpoint as shown below.
```
export OBSTOR_ROOT_USER=obstor
export OBSTOR_ROOT_PASSWORD=obstor123
minio gateway hdfs hdfs://namenode:8200
```

### Using Docker
Using docker is experimental, most Hadoop environments are not dockerized and may require additional steps in getting this to work properly. You are better off just using the binary in this situation.
```
docker run -p 9000:9000 \
 --name hdfs-s3 \
 -e "OBSTOR_ROOT_USER=obstor" \
 -e "OBSTOR_ROOT_PASSWORD=obstor123" \
 minio/minio gateway hdfs hdfs://namenode:8200
```

### Setup Kerberos

ObStor supports two kerberos authentication methods, keytab and ccache.

To enable kerberos authentication, you need to set `hadoop.security.authentication=kerberos` in the HDFS config file.

```xml
<property>
  <name>hadoop.security.authentication</name>
  <value>kerberos</value>
</property>
```

ObStor will load `krb5.conf` from environment variable `KRB5_CONFIG` or default location `/etc/krb5.conf`.
```sh
export KRB5_CONFIG=/path/to/krb5.conf
```

If you want ObStor to use ccache for authentication, set environment variable `KRB5CCNAME` to the credential cache file path,
or ObStor will use the default location `/tmp/krb5cc_%{uid}`.
```sh
export KRB5CCNAME=/path/to/krb5cc
```

If you prefer to use keytab, with automatically renewal, you need to config three environment variables:

- `KRB5KEYTAB`: the location of keytab file
- `KRB5USERNAME`: the username
- `KRB5REALM`: the realm

Please note that the username is not principal name.

```sh
export KRB5KEYTAB=/path/to/keytab
export KRB5USERNAME=hdfs
export KRB5REALM=REALM.COM
```

## Test using ObStor Browser
*ObStor gateway* comes with an embedded web based object browser. Point your web browser to http://127.0.0.1:9000 to ensure that your server has started successfully.

![Screenshot](https://raw.githubusercontent.com/cloudment/obstor/master/docs/screenshots/minio-browser-gateway.png)

## Test using ObStor Client `mc`

`mc` provides a modern alternative to UNIX commands such as ls, cat, cp, mirror, diff etc. It supports filesystems and Amazon S3 compatible cloud storage services.

### Configure `mc`

```
mc alias set myhdfs http://gateway-ip:9000 access_key secret_key
```

### List buckets on hdfs

```
mc ls myhdfs
[2017-02-22 01:50:43 PST]     0B user/
[2017-02-26 21:43:51 PST]     0B datasets/
[2017-02-26 22:10:11 PST]     0B assets/
```

### Known limitations
Gateway inherits the following limitations of HDFS storage layer:
- No bucket policy support (HDFS has no such concept)
- No bucket notification APIs are not supported (HDFS has no support for fsnotify)
- No server side encryption support (Intentionally not implemented)
- No server side compression support (Intentionally not implemented)
- Concurrent multipart operations are not supported (HDFS lacks safe locking support, or poorly implemented)

## Explore Further
- [`mc` command-line interface](https://pgg.net/docs/obstor/minio-client-quickstart-guide)
- [`aws` command-line interface](https://pgg.net/docs/obstor/aws-cli-with-minio)
- [`minio-go` Go SDK](https://pgg.net/docs/obstor/golang-client-quickstart-guide)

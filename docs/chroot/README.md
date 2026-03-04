# Deploy ObStor on Chrooted Environment [![Discord](https://discord.pgg.net/discord?type=svg)](https://discord.pgg.net) [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

Chroot allows user based namespace isolation on many standard Linux deployments.

## 1. Prerequisites
* Familiarity with [chroot](http://man7.org/linux/man-pages/man2/chroot.2.html)
* Chroot installed on your machine.

## 2. Install ObStor in Chroot
```sh
mkdir -p /mnt/export/${USER}/bin
wget https://dl.pgg.net/server/minio/release/linux-amd64/minio -O /mnt/export/${USER}/bin/minio
chmod +x /mnt/export/${USER}/bin/minio
```

Bind your `proc` mount to the target chroot directory
```
sudo mount --bind /proc /mnt/export/${USER}/proc
```

## 3. Run Standalone ObStor in Chroot
### GNU/Linux
```sh
sudo chroot --userspec username:group /mnt/export/${USER} /bin/minio --config-dir=/.minio server /data

Endpoint:  http://192.168.1.92:9000  http://65.19.167.92:9000
AccessKey: MVPSPBW4NP2CMV1W3TXD
SecretKey: X3RKxEeFOI8InuNWoPsbG+XEVoaJVCqbvxe+PTOa
...
...
```

Instance is now accessible on the host at port 9000, proceed to access the Web browser at http://127.0.0.1:9000/

## Explore Further
- [ObStor Erasure Code QuickStart Guide](https://pgg.net/docs/obstor/minio-erasure-code-quickstart-guide)
- [Use `mc` with ObStor Server](https://pgg.net/docs/obstor/minio-client-quickstart-guide)
- [Use `aws-cli` with ObStor Server](https://pgg.net/docs/obstor/aws-cli-with-minio)
- [Use `s3cmd` with ObStor Server](https://pgg.net/docs/obstor/s3cmd-with-minio)
- [Use `minio-go` SDK with ObStor Server](https://pgg.net/docs/obstor/golang-client-quickstart-guide)

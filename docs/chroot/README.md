# Deploy Obstor on Chrooted Environment

Chroot allows user based namespace isolation on many standard Linux deployments.

## 1. Prerequisites
* Familiarity with [chroot](http://man7.org/linux/man-pages/man2/chroot.2.html)
* Chroot installed on your machine.

## 2. Install Obstor in Chroot
```sh
mkdir -p /mnt/export/${USER}/bin
wget https://dl.pgg.net/packages/obstor/release/linux-amd64/obstor -O /mnt/export/${USER}/bin/obstor
chmod +x /mnt/export/${USER}/bin/obstor
```

Bind your `proc` mount to the target chroot directory
```
sudo mount --bind /proc /mnt/export/${USER}/proc
```

## 3. Run Standalone Obstor in Chroot
### GNU/Linux
```sh
sudo chroot --userspec username:group /mnt/export/${USER} /bin/obstor --config-dir=/.obstor server /data

Endpoint:  http://192.168.1.92:9000  http://65.19.167.92:9000
AccessKey: MVPSPBW4NP2CMV1W3TXD
SecretKey: X3RKxEeFOI8InuNWoPsbG+XEVoaJVCqbvxe+PTOa
...
...
```

Instance is now accessible on the host at port 9000, proceed to access the Web browser at http://127.0.0.1:9000/

## Explore Further
- [Obstor Erasure Code QuickStart Guide](https://obstor.net/docs/obstor-erasure-code-quickstart-guide)
- [Use `mc` with Obstor Server](https://obstor.net/docs/obstor-client-quickstart-guide)
- [Use `aws-cli` with Obstor Server](https://obstor.net/docs/aws-cli-with-obstor)
- [Use `s3cmd` with Obstor Server](https://obstor.net/docs/s3cmd-with-obstor)
- [Use `minio-go` SDK with Obstor Server](https://obstor.net/docs/golang-client-quickstart-guide)

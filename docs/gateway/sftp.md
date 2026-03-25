# Obstor SFTP Gateway

Obstor includes a built-in SFTP gateway that allows any standard SFTP client to access your object storage. This enables legacy applications, automated file transfers, and users who prefer SFTP workflows to interact with S3-compatible storage seamlessly.

## Features

- **Native SFTP protocol support**: no external SFTP server required
- **S3-compatible backend**: files uploaded via SFTP are stored as S3 objects
- **Standard authentication**: uses Obstor IAM credentials for SFTP access
- **Bucket mapping**: each SFTP user's home directory maps to a bucket
- **Large file support**: handles files of any size with multipart uploads

## Quick Start

### 1. Enable SFTP

Start Obstor with the SFTP flag enabled:

```sh
export OBSTOR_ROOT_USER=admin
export OBSTOR_ROOT_PASSWORD=password
obstor server /data --sftp-address :2022
```

The SFTP server will listen on port `2022` by default.

### 2. Connect with an SFTP Client

Use any standard SFTP client to connect:

```sh
sftp -P 2022 admin@localhost
```

Or use the `sftp://` protocol URL:

```
sftp://admin:password@localhost:2022
```

### 3. Transfer Files

Once connected, standard SFTP commands work as expected:

```sh
sftp> ls
sftp> put myfile.txt mybucket/myfile.txt
sftp> get mybucket/myfile.txt localfile.txt
sftp> mkdir mybucket
sftp> rm mybucket/myfile.txt
```

## Configuration

### SFTP Port

Configure the SFTP listening address:

```sh
obstor server /data --sftp-address :2222
```

### SSH Host Keys

Obstor auto-generates SSH host keys on first startup. To use custom host keys:

```sh
obstor server /data --sftp-address :2022 --sftp-ssh-key /path/to/ssh_host_rsa_key
```

### TLS with SFTP

SFTP uses SSH protocol encryption by default. For additional TLS encryption on the S3 API:

```sh
obstor server /data \
  --sftp-address :2022 \
  --certs-dir /path/to/certs
```

## Access Control

SFTP access uses the same IAM system as the S3 API:

- **Root credentials**: full access to all buckets
- **IAM users**: access controlled by IAM policies
- **Service accounts**: programmatic SFTP access for automation

### Creating an SFTP User

```sh
mc admin user add myminio sftpuser sftppassword
mc admin policy attach myminio readwrite --user sftpuser
```

The user can then connect via SFTP:

```sh
sftp -P 2022 sftpuser@your-server
```

## Client Compatibility

Obstor's SFTP gateway works with all standard SFTP clients:

| Client | Platform | Status |
|--------|----------|--------|
| OpenSSH sftp | Linux, macOS | Supported |
| WinSCP | Windows | Supported |
| FileZilla | Cross-platform | Supported |
| Cyberduck | macOS, Windows | Supported |
| rclone | Cross-platform | Supported |
| Paramiko (Python) | Cross-platform | Supported |

## Limitations

- SFTP does not support S3-specific features like object versioning, lifecycle rules, or multipart listing
- Symbolic links are not supported
- File permissions are mapped to S3 ACLs where possible
- Maximum concurrent SFTP connections depend on server resources

## Explore Further

- [Obstor Distributed Mode](../distributed/README.md) - deploy across multiple nodes
- [TLS Configuration](../tls/README.md) - secure your deployment
- [IAM & Policies](../multi-user/README.md) - manage user access
- [S3 Gateway](./s3.md) - S3 API compatibility details

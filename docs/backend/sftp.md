# Obstor SFTP Backend

Obstor SFTP Backend adds S3 and [other supported protocol](../protocols/README.md) compatibility to any existing SFTP server. Use your SFTP server as the storage backend while accessing data through the S3 API and other supported protocols.

> Looking for the built-in SFTP server? See [SFTP Protocol](../protocols/sftp.md).

## Run Obstor Backend for SFTP Storage

### Using Binary

```sh
export OBSTOR_ROOT_USER=accesskey
export OBSTOR_ROOT_PASSWORD=secretkey
export OBSTOR_BACKEND_SFTP_USER=sftpuser
export OBSTOR_BACKEND_SFTP_PASSWORD=sftppassword
obstor backend sftp sftp-server:22/data
```

### Using SSH Key Authentication

```sh
export OBSTOR_ROOT_USER=accesskey
export OBSTOR_ROOT_PASSWORD=secretkey
export OBSTOR_BACKEND_SFTP_USER=sftpuser
export OBSTOR_BACKEND_SFTP_KEY=/path/to/id_rsa
obstor backend sftp sftp-server:22/data
```

### Using Docker

```sh
docker run -p 9000:9000 \
  -e "OBSTOR_ROOT_USER=accesskey" \
  -e "OBSTOR_ROOT_PASSWORD=secretkey" \
  -e "OBSTOR_BACKEND_SFTP_USER=sftpuser" \
  -e "OBSTOR_BACKEND_SFTP_PASSWORD=sftppassword" \
  ghcr.io/cloudment/obstor backend sftp sftp-server:22/data
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `OBSTOR_BACKEND_SFTP_USER` | SSH username for the remote SFTP server | Yes |
| `OBSTOR_BACKEND_SFTP_PASSWORD` | SSH password (if not using key auth) | One of password/key |
| `OBSTOR_BACKEND_SFTP_KEY` | Path to SSH private key file | One of password/key |

### Endpoint Format

The SFTP endpoint follows this format:

```
host:port/base-path
```

- `host:port` - SFTP server address (port defaults to 22 if omitted)
- `/base-path` - Optional base directory on the SFTP server (defaults to `/`)

Examples:
```
sftp-server:22          # root directory on port 22
sftp-server:2222/data   # /data directory on port 2222
192.168.1.100:22/mnt    # /mnt directory on a specific IP
```

### Bucket Mapping

Top-level directories on the SFTP server are exposed as S3 buckets. Files within those directories become S3 objects:

```
SFTP Server                    S3 API
/data/                    -->  (root)
/data/photos/             -->  photos (bucket)
/data/photos/sunset.jpg   -->  photos/sunset.jpg (object)
/data/backups/            -->  backups (bucket)
/data/backups/db.sql.gz   -->  backups/db.sql.gz (object)
```

## Test Using Obstor Client `mc`

```sh
mc alias set mysftp http://localhost:9000 accesskey secretkey

mc ls mysftp
mc mb mysftp/newbucket
mc cp myfile.txt mysftp/newbucket/myfile.txt
```

## Known Limitations

- No bucket policy support
- No bucket notification support
- No server-side encryption
- No server-side copy (copy operations read from source and write to destination)
- No compression support
- No multipart upload support (handled transparently by the backend)
- Delete bucket only works on empty directories

## Explore Further

- [Supported Protocols](../protocols/README.md) - S3, SFTP, and more
- [Obstor Distributed Mode](../distributed/README.md)
- [TLS Configuration](../tls/README.md)
- [IAM & Policies](../multi-user/README.md)
- [SFTP Protocol](../protocols/sftp.md) - built-in SFTP server for accessing Obstor via SFTP

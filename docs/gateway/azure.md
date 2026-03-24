# Obstor Azure Gateway [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord)
Obstor Gateway adds Amazon S3 compatibility to Microsoft Azure Blob Storage.

## Run Obstor Gateway for Microsoft Azure Blob Storage
### Using Docker
```
docker run -p 9000:9000 --name azure-s3 \
 -e "OBSTOR_ROOT_USER=azurestorageaccountname" \
 -e "OBSTOR_ROOT_PASSWORD=azurestorageaccountkey" \
 ghcr.io/cloudment/obstor gateway azure
```

### Using Binary
```
export OBSTOR_ROOT_USER=azureaccountname
export OBSTOR_ROOT_PASSWORD=azureaccountkey
obstor gateway azure
```
## Test using Browser Dashboard
Obstor Gateway comes with an embedded web based object browser. Point your web browser to http://127.0.0.1:9000 to ensure that your server has started successfully.

![Screenshot](https://raw.githubusercontent.com/cloudment/obstor/main/docs/screenshots/dashboard.png)
## Test using Obstor Client `mc`
`mc` provides a modern alternative to UNIX commands such as ls, cat, cp, mirror, diff etc. It supports filesystems and Amazon S3 compatible cloud storage services.

### Configure `mc`
```
mc alias set myazure http://gateway-ip:9000 azureaccountname azureaccountkey
```

### List containers on Microsoft Azure
```
mc ls myazure
[2017-02-22 01:50:43 PST]     0B ferenginar/
[2017-02-26 21:43:51 PST]     0B my-container/
[2017-02-26 22:10:11 PST]     0B test-container1/
```

### Use custom access/secret keys

If you do not want to share the credentials of the Azure blob storage with your users/applications, you can set the original credentials in the shell environment using `AZURE_STORAGE_ACCOUNT` and `AZURE_STORAGE_KEY` variables and assign different access/secret keys to `OBSTOR_ROOT_USER` and `OBSTOR_ROOT_PASSWORD`.

### Known limitations
Gateway inherits the following Azure limitations:

- Only read-only bucket policy supported at bucket level, all other variations will return API Notimplemented error.
- Bucket names with "." in the bucket name are not supported.
- Non-empty buckets get removed on a DeleteBucket() call.
- _List Multipart Uploads_ always returns empty list.

Other limitations:

- Bucket notification APIs are not supported.

## Explore Further
- [`mc` command-line interface](https://obstor.net/docs/obstor-client-quickstart-guide)
- [`aws` command-line interface](https://obstor.net/docs/aws-cli-with-obstor)
- [`minio-go` Go SDK](https://obstor.net/docs/golang-client-quickstart-guide)

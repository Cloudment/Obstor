# Obstor GCS Gateway [![Discord](https://pgg.net/discord?type=svg)](https://pgg.net/discord)

Obstor GCS Gateway allows you to access Google Cloud Storage (GCS) with Amazon S3-compatible APIs

- [Run Obstor Gateway for GCS](#run-obstor-gateway-for-gcs)
- [Test Using Obstor Browser](#test-using-obstor-browser)
- [Test Using Obstor Client](#test-using-obstor-client)

## <a name="run-obstor-gateway-for-gcs"></a>1. Run Obstor Gateway for GCS

### 1.1 Create a Service Account key for GCS and get the Credentials File
1. Navigate to the [API Console Credentials page](https://console.developers.google.com/project/_/apis/credentials).
2. Select a project or create a new project. Note the project ID.
3. Select the **Create credentials** dropdown on the **Credentials** page, and click **Service account key**.
4. Select **New service account** from the **Service account** dropdown.
5. Populate the **Service account name** and **Service account ID**.
6. Click the dropdown for the **Role** and choose **Storage** > **Storage Admin** *(Full control of GCS resources)*.
7. Click the **Create** button to download a credentials file and rename it to `credentials.json`.

**Note:** For alternate ways to set up *Application Default Credentials*, see [Setting Up Authentication for Server to Server Production Applications](https://developers.google.com/identity/protocols/application-default-credentials).

### 1.2 Run Obstor GCS Gateway Using Docker
```sh
docker run -p 9000:9000 --name gcs-s3 \
 -v /path/to/credentials.json:/credentials.json \
 -e "GOOGLE_APPLICATION_CREDENTIALS=/credentials.json" \
 -e "OBSTOR_ROOT_USER=obstoraccountname" \
 -e "OBSTOR_ROOT_PASSWORD=obstoraccountkey" \
 obstor/obstor gateway gcs yourprojectid
```

### 1.3 Run Obstor GCS Gateway Using the Obstor Binary

```sh
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
export OBSTOR_ROOT_USER=obstoraccesskey
export OBSTOR_ROOT_PASSWORD=obstorsecretkey
obstor gateway gcs yourprojectid
```

## <a name="test-using-obstor-browser"></a>2. Test Using Obstor Browser

Obstor Gateway comes with an embedded web-based object browser that outputs content to http://127.0.0.1:9000. To test that Obstor Gateway is running, open a web browser, navigate to http://127.0.0.1:9000, and ensure that the object browser is displayed.

![Screenshot](https://github.com/cloudment/obstor/blob/master/docs/screenshots/obstor-browser-gateway.png?raw=true)

## <a name="test-using-obstor-client"></a>3. Test Using Obstor Client

Obstor Client is a command-line tool called `mc` that provides UNIX-like commands for interacting with the server (e.g. ls, cat, cp, mirror, diff, find, etc.).  `mc` supports file systems and Amazon S3-compatible cloud storage services (AWS Signature v2 and v4).

### 3.1 Configure the Gateway using Obstor Client

Use the following command to configure the gateway:

```sh
mc alias set mygcs http://gateway-ip:9000 minioaccesskey miniosecretkey
```

### 3.2 List Containers on GCS

Use the following command to list the containers on GCS:

```sh
mc ls mygcs
```

A response similar to this one should be displayed:

```
[2017-02-22 01:50:43 PST]     0B ferenginar/
[2017-02-26 21:43:51 PST]     0B my-container/
[2017-02-26 22:10:11 PST]     0B test-container1/
```

### 3.3 Known limitations
Obstor Gateway has the following limitations when used with GCS:

* It only supports read-only and write-only bucket policies at the bucket level; all other variations will return `API Not implemented`.
* The `List Multipart Uploads` and `List Object parts` commands always return empty lists. Therefore, the client must store all of the parts that it has uploaded and use that information when invoking the `_Complete Multipart Upload` command.

Other limitations:

* Bucket notification APIs are not supported.

## <a name="explore-further"></a>4. Explore Further
- [`mc` command-line interface](https://pgg.net/docs/obstor/obstor-client-quickstart-guide)
- [`aws` command-line interface](https://pgg.net/docs/obstor/aws-cli-with-obstor)
- [`minio-go` Go SDK](https://pgg.net/docs/obstor/golang-client-quickstart-guide)

# Obstor File Browser

``Obstor Browser`` provides minimal set of UI to manage buckets and objects on ``obstor`` server. ``Obstor Browser`` is written in javascript and released under [Apache 2.0 License](./LICENSE).


## Installation

### Install node
```sh
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.34.0/install.sh | bash
exec -l $SHELL
nvm install stable
```

### Install node dependencies
```sh
npm install
```

## Generating Assets

```sh
npm run release
```

This generates `production` in the current directory.


## Run Obstor Browser with live reload

### Run Obstor Browser with live reload

```sh
npm run dev
```

Open [http://localhost:8080/obstor/](http://localhost:8080/obstor/) in your browser to play with the application.

### Run Obstor Browser with live reload on custom port

Edit `browser/webpack.config.js`

```diff
diff --git a/browser/webpack.config.js b/browser/webpack.config.js
index 3ccdaba..9496c56 100644
--- a/browser/webpack.config.js
+++ b/browser/webpack.config.js
@@ -58,6 +58,7 @@ var exports = {
     historyApiFallback: {
       index: '/obstor/'
     },
+    port: 8888,
     proxy: {
       '/obstor/webrpc': {
         target: 'http://localhost:9000',
@@ -97,7 +98,7 @@ var exports = {
 if (process.env.NODE_ENV === 'dev') {
   exports.entry = [
     'webpack/hot/dev-server',
-    'webpack-dev-server/client?http://localhost:8080',
+    'webpack-dev-server/client?http://localhost:8888',
     path.resolve(__dirname, 'app/index.js')
   ]
 }
```

```sh
npm run dev
```

Open [http://localhost:8888/obstor/](http://localhost:8888/obstor/) in your browser to play with the application.

### Run Obstor Browser with live reload on any IP

Edit `browser/webpack.config.js`

```diff
diff --git a/browser/webpack.config.js b/browser/webpack.config.js
index 8bdbba53..139f6049 100644
--- a/browser/webpack.config.js
+++ b/browser/webpack.config.js
@@ -71,6 +71,7 @@ var exports = {
     historyApiFallback: {
       index: '/obstor/'
     },
+    host: '0.0.0.0',
     proxy: {
       '/obstor/webrpc': {
         target: 'http://localhost:9000',
```

```sh
npm run dev
```

Open [http://IP:8080/obstor/](http://IP:8080/obstor/) in your browser to play with the application.


## Run tests

    npm run test


## Docker development environment

This approach will download the sources on your machine such that you are able to use your IDE or editor of choice.
A Docker container will be used in order to provide a controlled build environment without messing with your host system.

### Prepare host system

Install [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) and [Docker](https://docs.docker.com/get-docker/).

### Development within container

Prepare and build container
```
git clone git@github.com:cloudment/obstor.git
cd obstor
docker build -t obstor-dev -f Dockerfile.dev.browser .
```

Run container, build and run core
```sh
docker run -it --rm --name obstor-dev -v "$PWD":/obstor obstor-dev

cd /obstor/browser
npm install
npm run release
cd /obstor
make
./obstor server /data
```
Note `Endpoint` IP (the one which is _not_ `127.0.0.1`), `AccessKey` and `SecretKey` (both default to `obstoradmin`) in order to enter them in the browser later.


Open another terminal.
Connect to container
```sh
docker exec -it obstor-dev bash
```

Apply patch to allow access from outside container
```sh
cd /obstor
git apply --ignore-whitespace <<EOF
diff --git a/browser/webpack.config.js b/browser/webpack.config.js
index 8bdbba53..139f6049 100644
--- a/browser/webpack.config.js
+++ b/browser/webpack.config.js
@@ -71,6 +71,7 @@ var exports = {
     historyApiFallback: {
       index: '/obstor/'
     },
+    host: '0.0.0.0',
     proxy: {
       '/obstor/webrpc': {
         target: 'http://localhost:9000',
EOF
```

Build and run frontend with auto-reload
```sh
cd /obstor/browser
npm install
npm run dev
```

Open [http://IP:8080/obstor/](http://IP:8080/obstor/) in your browser to play with the application.


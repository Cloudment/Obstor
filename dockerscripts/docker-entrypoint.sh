#!/bin/sh
#
# MinIO Cloud Storage, (C) 2019 MinIO, Inc.
# PGG Obstor, (C) 2021-2026 PGG, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# If command starts with an option, prepend obstor.
if [ "${1}" != "obstor" ]; then
    if [ -n "${1}" ]; then
        set -- obstor "$@"
    fi
fi

docker_secrets_env() {
    if [ -f "$OBSTOR_ROOT_USER_FILE" ]; then
        ROOT_USER_FILE="$OBSTOR_ROOT_USER_FILE"
    else
        ROOT_USER_FILE="/run/secrets/$OBSTOR_ROOT_USER_FILE"
    fi
    if [ -f "$OBSTOR_ROOT_PASSWORD_FILE" ]; then
        SECRET_KEY_FILE="$OBSTOR_ROOT_PASSWORD_FILE"
    else
        SECRET_KEY_FILE="/run/secrets/$OBSTOR_ROOT_PASSWORD_FILE"
    fi

    if [ -f "$ROOT_USER_FILE" ] && [ -f "$SECRET_KEY_FILE" ]; then
        if [ -f "$ROOT_USER_FILE" ]; then
            OBSTOR_ROOT_USER="$(cat "$ROOT_USER_FILE")"
            export OBSTOR_ROOT_USER
        fi
        if [ -f "$SECRET_KEY_FILE" ]; then
            OBSTOR_ROOT_PASSWORD="$(cat "$SECRET_KEY_FILE")"
            export OBSTOR_ROOT_PASSWORD
        fi
    fi
}

## Set KMS_MASTER_KEY from docker secrets if provided
docker_kms_encryption_env() {
    if [ -f "$OBSTOR_KMS_SECRET_KEY_FILE" ]; then
        KMS_SECRET_KEY_FILE="$OBSTOR_KMS_SECRET_KEY_FILE"
    else
        KMS_SECRET_KEY_FILE="/run/secrets/$OBSTOR_KMS_SECRET_KEY_FILE"
    fi

    if [ -f "$KMS_SECRET_KEY_FILE" ]; then
        OBSTOR_KMS_SECRET_KEY="$(cat "$KMS_SECRET_KEY_FILE")"
        export OBSTOR_KMS_SECRET_KEY
    fi
}

# su-exec to requested user, if service cannot run exec will fail.
docker_switch_user() {
    if [ -n "${OBSTOR_USERNAME}" ] && [ -n "${OBSTOR_GROUPNAME}" ]; then
        if [ -n "${OBSTOR_UID}" ] && [ -n "${OBSTOR_GID}" ]; then
            groupadd -g "$OBSTOR_GID" "$OBSTOR_GROUPNAME" && \
                useradd -u "$OBSTOR_UID" -g "$OBSTOR_GROUPNAME" "$OBSTOR_USERNAME"
        else
            groupadd "$OBSTOR_GROUPNAME" && \
                useradd -g "$OBSTOR_GROUPNAME" "$OBSTOR_USERNAME"
        fi
        exec setpriv --reuid="${OBSTOR_USERNAME}" --regid="${OBSTOR_GROUPNAME}" --keep-groups "$@"
    else
        exec "$@"
    fi
}

# Start frontend if available.
start_frontend() {
    if [ -f /opt/frontend/server.js ]; then
        FRONTEND_PORT=9001
        for arg in "$@"; do
            case "$arg" in
                --frontend-address=*) FRONTEND_PORT="${arg#*=:}" ;;
            esac
        done
        PORT=$FRONTEND_PORT \
        OBSTOR_ENDPOINT=${OBSTOR_ENDPOINT:-http://127.0.0.1:9000} \
        OBSTOR_HOST=${OBSTOR_HOST:-127.0.0.1:9000} \
        node /opt/frontend/server.js &
    fi
}

## Load secrets
docker_secrets_env
docker_kms_encryption_env

## Start frontend
start_frontend "$@"

## Switch to user and exec obstor
docker_switch_user "$@"

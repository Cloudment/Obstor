/*
 * MinIO Cloud Storage, (C) 2019 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

// Config value separator
const (
	ValueSeparator = ","
)

// Top level common ENVs
const (
	EnvAccessKey    = "OBSTOR_ACCESS_KEY"
	EnvSecretKey    = "OBSTOR_SECRET_KEY"
	EnvRootUser     = "OBSTOR_ROOT_USER"
	EnvRootPassword = "OBSTOR_ROOT_PASSWORD"

	EnvBrowser    = "OBSTOR_BROWSER"
	EnvDomain     = "OBSTOR_DOMAIN"
	EnvRegionName = "OBSTOR_REGION_NAME"
	EnvPublicIPs  = "OBSTOR_PUBLIC_IPS"
	EnvFSOSync    = "OBSTOR_FS_OSYNC"
	EnvArgs       = "OBSTOR_ARGS"
	EnvDNSWebhook = "OBSTOR_DNS_WEBHOOK_ENDPOINT"

	EnvUpdate = "OBSTOR_UPDATE"

	EnvKMSMasterKey  = "OBSTOR_KMS_MASTER_KEY" // legacy
	EnvKMSSecretKey  = "OBSTOR_KMS_SECRET_KEY"
	EnvKESEndpoint   = "OBSTOR_KMS_KES_ENDPOINT"
	EnvKESKeyName    = "OBSTOR_KMS_KES_KEY_NAME"
	EnvKESClientKey  = "OBSTOR_KMS_KES_KEY_FILE"
	EnvKESClientCert = "OBSTOR_KMS_KES_CERT_FILE"
	EnvKESServerCA   = "OBSTOR_KMS_KES_CAPATH"

	EnvEndpoints = "OBSTOR_ENDPOINTS" // legacy
	EnvWorm      = "OBSTOR_WORM"      // legacy
	EnvRegion    = "OBSTOR_REGION"    // legacy
)

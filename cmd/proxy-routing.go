/*
 * MinIO Cloud Storage, (C) 2016 MinIO, Inc.
 * PGG Obstor, (C) 2021-2026 PGG, Inc.
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

package cmd

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

func parseRequestToken(token string) (subToken string, nodeIndex int) {
	if token == "" {
		return token, -1
	}
	i := strings.Index(token, "@")
	if i < 0 {
		return token, -1
	}
	nodeIndex, err := strconv.Atoi(token[i+1:])
	if err != nil {
		return token, -1
	}
	subToken = token[:i]
	return subToken, nodeIndex
}

func proxyRequestByToken(ctx context.Context, w http.ResponseWriter, r *http.Request, token string) (string, bool) {
	subToken, nodeIndex := parseRequestToken(token)
	if nodeIndex > 0 {
		return subToken, proxyRequestByNodeIndex(ctx, w, r, nodeIndex)
	}
	return subToken, false
}

func proxyRequestByNodeIndex(ctx context.Context, w http.ResponseWriter, r *http.Request, index int) (success bool) {
	if len(globalProxyEndpoints) == 0 {
		return false
	}
	if index < 0 || index >= len(globalProxyEndpoints) {
		return false
	}
	ep := globalProxyEndpoints[index]
	if ep.IsLocal {
		return false
	}
	return proxyRequest(ctx, w, r, ep)
}

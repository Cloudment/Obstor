/*
 * MinIO Cloud Storage, (C) 2019 MinIO, Inc.
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
	"net/http"
	"sync/atomic"
)

// RequestStats - counts for Get and Head requests
type RequestStats struct {
	Get  uint64 `json:"Get"`
	Head uint64 `json:"Head"`
	Put  uint64 `json:"Put"`
	Post uint64 `json:"Post"`
}

// IncBytesReceived - Increase total bytes received from backend storage
func (s *BackendMetrics) IncBytesReceived(n uint64) {
	atomic.AddUint64(&s.bytesReceived, n)
}

// GetBytesReceived - Get total bytes received from backend storage
func (s *BackendMetrics) GetBytesReceived() uint64 {
	return atomic.LoadUint64(&s.bytesReceived)
}

// IncBytesSent - Increase total bytes sent to backend storage
func (s *BackendMetrics) IncBytesSent(n uint64) {
	atomic.AddUint64(&s.bytesSent, n)
}

// GetBytesSent - Get total bytes received from backend storage
func (s *BackendMetrics) GetBytesSent() uint64 {
	return atomic.LoadUint64(&s.bytesSent)
}

// IncRequests - Increase request count sent to backend storage by 1
func (s *BackendMetrics) IncRequests(method string) {
	// Only increment for Head & Get requests, else no op
	switch method {
	case http.MethodGet:
		atomic.AddUint64(&s.requestStats.Get, 1)
	case http.MethodHead:
		atomic.AddUint64(&s.requestStats.Head, 1)
	case http.MethodPut:
		atomic.AddUint64(&s.requestStats.Put, 1)
	case http.MethodPost:
		atomic.AddUint64(&s.requestStats.Post, 1)
	}
}

// GetRequests - Get total number of Get & Headrequests sent to backend storage
func (s *BackendMetrics) GetRequests() RequestStats {
	return s.requestStats
}

// NewMetrics - Prepare new BackendMetrics structure
func NewMetrics() *BackendMetrics {
	return &BackendMetrics{}
}

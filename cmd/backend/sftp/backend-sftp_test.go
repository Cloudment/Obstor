/*
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

package sftp

import (
	"fmt"
	"testing"
)

func TestParseSFTPEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		host     string
		port     string
		basePath string
	}{
		{"host:22", "host", "22", "/"},
		{"host:22/data", "host", "22", "/data"},
		{"host:2222/mnt/storage", "host", "2222", "/mnt/storage"},
		{"192.168.1.100:22", "192.168.1.100", "22", "/"},
		{"192.168.1.100:22/backup", "192.168.1.100", "22", "/backup"},
		{"myserver", "myserver", "22", "/"},
		{"myserver/data", "myserver", "22", "/data"},
	}

	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			host, port, basePath := parseSFTPEndpoint(tt.endpoint)
			if host != tt.host {
				t.Errorf("host: got %q, want %q", host, tt.host)
			}
			if port != tt.port {
				t.Errorf("port: got %q, want %q", port, tt.port)
			}
			if basePath != tt.basePath {
				t.Errorf("basePath: got %q, want %q", basePath, tt.basePath)
			}
		})
	}
}

func TestSftpIsValidBucketName(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"mybucket", true},
		{"my-bucket", true},
		{"my.bucket", true},
		{"123bucket", true},
		{"My-Bucket", false},    // uppercase not allowed in strict mode
		{".obstor.sys", false},  // starts with dot
		{"a", false},            // too short
		{"ab", false},           // too short
		{"abc", true},           // minimum length
		{"-bucket", false},      // starts with hyphen
		{"bucket-", false},      // ends with hyphen
		{"bucket..name", false}, // consecutive dots
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sftpIsValidBucketName(tt.name)
			if got != tt.valid {
				t.Errorf("sftpIsValidBucketName(%q) = %v, want %v", tt.name, got, tt.valid)
			}
		})
	}
}

func TestIsReservedOrInvalidBucket(t *testing.T) {
	tests := []struct {
		name     string
		strict   bool
		reserved bool
	}{
		{".obstor.sys", false, true},
		{"obstor", false, true},
		{"mybucket", false, false},
		{"my-bucket", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isReservedOrInvalidBucket(tt.name, tt.strict)
			if got != tt.reserved {
				t.Errorf("isReservedOrInvalidBucket(%q, %v) = %v, want %v", tt.name, tt.strict, got, tt.reserved)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected bool
	}{
		{"nil", "", false},
		{"not empty", "directory is not empty", true},
		{"not empty alt", "not empty", true},
		{"other error", "permission denied", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.errMsg != "" {
				err = fmt.Errorf("%s", tt.errMsg)
			}
			got := isNotEmpty(err)
			if got != tt.expected {
				t.Errorf("isNotEmpty(%v) = %v, want %v", err, got, tt.expected)
			}
		})
	}
}

/*
 * MinIO Cloud Storage, (C) 2016, 2017 MinIO, Inc.
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
	"os"
	"testing"
)

// Test printing Backend common message.
func TestPrintBackendCommonMessage(t *testing.T) {
	obj, fsDir, err := prepareFS()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fsDir)
	if err = newTestConfig(globalObstorDefaultRegion, obj); err != nil {
		t.Fatal(err)
	}

	apiEndpoints := []string{"http://127.0.0.1:9000"}
	printBackendCommonMsg(apiEndpoints)
}

// Test print backend startup message.
func TestPrintBackendStartupMessage(t *testing.T) {
	obj, fsDir, err := prepareFS()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fsDir)
	if err = newTestConfig(globalObstorDefaultRegion, obj); err != nil {
		t.Fatal(err)
	}

	apiEndpoints := []string{"http://127.0.0.1:9000"}
	printBackendStartupMessage(apiEndpoints, "azure")
}

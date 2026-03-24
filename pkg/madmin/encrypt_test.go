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
 *
 */

package madmin

import (
	"bytes"
	"fmt"
	"testing"
)

var encryptDataTests = []struct {
	Password string
	Data     []byte
}{
	{Password: "", Data: nil},
	{Password: "", Data: make([]byte, 256)},
	{Password: `xPl.8/rhR"Q_1xLt`, Data: make([]byte, 32)},
	{Password: "m69?yz4W-!k+7p0", Data: make([]byte, 1024*1024)},
	{Password: `7h5oU4$te{;K}fgqlI^]`, Data: make([]byte, 256)},
}

func TestEncryptData(t *testing.T) {
	for i, test := range encryptDataTests {
		i, test := i, test
		t.Run(fmt.Sprintf("Test-%d", i), func(t *testing.T) {
			ciphertext, err := EncryptData(test.Password, test.Data)
			if err != nil {
				t.Fatalf("Failed to encrypt data: %v", err)
			}

			plaintext, err := DecryptData(test.Password, bytes.NewReader(ciphertext))
			if err != nil {
				t.Fatalf("Failed to decrypt data: %v", err)
			}
			if !bytes.Equal(plaintext, test.Data) {
				t.Fatal("Decrypt plaintext does not match origin data")
			}
		})
	}
}

var decryptDataTests = []struct {
	Password string
	Data     []byte
}{
	{Password: "", Data: nil},
	{Password: "", Data: make([]byte, 256)},
	{Password: `xPl.8/rhR"Q_1xLt`, Data: make([]byte, 32)},
	{Password: `7h5oU4$te{;K}fgqlI^]`, Data: make([]byte, 256)},
}

func TestDecryptData(t *testing.T) {
	for i, test := range decryptDataTests {
		i, test := i, test
		t.Run(fmt.Sprintf("Test-%d", i), func(t *testing.T) {
			ciphertext, err := EncryptData(test.Password, test.Data)
			if err != nil {
				t.Fatalf("Failed to encrypt data: %v", err)
			}
			plaintext, err := DecryptData(test.Password, bytes.NewReader(ciphertext))
			if err != nil {
				t.Fatalf("Failed to decrypt data: %v", err)
			}
			if !bytes.Equal(plaintext, test.Data) {
				t.Fatal("Decrypted data does not match original")
			}
		})
	}
}

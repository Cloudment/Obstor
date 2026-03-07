/*
 * MinIO Cloud Storage, (C) 2016-2020 MinIO, Inc.
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
	"fmt"
)

// ReadBlock reads a block from the first available source node.
// It tries each source in order (prefer local disks first) and
// returns on the first successful read. No reconstruction is needed
// since every copy is identical.
func (br *BlockReplicator) ReadBlock(ctx context.Context, blockHash string, sources []StorageAPI) ([]byte, error) {
	var lastErr error
	for _, disk := range sources {
		if disk == nil {
			continue
		}
		data, err := disk.ReadBlock(ctx, blockHash)
		if err != nil {
			lastErr = err
			continue
		}
		// Verify integrity.
		if err := verifyBlockHash(data, blockHash); err != nil {
			lastErr = err
			continue
		}
		return data, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("block %s: no available source disks", blockHash)
}

// HasBlock checks if at least one source has the block.
func (br *BlockReplicator) HasBlock(ctx context.Context, blockHash string, sources []StorageAPI) bool {
	for _, disk := range sources {
		if disk == nil {
			continue
		}
		has, err := disk.HasBlock(ctx, blockHash)
		if err == nil && has {
			return true
		}
	}
	return false
}

// CountBlockCopies returns how many sources have the block.
func (br *BlockReplicator) CountBlockCopies(ctx context.Context, blockHash string, sources []StorageAPI) int {
	count := 0
	for _, disk := range sources {
		if disk == nil {
			continue
		}
		has, err := disk.HasBlock(ctx, blockHash)
		if err == nil && has {
			count++
		}
	}
	return count
}

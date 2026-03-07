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
)

// HealBlock copies a block from any source that has it to all targets
// that are missing it. Returns the number of copies repaired.
func (br *BlockReplicator) HealBlock(ctx context.Context, blockHash string, sources []StorageAPI, targets []StorageAPI) (int, error) {
	// Read the block from any available source.
	data, err := br.ReadBlock(ctx, blockHash, sources)
	if err != nil {
		return 0, err
	}

	repaired := 0
	for _, disk := range targets {
		if disk == nil {
			continue
		}
		has, _ := disk.HasBlock(ctx, blockHash)
		if has {
			continue
		}
		if err := disk.WriteBlock(ctx, blockHash, data); err == nil {
			repaired++
		}
	}
	return repaired, nil
}

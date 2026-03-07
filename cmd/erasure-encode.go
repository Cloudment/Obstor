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
	"sync"
)

// WriteBlock writes the same block data to multiple target storage nodes
// in parallel. It returns nil if at least writeQuorum targets succeed.
func (br *BlockReplicator) WriteBlock(ctx context.Context, blockHash string, data []byte, targets []StorageAPI, writeQuorum int) error {
	errs := make([]error, len(targets))
	var wg sync.WaitGroup

	for i, disk := range targets {
		if disk == nil {
			errs[i] = errDiskNotFound
			continue
		}
		wg.Add(1)
		go func(i int, disk StorageAPI) {
			defer wg.Done()
			errs[i] = disk.WriteBlock(ctx, blockHash, data)
		}(i, disk)
	}
	wg.Wait()

	var successes int
	for _, err := range errs {
		if err == nil {
			successes++
		}
	}
	if successes >= writeQuorum {
		return nil
	}
	return reduceWriteQuorumErrs(ctx, errs, objectOpIgnoredErrs, writeQuorum)
}

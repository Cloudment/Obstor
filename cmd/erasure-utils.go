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
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/cloudment/obstor/cmd/logger"
)

// writeDataBlocks writes data blocks to dst. Legacy function kept for test compatibility.
func writeDataBlocks(ctx context.Context, dst io.Writer, enBlocks [][]byte, dataBlocks int, offset int64, length int64) (int64, error) {
	if offset < 0 || length < 0 {
		logger.LogIf(ctx, errUnexpected)
		return 0, errUnexpected
	}
	if len(enBlocks) < dataBlocks {
		return 0, fmt.Errorf("too few blocks: have %d, need %d", len(enBlocks), dataBlocks)
	}

	write := length
	var totalWritten int64
	for _, block := range enBlocks[:dataBlocks] {
		if offset >= int64(len(block)) {
			offset -= int64(len(block))
			continue
		} else {
			block = block[offset:]
			offset = 0
		}
		if write < int64(len(block)) {
			n, err := io.Copy(dst, bytes.NewReader(block[:write]))
			if err != nil {
				return 0, err
			}
			totalWritten += n
			break
		}
		n, err := io.Copy(dst, bytes.NewReader(block))
		if err != nil {
			return 0, err
		}
		write -= n
		totalWritten += n
	}
	return totalWritten, nil
}

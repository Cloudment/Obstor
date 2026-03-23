/*
 * MinIO Cloud Storage, (C) 2017 MinIO, Inc.
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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path"
)

// Erasure - legacy erasure coding details. Kept for backward compatibility
// with objects written using Reed-Solomon encoding.
type Erasure struct {
	dataBlocks, parityBlocks int
	blockSize                int64
}

// NewErasure creates a new Erasure for legacy erasure-coded reads.
func NewErasure(ctx context.Context, dataBlocks, parityBlocks int, blockSize int64) (e Erasure, err error) {
	if dataBlocks <= 0 || parityBlocks < 0 {
		return e, fmt.Errorf("invalid erasure config: data=%d parity=%d", dataBlocks, parityBlocks)
	}
	return Erasure{
		dataBlocks:   dataBlocks,
		parityBlocks: parityBlocks,
		blockSize:    blockSize,
	}, nil
}

// ShardSize returns actual shard size from erasure blockSize.
func (e *Erasure) ShardSize() int64 {
	return ceilFrac(e.blockSize, int64(e.dataBlocks))
}

// ShardFileSize returns final erasure size from original size.
func (e *Erasure) ShardFileSize(totalLength int64) int64 {
	if totalLength == 0 {
		return 0
	}
	if totalLength == -1 {
		return -1
	}
	numShards := totalLength / e.blockSize
	lastBlockSize := totalLength % e.blockSize
	lastShardSize := ceilFrac(lastBlockSize, int64(e.dataBlocks))
	return numShards*e.ShardSize() + lastShardSize
}

// ShardFileOffset returns the effective offset where erasure reading begins.
func (e *Erasure) ShardFileOffset(startOffset, length, totalLength int64) int64 {
	shardSize := e.ShardSize()
	shardFileSize := e.ShardFileSize(totalLength)
	endShard := (startOffset + length) / e.blockSize
	tillOffset := endShard*shardSize + shardSize
	if tillOffset > shardFileSize {
		tillOffset = shardFileSize
	}
	return tillOffset
}

// Encode reads from src, erasure-encodes the data and writes to the writers.
// This is a stub that supports legacy reads. New writes use block replication.
func (e *Erasure) Encode(ctx context.Context, src io.Reader, writers []io.Writer, buf []byte, quorum int) (total int64, err error) {
	return 0, fmt.Errorf("erasure encoding is no longer supported; use block replication")
}

// Decode reads from readers, reconstructs data and writes to the writer.
// This is a stub that supports legacy reads. New reads use block replication.
func (e *Erasure) Decode(ctx context.Context, writer io.Writer, readers []io.ReaderAt, offset, length, totalLength int64, prefer []bool) (written int64, derr error) {
	return 0, fmt.Errorf("erasure decoding is no longer supported; use block replication")
}

// DecodeDataBlocks is a stub for legacy compatibility.
func (e *Erasure) DecodeDataBlocks(data [][]byte) error {
	return fmt.Errorf("erasure decoding is no longer supported; use block replication")
}

// EncodeData is a stub for legacy compatibility.
func (e *Erasure) EncodeData(ctx context.Context, data []byte) ([][]byte, error) {
	return nil, fmt.Errorf("erasure encoding is no longer supported; use block replication")
}

// Heal is a stub for legacy compatibility.
func (e *Erasure) Heal(ctx context.Context, readers []io.ReaderAt, writers []io.Writer, size int64) error {
	return fmt.Errorf("erasure healing is no longer supported; use block replication")
}

// DecodeDataAndParityBlocks is a stub for legacy compatibility.
func (e *Erasure) DecodeDataAndParityBlocks(ctx context.Context, data [][]byte) error {
	return fmt.Errorf("erasure decoding is no longer supported; use block replication")
}

// BlockRef identifies a content-addressed data block within an object.
type BlockRef struct {
	Hash  string `json:"h" msg:"h"` // SHA-256 hex digest
	Size  int64  `json:"s" msg:"s"` // block data length
	Index int    `json:"i" msg:"i"` // ordinal position in object
}

// BlockReplicator handles chunking data into content-addressed blocks
// and coordinating N-copy writes/reads across storage nodes.
type BlockReplicator struct {
	blockSize    int64
	replicaCount int
}

// NewBlockReplicator creates a new BlockReplicator.
func NewBlockReplicator(blockSize int64, replicaCount int) *BlockReplicator {
	return &BlockReplicator{
		blockSize:    blockSize,
		replicaCount: replicaCount,
	}
}

// ChunkAndHash reads from src in blockSize chunks, returning BlockRefs
// and the raw data for each block. The caller is responsible for writing
// the block data to storage nodes.
func (br *BlockReplicator) ChunkAndHash(src io.Reader) ([]BlockRef, [][]byte, error) {
	var refs []BlockRef
	var blocks [][]byte
	buf := make([]byte, br.blockSize)
	idx := 0

	for {
		n, err := io.ReadFull(src, buf)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, nil, err
		}
		if n == 0 {
			break
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		hash := sha256.Sum256(data)
		hashHex := hex.EncodeToString(hash[:])

		refs = append(refs, BlockRef{
			Hash:  hashHex,
			Size:  int64(n),
			Index: idx,
		})
		blocks = append(blocks, data)
		idx++

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}

	return refs, blocks, nil
}

// BlockCount returns the number of blocks needed for a given total size.
func (br *BlockReplicator) BlockCount(totalSize int64) int {
	if totalSize <= 0 {
		return 0
	}
	n := totalSize / br.blockSize
	if totalSize%br.blockSize != 0 {
		n++
	}
	return int(n)
}

// BlockRange returns the start and end block indices that overlap
// with the given byte range [offset, offset+length).
func (br *BlockReplicator) BlockRange(offset, length, totalSize int64) (startBlock, endBlock int) {
	if length <= 0 || totalSize <= 0 {
		return 0, 0
	}
	startBlock = int(offset / br.blockSize)
	end := offset + length - 1
	if end >= totalSize {
		end = totalSize - 1
	}
	endBlock = int(end / br.blockSize)
	return startBlock, endBlock
}

// blockStoragePath returns the on-disk path for a block hash using
// a two-level directory structure for filesystem scalability.
// e.g. "blocks/ab/cd/abcd1234..."
func blockStoragePath(hash string) string {
	if len(hash) < 4 {
		return path.Join("blocks", hash)
	}
	return path.Join("blocks", hash[:2], hash[2:4], hash)
}

// hashBlockData computes the SHA-256 hex digest of a data slice.
func hashBlockData(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// verifyBlockHash checks that data matches the expected hash.
func verifyBlockHash(data []byte, expectedHash string) error {
	actual := hashBlockData(data)
	if actual != expectedHash {
		return fmt.Errorf("block hash mismatch: expected %s, got %s", expectedHash, actual)
	}
	return nil
}

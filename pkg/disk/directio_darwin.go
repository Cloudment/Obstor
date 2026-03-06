/*
 * Minio Cloud Storage, (C) 2019-2020 Minio, Inc.
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

package disk

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

// alignment returns the alignment of the block device.
const alignment = 4096

// OpenFileDirectIO - bypass kernel cache.
// On macOS, O_DIRECT is not available, so open normally then set F_NOCACHE.
func OpenFileDirectIO(filePath string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(filePath, flag, perm)
	if err != nil {
		return f, err
	}
	// F_NOCACHE turns off the file system cache for this file.
	_, err = unix.FcntlInt(f.Fd(), unix.F_NOCACHE, 1)
	if err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

// DisableDirectIO - disables directio mode.
func DisableDirectIO(f *os.File) error {
	fd := f.Fd()
	_, err := unix.FcntlInt(fd, unix.F_NOCACHE, 0)
	return err
}

// AlignedBlock returns a block of the given size aligned to the device block size.
func AlignedBlock(BlockSize int) []byte {
	block := make([]byte, BlockSize+alignment)
	a := alignment - int(uintptr(unsafe.Pointer(&block[0]))&uintptr(alignment-1))
	return block[a : a+BlockSize]
}

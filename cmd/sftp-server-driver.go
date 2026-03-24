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

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/cloudment/obstor/pkg/hash"
	xsftp "github.com/pkg/sftp"
)

// 5GB file max for uploads
const sftpMaxFileSize = 5 << 30

// sftpDriver implements sftp.FileReader, sftp.FileWriter, sftp.FileCmder, sftp.FileLister.
type sftpDriver struct {
	accessKey string
}

func newSFTPDriver(accessKey string) *sftpDriver {
	return &sftpDriver{accessKey: accessKey}
}

// parseSFTPPath splits an SFTP path into bucket and object key.
// "/" -> ("", "")
// "/bucket" -> ("bucket", "")
// "/bucket/key" -> ("bucket", "key")
// "/bucket/dir/file" -> ("bucket", "dir/file")
func parseSFTPPath(p string) (bucket, object string) {
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return "", ""
	}
	parts := strings.SplitN(p, "/", 2)
	bucket = parts[0]
	if len(parts) > 1 {
		object = parts[1]
	}
	return
}

func (d *sftpDriver) getObjectLayer() (ObjectLayer, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return nil, fmt.Errorf("object layer not initialized")
	}
	return objAPI, nil
}

// Fileread implements sftp.FileReader.
func (d *sftpDriver) Fileread(r *xsftp.Request) (io.ReaderAt, error) {
	bucket, object := parseSFTPPath(r.Filepath)
	if bucket == "" || object == "" {
		return nil, os.ErrInvalid
	}

	objAPI, err := d.getObjectLayer()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	oi, err := objAPI.GetObjectInfo(ctx, bucket, object, ObjectOptions{})
	if err != nil {
		return nil, sftpErrorMap(err)
	}

	return &sftpFileReader{
		bucket: bucket,
		object: object,
		size:   oi.Size,
		objAPI: objAPI,
	}, nil
}

// sftpFileReader implements io.ReaderAt by issuing ranged GetObject calls
// instead of loading the entire object into memory.
type sftpFileReader struct {
	bucket string
	object string
	size   int64
	objAPI ObjectLayer
}

func (fr *sftpFileReader) ReadAt(p []byte, off int64) (int, error) {
	if off >= fr.size {
		return 0, io.EOF
	}

	end := off + int64(len(p)) - 1
	if end >= fr.size {
		end = fr.size - 1
	}

	ctx := context.Background()
	reader, err := fr.objAPI.GetObjectNInfo(ctx, fr.bucket, fr.object,
		&HTTPRangeSpec{Start: off, End: end}, nil, readLock, ObjectOptions{})
	if err != nil {
		return 0, sftpErrorMap(err)
	}
	defer reader.Close()

	n, err := io.ReadFull(reader, p[:end-off+1])
	if err == io.ErrUnexpectedEOF {
		err = io.EOF
	}
	return n, err
}

// Filewrite implements sftp.FileWriter.
func (d *sftpDriver) Filewrite(r *xsftp.Request) (io.WriterAt, error) {
	bucket, object := parseSFTPPath(r.Filepath)
	if bucket == "" || object == "" {
		return nil, os.ErrInvalid
	}

	return &sftpFileWriter{
		bucket: bucket,
		object: object,
		driver: d,
		buf:    &bytes.Buffer{},
	}, nil
}

// Filecmd implements sftp.FileCmder.
func (d *sftpDriver) Filecmd(r *xsftp.Request) error {
	ctx := context.Background()

	objAPI, err := d.getObjectLayer()
	if err != nil {
		return err
	}

	switch r.Method {
	case "Setstat":
		return nil // No-op, S3 doesn't support chmod/chown/chtimes.

	case "Rename":
		srcBucket, srcObject := parseSFTPPath(r.Filepath)
		dstBucket, dstObject := parseSFTPPath(r.Target)
		if srcBucket == "" || srcObject == "" || dstBucket == "" || dstObject == "" {
			return os.ErrInvalid
		}
		srcInfo, err := objAPI.GetObjectInfo(ctx, srcBucket, srcObject, ObjectOptions{})
		if err != nil {
			return sftpErrorMap(err)
		}
		if _, err := objAPI.CopyObject(ctx, srcBucket, srcObject, dstBucket, dstObject, srcInfo, ObjectOptions{}, ObjectOptions{}); err != nil {
			return sftpErrorMap(err)
		}
		if _, err := objAPI.DeleteObject(ctx, srcBucket, srcObject, ObjectOptions{}); err != nil {
			return sftpErrorMap(err)
		}
		return nil

	case "Rmdir":
		bucket, prefix := parseSFTPPath(r.Filepath)
		if bucket == "" {
			return os.ErrPermission
		}
		if prefix == "" {
			// Deleting a bucket.
			if err := objAPI.DeleteBucket(ctx, bucket, false); err != nil {
				return sftpErrorMap(err)
			}
			return nil
		}
		// Delete directory marker if it exists.
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		objAPI.DeleteObject(ctx, bucket, prefix, ObjectOptions{})
		return nil

	case "Mkdir":
		bucket, prefix := parseSFTPPath(r.Filepath)
		if bucket == "" {
			return os.ErrInvalid
		}
		if prefix == "" {
			// Creating a bucket.
			return sftpErrorMap(objAPI.MakeBucketWithLocation(ctx, bucket, BucketOptions{}))
		}
		// Create a directory marker object.
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		hashReader, err := hash.NewReader(bytes.NewReader(nil), 0, "", "", 0)
		if err != nil {
			return err
		}
		_, err = objAPI.PutObject(ctx, bucket, prefix, NewPutObjReader(hashReader), ObjectOptions{})
		return sftpErrorMap(err)

	case "Remove":
		bucket, object := parseSFTPPath(r.Filepath)
		if bucket == "" || object == "" {
			return os.ErrInvalid
		}
		_, err := objAPI.DeleteObject(ctx, bucket, object, ObjectOptions{})
		return sftpErrorMap(err)

	case "Symlink":
		return fmt.Errorf("symlinks are not supported")

	case "Link":
		return fmt.Errorf("hard links are not supported")
	}

	return fmt.Errorf("unsupported command: %s", r.Method)
}

// Filelist implements sftp.FileLister.
func (d *sftpDriver) Filelist(r *xsftp.Request) (xsftp.ListerAt, error) {
	ctx := context.Background()

	objAPI, err := d.getObjectLayer()
	if err != nil {
		return nil, err
	}

	bucket, prefix := parseSFTPPath(r.Filepath)

	switch r.Method {
	case "List":
		if bucket == "" {
			// List buckets.
			buckets, err := objAPI.ListBuckets(ctx)
			if err != nil {
				return nil, sftpErrorMap(err)
			}
			entries := make([]os.FileInfo, 0, len(buckets))
			for _, b := range buckets {
				entries = append(entries, &sftpFileInfo{
					name:    b.Name,
					size:    0,
					mode:    os.ModeDir | 0755,
					modTime: b.Created,
					isDir:   true,
				})
			}
			return listerAt(entries), nil
		}

		// List objects in bucket with prefix.
		if prefix != "" && !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		result, err := objAPI.ListObjects(ctx, bucket, prefix, "", "/", 10000)
		if err != nil {
			return nil, sftpErrorMap(err)
		}

		entries := make([]os.FileInfo, 0, len(result.Prefixes)+len(result.Objects))
		// Directories (common prefixes).
		for _, p := range result.Prefixes {
			name := strings.TrimPrefix(p, prefix)
			name = strings.TrimSuffix(name, "/")
			if name == "" {
				continue
			}
			entries = append(entries, &sftpFileInfo{
				name:    name,
				size:    0,
				mode:    os.ModeDir | 0755,
				modTime: time.Time{},
				isDir:   true,
			})
		}
		// Objects.
		for _, obj := range result.Objects {
			name := strings.TrimPrefix(obj.Name, prefix)
			if name == "" || strings.HasSuffix(name, "/") {
				continue // skip directory markers
			}
			entries = append(entries, &sftpFileInfo{
				name:    name,
				size:    obj.Size,
				mode:    0644,
				modTime: obj.ModTime,
				isDir:   false,
			})
		}
		return listerAt(entries), nil

	case "Stat":
		if bucket == "" {
			// Stat root.
			return listerAt([]os.FileInfo{&sftpFileInfo{
				name:    "/",
				size:    0,
				mode:    os.ModeDir | 0755,
				modTime: time.Now(),
				isDir:   true,
			}}), nil
		}
		if prefix == "" {
			// Stat a bucket.
			bi, err := objAPI.GetBucketInfo(ctx, bucket)
			if err != nil {
				return nil, sftpErrorMap(err)
			}
			return listerAt([]os.FileInfo{&sftpFileInfo{
				name:    bi.Name,
				size:    0,
				mode:    os.ModeDir | 0755,
				modTime: bi.Created,
				isDir:   true,
			}}), nil
		}
		// Stat an object.
		oi, err := objAPI.GetObjectInfo(ctx, bucket, prefix, ObjectOptions{})
		if err != nil {
			// Maybe it's a directory prefix.
			dirPrefix := prefix
			if !strings.HasSuffix(dirPrefix, "/") {
				dirPrefix += "/"
			}
			result, listErr := objAPI.ListObjects(ctx, bucket, dirPrefix, "", "/", 1)
			if listErr != nil {
				return nil, sftpErrorMap(err)
			}
			if len(result.Objects) > 0 || len(result.Prefixes) > 0 {
				name := strings.TrimSuffix(prefix, "/")
				if idx := strings.LastIndex(name, "/"); idx >= 0 {
					name = name[idx+1:]
				}
				return listerAt([]os.FileInfo{&sftpFileInfo{
					name:    name,
					size:    0,
					mode:    os.ModeDir | 0755,
					modTime: time.Time{},
					isDir:   true,
				}}), nil
			}
			return nil, os.ErrNotExist
		}
		name := oi.Name
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		return listerAt([]os.FileInfo{&sftpFileInfo{
			name:    name,
			size:    oi.Size,
			mode:    0644,
			modTime: oi.ModTime,
			isDir:   false,
		}}), nil

	case "Readlink":
		return nil, fmt.Errorf("readlink is not supported")
	}

	return nil, fmt.Errorf("unsupported list method: %s", r.Method)
}

// sftpFileWriter buffers writes and uploads on Close.
type sftpFileWriter struct {
	bucket string
	object string
	driver *sftpDriver
	buf    *bytes.Buffer
	offset int64
}

func (w *sftpFileWriter) WriteAt(p []byte, off int64) (n int, err error) {
	end := off + int64(len(p))
	if end > sftpMaxFileSize {
		return 0, fmt.Errorf("file size exceeds maximum allowed (%d bytes)", sftpMaxFileSize)
	}
	// Ensure buffer is large enough.
	if end > int64(w.buf.Len()) {
		w.buf.Grow(int(end) - w.buf.Len())
		b := w.buf.Bytes()[:end]
		// Zero-fill any gap between old length and write offset.
		w.buf.Reset()
		w.buf.Write(b)
	}
	copy(w.buf.Bytes()[off:], p)
	if end > w.offset {
		w.offset = end
	}
	return len(p), nil
}

func (w *sftpFileWriter) Close() error {
	objAPI, err := w.driver.getObjectLayer()
	if err != nil {
		return err
	}

	data := w.buf.Bytes()[:w.offset]
	reader := bytes.NewReader(data)
	hashReader, err := hash.NewReader(reader, int64(len(data)), "", "", int64(len(data)))
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = objAPI.PutObject(ctx, w.bucket, w.object, NewPutObjReader(hashReader), ObjectOptions{})
	return sftpErrorMap(err)
}

// listerAt implements sftp.ListerAt.
type listerAt []os.FileInfo

func (l listerAt) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l)) {
		return 0, io.EOF
	}
	n := copy(ls, l[offset:])
	if n+int(offset) >= len(l) {
		return n, io.EOF
	}
	return n, nil
}

// sftpFileInfo implements os.FileInfo.
type sftpFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (fi *sftpFileInfo) Name() string       { return fi.name }
func (fi *sftpFileInfo) Size() int64        { return fi.size }
func (fi *sftpFileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *sftpFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *sftpFileInfo) IsDir() bool        { return fi.isDir }
func (fi *sftpFileInfo) Sys() interface{}   { return nil }

// sftpErrorMap maps object layer errors to OS-level errors that SFTP clients understand.
func sftpErrorMap(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case isErrObjectNotFound(err), isErrVersionNotFound(err):
		return os.ErrNotExist
	case isErrBucketNotFound(err):
		return os.ErrNotExist
	}
	switch err.(type) {
	case BucketAlreadyExists, BucketAlreadyOwnedByYou:
		return os.ErrExist
	case BucketNotEmpty:
		return fmt.Errorf("directory not empty")
	}
	return err
}

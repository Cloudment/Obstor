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
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	obstor "github.com/cloudment/obstor/cmd"
	"github.com/cloudment/obstor/pkg/hash"
	xsftp "github.com/pkg/sftp"
)

// 5GB file max for uploads
const sftpMaxFileSize = 5 << 30

// sftp.FileReader, sftp.FileWriter, sftp.FileCmder, sftp.FileLister.
type sftpDriver struct {
	accessKey string
}

func newSFTPDriver(accessKey string) *sftpDriver {
	return &sftpDriver{accessKey: accessKey}
}

// Split an SFTP path into bucket and object key.
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

func (d *sftpDriver) getObjectLayer() (obstor.ObjectLayer, error) {
	objAPI := obstor.NewObjectLayerFn()
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
	if err := obstor.CheckSFTPAccess(d.accessKey, bucket, object, obstor.SFTPActionGetObject); err != nil {
		return nil, os.ErrPermission
	}

	objAPI, err := d.getObjectLayer()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	oi, err := objAPI.GetObjectInfo(ctx, bucket, object, obstor.ObjectOptions{})
	if err != nil {
		return nil, sftpErrorMap(err)
	}
	if oi.Size > sftpMaxFileSize {
		return nil, fmt.Errorf("object size %d exceeds SFTP download limit (%d bytes)", oi.Size, sftpMaxFileSize)
	}

	reader, err := objAPI.GetObjectNInfo(ctx, bucket, object, nil, nil, obstor.ReadLock, obstor.ObjectOptions{})
	if err != nil {
		return nil, sftpErrorMap(err)
	}

	data, err := io.ReadAll(reader)
	_ = reader.Close()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

// Filewrite implements sftp.FileWriter.
func (d *sftpDriver) Filewrite(r *xsftp.Request) (io.WriterAt, error) {
	bucket, object := parseSFTPPath(r.Filepath)
	if bucket == "" || object == "" {
		return nil, os.ErrInvalid
	}
	if err := obstor.CheckSFTPAccess(d.accessKey, bucket, object, obstor.SFTPActionPutObject); err != nil {
		return nil, os.ErrPermission
	}

	return &sftpFileWriter{
		bucket: bucket,
		object: object,
		driver: d,
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
		if err := obstor.CheckSFTPAccess(d.accessKey, srcBucket, srcObject, obstor.SFTPActionGetObject); err != nil {
			return os.ErrPermission
		}
		if err := obstor.CheckSFTPAccess(d.accessKey, dstBucket, dstObject, obstor.SFTPActionPutObject); err != nil {
			return os.ErrPermission
		}
		if err := obstor.CheckSFTPAccess(d.accessKey, srcBucket, srcObject, obstor.SFTPActionDeleteObject); err != nil {
			return os.ErrPermission
		}
		// CopyObject requires PutObjReader which GetObjectInfo doesnt have
		reader, err := objAPI.GetObjectNInfo(ctx, srcBucket, srcObject, nil, nil, obstor.ReadLock, obstor.ObjectOptions{})
		if err != nil {
			return sftpErrorMap(err)
		}
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			return err
		}
		hashReader, err := hash.NewReader(bytes.NewReader(data), int64(len(data)), "", "", int64(len(data)))
		if err != nil {
			return err
		}
		if _, err := objAPI.PutObject(ctx, dstBucket, dstObject, obstor.NewPutObjReader(hashReader), obstor.ObjectOptions{}); err != nil {
			return sftpErrorMap(err)
		}
		if _, err := objAPI.DeleteObject(ctx, srcBucket, srcObject, obstor.ObjectOptions{}); err != nil {
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
			if err := obstor.CheckSFTPAccess(d.accessKey, bucket, "", obstor.SFTPActionDeleteBucket); err != nil {
				return os.ErrPermission
			}
			if err := objAPI.DeleteBucket(ctx, bucket, false); err != nil {
				return sftpErrorMap(err)
			}
			return nil
		}
		if err := obstor.CheckSFTPAccess(d.accessKey, bucket, prefix, obstor.SFTPActionDeleteObject); err != nil {
			return os.ErrPermission
		}
		// Delete directory marker if it exists.
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		_, _ = objAPI.DeleteObject(ctx, bucket, prefix, obstor.ObjectOptions{})
		return nil

	case "Mkdir":
		bucket, prefix := parseSFTPPath(r.Filepath)
		if bucket == "" {
			return os.ErrInvalid
		}
		if prefix == "" {
			// Creating a bucket.
			if err := obstor.CheckSFTPAccess(d.accessKey, bucket, "", obstor.SFTPActionCreateBucket); err != nil {
				return os.ErrPermission
			}
			return sftpErrorMap(objAPI.MakeBucketWithLocation(ctx, bucket, obstor.BucketOptions{}))
		}
		if err := obstor.CheckSFTPAccess(d.accessKey, bucket, prefix, obstor.SFTPActionPutObject); err != nil {
			return os.ErrPermission
		}
		// Create a directory marker object.
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		hashReader, err := hash.NewReader(bytes.NewReader(nil), 0, "", "", 0)
		if err != nil {
			return err
		}
		_, err = objAPI.PutObject(ctx, bucket, prefix, obstor.NewPutObjReader(hashReader), obstor.ObjectOptions{})
		return sftpErrorMap(err)

	case "Remove":
		bucket, object := parseSFTPPath(r.Filepath)
		if bucket == "" || object == "" {
			return os.ErrInvalid
		}
		if err := obstor.CheckSFTPAccess(d.accessKey, bucket, object, obstor.SFTPActionDeleteObject); err != nil {
			return os.ErrPermission
		}
		_, err := objAPI.DeleteObject(ctx, bucket, object, obstor.ObjectOptions{})
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
				if err := obstor.CheckSFTPAccess(d.accessKey, b.Name, "", obstor.SFTPActionListBucket); err != nil {
					continue
				}
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
		if err := obstor.CheckSFTPAccess(d.accessKey, bucket, "", obstor.SFTPActionListBucket); err != nil {
			return nil, os.ErrPermission
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
		if err := obstor.CheckSFTPAccess(d.accessKey, bucket, prefix, obstor.SFTPActionListBucket); err != nil {
			return nil, os.ErrPermission
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
		oi, err := objAPI.GetObjectInfo(ctx, bucket, prefix, obstor.ObjectOptions{})
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

// Buffer writes and uploads on Close.
type sftpFileWriter struct {
	mu     sync.Mutex
	bucket string
	object string
	driver *sftpDriver
	buf    []byte
}

func (w *sftpFileWriter) WriteAt(p []byte, off int64) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	end := off + int64(len(p))
	if end > sftpMaxFileSize {
		return 0, fmt.Errorf("file size exceeds maximum allowed (%d bytes)", sftpMaxFileSize)
	}
	// Grow buffer if needed
	if end > int64(len(w.buf)) {
		grown := make([]byte, end)
		copy(grown, w.buf)
		w.buf = grown
	}
	copy(w.buf[off:], p)
	return len(p), nil
}

func (w *sftpFileWriter) Close() error {
	objAPI, err := w.driver.getObjectLayer()
	if err != nil {
		return err
	}

	reader := bytes.NewReader(w.buf)
	hashReader, err := hash.NewReader(reader, int64(len(w.buf)), "", "", int64(len(w.buf)))
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = objAPI.PutObject(ctx, w.bucket, w.object, obstor.NewPutObjReader(hashReader), obstor.ObjectOptions{})
	return sftpErrorMap(err)
}

// Implement sftp.ListerAt
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

// Implement os.FileInfo
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

// Map object layer errors to OS-level errors that SFTP clients understand.
func sftpErrorMap(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case obstor.IsErrObjectNotFound(err), obstor.IsErrVersionNotFound(err):
		return os.ErrNotExist
	case obstor.IsErrBucketNotFound(err):
		return os.ErrNotExist
	}
	switch err.(type) {
	case obstor.BucketAlreadyExists, obstor.BucketAlreadyOwnedByYou:
		return os.ErrExist
	case obstor.BucketNotEmpty:
		return fmt.Errorf("directory not empty")
	}
	return err
}

package atime

import (
	"os"
	"syscall"
	"time"
)

// Get returns the access time of the file described by the given FileInfo.
func Get(fi os.FileInfo) time.Time {
	d := fi.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, d.LastAccessTime.Nanoseconds())
}

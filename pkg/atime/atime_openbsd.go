package atime

import (
	"os"
	"syscall"
	"time"
)

// Get returns the access time of the file described by the given FileInfo.
func Get(fi os.FileInfo) time.Time {
	st := fi.Sys().(*syscall.Stat_t)
	return time.Unix(st.Atimespec.Sec, st.Atimespec.Nsec)
}

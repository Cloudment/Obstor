//go:build darwin || freebsd || netbsd || openbsd

package cmd

import "syscall"

func globalSync() {
	_ = syscall.Sync()
}

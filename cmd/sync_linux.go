//go:build linux && !appengine

package cmd

import "syscall"

func globalSync() {
	syscall.Sync()
}

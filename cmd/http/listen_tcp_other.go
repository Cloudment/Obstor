// +build !linux,!windows

package http

func setTCPOptions(fd int) {
	// TCP_FASTOPEN and TCP_DEFER_ACCEPT are Linux-specific.
	// On other Unix platforms, these are not available.
}

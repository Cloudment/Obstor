package http

import "golang.org/x/sys/unix"

func setTCPOptions(fd int) {
	_ = unix.SetsockoptInt(fd, unix.SOL_TCP, unix.TCP_FASTOPEN, 1)
	_ = unix.SetsockoptInt(fd, unix.SOL_TCP, unix.TCP_DEFER_ACCEPT, 1)
}

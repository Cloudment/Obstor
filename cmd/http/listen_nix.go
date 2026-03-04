// +build linux darwin dragonfly freebsd netbsd openbsd rumprun

/*
 * MinIO Cloud Storage, (C) 2018 MinIO, Inc.
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

package http

import (
	"context"
	"net"
	"syscall"
)

var listenCfg = net.ListenConfig{
	Control: func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			// Enable SO_REUSEADDR
			_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			// TCP_FASTOPEN and TCP_DEFER_ACCEPT are set via platform-specific files.
			setTCPOptions(int(fd))
		})
	},
}

// Unix listener with special TCP options.
var listen = func(network, addr string) (net.Listener, error) {
	return listenCfg.Listen(context.Background(), network, addr)
}
var fallbackListen = net.Listen

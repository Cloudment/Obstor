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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/cloudment/obstor/cmd/config/sftp"
	"github.com/cloudment/obstor/cmd/logger"
	xsftp "github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var globalSFTPConfig sftp.Config

func startSFTPServer(cfg sftp.Config) {
	globalSFTPConfig = cfg

	sshConfig := &ssh.ServerConfig{
		PasswordCallback: sftpPasswordCallback,
	}

	hostKey, err := loadOrGenerateHostKey(cfg.HostKeyPath)
	if err != nil {
		logger.Fatal(err, "Unable to load/generate SFTP host key")
	}
	sshConfig.AddHostKey(hostKey)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Fatal(err, "Unable to start SFTP server on %s", cfg.Address)
	}

	logger.Info("SFTP server listening on %s", cfg.Address)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				logger.LogIf(GlobalContext, fmt.Errorf("sftp accept error: %w", err))
				continue
			}
			go handleSFTPConnection(conn, sshConfig)
		}
	}()
}

func sftpPasswordCallback(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	accessKey := c.User()
	secretKey := string(pass)

	// Check root credentials.
	if accessKey == globalActiveCred.AccessKey && globalActiveCred.SecretKey == secretKey {
		return &ssh.Permissions{
			Extensions: map[string]string{
				"accessKey": accessKey,
			},
		}, nil
	}

	// Check IAM users.
	if globalIAMSys != nil {
		cred, ok := globalIAMSys.GetUser(accessKey)
		if ok && cred.SecretKey == secretKey && cred.IsValid() {
			return &ssh.Permissions{
				Extensions: map[string]string{
					"accessKey": accessKey,
				},
			}, nil
		}
	}

	return nil, fmt.Errorf("authentication failed for user %s", accessKey)
}

func handleSFTPConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		logger.LogIf(GlobalContext, fmt.Errorf("sftp ssh handshake error from %s: %w", conn.RemoteAddr(), err))
		return
	}
	defer sshConn.Close()

	logger.Info("SFTP connection from %s (user: %s)", conn.RemoteAddr(), sshConn.User())

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			logger.LogIf(GlobalContext, fmt.Errorf("sftp channel accept error: %w", err))
			return
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := req.Type == "subsystem" && len(req.Payload) > 4 && string(req.Payload[4:]) == "sftp"
				req.Reply(ok, nil)
			}
		}(requests)

		driver := newSFTPDriver(sshConn.User())
		server := xsftp.NewRequestServer(channel, xsftp.Handlers{
			FileGet:  driver,
			FilePut:  driver,
			FileCmd:  driver,
			FileList: driver,
		})

		if err := server.Serve(); err != nil {
			if err != io.EOF {
				logger.LogIf(GlobalContext, fmt.Errorf("sftp server error for user %s: %w", sshConn.User(), err))
			}
		}
		server.Close()
	}
}

func loadOrGenerateHostKey(hostKeyPath string) (ssh.Signer, error) {
	if hostKeyPath != "" {
		keyBytes, err := os.ReadFile(hostKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read host key %s: %w", hostKeyPath, err)
		}
		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse host key %s: %w", hostKeyPath, err)
		}
		return signer, nil
	}

	// Auto-generate and persist a host key in the certs directory.
	autoKeyPath := filepath.Join(globalCertsDir.Get(), "sftp_host_key")
	if keyBytes, err := os.ReadFile(autoKeyPath); err == nil {
		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err == nil {
			return signer, nil
		}
	}

	logger.Info("Generating SFTP host key at %s", autoKeyPath)
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate host key: %w", err)
	}

	keyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	if err := os.WriteFile(autoKeyPath, keyBytes, 0600); err != nil {
		return nil, fmt.Errorf("failed to save host key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated host key: %w", err)
	}
	return signer, nil
}

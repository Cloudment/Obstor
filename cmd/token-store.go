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
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

const (
	// Maximum number of active tokens before rejecting new ones.
	maxPresignedTokens = 100000
	// Cleanup interval for expired tokens.
	tokenCleanupInterval = 5 * time.Minute
)

var errTokenStoreFull = errors.New("presigned token store is full")

// PresignedTokenStore tracks issued presigned URL tokens to enforce one-time use.
type PresignedTokenStore struct {
	mu     sync.Mutex
	tokens map[string]time.Time
	done   chan struct{}
}

var globalTokenStore = newTokenStore()

func newTokenStore() *PresignedTokenStore {
	s := &PresignedTokenStore{
		tokens: make(map[string]time.Time),
		done:   make(chan struct{}),
	}
	go s.cleanupLoop()
	return s
}

func (s *PresignedTokenStore) cleanupLoop() {
	ticker := time.NewTicker(tokenCleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.cleanup()
		}
	}
}

// Stop and cleanup goroutine.
func (s *PresignedTokenStore) Stop() {
	close(s.done)
}

// Issue creates a new one-time token with the given TTL.
func (s *PresignedTokenStore) Issue(ttl time.Duration) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.tokens) >= maxPresignedTokens {
		return "", errTokenStoreFull
	}

	s.tokens[token] = time.Now().Add(ttl)
	return token, nil
}

// Consume validates and removes a token.
// Returns true if the token was valid and unused.
func (s *PresignedTokenStore) Consume(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Timing-safe lookup
	var found bool
	var foundKey string
	for k, expiry := range s.tokens {
		if subtle.ConstantTimeCompare([]byte(k), []byte(token)) == 1 {
			if time.Now().Before(expiry) {
				found = true
				foundKey = k
			}
			break
		}
	}

	if found {
		delete(s.tokens, foundKey)
	}
	return found
}

func (s *PresignedTokenStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, expiry := range s.tokens {
		if now.After(expiry) {
			delete(s.tokens, token)
		}
	}
}

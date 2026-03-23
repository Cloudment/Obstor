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

package replication

// ConsistencyMode determines read/write quorum behavior.
type ConsistencyMode string

const (
	ConsistencyModeConsistent ConsistencyMode = "consistent"
	ConsistencyModeDegraded   ConsistencyMode = "degraded"
	ConsistencyModeDangerous  ConsistencyMode = "dangerous"
)

// Config represents the zone-aware block replication configuration.
type Config struct {
	ReplicationFactor int             `json:"replication_factor"`
	Consistency       ConsistencyMode `json:"consistency_mode"`
	BlockSize         int64           `json:"block_size"`
	Zone              string          `json:"zone"`
	NodeCapacity      int64           `json:"node_capacity"`
	SyncQueueSize     int             `json:"sync_queue_size"`
	SyncWorkers       int             `json:"sync_workers"`
	MaxSyncRetries    int             `json:"max_sync_retries"`
}

// WriteQuorum returns the minimum copies that must be written
// before acknowledging success.
func (c Config) WriteQuorum() int {
	switch c.Consistency {
	case ConsistencyModeDangerous:
		return 1
	default:
		if c.ReplicationFactor <= 1 {
			return 1
		}
		if c.ReplicationFactor == 2 {
			return 2
		}
		return (c.ReplicationFactor / 2) + 1
	}
}

// ReadQuorum returns the minimum copies that must respond
// to satisfy a read request.
func (c Config) ReadQuorum() int {
	switch c.Consistency {
	case ConsistencyModeDangerous, ConsistencyModeDegraded:
		return 1
	default:
		if c.ReplicationFactor <= 2 {
			return 1
		}
		return (c.ReplicationFactor / 2) + 1
	}
}

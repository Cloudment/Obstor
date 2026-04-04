/*
 * MinIO Cloud Storage, (C) 2019-2024 MinIO, Inc.
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

// Package jstream provides a compatibility layer that replaces
// github.com/bcicen/jstream using only the Go standard library's
// encoding/json package.
//
// It reimplements the subset of jstream used by Obstor:
//   - KV, KVS types (ordered key-value pairs with custom JSON marshaling)
//   - Decoder with depth-based filtering and ObjectAsKVS mode
//   - MetaValue with ValueType discrimination
package jstream

import (
	"encoding/json"
	"io"
	"sync"
)

// ValueType indicates the JSON type of a decoded value.
type ValueType int

const (
	// Unknown is the zero value.
	Unknown ValueType = iota
	// Null is a JSON null.
	Null
	// Bool is a JSON boolean.
	Bool
	// Number is a JSON number (float64 or json.Number).
	Number
	// String is a JSON string.
	String
	// Array is a JSON array.
	Array
	// Object is a JSON object (decoded as KVS when ObjectAsKVS is set).
	Object
)

// KV is an ordered key-value pair from a JSON object.
type KV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// KVS is an ordered slice of key-value pairs that represents a JSON
// object while preserving key order. It implements json.Marshaler so
// that encoding/json.Encoder produces a JSON object (not an array of
// KV structs).
type KVS []KV

// MarshalJSON encodes KVS as a JSON object with keys in order.
func (kvs KVS) MarshalJSON() ([]byte, error) {
	// Estimate capacity: 2 for braces + per-entry overhead.
	buf := make([]byte, 0, 128)
	buf = append(buf, '{')
	for i, kv := range kvs {
		if i > 0 {
			buf = append(buf, ',')
		}
		// Marshal key
		keyBytes, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		buf = append(buf, keyBytes...)
		buf = append(buf, ':')
		// Marshal value
		valBytes, err := json.Marshal(kv.Value)
		if err != nil {
			return nil, err
		}
		buf = append(buf, valBytes...)
	}
	buf = append(buf, '}')
	return buf, nil
}

// MetaValue holds a decoded JSON value along with metadata.
type MetaValue struct {
	// Value is the decoded JSON value. For objects decoded with
	// ObjectAsKVS it will be of type KVS. Arrays become
	// []interface{}, strings become string, numbers become float64,
	// bools become bool, and null becomes nil.
	Value interface{}

	// ValueType indicates the JSON type of Value.
	ValueType ValueType

	// Offset is the approximate byte offset in the stream where
	// this value started. It is a best-effort approximation since
	// encoding/json.Decoder does not expose exact offsets.
	Offset int

	// Length is kept for API compatibility but is always 0 in this
	// implementation.
	Length int
}

// Decoder reads a JSON stream and emits values at a specified depth.
type Decoder struct {
	reader      io.Reader
	depth       int
	objectAsKVS bool
	err         error
	mu          sync.Mutex
	ch          chan *MetaValue
	done        chan struct{}
}

// NewDecoder creates a Decoder that emits values found at the given
// nesting depth. Depth 0 means top-level values (the same semantics
// as bcicen/jstream).
func NewDecoder(r io.Reader, depth int) *Decoder {
	return &Decoder{
		reader: r,
		depth:  depth,
	}
}

// ObjectAsKVS configures the decoder to represent JSON objects as KVS
// (ordered key-value slices) instead of map[string]interface{}.
// Returns the decoder for chaining.
func (d *Decoder) ObjectAsKVS() *Decoder {
	d.objectAsKVS = true
	return d
}

// Stream starts decoding in a background goroutine and returns a
// channel of *MetaValue. The channel is closed when the stream ends
// (either EOF or error). Use Err() after the channel is drained to
// check for errors.
func (d *Decoder) Stream() chan *MetaValue {
	d.ch = make(chan *MetaValue, 128)
	d.done = make(chan struct{})
	go d.run()
	return d.ch
}

// Err returns the first non-EOF error encountered during decoding.
func (d *Decoder) Err() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.err
}

func (d *Decoder) setErr(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.err == nil {
		d.err = err
	}
}

// run drives the decoding loop.
func (d *Decoder) run() {
	defer close(d.ch)
	defer close(d.done)

	dec := json.NewDecoder(d.reader)
	dec.UseNumber()

	// We walk the token stream, tracking nesting depth and emitting
	// complete values when we reach the target depth.
	d.decode(dec, 0)
}

// decode recursively walks the JSON token stream. currentDepth is
// the depth of the *container* we are currently inside (0 = top level).
// When currentDepth == d.depth, each value encountered is emitted.
func (d *Decoder) decode(dec *json.Decoder, currentDepth int) {
	for dec.More() {
		if currentDepth == d.depth {
			// Read and emit a complete value at target depth.
			val, vt, err := d.readValue(dec)
			if err != nil {
				if err != io.EOF {
					d.setErr(err)
				}
				return
			}
			mv := &MetaValue{
				Value:     val,
				ValueType: vt,
				Offset:    int(dec.InputOffset()),
			}
			d.ch <- mv
		} else {
			// We need to descend deeper. Read the next token.
			tok, err := dec.Token()
			if err != nil {
				if err != io.EOF {
					d.setErr(err)
				}
				return
			}
			switch tok := tok.(type) {
			case json.Delim:
				delim := tok
				switch delim {
				case '{':
					d.decodeObject(dec, currentDepth+1)
				case '[':
					d.decodeArray(dec, currentDepth+1)
				case '}', ']':
					// End of current container; return to parent.
					return
				}
			default:
				// Scalar at a depth less than target; skip it.
			}
		}
	}
}

// decodeObject handles tokens inside a JSON object (after the opening '{').
// The closing '}' is consumed before returning.
func (d *Decoder) decodeObject(dec *json.Decoder, currentDepth int) {
	for dec.More() {
		// Read the key token.
		tok, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				d.setErr(err)
			}
			return
		}
		// The key should be a string.
		_, ok := tok.(string)
		if !ok {
			// Could be closing delim if More() was wrong; handle gracefully.
			return
		}

		if currentDepth == d.depth {
			// Emit the value.
			val, vt, err := d.readValue(dec)
			if err != nil {
				if err != io.EOF {
					d.setErr(err)
				}
				return
			}
			mv := &MetaValue{
				Value:     val,
				ValueType: vt,
				Offset:    int(dec.InputOffset()),
			}
			d.ch <- mv
		} else if currentDepth < d.depth {
			// Need to go deeper.
			tok2, err := dec.Token()
			if err != nil {
				if err != io.EOF {
					d.setErr(err)
				}
				return
			}
			switch tok2 := tok2.(type) {
			case json.Delim:
				delim := tok2
				switch delim {
				case '{':
					d.decodeObject(dec, currentDepth+1)
				case '[':
					d.decodeArray(dec, currentDepth+1)
				}
			default:
				// Scalar value at wrong depth; skip.
			}
		} else {
			// currentDepth > d.depth; should not happen if called correctly.
			// Skip the value.
			d.skipValue(dec)
		}
	}
	// Consume closing '}'.
	dec.Token() // nolint: errcheck
}

// decodeArray handles tokens inside a JSON array (after the opening '[').
// The closing ']' is consumed before returning.
func (d *Decoder) decodeArray(dec *json.Decoder, currentDepth int) {
	for dec.More() {
		if currentDepth == d.depth {
			val, vt, err := d.readValue(dec)
			if err != nil {
				if err != io.EOF {
					d.setErr(err)
				}
				return
			}
			mv := &MetaValue{
				Value:     val,
				ValueType: vt,
				Offset:    int(dec.InputOffset()),
			}
			d.ch <- mv
		} else if currentDepth < d.depth {
			tok, err := dec.Token()
			if err != nil {
				if err != io.EOF {
					d.setErr(err)
				}
				return
			}
			switch tok := tok.(type) {
			case json.Delim:
				delim := tok
				switch delim {
				case '{':
					d.decodeObject(dec, currentDepth+1)
				case '[':
					d.decodeArray(dec, currentDepth+1)
				}
			default:
				// Scalar at wrong depth; skip.
			}
		} else {
			d.skipValue(dec)
		}
	}
	// Consume closing ']'.
	dec.Token() // nolint: errcheck
}

// readValue reads a complete JSON value from the decoder. If
// objectAsKVS is enabled, objects are decoded as KVS.
func (d *Decoder) readValue(dec *json.Decoder) (interface{}, ValueType, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, Unknown, err
	}

	switch v := tok.(type) {
	case json.Delim:
		switch v {
		case '{':
			return d.readObjectValue(dec)
		case '[':
			return d.readArrayValue(dec)
		default:
			// Unexpected closing delimiter.
			return nil, Unknown, nil
		}
	case json.Number:
		// Convert json.Number to float64 to match jstream behavior.
		f, err := v.Float64()
		if err != nil {
			return nil, Unknown, err
		}
		return f, Number, nil
	case string:
		return v, String, nil
	case bool:
		return v, Bool, nil
	case nil:
		return nil, Null, nil
	default:
		return v, Unknown, nil
	}
}

// readObjectValue reads object contents after the opening '{' and
// returns the value. If objectAsKVS is true, returns KVS; otherwise
// returns map[string]interface{}.
func (d *Decoder) readObjectValue(dec *json.Decoder) (interface{}, ValueType, error) {
	if d.objectAsKVS {
		kvs := make(KVS, 0, 8)
		for dec.More() {
			tok, err := dec.Token()
			if err != nil {
				return nil, Unknown, err
			}
			key, ok := tok.(string)
			if !ok {
				return nil, Unknown, nil
			}
			val, _, err := d.readValue(dec)
			if err != nil {
				return nil, Unknown, err
			}
			kvs = append(kvs, KV{Key: key, Value: val})
		}
		// Consume closing '}'.
		dec.Token() // nolint: errcheck
		return kvs, Object, nil
	}

	m := make(map[string]interface{})
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return nil, Unknown, err
		}
		key, ok := tok.(string)
		if !ok {
			return nil, Unknown, nil
		}
		val, _, err := d.readValue(dec)
		if err != nil {
			return nil, Unknown, err
		}
		m[key] = val
	}
	// Consume closing '}'.
	dec.Token() // nolint: errcheck
	return m, Object, nil
}

// readArrayValue reads array contents after the opening '['.
func (d *Decoder) readArrayValue(dec *json.Decoder) (interface{}, ValueType, error) {
	var arr []interface{}
	for dec.More() {
		val, _, err := d.readValue(dec)
		if err != nil {
			return nil, Unknown, err
		}
		arr = append(arr, val)
	}
	// Consume closing ']'.
	dec.Token() // nolint: errcheck
	if arr == nil {
		arr = []interface{}{}
	}
	return arr, Array, nil
}

// skipValue reads and discards a complete JSON value from the decoder.
func (d *Decoder) skipValue(dec *json.Decoder) {
	tok, err := dec.Token()
	if err != nil {
		return
	}
	switch v := tok.(type) {
	case json.Delim:
		switch v {
		case '{':
			d.skipObject(dec)
		case '[':
			d.skipArray(dec)
		}
	}
	// Scalars are already consumed.
}

func (d *Decoder) skipObject(dec *json.Decoder) {
	for dec.More() {
		// Read key.
		_, err := dec.Token()
		if err != nil {
			return
		}
		// Skip value.
		d.skipValue(dec)
	}
	_, _ = dec.Token() // consume '}'
}

func (d *Decoder) skipArray(dec *json.Decoder) {
	for dec.More() {
		d.skipValue(dec)
	}
	_, _ = dec.Token() // consume ']'
}

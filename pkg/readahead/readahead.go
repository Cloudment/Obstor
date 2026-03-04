// Package readahead provides a goroutine-based readahead reader that
// pre-fills buffers from an underlying reader in a background goroutine
// to reduce read latency for sequential workloads.
package readahead

import (
	"errors"
	"io"
	"sync"
)

const (
	defaultBuffers = 4
	defaultBufSize = 1 << 20 // 1 MiB
)

// buffer holds a single pre-read chunk and any error encountered while filling it.
type buffer struct {
	data []byte
	n    int
	err  error
}

// reader implements io.ReadCloser with goroutine-based readahead.
type reader struct {
	closer io.Closer    // non-nil when wrapping an io.ReadCloser
	bufs   chan buffer   // filled buffers sent from the background goroutine
	cur    buffer        // current buffer being consumed by Read
	off    int           // read offset within cur.data[:cur.n]
	done   chan struct{}  // closed to signal the background goroutine to stop
	once   sync.Once     // ensures Close is idempotent
	err    error         // sticky error returned after all buffered data is drained
}

// NewReader creates a readahead reader with default settings (4 buffers, 1 MiB each).
func NewReader(r io.Reader) (io.ReadCloser, error) {
	return NewReaderSize(r, defaultBuffers, defaultBufSize)
}

// NewReaderSize creates a readahead reader with the specified number of buffers
// and buffer size. The background goroutine pre-reads up to `buffers` chunks of
// `bufSize` bytes each, keeping them ready for consumption by Read.
func NewReaderSize(r io.Reader, buffers, bufSize int) (io.ReadCloser, error) {
	if r == nil {
		return nil, errors.New("readahead: reader is nil")
	}
	if buffers < 1 {
		buffers = defaultBuffers
	}
	if bufSize < 1 {
		bufSize = defaultBufSize
	}

	ra := &reader{
		bufs: make(chan buffer, buffers),
		done: make(chan struct{}),
	}

	go ra.fill(r, buffers, bufSize)
	return ra, nil
}

// NewReadCloser creates a readahead ReadCloser with default settings.
// When the returned ReadCloser is closed it also closes the underlying rc.
func NewReadCloser(rc io.ReadCloser) (io.ReadCloser, error) {
	return NewReadCloserSize(rc, defaultBuffers, defaultBufSize)
}

// NewReadCloserSize creates a readahead ReadCloser with custom buffer settings.
// When the returned ReadCloser is closed it also closes the underlying rc.
func NewReadCloserSize(rc io.ReadCloser, buffers, bufSize int) (io.ReadCloser, error) {
	if rc == nil {
		return nil, errors.New("readahead: reader is nil")
	}
	if buffers < 1 {
		buffers = defaultBuffers
	}
	if bufSize < 1 {
		bufSize = defaultBufSize
	}

	ra := &reader{
		closer: rc,
		bufs:   make(chan buffer, buffers),
		done:   make(chan struct{}),
	}

	go ra.fill(rc, buffers, bufSize)
	return ra, nil
}

// fill runs in a background goroutine. It reads from r into pre-allocated
// buffers and sends them on the bufs channel. It stops when r returns an
// error (including io.EOF) or when done is closed.
func (ra *reader) fill(r io.Reader, buffers, bufSize int) {
	defer close(ra.bufs)

	// Pre-allocate a pool of byte slices that rotate through the channel.
	pool := make([][]byte, buffers)
	for i := range pool {
		pool[i] = make([]byte, bufSize)
	}

	idx := 0
	for {
		buf := pool[idx%buffers]
		n, err := io.ReadFull(r, buf)

		// Send whatever we got, even on error.
		b := buffer{data: buf, n: n, err: err}

		select {
		case ra.bufs <- b:
		case <-ra.done:
			return
		}

		// io.ErrUnexpectedEOF from ReadFull means the reader returned
		// fewer bytes than bufSize but did provide some data; the real
		// underlying error is io.EOF (short final chunk).
		if err != nil {
			return
		}
		idx++
	}
}

// Read satisfies io.Reader. It consumes pre-filled buffers delivered by the
// background goroutine.
func (ra *reader) Read(p []byte) (int, error) {
	for {
		// Serve from the current buffer first.
		if ra.off < ra.cur.n {
			n := copy(p, ra.cur.data[ra.off:ra.cur.n])
			ra.off += n
			return n, nil
		}

		// Current buffer exhausted. If we already saw a terminal error, return it.
		if ra.err != nil {
			return 0, ra.err
		}

		// Fetch the next buffer from the background goroutine.
		b, ok := <-ra.bufs
		if !ok {
			// Channel closed with no more buffers and no stored error means EOF.
			ra.err = io.EOF
			return 0, io.EOF
		}
		ra.cur = b
		ra.off = 0

		// Translate errors: io.ErrUnexpectedEOF from ReadFull means the
		// underlying reader hit EOF mid-buffer, which is normal for the
		// last chunk. We surface io.EOF after the remaining data is consumed.
		if b.err != nil {
			if errors.Is(b.err, io.ErrUnexpectedEOF) || errors.Is(b.err, io.EOF) {
				ra.err = io.EOF
			} else {
				ra.err = b.err
			}
		}

		// Loop back to serve data from this new buffer.
	}
}

// Close signals the background goroutine to stop, drains remaining buffers,
// and (if the underlying reader is an io.ReadCloser) closes it.
func (ra *reader) Close() error {
	var err error
	ra.once.Do(func() {
		// Signal the filler goroutine to stop.
		close(ra.done)
		// Drain remaining buffers so the goroutine is not blocked on send.
		for range ra.bufs {
		}
		// Close the underlying reader if applicable.
		if ra.closer != nil {
			err = ra.closer.Close()
		}
	})
	return err
}

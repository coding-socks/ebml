package ebmltext

import (
	"bytes"
	"io"
)

// A byteReader implements a sliding window over an io.Reader.
type byteReader struct {
	data   []byte
	offset int
	r      io.ReadSeeker
	err    error
}

// release discards n bytes from the front of the window.
func (b *byteReader) release(n int) {
	b.offset += n
}

// window returns the current window.
// The window is invalidated by calls to release or extend.
func (b *byteReader) window() []byte {
	return b.data[b.offset:]
}

// tuning constants for byteReader.extend.
const (
	newBufferSize = 1024
	minReadSize   = 16
)

// extend extends the window with data from the underlying reader.
func (b *byteReader) extend() int {
	if b.err != nil {
		return 0
	}

	remaining := len(b.data) - b.offset
	if remaining <= 0 {
		b.reset()
		remaining = 0
	}
	if remaining >= minReadSize {
		// nothing to do, enough data remaining.
		return 0
	} else if cap(b.data)-remaining >= minReadSize {
		// buffer has enough space if we move the data to the front.
		b.compact()
	} else {
		// otherwise, we must allocate/extend a new buffer
		b.grow()
	}
	remaining += b.offset
	n, err := b.r.Read(b.data[remaining:cap(b.data)])
	// reduce length to the existing plus the data we read.
	b.data = b.data[:remaining+n]
	b.err = err
	return n
}

// grow grows the buffer, moving the active data to the front.
func (b *byteReader) grow() {
	buf := make([]byte, max(cap(b.data)*2, newBufferSize))
	copy(buf, b.data[b.offset:])
	b.data = buf
	b.offset = 0
}

// compact moves the active data to the front of the buffer.
func (b *byteReader) compact() {
	copy(b.data, b.data[b.offset:])
	b.offset = 0
}

func (b *byteReader) reset() {
	b.data = b.data[:0]
	b.offset = 0
	b.err = nil
}

func (b *byteReader) Read(buf []byte) (int, error) {
	window := b.window()
	n, err := io.MultiReader(bytes.NewReader(window), b.r).Read(buf)
	b.release(n)
	return n, err
}

func (b *byteReader) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekCurrent && int(offset) < (len(b.data)-b.offset) {
		b.release(int(offset))

		n, err := b.r.Seek(0, whence)
		if err != nil {
			return n, err
		}
		n -= int64(len(b.data) - b.offset)
		return n, err
	}
	if whence == io.SeekCurrent && offset != 0 {
		offset -= int64(len(b.data) - b.offset)
		b.reset()
	}
	if whence != io.SeekCurrent && len(b.data) != 0 {
		b.reset()
	}
	return b.r.Seek(offset, whence)
}

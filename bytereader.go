package rc

import "io"

type byteReader struct {
	io.Reader
}

// ReadByte reads a single byte from the underlying io.Reader. It will call Read() until
// it either succeeds or fails.
func (b *byteReader) ReadByte() (byte, error) {
	var s [1]byte
	for {
		n, err := b.Reader.Read(s[:])
		if n == 1 {
			// if we did read something, ignore error and return it
			return s[0], nil
		}
		if err != nil {
			// if there was an error, return it
			return 0, err
		}
	}
}

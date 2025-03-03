package rc

import (
	"fmt"
	"io"
	"path"
	"runtime"
)

type ByteAndReadReader interface {
	io.ByteReader
	io.Reader
}

type ReadCounter struct {
	parent ByteAndReadReader
	ctx    string
	cnt    int64
}

// ReadByte reads a single byte from the underlying stream, incrementing the counter
// by one on success.
func (rc *ReadCounter) ReadByte() (byte, error) {
	res, err := rc.parent.ReadByte()
	if err == nil {
		rc.cnt += 1
	}
	return res, err
}

// Read reads a buffer from the underlying stream, incrementing the counter accordingly
func (rc *ReadCounter) Read(p []byte) (n int, err error) {
	res, err := rc.parent.Read(p)
	rc.cnt += int64(res)
	return res, err
}

// Ret64 returns a value appropriate for ReaderAs
func (rc *ReadCounter) Ret64() (int64, error) {
	return rc.cnt, nil
}

// Ret returns a value appropriate for Read and others
func (rc *ReadCounter) Ret() (int, error) {
	return int(rc.cnt), nil
}

// Error64 returns the given error with the appropriate context
func (rc *ReadCounter) Error64(err error) (int64, error) {
	return rc.cnt, fmt.Errorf("in %s: %w", rc.ctx, err)
}

func (rc *ReadCounter) Error(err error) (int, error) {
	return int(rc.cnt), fmt.Errorf("in %s: %w", rc.ctx, err)
}

func (rc *ReadCounter) ReadFull(buf []byte) error {
	_, err := io.ReadFull(rc, buf)
	return err
}

// New returns a new instance of ReadCounter with a zero counter
func New(r io.Reader) *ReadCounter {
	// func Caller(skip int) (pc uintptr, file string, line int, ok bool)
	pc, fn, ln, ok := runtime.Caller(1)
	var ctx string
	if ok {
		f := runtime.FuncForPC(pc)
		ctx = fmt.Sprintf("%s at %s:%d", f.Name(), path.Base(fn), ln)
	} else {
		ctx = fmt.Sprintf("%s:%d", path.Base(fn), ln)
	}
	switch o := r.(type) {
	case *ReadCounter:
		return &ReadCounter{parent: o, ctx: ctx}
	case ByteAndReadReader:
		return &ReadCounter{parent: o, ctx: ctx}
	case io.Reader:
		return &ReadCounter{parent: &byteReader{o}, ctx: ctx}
	default:
		panic("object cannot be handled as ByteAndReadReader")
	}
}

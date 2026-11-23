package multihash

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"runtime"
	"testing"
	"testing/iotest"
)

type evilReader struct {
	buffer []byte
}

func (er *evilReader) Read(buf []byte) (int, error) {
	n := copy(buf, er.buffer)
	er.buffer = er.buffer[n:]
	var err error
	if len(er.buffer) == 0 {
		err = io.EOF
	}
	return n, err
}

func TestEvilReader(t *testing.T) {
	emptyHash, err := Sum(nil, IDENTITY, 0)
	if err != nil {
		t.Fatal(err)
	}
	r := NewReader(&evilReader{emptyHash})
	h, err := r.ReadMultihash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(h, []byte(emptyHash)) {
		t.Fatal(err)
	}
	h, err = r.ReadMultihash()
	if len([]byte(h)) > 0 || err != io.EOF {
		t.Fatal("expected end of file")
	}
}

func TestReader(t *testing.T) {

	var buf bytes.Buffer

	for _, tc := range testCases {
		m, err := tc.Multihash()
		if err != nil {
			t.Fatal(err)
		}

		buf.Write([]byte(m))
	}

	r := NewReader(&buf)

	for _, tc := range testCases {
		h, err := tc.Multihash()
		if err != nil {
			t.Fatal(err)
		}

		h2, err := r.ReadMultihash()
		if err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(h, h2) {
			t.Error("h and h2 should be equal")
		}
	}
}

func TestWriter(t *testing.T) {

	var buf bytes.Buffer
	w := NewWriter(&buf)

	for _, tc := range testCases {
		m, err := tc.Multihash()
		if err != nil {
			t.Error(err)
			continue
		}

		if err := w.WriteMultihash(m); err != nil {
			t.Error(err)
			continue
		}

		buf2 := make([]byte, len(m))
		if _, err := io.ReadFull(&buf, buf2); err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(m, buf2) {
			t.Error("m and buf2 should be equal")
		}
	}
}

func TestReadMultihashOverlongLength(t *testing.T) {
	// A sha2-256 header claiming a 256MiB digest.
	header := binary.AppendUvarint(nil, SHA2_256)
	header = binary.AppendUvarint(header, 1<<28)

	readers := map[string]func([]byte) io.Reader{
		"bytes.Reader": func(b []byte) io.Reader { return bytes.NewReader(b) },
		"bufio.Reader": func(b []byte) io.Reader { return bufio.NewReader(bytes.NewReader(b)) },
		"evilReader":   func(b []byte) io.Reader { return &evilReader{b} },
	}
	for name, newReader := range readers {
		t.Run(name, func(t *testing.T) {
			for _, tail := range [][]byte{nil, []byte("partial")} {
				var before, after runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&before)
				_, err := NewReader(newReader(append(header, tail...))).ReadMultihash()
				runtime.ReadMemStats(&after)

				want := io.ErrUnexpectedEOF
				if len(tail) == 0 {
					want = io.EOF
				}
				if !errors.Is(err, want) {
					t.Fatalf("tail %q: expected %v, got %v", tail, want, err)
				}
				if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 1<<20 {
					t.Fatalf("tail %q: allocated %d bytes for a truncated 256MiB digest", tail, allocated)
				}
			}
		})
	}
}

func TestReadMultihashMultiChunk(t *testing.T) {
	data := make([]byte, 3*digestPrealloc+17)
	for i := range data {
		data[i] = byte(i)
	}
	want, err := Sum(data, IDENTITY, -1)
	if err != nil {
		t.Fatal(err)
	}

	readers := map[string]io.Reader{
		"bytes.Reader": bytes.NewReader(want),
		"OneByte":      iotest.OneByteReader(bytes.NewReader(want)),
		"DataErr":      iotest.DataErrReader(bytes.NewReader(want)),
	}
	for name, r := range readers {
		t.Run(name, func(t *testing.T) {
			got, err := NewReader(r).ReadMultihash()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatal("round trip changed the multihash")
			}
		})
	}

	// Truncated after the preallocated region.
	_, err = NewReader(bytes.NewReader(want[:len(want)-digestPrealloc])).ReadMultihash()
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected %v, got %v", io.ErrUnexpectedEOF, err)
	}
}

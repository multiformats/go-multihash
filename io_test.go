package multihash

import (
	"bytes"
	"encoding/binary"
	"io"
	"runtime"
	"testing"
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

// TestReadMultihashOverlongLength checks that a digest length larger than the data
// behind it does not size the buffer before that data has arrived, and that a
// length matching the data is still accepted. The assertion is on allocation
// volume, not merely on getting an error: an overlong length errors either way,
// once io.ReadFull runs out of input.
func TestReadMultihashOverlongLength(t *testing.T) {
	// sha2-256 code, then a length of 2^28, then nothing.
	var hostile []byte
	hostile = binary.AppendUvarint(hostile, SHA2_256)
	hostile = binary.AppendUvarint(hostile, 1<<28)

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	_, err := NewReader(bytes.NewReader(hostile)).ReadMultihash()
	runtime.ReadMemStats(&after)

	if err == nil {
		t.Fatal("expected an error for a length with no digest behind it")
	}
	if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 1<<20 {
		t.Fatalf("reading a multihash that claims a 256MiB digest allocated %d bytes", allocated)
	}

	// The same shape with a length that the data does satisfy must still parse.
	good, err := Sum([]byte("foo"), SHA2_256, -1)
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewReader(bytes.NewReader(good)).ReadMultihash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, good) {
		t.Fatalf("round trip changed the multihash: %x != %x", got, good)
	}
}

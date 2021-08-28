// Generate hasher input that matches BLAKE3's test vectors.
//
// See:
//
// https://github.com/BLAKE3-team/BLAKE3/blob/080b3330159a19407dddb407dc917925ac40c4d3/test_vectors/test_vectors.json
package main

import (
	"fmt"
	"os"
	"strconv"
)

var usage = `usage: %s n
Generate n bytes of input for blake3 test vector.
`

const BYTE_LOOP = 251

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}

	c, err := strconv.ParseUint(os.Args[1], 0, 32)
	if err != nil {
		die(err)
	}

	var b [BYTE_LOOP]uint8
	for i := range b {
		b[i] = uint8(i)
	}

	for c > 0 {
		var w uint64 = BYTE_LOOP
		if c < BYTE_LOOP {
			w = c
		}

		_, err := os.Stdout.Write(b[:w])
		if err != nil {
			die(err)
		}

		c -= w
	}
}

func die(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

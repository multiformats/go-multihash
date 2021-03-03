package multihash

import (
	"errors"
	"fmt"
)

// ErrSumNotSupported is returned when the Sum function code is not implemented
var ErrSumNotSupported = errors.New("Function not implemented. Complain to lib maintainer.")

var ErrLenTooLarge = errors.New("requested length was too large for digest")

// Sum obtains the cryptographic sum of a given buffer. The length parameter
// indicates the length of the resulting digest and passing a negative value
// use default length values for the selected hash function.
func Sum(data []byte, code uint64, length int) (Multihash, error) {
	// Get the algorithm.
	hasher, err := GetHasher(code)
	if err != nil {
		return nil, err
	}

	// Feed data in.
	hasher.Write(data)

	// Compute hash.
	//  Use a fixed size array here: should keep things on the stack.
	var space [64]byte
	sum := hasher.Sum(space[0:0])

	// Deal with any truncation.
	if length < 0 {
		var ok bool
		length, ok = DefaultLengths[code]
		if !ok {
			return nil, fmt.Errorf("no default length for code %d", code)
		}
	}
	if len(sum) < length {
		return nil, ErrLenTooLarge
	}
	if length >= 0 {
		sum = sum[:length]
	}

	// Put the multihash metainfo bytes at the front of the buffer.
	// FIXME: this does many avoidable allocations, but it's the shape of the Encode method arguments that forces this.
	return Encode(sum, code)
}

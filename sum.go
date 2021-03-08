package multihash

import (
	"errors"
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

	// Compute final hash.
	//  A new slice is allocated.  FUTURE: see other comment below about allocation, and review together with this line to try to improve.
	sum := hasher.Sum(nil)

	// Deal with any truncation.
	if length < 0 {
		length = hasher.Size()
	}
	if len(sum) < length {
		return nil, ErrLenTooLarge
	}
	if length >= 0 {
		sum = sum[:length]
	}

	// Put the multihash metainfo bytes at the front of the buffer.
	//  FUTURE: try to improve allocations here.  Encode does several which are probably avoidable, but it's the shape of the Encode method arguments that forces this.
	return Encode(sum, code)
}

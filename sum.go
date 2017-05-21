package multihash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"

	keccak "gx/ipfs/QmQPWTeQJnJE7MYu6dJTiNTQRNuqBr41dis6UgY6Uekmgd/keccakpg"
	blake2b "gx/ipfs/QmaPHkZLbQQbvcyavn8q1GFHg6o6yeceyHFSJ3Pjf3p3TQ/go-crypto/blake2b"
	blake2s "gx/ipfs/QmaPHkZLbQQbvcyavn8q1GFHg6o6yeceyHFSJ3Pjf3p3TQ/go-crypto/blake2s"
	sha3 "gx/ipfs/QmaPHkZLbQQbvcyavn8q1GFHg6o6yeceyHFSJ3Pjf3p3TQ/go-crypto/sha3"
	"gx/ipfs/QmfJHywXQu98UeZtGJBQrPAR6AtmDjjbe3qjTo9piXHPnx/murmur3"
	skein "leb.io/hashland/skein" // could use https://github.com/wernerd/Skein3Fish/tree/master/go as well
)

// ErrSumNotSupported is returned when the Sum function code is not implemented
var ErrSumNotSupported = errors.New("Function not implemented. Complain to lib maintainer.")

// Sum obtains the cryptographic sum of a given buffer. The length parameter
// indicates the length of the resulting digest and passing a negative value
// use default length values for the selected hash function.
func Sum(data []byte, code uint64, length int) (Multihash, error) {
	m := Multihash{}
	err := error(nil)
	if !ValidCode(code) {
		return m, fmt.Errorf("invalid multihash code %d", code)
	}

	if length < 0 {
		var ok bool
		length, ok = DefaultLengths[code]
		if !ok {
			return m, fmt.Errorf("no default length for code %d", code)
		}
	}

	var d []byte
	switch {
	case isBlake2s(code):
		olen := code - BLAKE2S_MIN + 1
		switch olen {
		case 32:
			out := blake2s.Sum256(data)
			d = out[:]
		default:

			// var sum [32]byte
			// var tmp [olen]byte
			// checkSum(&sum, olen, data)
			// copy(tmp[:], sum[:olen])
			// d = tmp

			return nil, fmt.Errorf("unsupported length for blake2s: %d", olen)
		}
	case isBlake2b(code):
		olen := code - BLAKE2B_MIN + 1
		switch olen {
		case 32:
			out := blake2b.Sum256(data)
			d = out[:]
		case 48:
			out := blake2b.Sum384(data)
			d = out[:]
		case 64:
			out := blake2b.Sum512(data)
			d = out[:]
		default:

			// var sum [64]byte
			// var tmp [olen]byte
			// checkSum(&sum, olen, data)
			// copy(tmp[:], sum[:olen])
			// d = tmp

			return nil, fmt.Errorf("unsupported length for blake2b: %d", olen)
		}
	case isSkein256(code):
		olen := code - SKEIN256_MIN + 1
		state, _ := skein.New(256, int(olen))
		out := state.Sum(data)
		d = out[:]
	case isSkein512(code):
		olen := code - SKEIN512_MIN + 1
		state, _ := skein.New(512, int(olen))
		out := state.Sum(data)
		d = out[:]
	case isSkein1024(code):
		olen := code - SKEIN1024_MIN + 1
		state, _ := skein.New(1024, int(olen))
		out := state.Sum(data)
		d = out[:]
	default:
		switch code {
		case ID:
			d = sumID(data)
		case SHA1:
			d = sumSHA1(data)
		case SHA2_256:
			d = sumSHA256(data)
		case SHA2_512:
			d = sumSHA512(data)
		case KECCAK_224:
			d = sumKeccak224(data)
		case KECCAK_256:
			d = sumKeccak256(data)
		case KECCAK_384:
			d = sumKeccak384(data)
		case KECCAK_512:
			d = sumKeccak512(data)
		case SHA3_224:
			d = sumSHA3_224(data)
		case SHA3_256:
			d = sumSHA3_256(data)
		case SHA3_384:
			d = sumSHA3_384(data)
		case SHA3_512:
			d = sumSHA3_512(data)
		case DBL_SHA2_256:
			d = sumSHA256(sumSHA256(data))
		case MURMUR3:
			d, err = sumMURMUR3(data)
		case SHAKE_128:
			d = sumSHAKE128(data)
		case SHAKE_256:
			d = sumSHAKE256(data)
		default:
			return m, ErrSumNotSupported
		}
	}
	if err != nil {
		return m, err
	}
	return Encode(d[0:length], code)
}

func isBlake2s(code uint64) bool {
	return code >= BLAKE2S_MIN && code <= BLAKE2S_MAX
}
func isBlake2b(code uint64) bool {
	return code >= BLAKE2B_MIN && code <= BLAKE2B_MAX
}

func isSkein256(code uint64) bool {
	return code >= SKEIN256_MIN && code <= SKEIN256_MAX
}

func isSkein512(code uint64) bool {
	return code >= SKEIN512_MIN && code <= SKEIN512_MAX
}

func isSkein1024(code uint64) bool {
	return code >= SKEIN1024_MIN && code <= SKEIN1024_MAX
}

func sumID(data []byte) []byte {
	return data
}

func sumSHA1(data []byte) []byte {
	a := sha1.Sum(data)
	return a[0:20]
}

func sumSHA256(data []byte) []byte {
	a := sha256.Sum256(data)
	return a[0:32]
}

func sumSHA512(data []byte) []byte {
	a := sha512.Sum512(data)
	return a[0:64]
}

func sumKeccak224(data []byte) []byte {
	h := keccak.New224()
	h.Write(data)
	return h.Sum(nil)
}

func sumKeccak256(data []byte) []byte {
	h := keccak.New256()
	h.Write(data)
	return h.Sum(nil)
}

func sumKeccak384(data []byte) []byte {
	h := keccak.New384()
	h.Write(data)
	return h.Sum(nil)
}

func sumKeccak512(data []byte) []byte {
	h := keccak.New512()
	h.Write(data)
	return h.Sum(nil)
}

func sumSHA3(data []byte) ([]byte, error) {
	h := sha3.New512()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func sumSHA3_512(data []byte) []byte {
	a := sha3.Sum512(data)
	return a[:]
}

func sumMURMUR3(data []byte) ([]byte, error) {
	number := murmur3.Sum32(data)
	bytes := make([]byte, 4)
	for i := range bytes {
		bytes[i] = byte(number & 0xff)
		number >>= 8
	}
	return bytes, nil
}

func sumSHAKE128(data []byte) []byte {
	bytes := make([]byte, 32)
	sha3.ShakeSum128(bytes, data)
	return bytes
}

func sumSHAKE256(data []byte) []byte {
	bytes := make([]byte, 64)
	sha3.ShakeSum256(bytes, data)
	return bytes
}

func sumSHA3_384(data []byte) []byte {
	a := sha3.Sum384(data)
	return a[:]
}

func sumSHA3_256(data []byte) []byte {
	a := sha3.Sum256(data)
	return a[:]
}

func sumSHA3_224(data []byte) []byte {
	a := sha3.Sum224(data)
	return a[:]
}

package multihash

import (
	"crypto/sha1"
	"crypto/sha512"
	"errors"
	"fmt"

	keccak "github.com/gxed/hashland/keccakpg"
	blake2b "github.com/minio/blake2b-simd"
	sha256 "github.com/minio/sha256-simd"
	murmur3 "github.com/spaolacci/murmur3"
	blake2s "golang.org/x/crypto/blake2s"
	sha3 "golang.org/x/crypto/sha3"
)

// ErrSumNotSupported is returned when the Sum function code is not implemented
var ErrSumNotSupported = errors.New("Function not implemented. Complain to lib maintainer.")

// funcTable maps multicodec values to hash functions.
var funcTable = make(map[uint64]func([]byte) []byte)

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
	// TODO: Consider how to register blake funcs also. This length based is a
	// special case which is more complex.
	case isBlake2s(code):
		olen := code - BLAKE2S_MIN + 1
		switch olen {
		case 32:
			out := blake2s.Sum256(data)
			d = out[:]
		default:
			return nil, fmt.Errorf("unsupported length for blake2s: %d", olen)
		}
	case isBlake2b(code):
		olen := uint8(code - BLAKE2B_MIN + 1)
		d = sumBlake2b(olen, data)
	default:
		switch code {
		case ID:
			d, err = sumID(data, length)
		default:
			hashFunc, ok := funcTable[code]
			if !ok {
				return m, ErrSumNotSupported
			}
			d = hashFunc(data)
		}
	}
	if err != nil {
		return m, err
	}
	if length >= 0 {
		d = d[:length]
	}
	return Encode(d, code)
}

func isBlake2s(code uint64) bool {
	return code >= BLAKE2S_MIN && code <= BLAKE2S_MAX
}
func isBlake2b(code uint64) bool {
	return code >= BLAKE2B_MIN && code <= BLAKE2B_MAX
}

func sumBlake2b(size uint8, data []byte) []byte {
	hasher, err := blake2b.New(&blake2b.Config{Size: size})
	if err != nil {
		panic(err)
	}

	if _, err := hasher.Write(data); err != nil {
		panic(err)
	}

	return hasher.Sum(nil)[:]
}

func sumID(data []byte, length int) ([]byte, error) {
	if length >= 0 && length != len(data) {
		return nil, fmt.Errorf("the length of the identity hash (%d) must be equal to the length of the data (%d)",
			length, len(data))

	}
	return data, nil
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

func sumSHA3_512(data []byte) []byte {
	a := sha3.Sum512(data)
	return a[:]
}

func sumMURMUR3(data []byte) []byte {
	number := murmur3.Sum32(data)
	bytes := make([]byte, 4)
	for i := range bytes {
		bytes[i] = byte(number & 0xff)
		number >>= 8
	}
	return bytes
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

func registerStdlibHashFuncs() {
	RegisterHashFunc(SHA1, sumSHA1)
	RegisterHashFunc(SHA2_512, sumSHA512)
}

func registerNonStdlibHashFuncs() {
	RegisterHashFunc(SHA2_256, sumSHA256)
	RegisterHashFunc(DBL_SHA2_256, func(data []byte) []byte {
		return sumSHA256(sumSHA256(data))
	})

	RegisterHashFunc(KECCAK_224, sumKeccak224)
	RegisterHashFunc(KECCAK_256, sumKeccak256)
	RegisterHashFunc(KECCAK_384, sumKeccak384)
	RegisterHashFunc(KECCAK_512, sumKeccak512)

	RegisterHashFunc(SHA3_224, sumSHA3_224)
	RegisterHashFunc(SHA3_256, sumSHA3_256)
	RegisterHashFunc(SHA3_384, sumSHA3_384)
	RegisterHashFunc(SHA3_512, sumSHA3_512)

	RegisterHashFunc(MURMUR3, sumMURMUR3)

	RegisterHashFunc(SHAKE_128, sumSHAKE128)
	RegisterHashFunc(SHAKE_256, sumSHAKE256)
}

func init() {
	registerStdlibHashFuncs()
	registerNonStdlibHashFuncs()
}

// RegisterHashFunc adds an entry to the package-level code -> hash func map.
func RegisterHashFunc(code uint64, hashFunc func([]byte) []byte) error {
	if !ValidCode(code) {
		return fmt.Errorf("code %v not valid", code)
	}

	_, ok := funcTable[code]
	if ok {
		return fmt.Errorf("hash func for code %v already registered", code)
	}

	funcTable[code] = hashFunc
	return nil
}

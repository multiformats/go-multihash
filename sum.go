package multihash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"

	blake2b "golang.org/x/crypto/blake2b"
	blake2s "golang.org/x/crypto/blake2s"
	sha3 "golang.org/x/crypto/sha3"
)

var ErrSumNotSupported = errors.New("Function not implemented. Complain to lib maintainer.")

func Sum(data []byte, code int, length int) (Multihash, error) {
	m := Multihash{}
	err := error(nil)
	if !ValidCode(code) {
		return m, fmt.Errorf("invalid multihash code %d", code)
	}

	var d []byte
	switch code {
	case SHA1:
		d = sumSHA1(data)
	case SHA2_256:
		d = sumSHA256(data)
	case SHA2_512:
		d = sumSHA512(data)
	case SHA3:
		d, err = sumSHA3(data)
	case BLAKE2B:
		d = sumBLAKE2B(data)
	case BLAKE2S:
		d = sumBLAKE2S(data)
	case BLAKE2B_256:
		d = sumBLAKE2B_256(data)
	default:
		return m, ErrSumNotSupported
	}

	if err != nil {
		return m, err
	}

	if length < 0 {
		var ok bool
		length, ok = DefaultLengths[code]
		if !ok {
			return m, fmt.Errorf("no default length for code %d", code)
		}
	}

	return Encode(d[0:length], code)
}

func sumSHA1(data []byte) []byte {
	a := sha1.Sum(data)
	return a[0:sha1.Size]
}

func sumSHA256(data []byte) []byte {
	a := sha256.Sum256(data)
	return a[0:sha256.Size]
}

func sumSHA512(data []byte) []byte {
	a := sha512.Sum512(data)
	return a[0:sha512.Size]
}

func sumSHA3(data []byte) ([]byte, error) {
	h := sha3.New512()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func sumBLAKE2B(data []byte) []byte {
	a := blake2b.Sum512(data)
	return a[0:blake2b.Size]
}

func sumBLAKE2S(data []byte) []byte {
	a := blake2s.Sum256(data)
	return a[0:blake2s.Size]
}

func sumBLAKE2B_256(data []byte) []byte {
	a := blake2b.Sum256(data)
	return a[0:blake2b.Size256]
}

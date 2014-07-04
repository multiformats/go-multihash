package multihash

import (
	"fmt"
)

var Names = map[string]int{
	"sha1":     0x11,
	"sha2-256": 0x12,
	"sha2-512": 0x13,
	"sha3":     0x14,
	"blake2b":  0x40,
	"blake2s":  0x41,
}

var Codes = map[int]string{
	0x11: "sha1",
	0x12: "sha2-256",
	0x13: "sha2-512",
	0x14: "sha3",
	0x40: "blake2b",
	0x41: "blake2s",
}

var DefaultLengths = map[int]int{
	0x11: 20,
	0x12: 32,
	0x13: 64,
	0x14: 64,
	0x40: 64,
	0x41: 32,
}

type DecodedMultihash struct {
	Code   int
	Name   string
	Length int
	Digest []byte
}

type Multihash []byte

func Cast(buf []byte) (Multihash, error) {
	dm, err := Decode(buf)
	if err != nil {
		return Multihash{}, err
	}

	if !ValidCode(dm.Code) {
		return Multihash{}, fmt.Errorf("unknown multihash code")
	}

	return Multihash(buf), nil
}

// Decodes a hash from the given Multihash.
func Decode(buf []byte) (*DecodedMultihash, error) {

	if len(buf) < 3 {
		return nil, fmt.Errorf("multihash too short. must be > 3 bytes.")
	}

	if len(buf) > 129 {
		return nil, fmt.Errorf("multihash too long. must be < 129 bytes.")
	}

	dm := &DecodedMultihash{
		Code:   int(uint8(buf[0])),
		Name:   Codes[int(uint8(buf[0]))],
		Length: int(uint8(buf[1])),
		Digest: buf[2:],
	}

	if len(dm.Digest) != dm.Length {
		return nil, fmt.Errorf("multihash length inconsistent: %v", dm)
	}

	return dm, nil
}

// Encodes a hash digest along with the specified function code.
// Note: the length is derived from the length of the digest itself.
func Encode(buf []byte, code int) ([]byte, error) {

	if !ValidCode(code) {
		return nil, fmt.Errorf("unknown multihash code")
	}

	if len(buf) > 127 {
		m := "multihash does not yet support digests longer than 127 bytes."
		return nil, fmt.Errorf(m)
	}

	pre := make([]byte, 2)
	pre[0] = byte(uint8(code))
	pre[1] = byte(uint8(len(buf)))
	return append(pre, buf...), nil
}

func EncodeName(buf []byte, name string) ([]byte, error) {
	return Encode(buf, Names[name])
}

// Checks whether a multihash code is valid.
func ValidCode(code int) bool {
	if AppCode(code) {
		return true
	}

	if _, ok := Codes[code]; ok {
		return true
	}

	return false
}

// Checks whether a multihash code is part of the App range.
func AppCode(code int) bool {
	return code > 0 && code < 0x10
}

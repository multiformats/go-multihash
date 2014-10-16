package multihash

import (
	"encoding/hex"
	"errors"
	"fmt"

	b58 "github.com/jbenet/go-base58"
)

// errors
var (
	ErrorUknownCode    = errors.New("unknown multihash code")
	ErrTooShort        = errors.New("multihash too short. must be > 3 bytes")
	ErrTooLong         = errors.New("multihash too long. must be < 129 bytes")
	ErrLenNotSupported = errors.New("multihash does not yet support digests longer than 127 bytes")
)

type ErrInconsistentLen struct {
	dm *DecodedMultihash
}

func (e ErrInconsistentLen) Error() string {
	return fmt.Sprintf("multihash length inconsistent: %v", e.dm)
}

// constants
const SHA1 = 0x11
const SHA2_256 = 0x12
const SHA2_512 = 0x13
const SHA3 = 0x14
const BLAKE2B = 0x40
const BLAKE2S = 0x41

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

func (m *Multihash) HexString() string {
	return hex.EncodeToString([]byte(*m))
}

func (m *Multihash) String() string {
	return m.HexString()
}

func FromHexString(s string) (Multihash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Multihash{}, err
	}

	return Cast(b)
}

func (m Multihash) B58String() string {
	return b58.Encode([]byte(m))
}

func FromB58String(s string) (m Multihash, err error) {
	// panic handler, in case we try accessing bytes incorrectly.
	defer func() {
		if e := recover(); e != nil {
			m = Multihash{}
			err = e.(error)
		}
	}()

	//b58 smells like it can panic...
	b := b58.Decode(s)
	return Cast(b)
}

func Cast(buf []byte) (Multihash, error) {
	dm, err := Decode(buf)
	if err != nil {
		return Multihash{}, err
	}

	if !ValidCode(dm.Code) {
		return Multihash{}, ErrorUknownCode
	}

	return Multihash(buf), nil
}

// Decodes a hash from the given Multihash.
func Decode(buf []byte) (*DecodedMultihash, error) {

	if len(buf) < 3 {
		return nil, ErrTooShort
	}

	if len(buf) > 129 {
		return nil, ErrTooLong
	}

	dm := &DecodedMultihash{
		Code:   int(uint8(buf[0])),
		Name:   Codes[int(uint8(buf[0]))],
		Length: int(uint8(buf[1])),
		Digest: buf[2:],
	}

	if len(dm.Digest) != dm.Length {
		return nil, ErrInconsistentLen{dm}
	}

	return dm, nil
}

// Encodes a hash digest along with the specified function code.
// Note: the length is derived from the length of the digest itself.
func Encode(buf []byte, code int) ([]byte, error) {

	if !ValidCode(code) {
		return nil, ErrorUknownCode
	}

	if len(buf) > 127 {
		return nil, ErrLenNotSupported
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
	return code >= 0 && code < 0x10
}

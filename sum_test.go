package multihash_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"sort"
	"testing"

	"github.com/multiformats/go-multihash"
	_ "github.com/multiformats/go-multihash/register/all"
)

type SumTestCase struct {
	code             uint64
	length           int
	input            string
	hex              string
	expectedSumError error
}

var sumTestCases = []SumTestCase{
	{multihash.IDENTITY, 3, "foo", "0003666f6f", nil},
	{multihash.IDENTITY, -1, "foofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoo", "0030666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f666f6f", nil},
	{multihash.SHA1, -1, "foo", "11140beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33", nil},
	{multihash.SHA1, 10, "foo", "110a0beec7b5ea3f0fdbc95d", nil},
	{multihash.SHA2_256, -1, "foo", "12202c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae", nil},
	{multihash.SHA2_256, 31, "foo", "121f2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7", nil},
	{multihash.SHA2_256, 32, "foo", "12202c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae", nil},
	{multihash.SHA2_256, 16, "foo", "12102c26b46b68ffc68ff99b453c1d304134", nil},
	{multihash.SHA2_512, -1, "foo", "1340f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7", nil},
	{multihash.SHA2_512, 32, "foo", "1320f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc663832", nil},
	{multihash.SHA3, 32, "foo", "14204bca2b137edc580fe50a88983ef860ebaca36c857b1f492839d6d7392452a63c", nil},
	{multihash.SHA3_512, 16, "foo", "14104bca2b137edc580fe50a88983ef860eb", nil},
	{multihash.SHA3_512, -1, "foo", "14404bca2b137edc580fe50a88983ef860ebaca36c857b1f492839d6d7392452a63c82cbebc68e3b70a2a1480b4bb5d437a7cba6ecf9d89f9ff3ccd14cd6146ea7e7", nil},
	{multihash.SHA3_224, -1, "beep boop", "171c0da73a89549018df311c0a63250e008f7be357f93ba4e582aaea32b8", nil},
	{multihash.SHA3_224, 16, "beep boop", "17100da73a89549018df311c0a63250e008f", nil},
	{multihash.SHA3_256, -1, "beep boop", "1620828705da60284b39de02e3599d1f39e6c1df001f5dbf63c9ec2d2c91a95a427f", nil},
	{multihash.SHA3_256, 16, "beep boop", "1610828705da60284b39de02e3599d1f39e6", nil},
	{multihash.SHA3_384, -1, "beep boop", "153075a9cff1bcfbe8a7025aa225dd558fb002769d4bf3b67d2aaf180459172208bea989804aefccf060b583e629e5f41e8d", nil},
	{multihash.SHA3_384, 16, "beep boop", "151075a9cff1bcfbe8a7025aa225dd558fb0", nil},
	{multihash.DBL_SHA2_256, 32, "foo", "5620c7ade88fc7a21498a6a5e5c385e1f68bed822b72aa63c4a9a48a02c2466ee29e", nil},
	{multihash.BLAKE2B_MAX, -1, "foo", "c0e40240ca002330e69d3e6b84a46a56a6533fd79d51d97a3bb7cad6c2ff43b354185d6dc1e723fb3db4ae0737e120378424c714bb982d9dc5bbd7a0ab318240ddd18f8d", nil},
	{multihash.BLAKE2B_MAX, 64, "foo", "c0e40240ca002330e69d3e6b84a46a56a6533fd79d51d97a3bb7cad6c2ff43b354185d6dc1e723fb3db4ae0737e120378424c714bb982d9dc5bbd7a0ab318240ddd18f8d", nil},
	{multihash.BLAKE2B_MAX - 32, -1, "foo", "a0e40220b8fe9f7f6255a6fa08f668ab632a8d081ad87983c77cd274e48ce450f0b349fd", nil},
	{multihash.BLAKE2B_MAX - 32, 32, "foo", "a0e40220b8fe9f7f6255a6fa08f668ab632a8d081ad87983c77cd274e48ce450f0b349fd", nil},
	{multihash.BLAKE2B_MAX - 19, -1, "foo", "ade4022dca82ab956d5885e3f5db10cca94182f01a6ca2c47f9f4228497dcc9f4a0121c725468b852a71ec21fcbeb725df", nil},
	{multihash.BLAKE2B_MAX - 19, 45, "foo", "ade4022dca82ab956d5885e3f5db10cca94182f01a6ca2c47f9f4228497dcc9f4a0121c725468b852a71ec21fcbeb725df", nil},
	{multihash.BLAKE2B_MAX - 16, -1, "foo", "b0e40230e629ee880953d32c8877e479e3b4cb0a4c9d5805e2b34c675b5a5863c4ad7d64bb2a9b8257fac9d82d289b3d39eb9cc2", nil},
	{multihash.BLAKE2B_MAX - 16, 48, "foo", "b0e40230e629ee880953d32c8877e479e3b4cb0a4c9d5805e2b34c675b5a5863c4ad7d64bb2a9b8257fac9d82d289b3d39eb9cc2", nil},
	{multihash.BLAKE2B_MIN + 19, -1, "foo", "94e40214983ceba2afea8694cc933336b27b907f90c53a88", nil},
	{multihash.BLAKE2B_MIN + 19, 20, "foo", "94e40214983ceba2afea8694cc933336b27b907f90c53a88", nil},
	{multihash.BLAKE2B_MIN, -1, "foo", "81e4020152", nil},
	{multihash.BLAKE2B_MIN, 1, "foo", "81e4020152", nil},
	{multihash.BLAKE2S_MAX, 32, "foo", "e0e4022008d6cad88075de8f192db097573d0e829411cd91eb6ec65e8fc16c017edfdb74", nil},
	{multihash.KECCAK_256, 32, "foo", "1b2041b1a0649752af1b28b3dc29a1556eee781e4a4c3a1f7f53f90fa834de098c4d", nil},
	{multihash.KECCAK_512, -1, "beep boop", "1d40e161c54798f78eba3404ac5e7e12d27555b7b810e7fd0db3f25ffa0c785c438331b0fbb6156215f69edf403c642e5280f4521da9bd767296ec81f05100852e78", nil},
	{multihash.SHAKE_128, 32, "foo", "1820f84e95cb5fbd2038863ab27d3cdeac295ad2d4ab96ad1f4b070c0bf36078ef08", nil},
	{multihash.SHAKE_256, 64, "foo", "19401af97f7818a28edfdfce5ec66dbdc7e871813816d7d585fe1f12475ded5b6502b7723b74e2ee36f2651a10a8eaca72aa9148c3c761aaceac8f6d6cc64381ed39", nil},
	{multihash.MD5, -1, "foo", "d50110acbd18db4cc2f85cedef654fccc4a4d8", nil},
	{multihash.BLAKE3, 32, "foo", "1e2004e0bb39f30b1a3feb89f536c93be15055482df748674b00d26e5a75777702e9", nil},
	{multihash.BLAKE3, 64, "foo", "1e4004e0bb39f30b1a3feb89f536c93be15055482df748674b00d26e5a75777702e9791074b7511b59d31c71c62f5a745689fa6c9497f68bdf1061fe07f518d410c0", nil},
	{multihash.BLAKE3, 128, "foo", "1e800104e0bb39f30b1a3feb89f536c93be15055482df748674b00d26e5a75777702e9791074b7511b59d31c71c62f5a745689fa6c9497f68bdf1061fe07f518d410c0b0c27f41b3cf083f8a7fdc67a877e21790515762a754a45dcb8a356722698a7af5ed2bb608983d5aa75d4d61691ef132efe8631ce0afc15553a08fffc60ee936", nil},
	{multihash.BLAKE3, -1, "foo", "1e2004e0bb39f30b1a3feb89f536c93be15055482df748674b00d26e5a75777702e9", nil},
	{multihash.BLAKE3, 129, "foo", "1e810104e0bb39f30b1a3feb89f536c93be15055482df748674b00d26e5a75777702e9791074b7511b59d31c71c62f5a745689fa6c9497f68bdf1061fe07f518d410c0b0c27f41b3cf083f8a7fdc67a877e21790515762a754a45dcb8a356722698a7af5ed2bb608983d5aa75d4d61691ef132efe8631ce0afc15553a08fffc60ee9369b", multihash.ErrLenTooLarge},
}

func TestSum(t *testing.T) {
	for _, tc := range sumTestCases {
		m1, err := multihash.FromHexString(tc.hex)
		if err != nil {
			t.Error(err)
			continue
		}

		m2, err := multihash.Sum([]byte(tc.input), tc.code, tc.length)
		if err != tc.expectedSumError {
			t.Error(tc.code, "sum failed or succeeded unexpectedly.", err)
			continue
		} else if err != nil {
			// test case was expected to fail at the sum, and failed as expected. stop the test.
			continue
		}

		if !bytes.Equal(m1, m2) {
			t.Error(tc.code, multihash.Codes[tc.code], "sum failed.", m1, m2)
			t.Error(hex.EncodeToString(m2))
		}

		s1 := m1.HexString()
		if s1 != tc.hex {
			t.Error("hex strings not the same")
		}

		s2 := m1.B58String()
		m3, err := multihash.FromB58String(s2)
		if err != nil {
			t.Error("failed to decode b58")
		} else if !bytes.Equal(m3, m1) {
			t.Error("b58 failing bytes")
		} else if s2 != m3.B58String() {
			t.Error("b58 failing string")
		}
	}
}

func BenchmarkSum(b *testing.B) {
	tc := sumTestCases[0]
	for i := 0; i < b.N; i++ {
		_, _ = multihash.Sum([]byte(tc.input), tc.code, tc.length)
	}
}

func BenchmarkBlake2B(b *testing.B) {
	sizes := []uint64{128, 129, 130, 255, 256, 257, 386, 512}
	for _, s := range sizes {
		func(si uint64) {
			b.Run(fmt.Sprintf("blake2b-%d", s), func(b *testing.B) {
				arr := []byte("test data for some hashing, this is broken")
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					m, err := multihash.Sum(arr, multihash.BLAKE2B_MIN+si/8-1, -1)
					if err != nil {
						b.Fatal(err)
					}
					runtime.KeepAlive(m)
				}
			})
		}(s)
	}
}

func TestSmallerLengthHashID(t *testing.T) {

	data := []byte("Identity hash input data.")
	dataLength := len(data)

	// Normal case: `length == len(data)`.
	_, err := multihash.Sum(data, multihash.IDENTITY, dataLength)
	if err != nil {
		t.Fatal(err)
	}

	// Unconstrained length (-1): also allowed.
	_, err = multihash.Sum(data, multihash.IDENTITY, -1)
	if err != nil {
		t.Fatal(err)
	}

	// Any other variation of those two scenarios should fail.
	for l := dataLength - 1; l >= 0; l-- {
		_, err = multihash.Sum(data, multihash.IDENTITY, l)
		if err == nil {
			t.Fatalf("identity hash of length %d smaller than data length %d didn't fail",
				l, dataLength)
		}
	}
}

func TestTooLargeLength(t *testing.T) {
	_, err := multihash.Sum([]byte("test"), multihash.SHA2_256, 33)
	if err != multihash.ErrLenTooLarge {
		t.Fatal("bad error", err)
	}
}

func TestBasicSum(t *testing.T) {
	for code, name := range multihash.Codes {
		_, err := multihash.Sum([]byte("test"), code, -1)
		switch {
		case errors.Is(err, multihash.ErrSumNotSupported):
		case err == nil:
		default:
			t.Errorf("unexpected error for %s: %s", name, err)
		}
	}
}

var Sink []byte

type codeNamePair struct {
	id   uint64
	name string
}

func BenchmarkSumAllLarge(b *testing.B) {
	var data [1024 * 1024 * 16]byte
	src := rand.New(rand.NewSource(0x4242424242424242))
	_, err := io.ReadFull(src, data[:])
	if err != nil {
		b.Fatal(err)
	}

	// Write to a slice to sort elements
	s := make([]codeNamePair, 0, len(multihash.Codes))
	for id, name := range multihash.Codes {
		s = append(s, codeNamePair{id, name})
	}

	sort.Slice(s, func(x, y int) bool {
		return s[x].id < s[y].id
	})

	for _, v := range s {
		h, e := multihash.GetHasher(v.id)
		if h == nil || e != nil {
			continue // Don't benchmark unsupported hashing functions
		}

		b.Run(v.name, func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			for i := b.N; i > 0; i-- {
				var err error
				Sink, err = multihash.Sum(data[:], v.id, -1)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

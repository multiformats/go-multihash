/*
	This package has no purpose except to perform registration of multihashes.

	It is meant to be used as a side-effecting import, e.g.

		import (
			_ "github.com/multiformats/go-multihash/register/blake2"
		)

	This package registers several multihashes for the blake2 family
	(both the 's' and the 'b' variants, and in a variety of sizes).
*/
package blake2

import (
	"hash"

	"github.com/minio/blake2b-simd"
	"golang.org/x/crypto/blake2s"

	"github.com/multiformats/go-multihash"
)

const (
	BLAKE2B_MIN = 0xb201
	BLAKE2B_MAX = 0xb240
	BLAKE2S_MIN = 0xb241
	BLAKE2S_MAX = 0xb260
)

func init() {
	// BLAKE2S
	// This package only enables support for 32byte (256 bit) blake2s.
	multihash.Register(BLAKE2S_MIN+31, func() hash.Hash { h, _ := blake2s.New256(nil); return h })

	// BLAKE2B
	// There's a whole range of these.
	for c := uint64(BLAKE2B_MIN); c <= BLAKE2B_MAX; c++ {
		size := int(c - BLAKE2B_MIN + 1)

		// special case these lengths to avoid allocations.
		switch size {
		case 32:
			multihash.Register(c, blake2b.New256)
			continue
		case 64:
			multihash.Register(c, blake2b.New512)
			continue
		}

		// Ok, allocate away.
		//  (The config object here being a pointer is a tad unfortunate,
		//   but we manage amortize it away by making them just once anyway.)
		cfg := &blake2b.Config{Size: uint8(size)}
		multihash.Register(c, func() hash.Hash {
			hasher, err := blake2b.New(cfg)
			if err != nil {
				panic(err)
			}
			return hasher
		})
	}
}

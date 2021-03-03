/*
	This package has no purpose except to perform registration of multihashes.

	It is meant to be used as a side-effecting import, e.g.

		import (
			_ "github.com/multiformats/go-multihash/register/murmur3"
		)

	This package registers multihashes for the murmur3 family.
*/
package murmur3

// import (
// 	"hash"
//
// 	"github.com/gxed/hashland/murmur3"
//
// 	"github.com/multiformats/go-multihash"
// )

func init() {
	// REVIEW: what go-multihash has done historically is New32, but this doesn't match what the multihash table says, which is 128!  Resolution needed.
	// REVIEW: I have also heard that something in ipfs unixfsv1 uses a murmur hash, but that is yet different than this.  Resolution needed.
	// REVIEW: these bit-twiddling things *are* in fact load-bearing somehow.  If you just return murmur3.New32 without this wrapper type, it produces different results.  Resolution needed.

	// multihash.Register(0x22, func() hash.Hash { return murmur3.New32() })

	// -or-, what previously existed:

	// number := murmur3.Sum32(data)
	// bytes := make([]byte, 4)
	// for i := range bytes {
	// 	bytes[i] = byte(number & 0xff)
	// 	number >>= 8
	// }
	// return bytes, nil
}

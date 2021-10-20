/*
	This package has no purpose except to perform registration of multihashes.

	It is meant to be used as a side-effecting import, e.g.

		import (
			_ "github.com/multiformats/go-multihash/register/murmur3"
		)

	This package registers multihashes for murmur3
*/
package murmur3

import (
	"hash"

	multihash "github.com/multiformats/go-multihash/core"
	"github.com/spaolacci/murmur3"
)

func init() {
	multihash.Register(multihash.MURMUR3X64_64, func() hash.Hash { h := murmur3.New64(); return h })
}

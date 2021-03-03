package multihash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// registry is a simple map which maps a multihash indicator number
// to a standard golang Hash interface.
//
// Multihash indicator numbers are reserved and described in
// https://github.com/multiformats/multicodec/blob/master/table.csv .
// The keys used in this map must match those reservations.
//
// Hashers which are available in the golang stdlib will be registered automatically.
// Others can be added using the Register function.
var registry = make(map[uint64]func() hash.Hash)

// Register adds a new hash to the set available from GetHasher and Sum.
//
// Register has a global effect and should only be used at package init time to avoid data races.
//
// The indicator code should be per the numbers reserved and described in
// https://github.com/multiformats/multicodec/blob/master/table.csv .
//
// If Register is called with the same indicator code more than once, the last call wins.
// In practice, this means that if an application has a strong opinion about what implementation to use for a certain hash
// (e.g., perhaps they want to override the sha256 implementation to use a special hand-rolled assembly variant
// rather than the stdlib one which is registered by default),
// then this can be done by making a Register call with that effect at init time in the application's main package.
// This should have the desired effect because the root of the import tree has its init time effect last.
func Register(indicator uint64, hasherFactory func() hash.Hash) {
	registry[indicator] = hasherFactory
}

// GetHasher returns a new hash.Hash according to the indicator code number provided.
//
// The indicator code should be per the numbers reserved and described in
// https://github.com/multiformats/multicodec/blob/master/table.csv .
//
// The actual hashers available are determined by what has been registered.
// The registry automatically contains those hashers which are available in the golang standard libraries
// (which includes md5, sha1, sha256, sha384, sha512, and the "identity" mulithash, among others).
// Other hash implementations can be made available by using the Register function.
// The 'go-mulithash/register/*' packages can also be imported to gain more common hash functions.
//
// If an error is returned, it will be ErrSumNotSupported.
func GetHasher(indicator uint64) (hash.Hash, error) {
	factory, exists := registry[indicator]
	if !exists {
		return nil, ErrSumNotSupported // REVIEW: it's unfortunate that this error doesn't say what code was missing.  Also "NotSupported" is a bit of a misnomer now.
	}
	return factory(), nil
}

func init() {
	Register(0x00, func() hash.Hash { return &identityMultihash{} })
	Register(0xd5, md5.New)
	Register(0x11, sha1.New)
	Register(0x12, sha256.New)
	Register(0x13, sha512.New)
	Register(0x1f, sha256.New224)
	Register(0x20, sha512.New384)
	Register(0x56, func() hash.Hash { return &doubleSha256{sha256.New()} })
}

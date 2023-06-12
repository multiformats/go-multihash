//go:build !go1.20

package multihash

import "testing"

func mustNotAllocateMore(_ *testing.T, _ float64, f func()) {
	// the compiler isn't able to detect our outlined stack allocation on before
	// 1.20 so let's not test for it. We don't mind if outdated versions are slightly slower.
	f()
}

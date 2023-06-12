//go:build go1.20

package multihash

import "testing"

func mustNotAllocateMore(t *testing.T, n float64, f func()) {
	t.Helper()
	if b := testing.AllocsPerRun(10, f); b > n {
		t.Errorf("it allocated %f max %f !", b, n)
	}
}

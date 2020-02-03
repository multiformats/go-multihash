package multihash

import (
	"crypto/rand"
	"errors"
	"testing"
)

func makeRandomMultihash(t *testing.T) Multihash {
	t.Helper()

	p := make([]byte, 256)
	_, err := rand.Read(p)
	if err != nil {
		t.Fatal(err)
	}

	m, err := Sum(p, SHA3, 4)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestSet(t *testing.T) {
	mhSet := NewSet()

	total := 10
	for i := 0; i < total; i++ {
		mhSet.Add(makeRandomMultihash(t))
	}

	m0 := makeRandomMultihash(t)

	if mhSet.Len() != total {
		t.Error("bad length")
	}

	if mhSet.Has(m0) {
		t.Error("m0 should not be in set")
	}

	mhSet.Add(m0)

	if !mhSet.Has(m0) {
		t.Error("m0 should be in set")
	}

	i := 0
	f := func(m Multihash) error {
		i++
		if i == 3 {
			return errors.New("3")
		}
		return nil
	}

	mhSet.ForEach(f)
	if i != 3 {
		t.Error("forEach should have run 3 times")
	}

	mhSet.Remove(m0)

	if mhSet.Len() != total {
		t.Error("an element should have been removed")
	}

	if mhSet.Has(m0) {
		t.Error("m0 should not be in set")
	}

	if !mhSet.Visit(m0) {
		t.Error("Visit() should return true when new element added")
	}

	all := mhSet.All()
	if len(all) != mhSet.Len() {
		t.Error("All() should return all")
	}
	for _, mh := range all {
		if !mhSet.Has(mh) {
			t.Error("element in All() not in set")
		}
	}
}

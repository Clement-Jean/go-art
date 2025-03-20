package art

import "golang.org/x/text/collate"

type BinaryComparableKey[K nodeKey] interface {
	Transform(K) ([]byte, []byte)
	Restore([]byte) K
}

type AlphabeticalOrderKey[K nodeKey] struct{}

func (ak AlphabeticalOrderKey[K]) Transform(k K) ([]byte, []byte) {
	b := []byte(string(k))
	return b, b
}
func (ak AlphabeticalOrderKey[K]) Restore(b []byte) K { return K(string(b)) }

var _ BinaryComparableKey[string] = AlphabeticalOrderKey[string]{}

type CollationOrderKey[K nodeKey] struct {
	c   *collate.Collator
	buf *collate.Buffer
	src K
}

func (ck CollationOrderKey[K]) Transform(k K) ([]byte, []byte) {
	ck.src = k
	return []byte(string(k)), ck.c.KeyFromString(ck.buf, string(k))
}
func (ck CollationOrderKey[K]) Restore(b []byte) K { return ck.src }

var _ BinaryComparableKey[string] = CollationOrderKey[string]{}

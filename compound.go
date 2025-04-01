package art

func NewCompoundTree[K any, V any](bck BinaryComparableKey[K]) Tree[K, V] {
	return &compoundSortedTree[K, V]{
		bck: bck,
	}
}

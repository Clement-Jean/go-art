package art

func NewSignedBinaryTree[K ints, V any]() Tree[K, V] {
	return &signedSortedTree[K, V]{}
}

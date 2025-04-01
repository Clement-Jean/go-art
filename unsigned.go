package art

func NewUnsignedBinaryTree[K uints, V any]() Tree[K, V] {
	return &unsignedSortedTree[K, V]{}
}

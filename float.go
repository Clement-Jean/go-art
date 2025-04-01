package art

func NewFloatBinaryTree[K floats, V any]() Tree[K, V] {
	return &floatSortedTree[K, V]{}
}

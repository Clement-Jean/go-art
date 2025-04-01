package art

func NewAlphaSortedTree[K chars, V any]() Tree[K, V] {
	return &alphaSortedTree[K, V]{}
}

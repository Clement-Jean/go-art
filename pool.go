package art

import "sync"

// TODO remove undefined and leaf
var nodePools [nodeKindMax]sync.Pool = [nodeKindMax]sync.Pool{
	{New: func() any { panic("shouldn't be possible!") }}, // nodeKindUndefined
	{New: func() any { panic("shouldn't be possible!") }}, // nodeKindLeaf
	{New: func() any { return new(node4) }},               // nodeKind4
	{New: func() any { return new(node16) }},              // nodeKind16
	{New: func() any { return new(node48) }},              // nodeKind48
	{New: func() any { return new(node256) }},             // nodeKind256
}

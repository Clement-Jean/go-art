package art

import "sync"

var nodePools [nodeKindLeaf]sync.Pool = [nodeKindLeaf]sync.Pool{
	{New: func() any { return new(node4) }},   // nodeKind4
	{New: func() any { return new(node16) }},  // nodeKind16
	{New: func() any { return new(node48) }},  // nodeKind48
	{New: func() any { return new(node256) }}, // nodeKind256
}

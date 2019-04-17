package functions

import (
	"sync"

	"honnef.co/go/tools/callgraph"
	"honnef.co/go/tools/callgraph/static"
	"honnef.co/go/tools/ssa"
)

type Description struct {
	Loops []Loop
}

type descriptionEntry struct {
	ready  chan struct{}
	result Description
}

type Descriptions struct {
	CallGraph *callgraph.Graph
	mu        sync.Mutex
	cache     map[*ssa.Function]*descriptionEntry
}

func NewDescriptions(prog *ssa.Program) *Descriptions {
	return &Descriptions{
		CallGraph: static.CallGraph(prog),
		cache:     map[*ssa.Function]*descriptionEntry{},
	}
}

func (d *Descriptions) Get(fn *ssa.Function) Description {
	d.mu.Lock()
	fd := d.cache[fn]
	if fd == nil {
		fd = &descriptionEntry{
			ready: make(chan struct{}),
		}
		d.cache[fn] = fd
		d.mu.Unlock()

		fd.result.Loops = findLoops(fn)

		close(fd.ready)
	} else {
		d.mu.Unlock()
		<-fd.ready
	}
	return fd.result
}
